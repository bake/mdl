package main

import (
	"context"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/url"
	"os"
	"path"
	"runtime"

	_ "github.com/bake/mri"
	"github.com/cheggaaa/pb/v3"
	"github.com/pkg/errors"
	"golang.org/x/sync/semaphore"
)

// client describes the interface that each supported site has to implement.
type client interface {
	// Match returns true if the given URL can be processed by this client.
	Match(url *url.URL) bool
	// Files return a slice of URLs of images.
	Files(url *url.URL) ([]string, error)
	// Download downloads an image and returns its read closer. This is usually
	// the body of a HTTP response. Useful if a site needs to set special headers.
	Download(url string) (io.ReadCloser, error)
}

type clients []client

func (cs clients) filter(url *url.URL) client {
	for _, c := range cs {
		if !c.Match(url) {
			continue
		}
		return c
	}
	return nil
}

func main() {
	out := flag.String("out", ".", "Download directory")
	worker := flag.Int("worker", runtime.NumCPU(), "Concurrent downloads")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s ", os.Args[0])
		fmt.Fprintf(os.Stderr, "[-out=%s] [-worker=%d] url\n", *out, *worker)
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	url, err := url.Parse(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse URL: %v\n", err)
		os.Exit(1)
	}
	cs := clients{
		newDongmanmanhua(),
		newMangaDex(),
		newMangaRock(),
	}
	c := cs.filter(url)
	if c == nil {
		fmt.Fprintf(os.Stderr, "the URL is not supported\n")
		os.Exit(1)
	}
	files, err := c.Files(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get files: %v\n", err)
		os.Exit(1)
	}

	sem := semaphore.NewWeighted(int64(*worker))
	ctx := context.Background()
	bar := pb.StartNew(len(files))
	defer bar.Finish()
	defer sem.Acquire(ctx, int64(*worker))
	for i, file := range files {
		i, file := i, file
		sem.Acquire(ctx, 1)
		go func() {
			defer sem.Release(1)
			defer bar.Increment()
			err := download(c, path.Join(*out, fmt.Sprintf("%04d.png", i)), file)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}()
	}
}

func download(c client, dst, src string) error {
	r, err := c.Download(src)
	if err != nil {
		return errors.Wrap(err, "cound not download image")
	}
	defer r.Close()
	img, _, err := image.Decode(r)
	if err != nil {
		return errors.Wrap(err, "could not decode image")
	}
	w, err := os.Create(dst)
	if err != nil {
		return errors.Wrap(err, "could not create file")
	}
	defer w.Close()
	if err := png.Encode(w, img); err != nil {
		return errors.Wrap(err, "could not encode image")
	}
	return nil
}
