<!--

TODO
- optional npm/github (with help if none)
- consistent error handling
- rate limiting
- make quantity filters flags
- control logging with flags
- log time since start
- add stars to github table
- remove argument spreads
- add godoc/documentation

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
  -github-owner  Comma-separated list of GitHub owners to query for owned repos.
  -github-token  Authentication token used for GitHub requests.
  -npm-owner     Comma-separated list of NPM owners to query for owned packages.

GitHub traffic can only be collected from repositories that your account has push access to.
```

```
+-------------------+-------+-----+--------+-------------------------------------------------------------+
|       REPO        | VIEWS | DAY | UNIQUE |                            LINK                             |
+-------------------+-------+-----+--------+-------------------------------------------------------------+
| okwolo            |   678 |   1 |    336 | https://github.com/okwolo/okwolo/graphs/traffic             |
| searx             |   361 |   1 |      9 | https://github.com/g-harel/searx/graphs/traffic             |
| okwolo-todomvc    |    19 |   0 |      2 | https://github.com/okwolo/okwolo-todomvc/graphs/traffic     |
| superpermutations |    13 |   0 |      4 | https://github.com/g-harel/superpermutations/graphs/traffic |
| cover-gen         |     5 |   0 |      2 | https://github.com/g-harel/cover-gen/graphs/traffic         |
| website           |     5 |   0 |      2 | https://github.com/okwolo/website/graphs/traffic            |
| edelweiss         |     4 |   0 |      2 | https://github.com/g-harel/edelweiss/graphs/traffic         |
| codewars          |     4 |   0 |      1 | https://github.com/g-harel/codewars/graphs/traffic          |
| backend-challenge |     1 |   0 |      1 | https://github.com/g-harel/backend-challenge/graphs/traffic |
| ence              |     1 |   0 |      1 | https://github.com/g-harel/ence/graphs/traffic              |
+-------------------+-------+-----+--------+-------------------------------------------------------------+
```
