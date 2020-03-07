package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/bake/tapas"
	"github.com/pkg/errors"
)

type tapasClient struct{ *tapas.Client }

func newTapasClient() *tapasClient {
	return &tapasClient{tapas.New()}
}

func (tapasClient) Match(url *url.URL) bool {
	return strings.HasSuffix(url.Hostname(), "tapas.io") &&
		strings.HasPrefix(url.Path, "/episode/")
}

func (c *tapasClient) Files(url *url.URL) ([]string, error) {
	cids := strings.Split(url.Path, "/")
	if len(cids) < 2 {
		return nil, errors.New("URL does not contain a chapter ID")
	}
	cid, err := strconv.Atoi(cids[2])
	if err != nil {
		return nil, errors.Wrapf(err, "could not convert chapter id %s", cids[2])
	}

	res, err := http.Get(url.String())
	if err != nil {
		return nil, errors.Wrap(err, "could not get chapter")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not read response body")
	}
	sids := regexp.MustCompile(`sid:\s*(\d+),`).FindSubmatch(body)
	if len(sids) < 2 {
		return nil, errors.New("could not find series ID")
	}
	sid, err := strconv.Atoi(string(sids[1]))
	if err != nil {
		return nil, errors.Wrapf(err, "could not convert series id %s", sids[1])
	}
	e, err := c.Episode(sid, cid)
	if err != nil {
		return nil, err
	}
	urls := make([]string, len(e.Contents))
	for i, url := range e.Contents {
		urls[i] = url.String()
	}
	return urls, nil
}

func (tapasClient) Authenticate(req *http.Request) {}
