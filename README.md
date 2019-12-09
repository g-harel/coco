<!--

TODO
- optional npm/github (with help if none)
- consistent error handling
- rate limiting
- control logging with flags
- add godoc/documentation
- add example usages in readme

 -->

# coco

> track project stats

## Usage

```
$ go get -u github.com/g-harel/coco
```

```
Usage: coco [flags...]

Flags:
  -github-owner string
        List of GitHub owners whose repos to query (comma separated).
  -github-today int
        Show repos if they have this quantity of views today.
  -github-token string
        GitHub API token
  -github-views int
        Show repos if they have this quantity of views.
  -npm-owner string
        List of NPM owners whose packages to query (comma separated).
  -npm-weekly int
        Show repos if they have this quantity of weekly downloads.

GitHub traffic can only be collected from repositories that the token grants push access to.
```

## License

[MIT](./LICENSE)
