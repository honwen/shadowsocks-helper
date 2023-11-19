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
		topDomainLists = append(topDomainLists, strings.Replace(topDomainLists[i], "https://", "https://gh-proxy.com/https://", -1))
		// topDomainLists = append(topDomainLists, strings.Replace(topDomainLists[i], "raw.githubusercontent.com", "raw.githubusercontents.com", -1))
		// topDomainLists = append(topDomainLists, strings.Replace(topDomainLists[i], "raw.githubusercontent.com", "cdn.staticaly.com/gh", -1))
		// topDomainLists = append(topDomainLists, strings.Replace(topDomainLists[i], "raw.githubusercontent.com", "raw.fastgit.org", -1))
	}

	originLen = len(communityDomainLists)
	for i := 0; i < originLen; i++ {
		communityDomainLists = append(communityDomainLists, strings.Replace(communityDomainLists[i], "https://", "https://gh-proxy.com/https://", -1))
		// communityDomainLists = append(communityDomainLists, strings.Replace(communityDomainLists[i], "raw.githubusercontent.com", "raw.githubusercontents.com", -1))
		// communityDomainLists = append(communityDomainLists, strings.Replace(communityDomainLists[i], "raw.githubusercontent.com", "cdn.staticaly.com", -1))
		// communityDomainLists = append(communityDomainLists, strings.Replace(communityDomainLists[i], "raw.githubusercontent.com", "raw.fastgit.org", -1))
	}

	originLen = len(officalGFWListURLs)
	for i := 0; i < originLen; i++ {
		officalGFWListURLs = append(officalGFWListURLs, strings.Replace(officalGFWListURLs[i], "https://", "https://gh-proxy.com/https://", -1))
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
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/whatsapp`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/matrix`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/signal`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/slack`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/telegram`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/rarbg`,

	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/category-porn`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/18comic`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/54647`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/anon-v`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/avmoo`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/bdsmhub`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/bilibili2`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/boboporn`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/bongacams`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/boylove`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/bttzyw`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/camwhores`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/cavporn`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/chatwhores`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/clips4sale`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/dlsite`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/dmm-porn`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/ehentai`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/erolabs`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/fans66`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/fansta`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/hentaivn`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/heyzo`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/hooligapps`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/illusion`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/illusion-nonofficial`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/jable`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/javbus`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/javcc`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/javdb`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/japonx`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/javwide`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/jkf`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/johren`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/justav`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/kemono`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/konachan`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/kubakuba`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/lethalhardcore`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/lisiku`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/metart`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/mindgeek-porn`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/moxing`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/netflav`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/nudevista`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/picacg`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/playboy`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/pornpros`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/shireyishunjian`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/smtiaojiaoshi`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/theporndude`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/thescoregroup`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/tokyo-toshokan`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/truyen-hentai`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/uu-chat`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/xhamster`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/xingkongwuxianmedia`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/xnxx`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/xvideos`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/youjizz`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/yunlaopo`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/zhimeishe`,

	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/category-pt`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/category-media`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/6park`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/9news`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/9to5`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/afp`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/aljazeera`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/ap`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/appledaily`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/asahi`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/bbc`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/bloomberg`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/boxun`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/c-span`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/cabletv`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/cbs`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/cnbc`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/cnn`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/dailymail`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/dw`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/economist`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/falungong`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/forbes`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/fox`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/ft`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/ftv`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/gannett`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/globalvoices`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/guo`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/hketgroup`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/hkopentv`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/huffpost`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/infowars`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/insider`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/kyodonews`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/londonreal`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/ltn`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/mainichi`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/manorama`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/manoto`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/matters`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/mediachinesegroup`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/medium`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/myradio`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/nbcuniversal`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/neowin`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/newscorp`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/newsmax`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/nikkan-gendai`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/nikkei`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/now`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/nytimes`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/oan`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/realclear`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/rferl`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/rthk`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/sankei`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/scmp`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/setn`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/singtaonewscorp`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/sky`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/theguardian`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/theinitium`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/thetype`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/thomsonreuters`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/tokyo-sports`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/tvb`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/tver`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/udn`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/voanews`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/voxmedia`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/wantmedia`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/watchout`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/ycombinator`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/yomiuri`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/zaobao`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/zee`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/category-novel`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/category-public-tracker`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/category-ipfs`,
	// `https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/category-dev`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/category-anticensorship`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/cloudflare-ipfs`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/ceno`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/citizenlab`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/epicbrowser`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/greatfire`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/haveibeenpwned`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/lantern`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/paskoocheh`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/softether`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/torproject`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/v2ray`,
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/vpngate`,
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
	`https://raw.githubusercontent.com/v2fly/domain-list-community/master/data/vilavpn`,
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
