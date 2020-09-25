# mkill

[![PkgGoDev](https://pkg.go.dev/badge/golang.design/x/mkill)](https://pkg.go.dev/golang.design/x/mkill) [![Go Report Card](https://goreportcard.com/badge/golang.design/x/mkill)](https://goreportcard.com/report/golang.design/x/mkill)
![mkill](https://github.com/golang-design/mkill/workflows/mkill/badge.svg?branch=master)

Package mkill limits the number of threads in a Go program, without crash the whole program.

```
import "golang.design/x/mkill"
```

## Quick Start

```
mkill.GOMAXTHREADS(10)
```

## License

MIT &copy; The [golang.design](https://golang.design) Authors