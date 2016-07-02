package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"io"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

const (
	officalGFWListURL = `https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt`
)

type GFWList struct {
	Domains []string
}

func ParseGFWList(failSafeSS []string) (*GFWList, error) {
	base64GFWList, err := wGetWithSSFailsafe(officalGFWListURL, failSafeSS)
	if err != nil {
		return nil, err
	}
	domains, err := autoProxy2Domains(base64GFWList)
	if err != nil {
		return nil, err
	}
	gfwList = &GFWList{}
	gfwList.Domains = domains
	return gfwList, nil
}

func (gfwList *GFWList) GenDnsmasqIPset(ipset string) string {
	var str bytes.Buffer
	var w io.Writer = &str
	for _, s := range gfwList.Domains {
		io.WriteString(w, "ipset=/")
		io.WriteString(w, s)
		io.WriteString(w, "/")
		io.WriteString(w, ipset)
		io.WriteString(w, "\n")
	}
	return str.String()
}

func (gfwList *GFWList) GenDnsmasqServer(dns string) string {
	var str bytes.Buffer
	var w io.Writer = &str
	for _, s := range gfwList.Domains {
		io.WriteString(w, "server=/")
		io.WriteString(w, s)
		io.WriteString(w, "/")
		io.WriteString(w, dns)
		io.WriteString(w, "\n")
	}
	return str.String()
}

func (gfwList GFWList) String() string {
	var str bytes.Buffer
	var w io.Writer = &str
	for _, s := range gfwList.Domains {
		io.WriteString(w, s)
		io.WriteString(w, "\n")
	}
	return str.String()
}

func autoProxy2Domains(base64Str string) (Domains []string, err error) {
	gfwlistDEC, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}
	// ---------------------------------gfwlist---------------------------------
	sites := make(map[string]struct{}, 0)
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
			sites[site] = struct{}{}
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
				sites[site] = struct{}{}
			}
		case strings.HasPrefix(s, "."):
			site := strings.Split(strings.Split(s[1:], "/")[0], "*")[0]
			sites[site] = struct{}{}
		case !strings.ContainsAny(s, "*"):
			site := strings.Split(s, "/")[0]
			if regexp.MustCompile(`^[a-zA-Z0-9\.\_\-]+$`).MatchString(site) {
				sites[site] = struct{}{}
			}
		}
	}
	if err = scanner.Err(); err != nil {
		return
	}
	// ---------------------------------gfwlist---------------------------------
	for s := range sites {
		Domains = append(Domains, s)
	}
	sort.Strings(Domains)
	return
}
