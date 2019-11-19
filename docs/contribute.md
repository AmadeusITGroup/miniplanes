
## Development environment
`miniplanes` is written in `go-lang` more and at least `go-lang` version `1.12.6` is needed. It could work with `1.11` too, but no effort is put into development to grant `1.11` compatibility.

### Dependecies
`miniplanes` makes use of `go mod`, more particulary `vendor`ed `go-mod`. So basically one has to set `GO111MODULE=on` before compiling it.
Compilation is performed through
```shell
$ go build --mod=vendor
```
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
