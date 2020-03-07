package main

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/bake/mangadex"
	"github.com/pkg/errors"
)

type mangaDexClient struct{ *mangadex.Client }

func newMangaDexClient() *mangaDexClient {
	return &mangaDexClient{mangadex.New()}
}

func (mangaDexClient) Match(url *url.URL) bool {
	return strings.HasSuffix(url.Hostname(), "mangadex.org") &&
		strings.HasPrefix(url.Path, "/chapter/")
}

func (md *mangaDexClient) Files(url *url.URL) ([]string, error) {
	id := strings.Split(url.Path, "/")
	chapter, err := md.Chapter(id[2])
	return chapter.Images(), errors.Wrapf(err, "could not get chapter %s", id)
}

func (mangaDexClient) Authenticate(req *http.Request) {}
