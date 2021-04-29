package main

import (
	"index/suffixarray"
	"sort"
	"strings"
)

func reverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

func tide(domains []string) (tidedDomains []string) {
	sort.Strings(domains)
	for i := range domains {
		domains[i] = reverseString(domains[i])
	}
	coredomains := []string{}
	idx := suffixarray.New([]byte(strings.Join(domains, "")))
	for i := range domains {
		offsets := idx.Lookup([]byte("."+domains[i]), 1)
		if len(offsets) > 0 {
			coredomains = append(coredomains, domains[i])
		}
	}
	for i := range domains {
		cnt := 0
		idx = suffixarray.New([]byte(domains[i]))
		for ii := range coredomains {
			offsets := idx.Lookup([]byte("."+coredomains[ii]), 1)
			cnt += len(offsets)
		}
		if cnt == 0 {
			tidedDomains = append(tidedDomains, domains[i])
		}
	}
	for i := range tidedDomains {
		tidedDomains[i] = reverseString(tidedDomains[i])
	}
	sort.Strings(tidedDomains)
	for i := range tidedDomains {
		tidedDomains[i] = reverseString(tidedDomains[i])
	}
	return
}
