package main

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

type dongmanClient struct {
	domain string
}

func newDongmanClient() *dongmanClient {
	return &dongmanClient{"https://www.dongmanmanhua.cn/"}
}

func (dongmanClient) Match(url *url.URL) bool {
	if !strings.HasSuffix(url.Hostname(), "dongmanmanhua.cn") {
		return false
	}
	vals := url.Query()
	return vals.Get("title_no") != "" && vals.Get("episode_no") != ""
}

func (d *dongmanClient) Files(url *url.URL) (urls []string, err error) {
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not generate new request")
	}
	d.Authenticate(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not get chapter")
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not read document")
	}
	doc.Find("img._images[data-url]").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("data-url")
		urls = append(urls, url)
	})
	return urls, nil
}

func (d *dongmanClient) Authenticate(req *http.Request) {
	req.Header.Add("Referer", d.domain)
}
