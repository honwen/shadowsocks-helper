package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/Workiva/go-datastructures/set"
)

const (
	officalGFWListURL   = `https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt`
	officalGoogleDomain = `https://www.google.com/supported_domains`
)

const regxIP = `(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)`

// GFWList contain GFWed domains
type GFWList []string

// ParseGFWList Parse GFWList
func ParseGFWList(base64GFWList, extraList string) (gfwlist GFWList, err error) {
	gfwlist, err = autoProxy2Domains(base64GFWList, extraList)
	return
}

// GenDnsmasqIPset with specify ipset
func (gfwList *GFWList) GenDnsmasqIPset(ipset string) string {
	var (
		str bytes.Buffer
		w   io.Writer = &str
	)
	for _, s := range *gfwList {
		io.WriteString(w, fmt.Sprintf("ipset=/%s/%s\n", s, ipset))
	}
	return str.String()
}

// GenDnsmasqServer with specify server
func (gfwList *GFWList) GenDnsmasqServer(dns string) string {
	var (
		str bytes.Buffer
		w   io.Writer = &str
	)
	for _, s := range *gfwList {
		io.WriteString(w, fmt.Sprintf("server=/%s/%s\n", s, dns))
	}
	return str.String()
}

func (gfwList GFWList) String() string {
	var (
		str bytes.Buffer
		w   io.Writer = &str
	)
	for _, s := range gfwList {
		io.WriteString(w, s)
		io.WriteString(w, "\n")
	}
	return str.String()
}

func autoProxy2Domains(base64Str, extraList string) (Domains []string, err error) {
	gfwlistDEC, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}
	// ---------------------------------gfwlist---------------------------------
	sites := set.New()
	gfwlistDEC = append(gfwlistDEC, []byte(initList)...)
	gfwlistDEC = append(gfwlistDEC, []byte(extraList)...)
	scanner := bufio.NewScanner(bytes.NewReader(gfwlistDEC))
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())

		if s == "" ||
			strings.HasPrefix(s, "[") ||
			strings.HasPrefix(s, "!") ||
			strings.HasPrefix(s, "||!") ||
			strings.HasPrefix(s, "@@") {
			continue
		}

		if unURL, err := url.QueryUnescape(s); err == nil {
			s = unURL
		}

		switch {
		case strings.HasPrefix(s, "||"):
			site := strings.Split(s[2:], "/")[0]
			switch {
			case strings.Contains(site, "*."):
				parts := strings.Split(site, "*.")
				site = parts[len(parts)-1]
			case strings.HasPrefix(site, "*"):
				parts := strings.SplitN(site, ".", 2)
				site = parts[len(parts)-1]
			}
			sites.Add(site)
		case strings.HasPrefix(s, "|https://"):
			fallthrough
		case strings.HasPrefix(s, "|http://"):
			if u, err := url.Parse(s[1:]); err == nil {
				site := u.Host
				switch {
				case strings.Contains(site, "*."):
					parts := strings.Split(site, "*.")
					site = parts[len(parts)-1]
				case strings.HasPrefix(site, "*"):
					parts := strings.SplitN(site, ".", 2)
					site = parts[len(parts)-1]
				}
				sites.Add(site)
			}
		case strings.HasPrefix(s, "."):
			site := strings.Split(strings.Split(s[1:], "/")[0], "*")[0]
			sites.Add(site)
		case !strings.ContainsAny(s, "*"):
			site := strings.Split(s, "/")[0]
			if regexp.MustCompile(`^[a-zA-Z0-9\.\_\-]+$`).MatchString(site) {
				sites.Add(site)
			}
		}
	}
	if err = scanner.Err(); err != nil {
		return
	}
	// ---------------------------------gfwlist---------------------------------
	for _, site := range sites.Flatten() {
		if strings.ContainsAny(site.(string), ".") {
			Domains = append(Domains, site.(string))
		}
	}

	// ----------------------------------sort-----------------------------------
	regexpISIP := regexp.MustCompile(regxIP)
	spDomains := []string{}
	for i := 0; i < len(Domains); {
		if strings.HasPrefix(Domains[i], "google.") || strings.HasPrefix(Domains[i], "blogspot.") ||
			regexpISIP.MatchString(Domains[i]) {
			spDomains = append(spDomains, Domains[i])
			Domains = append(Domains[:i], Domains[i+1:]...)
		} else {
			Domains[i] = reverseString(Domains[i])
			i++
		}
	}
	sort.Strings(spDomains)
	sort.Strings(Domains)
	for i := range Domains {
		Domains[i] = reverseString(Domains[i])
	}
	Domains = append(spDomains, Domains...)
	return
}

func reverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

const initList = `
otcbtc.com
wikipedia.org
whatsapp.com
whatsapp.net
twimg.edgesuite.net
`
