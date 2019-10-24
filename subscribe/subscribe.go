package subscribe

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/nadoo/glider/proxy"

	// Register proxys
	_ "github.com/nadoo/glider/proxy/ssr"
)

var (
	reSSR     = regexp.MustCompile(`ssr://(\S+)`)
	reSSRInfo = regexp.MustCompile(`(\S+):(\S+):(\S+):(\S+):(\S+):([^/]+)(/?\S+)?`)
)

func base64Decode(str string, enc *base64.Encoding) string {
	rawBytes, err := enc.DecodeString(str)
	if err != nil {
		// handle err
		log.Fatal(str, err)
	}
	return string(rawBytes)
}

// URL2URIs gen SSR-URIs
func URL2URIs(uri string) (ssrURIs, ssrRemarks []string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		// handle err
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		// handle err
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle err
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	rawData, err := base64.StdEncoding.DecodeString(bodyString)
	if err != nil {
		// handle err
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(rawData))
	for scanner.Scan() {
		raw := reSSR.ReplaceAllString(strings.TrimSpace(scanner.Text()), `$1`)
		raw += strings.Repeat(`=`, (4-(len(raw)%4))%4)
		uriBytes, err := base64.URLEncoding.DecodeString(raw)
		if err != nil {
			// handle err
			log.Fatal(err)
		}
		uri = string(uriBytes)

		method := reSSRInfo.ReplaceAllString(uri, `$4`)

		password := reSSRInfo.ReplaceAllString(uri, `$6`)
		password += strings.Repeat(`=`, (4-(len(password)%4))%4)
		password = base64Decode(password, base64.StdEncoding)

		server := reSSRInfo.ReplaceAllString(uri, `$1:$2`)

		u, err := url.Parse(reSSRInfo.ReplaceAllString(uri, `$7`))
		if err != nil {
			log.Fatal(err)
		}
		params, err := url.ParseQuery(u.RawQuery)
		if err != nil {
			log.Fatal(err)
		}
		obfs := reSSRInfo.ReplaceAllString(uri, `$5`)
		obfsparam := params["obfsparam"][0]
		obfsparam += strings.Repeat(`=`, (4-(len(obfsparam)%4))%4)
		obfsparam = base64Decode(obfsparam, base64.StdEncoding)
		if strings.Contains(obfsparam, `,`) {
			obfsparam = (strings.Split(obfsparam, `,`))[0]
		}

		proto := reSSRInfo.ReplaceAllString(uri, `$3`)
		protoparam := params["protoparam"][0]
		protoparam += strings.Repeat(`=`, (4-(len(protoparam)%4))%4)
		protoparam = base64Decode(protoparam, base64.StdEncoding)

		remark := params["remarks"][0]
		remark += strings.Repeat(`=`, (4-(len(remark)%4))%4)
		remark = base64Decode(remark, base64.URLEncoding)
		remark = regexp.MustCompile(`\s`).ReplaceAllString(remark, "")

		ssr := fmt.Sprintf(
			`ssr://%s:%s@%s/?protocol=%s&protocol_param=%s&obfs=%s&obfs_param=%s`,
			method, password, server, proto, protoparam, obfs, obfsparam,
		)

		ssrURIs = append(ssrURIs, ssr)
		ssrRemarks = append(ssrRemarks, remark)
		// fmt.Printf("%s:\n%s\n\n", remark, ssr)
	}
	return
}

// WGetRawFastBySSRProxy wget over ssr proxy
func WGetRawFastBySSRProxy(urlAddr, ssrProxyAddr string, timeout time.Duration) ([]byte, time.Duration, error) {
	request, _ := http.NewRequest("GET", urlAddr, nil)
	direct, _ := proxy.NewDirect("")
	proxy, err := proxy.DialerFromURL(ssrProxyAddr, direct)
	if err != nil {
		log.Fatal(err)
		return nil, 0, err
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy:                 nil,
			Dial:                  proxy.Dial,
			ResponseHeaderTimeout: 3 * time.Second,
			IdleConnTimeout:       timeout,
			DisableKeepAlives:     true,
		},
		Timeout: timeout,
	}
	t0 := time.Now()
	resp, err := client.Do(request)
	if err != nil {
		return nil, 0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return body, time.Now().Sub(t0), err
}
