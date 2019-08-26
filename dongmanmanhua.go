package main

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

type dongmanmanhua struct{}

func newDongmanmanhua() *dongmanmanhua {
	return &dongmanmanhua{}
}

func (dongmanmanhua) Match(url *url.URL) bool {
	if !strings.HasSuffix(url.Hostname(), "dongmanmanhua.cn") {
		return false
	}
	vals := url.Query()
	return vals.Get("title_no") != "" && vals.Get("episode_no") != ""
}

func (d *dongmanmanhua) Files(url *url.URL) (urls []string, err error) {
	res, err := d.get(url.String())
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

func (d *dongmanmanhua) Download(url string) (io.ReadCloser, error) {
	res, err := d.get(url)
	return res.Body, err
}

// get sends a GET request and sets its referer to the requested URL.
// Dongmanmanhua does not accept requests without a referer from their own site.
func (dongmanmanhua) get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create request")
	}
	req.Header.Add("Referer", url)
	return http.DefaultClient.Do(req)
}
