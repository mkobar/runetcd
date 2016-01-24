DISCLAIMER: I dumped a lot of code without much documentation and optimization.
It needs to be cleaned up and better organized. Further work will continue at
https://github.com/gyuho/runetcd.

<br>

## runetcd [![Build Status](https://img.shields.io/travis/gophergala2016/runetcd.svg?style=flat-square)](https://travis-ci.org/gophergala2016/runetcd) [![Godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/gophergala2016/runetcd)

runetcd runs, demos [`etcd`](https://github.com/coreos/etcd) with CLIs and dashboards.

- [Why](#why)
- [**_10-second Demo_**](#10-second-demo)
- [Dashboard for production](#dashboard-for-production)
- [Credits](#credits)
- [Progress](#progress)

[↑ top](#runetcd--)
<br><br>


## Why

- http://try.redis.io
- http://play.golang.org

So must `etcd` be easy to try and play.

[↑ top](#runetcd--)
<br><br>


## 10-second Demo

Public demo here https://runetcd.io.

```

```

It runs exactly the same `etcd` cluster as in production. Each client launches
N number of machines and writes to the distributed database. `runetcd.io` has
limits on resources that can be used. To experience full-powered `etcd`, please
run CLI locally. `runetcd` runs executable `etcd` binaries. Make sure to have them
ready in your local machine. And here's how:

```

```

It's that easy! Just binary, nothing else.

[↑ top](#runetcd--)
<br><br>


## Dashboard for production

You can use this as an `etcd` dashboard:

```

```

[↑ top](#runetcd--)
<br><br>


## Credits

- https://github.com/coreos/etcd
- https://github.com/rakyll/statik
- https://github.com/mattn/goreman
- https://github.com/gyuho/psn

[↑ top](#runetcd--)
<br><br>


## Progress

<br>
##### 5:30PM Saturday, January 23, 2016

Got some realtime updates, logs up and running!

![changelog_01](./demos/changelog_01)

<br>
##### 3AM Saturday, January 23, 2016

CLI demo done.

![changelog_00](./demos/changelog_00)

[↑ top](#runetcd--)
<br><br>
