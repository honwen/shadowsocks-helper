package main

import (
	_ "embed"
	"regexp"
	"sort"
	"strings"

	"github.com/honwen/golibs/cip"
	"github.com/honwen/golibs/domain"
)

func init() {
	originLen := len(topDomainLists)
	for i := 0; i < originLen; i++ {
		topDomainLists = append(topDomainLists, strings.Replace(topDomainLists[i], "https://", "https://ghproxy.com/https://", -1))
		// topDomainLists = append(topDomainLists, strings.Replace(topDomainLists[i], "raw.githubusercontent.com", "raw.githubusercontents.com", -1))
		// topDomainLists = append(topDomainLists, strings.Replace(topDomainLists[i], "raw.githubusercontent.com", "cdn.staticaly.com/gh", -1))
		// topDomainLists = append(topDomainLists, strings.Replace(topDomainLists[i], "raw.githubusercontent.com", "raw.fastgit.org", -1))
	}

	originLen = len(communityDomainLists)
	for i := 0; i < originLen; i++ {
		communityDomainLists = append(communityDomainLists, strings.Replace(communityDomainLists[i], "https://", "https://ghproxy.com/https://", -1))
		// communityDomainLists = append(communityDomainLists, strings.Replace(communityDomainLists[i], "raw.githubusercontent.com", "raw.githubusercontents.com", -1))
		// communityDomainLists = append(communityDomainLists, strings.Replace(communityDomainLists[i], "raw.githubusercontent.com", "cdn.staticaly.com", -1))
		// communityDomainLists = append(communityDomainLists, strings.Replace(communityDomainLists[i], "raw.githubusercontent.com", "raw.fastgit.org", -1))
	}

	originLen = len(officalGFWListURLs)
	for i := 0; i < originLen; i++ {
		officalGFWListURLs = append(officalGFWListURLs, strings.Replace(officalGFWListURLs[i], "https://", "https://ghproxy.com/https://", -1))
		// officalGFWListURLs = append(officalGFWListURLs, strings.Replace(officalGFWListURLs[i], "raw.githubusercontent.com", "raw.githubusercontents.com", -1))
		// officalGFWListURLs = append(officalGFWListURLs, strings.Replace(officalGFWListURLs[i], "raw.githubusercontent.com", "cdn.staticaly.com", -1))
		// officalGFWListURLs = append(officalGFWListURLs, strings.Replace(officalGFWListURLs[i], "raw.githubusercontent.com", "raw.fastgit.org", -1))
	}

}

var topDomainLists = []string{
	`https://raw.githubusercontent.com/Loyalsoldier/clash-rules/release/google.txt`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/google`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/youtube`,
	// `https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/google-ads`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/google-registry`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/google-scholar`,
}

var communityDomainLists = []string{
	`https://raw.githubusercontent.com/pexcn/daily/gh-pages/gfwlist/gfwlist.txt`,
	`https://raw.githubusercontent.com/pexcn/daily-extras/master/gfwlist-extras.txt`,
	`https://raw.githubusercontent.com/Loyalsoldier/surge-rules/release/gfw.txt`,
	`https://raw.githubusercontent.com/poctopus/gfwlist-plus/master/gfwlist-plus.txt`,

	`https://raw.githubusercontent.com/hy-cs/subTemplate/master/rules/Proxy.list`,
	`https://raw.githubusercontent.com/blackmatrix7/ios_rule_script/master/rule/Surge/Global/Global_Domain.list`,

	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/openai`,

	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/android`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/blogspot`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/dart`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/fastlane`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/flutter`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/golang`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/kaggle`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/polymer`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/v8`,

	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/github`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/stackexchange`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/coursera`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/ted`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/akamai`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/cdn77`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/cloudflare`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/fastly`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/imgix`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/limelight`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/pagecdn`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/quantil`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/stackpath`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/clubhouse`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/discord`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/freenode`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/intercom`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/line`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/matrix`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/signal`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/slack`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/telegram`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/rarbg`,

	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/category-porn`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/bongacams`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/bttzyw`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/ehentai`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/heyzo`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/jable`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/javbus`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/javfinder`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/javwide`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/livejasmin`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/metart`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/pornhub`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/scoreland`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/uu-chat`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/xhamster`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/xnxx`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/xvideos`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/youjizz`,

	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/category-pt`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/category-media`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/category-novel`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/category-public-tracker`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/category-ipfs`,
	// `https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/category-dev`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/category-anticensorship`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/category-vpnservices`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/928plus`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/boslife`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/cylink`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/dlercloud`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/duyaoss`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/hitun`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/justmysocks`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/n3ro`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/nexitally`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/paofuyun`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/shadowsockscom`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/squirrelvpn`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/ssrcloud`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/surflite`,
}

var officalGFWListURLs = []string{
	`https://raw.githubusercontent.com/Loukky/gfwlist-by-loukky/master/gfwlist.txt`,
	`https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt`,
	// `https://raw.githubusercontent.com/YW5vbnltb3Vz/domain-list-community/release/gfwlist.txt`,
}

const regxIP = `(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)`

func filter(s []string, filterFunc func(string) bool) (after []string, fliterd []string) {
	for i := range s {
		if filterFunc(s[i]) {
			fliterd = append(fliterd, s[i])
		} else {
			after = append(after, s[i])
		}
	}
	return
}

func customsSort(domains []string, topDomains []string) []string {
	domains = domain.Sort(domains)
	regxIPv4 := regexp.MustCompile(cip.RegxIPv4)
	domains, ips := filter(domains, func(s string) bool {
		return regxIPv4.MatchString(s)
	})
	filterStr := "|" + strings.Join(topDomains, "|") + "|"
	domains, extraTops := filter(domains, func(s string) bool {
		if strings.HasPrefix(s, "google.") || strings.HasPrefix(s, "blogspot.") {
			return true
		}
		return strings.Contains(filterStr, "."+s+"|")
	})
	topDomains = domain.Sort(append(topDomains, extraTops...))

	sort.Strings(ips)
	sort.Strings(topDomains)
	domains = append(ips, append(topDomains, domains...)...)
	return domains
}

//go:embed gfwlist.init.txt
var initList string
