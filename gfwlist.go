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
	// officalGFWListURL = `https://raw.sevencdn.com/Loukky/gfwlist-by-loukky/master/gfwlist.txt`
	officalGFWListURL = `https://raw.githubusercontent.com/Loukky/gfwlist-by-loukky/master/gfwlist.txt`
	// officalGFWListURL   = `https://raw.sevencdn.com/gfwlist/gfwlist/master/gfwlist.txt`
	officalGoogleDomain = `https://www.google.com/supported_domains`
)

const regxIP = `(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)`

// GFWList contain GFWed domains
type GFWList []string

// ParseGFWList Parse GFWList
func ParseGFWList(base64GFWList, extraList string, noInitList bool) (gfwlist GFWList, err error) {
	gfwlist, err = autoProxy2Domains(base64GFWList, extraList, noInitList)
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

func autoProxy2Domains(base64Str, extraList string, noInitList bool) (Domains []string, err error) {
	gfwlistDEC, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}
	// ---------------------------------gfwlist---------------------------------
	sites := set.New()
	if !noInitList {
		gfwlistDEC = append(gfwlistDEC, []byte(initList)...)
		gfwlistDEC = append(gfwlistDEC, []byte(extraInitList)...)
	}
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
			sites.Add(trimSite(site))
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
				sites.Add(trimSite(site))
			}
		case strings.HasPrefix(s, "."):
			site := strings.Split(strings.Split(s[1:], "/")[0], "*")[0]
			if !strings.HasSuffix(site, ".") {
				sites.Add(trimSite(site))
			}
		case !strings.ContainsAny(s, "*"):
			site := strings.Split(s, "/")[0]
			if regexp.MustCompile(`^[a-zA-Z0-9\.\_\-]+$`).MatchString(site) {
				sites.Add(trimSite(site))
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
	Domains = append(spDomains, tide(Domains)...)
	return
}

func trimSite(s string) string {
	site := strings.TrimSpace(s)
	switch {
	case strings.HasSuffix(site, `USA`):
		site = site[0:strings.Index(site, `USA`)]
	case strings.HasSuffix(site, `--`):
		site = site[0:strings.Index(site, `--`)]
	}
	return site
}

const initList = `
000webhost.com
0zmedia.com
adobe.com
akamaiedge.net
akamai.net
amazonaws.com
ampproject.org
androidapksfree.com
apkcombo.com
apple.news
apps.evozi.com
aws.amazon.com
azureedge.net
backpackers.com.tw
bing.com
bitfinex.com
buzzfeed.com
china-internet-exchange.com
clockwise.ee
cloudfront.net
cnbc.com
coindesk.com
coinsquare.io
cryptocompare.com
cryptographyengineering.com
dailysabah.com
dlercloud.com
dropboxstatic.com
e-classical.com.tw
edgekey.net
edgesuite.net
edx-cdn.org
etherscan.io
eurecom.fr
facebook.com
fastly.net
firstpost.com
fontawesome.com
gdax.com
ggpht.com
github.com
github.io
githubusercontent.com
gnews.org
googleapis.com
googlecode.com
google.com
google.com.hk
googledrive.com
googlefiber.net
googlepages.com
googlesource.com
googleusercontent.com
googlevideo.com
greasyfork.org
gstatic.com
hkmap.live
invidio.us
kknews.cc
lihkg.com
limbopro.xyz
mastodon.xyz
maying.co
mergersandinquisitions.com
metafilter.com
nbcnews.com
newyorker.com
nexitally.com
nexitcore.com
nutaq.com
openairinterface.org
otcbtc.com
oxfordscholarship.com
pornstarbyface.com
project-syndicate.org
sinoinsider.com
skype.com
sleazyfork.org
streamate.com
sublimetext.com
t66y.com
telegra.ph
textnow.com
textnow.me
theatlantic.com
timtales.com
trouter.io
twimg.edgesuite.net
twitter.com
uploaded.net
userstyles.org
v2fly.org
v2raycn.com
weverse.io
whatsapp.com
whatsapp.net
wikipedia.org
wsj.net
yande.re
ycombinator.com
youtu.be
youtube.com
ytimg.com
ees.elsevier.com
archiveofourown.org
advertisercommunity.com
uptodown.com
ifixit.com
wireguard.com
clubhouseapi.com
joinclubhouse.com
edia.prod.mdn.mozit.cloud
eveloper.mozilla.org
hinapower.csis.org
dn.mozillademos.org
tm.tnt-ea.com
isco.evergage.com
pstore.corpmerchandise.com
dn77.scoreuniverse.com
mperial.insendi.com
urofensk-prod-env.eu-west-1.elasticbeanstalk.com
andom.zendesk.com
oss-cn-hongkong.aliyuncs.com
zuremarketplace.microsoft.com
zure.microsoft.com
eveloper.microsoft.com
stripst.com
lizzard.nefficient.co.kr
nteractive-examples.mdn.mozilla.net
omtrdc.net
cloudapp.net
akadns.net
rightcove.imgix.net
`
