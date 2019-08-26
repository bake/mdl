package main

import (
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/bake/mangarock"
	"github.com/pkg/errors"
)

type mangaRock struct {
	*mangarock.Client
	reg *regexp.Regexp
}

func newMangaRock() *mangaRock {
	return &mangaRock{
		mangarock.New(),
		regexp.MustCompile("^/manga/(mrs-serie-[0-9]+)/chapter/(mrs-chapter-[0-9]+)"),
	}
}

func (mr *mangaRock) Match(url *url.URL) bool {
	return strings.HasSuffix(url.Hostname(), "mangarock.com") &&
		mr.reg.MatchString(url.Path)
}

func (mr *mangaRock) Files(url *url.URL) ([]string, error) {
	ids := mr.reg.FindStringSubmatch(url.Path)
	chapter, err := mr.Chapter(ids[1], ids[2])
	if err != nil {
		return nil, errors.Wrapf(err, "could not get chapter %s of %s", ids[2], ids[1])
	}
	return chapter.Pages, nil
}

func (mangaRock) Download(url string) (io.ReadCloser, error) {
	res, err := http.Get(url)
	return res.Body, err
}
