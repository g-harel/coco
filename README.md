<!--

TODO
- optional npm/github
- parallelization utility
- consistent error handling
- rate limiting
- request user packages in parallel

 -->

# coco

> track repo traffic through the github api

Repositories with no views (in the past 14 days) will not get displayed.

## Usage

```
$ go get -u github.com/g-harel/coco
```

```
Usage: coco [USER]...
List repository traffic for USER(s).

Examples:
  coco username
  coco orgname username

Looks for GitHub api key in GITHUB_API_TOKEN environment variable.
Traffic can only be collected from repositories that your account has push access to.
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
