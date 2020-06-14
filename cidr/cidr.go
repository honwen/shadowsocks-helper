package cidr

import (
	"log"
	"regexp"
	"strings"

	goquery "github.com/PuerkitoBio/goquery"
)

const (
	googleAS = "https://whois.ipip.net/AS15169"
)

var (
	regexCDIR = regexp.MustCompile(`^(?:(?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}(?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\/([1-9]|[1-2]\d|3[0-2])$`)
)

const googleExtra = `
194.9.24.64/27
4.35.153.192/27
203.210.7.136/29
203.210.6.136/29
64.15.112.0/20
203.208.32.0/19
`

func Google() (cidrs []string) {
	docEachCurrency, err := goquery.NewDocument(googleAS)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range strings.Split(googleExtra, "\n") {
		if regexCDIR.MatchString(v) {
			cidrs = append(cidrs, v)
		}
	}
	docEachCurrency.Find(`td>a`).Each(func(i int, s *goquery.Selection) {
		v := s.Text()
		if regexCDIR.MatchString(v) {
			cidrs = append(cidrs, v)
		}
	})
	// log.Printf("%+v", cidrs)
	return
}
