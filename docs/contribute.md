
## Dependecies

`miniplanes` development environment depends on two external tools:
* `swagger`: supplied by `go-swagger`
* `go-bindata`: supplied by `go-bindata`

in order to build `miniplanes` executable you need both. You can install on your development environment via:

```shell
$ go get -u github.com/go-swagger/go-swagger/cmd/swagger
$ go get -u github.com/jteeuwen/go-bindata/...
```

## How to populate DB

To start having `miniplanes` running you must have DB populated. Some data come statically from [openflights.org](http://www.openflights.org/data.html), the `schedule` is generated from a small tool bundled with `miniplanes`.

### In Kube/Minishift

### Locally
