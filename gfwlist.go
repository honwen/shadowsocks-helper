package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"index/suffixarray"
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
	gfwlistDEC = append(gfwlistDEC, []byte(extraInitList)...)
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
			if !strings.HasSuffix(site, ".") {
				sites.Add(site)
			}
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
	coreDomains := []string{}
	idx := suffixarray.New([]byte(strings.Join(Domains, "")))
	for i := range Domains {
		offsets := idx.Lookup([]byte("."+Domains[i]), 1)
		if len(offsets) > 0 {
			coreDomains = append(coreDomains, Domains[i])
		}
	}
	tidedDomains := []string{}
	for i := range Domains {
		cnt := 0
		idx = suffixarray.New([]byte(Domains[i]))
		for ii := range coreDomains {
			offsets := idx.Lookup([]byte("."+coreDomains[ii]), 1)
			cnt += len(offsets)
		}
		if cnt == 0 {
			tidedDomains = append(tidedDomains, Domains[i])
		}
	}
	for i := range tidedDomains {
		tidedDomains[i] = reverseString(tidedDomains[i])
	}
	sort.Strings(tidedDomains)
	for i := range tidedDomains {
		tidedDomains[i] = reverseString(tidedDomains[i])
	}
	Domains = append(spDomains, tidedDomains...)
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
000webhost.com
0zmedia.com
adobe.com
akamaiedge.net
akamaihd.net
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
`

// curl -skSL https://raw.githubusercontent.com/lhie1/Rules/master/Shadowrocket/Complete.conf | grep -E '^DOMAIN-SUFFIX.*,PROXY$' | awk -F',' '{print $2}'
const extraInitList = `1e100.net
2mdn.net
2o7.net
4everProxy.com
4shared.com
4sqi.net
9to5mac.com
abc.xyz
abema.io
abema.tv
abpchina.org
accountkit.com
adblockplus.org
adobe.com
adobedtm.com
aerisapi.com
airtable.com
aiv-cdn.net
akamaihd.net
akamai.net
akamaized.net
alfredapp.com
allconnected.co
amazonaws.com
amazonaws.co.uk
amazon.co
amazon.co.jp
amazon.com
amazonvideo.com
ameba.jp
ampproject.com
ampproject.net
ampproject.org
anaconda.com
ancsconf.org
android.com
androidify.com
android-x86.org
angularjs.org
anthonycalzadilla.com
aolcdn.com
aol.com
api.amplitude.com
apigee.com
api.mixpanel.com
api.termius.com
api.tiktokv.com
api.urbandictionary.com
apk-dl.com
apkpure.com
appdownloader.net
apple-dns.net
app-measurement.com
appshopper.com
app.smartmailcloud.com
appspot.com
arcgis.com
archive.is
archive.org
archives.gov
armorgames.com
aspnetcdn.com
async.be
att.com
autodraw.com
avgle.com
awsstatic.com
azure.com
azureedge.net
azurewebsites.net
badoo.com
bahamut.com.tw
bamgrid.com
bandisoft.com
bandwagonhost.com
bbc.co.uk
bbci.co.uk
bbtoystore.com
betvictor.com
bigsound.org
bing.com
bing.net
bintray.com
bitbucket.org
bitcointalk.org
bit.com
bit.do
bit.ly
bitshare.com
bkrtx.com
blogblog.com
blogcdn.com
blog.com
blogger.com
bloglovin.com
blogsmithmedia.com
blogspot.hk
bloomberg.cn
bloomberg.com
books.com.tw
boomtrain.com
botanwang.com
box.com
box.net
boxun.com
cachefly.net
cbc.ca
cdn.angruo.com
cdn.segment.com
cdnst.net
celestrak.com
census.gov
certificate-transparency.org
chinadigitaltimes.net
chinatimes.com
chromecast.com
chrome.com
chromeexperiments.com
chromercise.com
chromestatus.com
chromium.org
cl.ly
cloudflare.com
cloudfront.net
cloudgarage.jp
cloudmagic.com
cmail19.com
cnet.com
cnn.com
cocoapods.org
comodoca.com
content.office.net
conviva.com
creativelab5.com
crittercism.com
culturalspot.org
culturedcode.com
cygames.jp
d151l6v8er5bdm.cloudfront.net
d3c7rimkq79yfu.cloudfront.net
danilo.to
daolan.net
dartlang.org
data-vocabulary.org
dayone.me
dazn-api.com
dazn.com
dazndn.com
db.tt
dcmilitary.com
deja.com
demdex.net
deskconnect.com
digisfera.com
digitaltrends.com
disconnect.me
disney.demdex.net
disneyplus.com
disney-plus.net
disq.us
disquscdn.com
disqus.com
dmm.co.jp
dnsimple.com
docker.com
doub.io
dowjones.com
d.pr
dribbble.com
droplr.com
dssott.com
duckduckgo.com
dueapp.com
dw.com
easybib.com
economist.com
edgecastcdn.net
edgedatg.com
edgekey.net
edgesuite.net
encoretvb.com
engadget.com
entrust.net
eurekavpt.com
evernote.com
execute-api.ap-southeast-1.amazonaws.com
execute-api.us-east-1.amazonaws.com
extmatrix.com
eyny.com
fabric.io
fanatical.com
fast.com
fastly.net
fastmail.com
fbcdn.net
fb.com
fb.me
fbsbx.com
fc2.com
feedburner.com
feedly.com
feedsportal.com
fiftythree.com
firebaseio.com
flexibits.com
flickr.com
flipboard.com
flipkart.com
flitto.com
flurry.com
fox.com
foxdcg.com
foxplus.com
freeopenProxy.com
fubo.tv
fullstory.com
fzlm.net
gabia.net
gamer.com.tw
garena.com
gate.hockeyapp.net
g.co
gcr.io
geni.us
getcloudapp.com
getfoxyProxy.org
get.how
getlantern.org
getmdl.io
getpricetag.com
gfw.press
gfx.ms
ggpht.com
ghostnoteapp.com
gitbook.com
git.io
globalsign.com
gmail.com
gmocloud.com
gmodules.com
go.com
godaddy.com
go.jp
golang.org
gongm.in
goodreaders.com
goodreads.com
goo.gl
googlevideo.com
gosetsuden.jp
gravatar.com
gstatic.cn
gstatic.com
gunsamerica.com
gvt0.com
gvt1.com
gvt2.com
gwtproject.org
happyon.jp
hayabusa.io
hboasia.com
hbo.com
hbogoasia.com
hbogoasia.hk
hbogo.com
hbonow.com
helpshift.com
heroku.com
hinet.net
hockeyapp.net
homedepot.com
hootsuite.com
howtoforge.com
html5rocks.com
hulu.com
huluim.com
hulu.jp
hulustream.com
humblebundle.com
iam.soy
i-cable.com
icoco.com
icons8.com
ift.tt
ifttt.com
imageshack.us
img.ly
imgur.com
imore.com
indazn.com
ingress.com
insder.co
instapaper.com
instructables.com
intercom.io
io.io
ipaddress.com
ipn.li
is.gd
ishowsapp.com
issuu.com
itgonglun.com
itun.es
ixquick.com
japonx.com
japonx.net
japonx.tv
japonx.vip
japronx.com
japronx.net
japronx.tv
japronx.vip
javbus.com
j.mp
joox.com
jshint.com
js.revsci.net
jtvnw.net
justgetflux.com
kakaocdn.net
kakao.co.kr
kakao.com
kat.cr
kenengba.com
keyhole.com
kfs.io
kik.com
kkbox.com
kkbox.com.tw
klip.me
kobobooks.com
kobo.com
leancloud.com
leetcode.com
lhie1.com
libsyn.com
licdn.com
lightboxcdn.com
like.com
line-apps.com
line-cdn.net
lin.ee
line.me
line.naver.jp
line-scdn.net
linetv.tw
linkedin.com
linode.com
lithium.com
littlehj.com
live.com
livefilestore.com
live.net
llnwd.net
localnetwork.uop
logmein.com
macid.co
macromedia.com
macrumors.com
madewithcode.com
manhuaren.com
marketwatch.com
mashable.com
material.io
mathjax.org
maven.org
medium.com
mega.co.nz
mega.nz
megaupload.com
messenger.com
microsofttranslator.com
mindnode.com
mlssoccer.com
mobile01.com
modmyi.com
moves-export.com
mp3buscador.com
msedge.net
mycnnews.com
myfontastic.com
mytvsuper.com
name.com
nasa.gov
ndr.de
netdna-cdn.com
netflix.com
netflix.net
netmarble.com
newipnow.com
nextmedia.com
nflxext.com
nflximg.com
nflximg.net
nflxso.net
nflxvideo.net
nhncorp.jp
nianticlabs.com
nih.gov
nintendo.com
nintendo.net
notion.so
novafile.com
nrk.no
nsstatic.net
nssurge.com
nyt.com
nytimes.com
nytimg.com
nytstyle.com
office365.com
office.com
office.net
omnigroup.com
onenote.com
ooklaserver.net
ooyala.com
openvpn.net
openwrt.org
optimizely.com
orkut.com
osha.gov
osxdaily.com
overcast.fm
ow.ly
paddleapi.com
paddle.com
pandora.com
panoramio.com
parallels.com
parse.com
pbs.org
pdfexpert.com
periscope.tv
phncdn.com
phprcdn.com
piaotian.net
picacomic.com
picasaweb.com
pinboard.in
ping.pe
pinimg.com
pinterest.com
pixelmator.com
pixiv.net
pixnet.net
playpcesor.com
playstation.com
playstation.net
playstationnetwork.com
pokemon.com
polymer-project.org
popo.tw
pornhub.com
pornhubpremium.com
potato.im
prfct.co
proxfree.com
psiphon3.com
ptt.cc
pubnub.com
pubu.com.tw
puffinbrowser.com
pushbullet.com
pushwoosh.com
pximg.net
quoracdn.net
quora.com
readingtimes.com.tw
readmoo.com
recaptcha.net
redd.it
reddit.com
redditmedia.com
reuters.com
rfi.fr
rileyguide.com
rime.im
rsf.org
rthk.hk
scdn.co
sciencedaily.com
sciencemag.org
scribd.com
search.com
servebom.com
sexinsex.net
sfx.ms
shadowsocks.org
shadowverse.jp
sharethis.com
shazam.com
shutterstock.com
sidelinesnews.com
simplenote.com
simp.ly
sketchappsources.com
skype.com
slack.com
slack-edge.com
slack-msgs.com
slideshare.net
smartdnsproxy.com
smartmailcloud.com
smh.com.au
snapchat.com
sndcdn.com
sockslist.net
sony.com
sonyentertainmentnetwork.com
soundcloud.com
sourceforge.net
sowers.org.hk
speedsmart.net
spike.com
spoti.fi
squarespace.com
ssa.gov
sstatic.net
stackoverflow.com
starp2p.com
startpage.com
state.gov
staticflickr.com
steamcommunity.com
steampowered.com
steamstatic.com
st.luluku.pw
storify.com
stumbleupon.com
sugarsync.com
supermariorun.com
surfeasy.com.au
surge.run
surrenderat20.net
sydneytoday.com
symauth.com
symcb.com
symcd.com
t66y.com
tablesgenerator.com
tabtter.jp
talk853.com
talkboxapp.com
talkonly.net
tapbots.com
tapbots.net
t.co
tdesktop.com
techcrunch.com
technorati.com
techsmith.com
teddysun.com
telegram.me
telegram.org
telegra.ph
tensorflow.org
textnow.me
thebobs.com
theinitium.com
thepiratebay.org
theplatform.com
theverge.com
thewgo.org
tiltbrush.com
time.com
timeinc.net
tinder.com
tiny.cc
tinychat.com
tinypic.com
tmblr.co
t.me
todoist.com
togetter.com
toggleable.com
tokyocn.com
tomshardware.com
torcn.com
torproject.org
torrentprivacy.com
torrentproject.se
torrentz.eu
traffichaus.com
trakt.tv
transparency.org
trello.com
trendsmap.com
trulyergonomic.com
trustasiassl.com
tt-rss.org
ttvnw.net
tumblr.co
tumblr.com
turbobit.net
tv.apple.com
tvb.com
tv.com
tweetdeck.com
tweetmarker.net
twimg.co
twimg.com
twitch.tv
twitpic.com
twitthat.com
twtkr.com
twttr.com
txmblr.com
typcn.com
typekit.net
typography.com
ubertags.com
ublock.org
ubnt.com
uchicago.edu
udn.com
ugo.com
uhdwallpapers.org
ulyssesapp.com
unblockdmm.com
unblocksites.co
unfiltered.news
unpo.org
unsplash.com
untraceable.us
uploaded.net
uplynk.com
uProxy.org
upwork.com
urchin.com
urlparser.com
usertrust.com
usgs.gov
usma.edu
uspto.gov
us.to
ustream.tv
v2ex.co
v2ex.com
v2ray.com
van001.com
vanpeople.com
vansky.com
vbstatic.co
venchina.com
venturebeat.com
veoh.com
verizonwireless.com
v.gd
viber.com
videomega.tv
vidinfo.org
vid.me
vimeocdn.com
vimeo.com
vimperator.org
vine.co
visibletweets.com
viu.com
viu.now.com
viu.tv
vivaldi.com
voachinese.com
vocativ.com
vox-cdn.com
vpnaccount.org
vpnbook.com
vpngate.net
vsco.co
vultr.com
vzw.com
w3schools.com
wattpad.com
waveprotocol.org
web2project.net
webfreer.com
weblagu.com
webmproject.org
webrtc.org
websnapr.com
webtype.com
webwarper.net
wenxuecity.com
westca.com
westpoint.edu
whatbrowser.org
wikileaks-forum.com
wikileaks.info
wikileaks.org
wikimedia.org
wikipedia.com
wikipedia.org
wn.com
wordpress.com
w.org
workflow.is
workflowy.com
worldcat.org
wow.com
wp.com
wsj.com
wsj.net
wwitv.com
xanga.com
xbox.com
xboxlive.com
xclient.info
xda-developers.com
xeeno.com
xiti.com
xteko.com
xuite.net
xvideos.com
yahooapis.com
yahoo.com
yasni.co.uk
yastatic.net
yeeyi.com
yesasia.com
yes-news.com
yidio.com
yimg.com
ying.com
yorkbbs.ca
youmaker.com
yourlisten.com
youtu.be
yoyo.org
ytimg.com
zacebook.com
zalmos.com
zaobao.com.sg
zeutch.com
zynamics.com
`
