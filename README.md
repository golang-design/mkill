# mkill

[![PkgGoDev](https://pkg.go.dev/badge/golang.design/x/mkill)](https://pkg.go.dev/golang.design/x/mkill) [![Go Report Card](https://goreportcard.com/badge/golang.design/x/mkill)](https://goreportcard.com/report/golang.design/x/mkill)
![mkill](https://github.com/golang-design/mkill/workflows/mkill/badge.svg?branch=master)

> [!WARNING]
> **This repository is archived and no longer maintained.**
>
> mkill capped the number of OS threads *reactively*: a background daemon
> polled `ps` once a second and reaped the surplus by spawning goroutines that
> lock an OS thread and exit. That approach is fragile (it shells out to `ps`,
> runs work on import, and kills threads it does not own).
>
> Use [`golang.design/x/thread`](https://github.com/golang-design/thread)
> instead. Its [`thread.Pool`](https://pkg.go.dev/golang.design/x/thread#Pool)
> caps threads *by construction* — work runs on a fixed set of owned threads,
> with no subprocess and no import-time side effects:
>
> ```go
> p := thread.NewPool(10)
> defer p.Terminate()
> p.Go(func() { /* runs on one of 10 owned threads */ })
> ```

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