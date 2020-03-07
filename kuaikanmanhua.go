package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type kuaikanClient struct{}

func newKuaikanClient() *kuaikanClient { return &kuaikanClient{} }

func (k *kuaikanClient) Match(url *url.URL) bool {
	return strings.HasSuffix(url.Hostname(), "kuaikanmanhua.com")
}

func (k *kuaikanClient) Files(url *url.URL) (urls []string, err error) {
	res, err := http.Get(url.String())
	if err != nil {
		return nil, errors.Wrap(err, "could not get chapter")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not read body")
	}
	raws := regexp.MustCompile(`url:"(http.+?\?sign.+?)"`).FindAllSubmatch(body, -1)
	for _, raw := range raws {
		url, err := strconv.Unquote("\"" + string(raw[1]) + "\"")
		if err != nil {
			return nil, errors.Wrap(err, "could not decode url")
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func (k *kuaikanClient) Authenticate(req *http.Request) {}
