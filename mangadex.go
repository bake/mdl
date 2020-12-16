package main

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/bake/mangadex/v2"
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
	parts := strings.Split(url.Path, "/")
	if len(parts) < 2 {
		return nil, errors.Errorf("chapter ID missing in URL")
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, errors.Wrap(err, "could not parse chapter ID")
	}
	chapter, err := md.Chapter(context.Background(), id, nil)
	return chapter.Images(), errors.Wrapf(err, "could not get chapter %s", id)
}

func (mangaDexClient) Authenticate(req *http.Request) {}
