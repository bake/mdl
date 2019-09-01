# mdl (manga downloader)

[![Go Report Card](https://goreportcard.com/badge/github.com/bake/mdl)](https://goreportcard.com/report/github.com/bake/mdl)

Think of youtube-dl but for mangas. And with a very, very small number of
supported sites. An alternative with more features is
[comics-downloader](https://github.com/Girbons/comics-downloader).

## Supported sites

- Dongmanmanhua (Chinese)
- MangaDex (English, [library](https://github.com/bake/mangadex))
- MangaRock (English, [library](https://github.com/bake/mangarock))

## Install

```bash
$ go get github.com/bake/mdl
$
```

## Usage

```bash
$ mdl -help
Usage: mdl [-out=.] [-worker=4] url
Flags:
  -out string
        Download directory (default ".")
  -worker int
        Concurrent downloads (default 4)
```
