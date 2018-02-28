# coco

> track repo traffic through the github api

All specified users' public repositories are queried.

Repositories with no views will not get displayed.

## Usage

```shell
$ go get -u github.com/g-harel/coco
```

```shell
$ coco -names="{username1,orgname1,...}" -token="{access_token}"

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
