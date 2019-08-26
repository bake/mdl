package main

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/bake/mangadex"
	"github.com/pkg/errors"
)

type mangaDex struct{ *mangadex.Client }

func newMangaDex() *mangaDex {
	return &mangaDex{mangadex.New()}
}

func (mangaDex) Match(url *url.URL) bool {
	return strings.HasSuffix(url.Hostname(), "mangadex.org") &&
		strings.HasPrefix(url.Path, "/chapter/")
}

func (md *mangaDex) Files(url *url.URL) ([]string, error) {
	id := strings.TrimPrefix(url.Path, "/chapter/")
	chapter, err := md.Chapter(id)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get chapter %s", id)
	}
	return chapter.Images(), nil
}

func (mangaDex) Download(url string) (io.ReadCloser, error) {
	res, err := http.Get(url)
	return res.Body, err
}
