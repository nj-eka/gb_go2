- start godoc
> godoc -http=:6060
- see documentation by opening url: `http://localhost:6060/pkg/github.com/nj-eka/gb_go2/lsn2/`

> go test -bench=. primes/* -benchmem -benchtime=1000x

    goos: linux
    goarch: amd64
    cpu: Intel(R) Core(TM) i7-4710HQ CPU @ 2.50GHz
    BenchmarkSieveWithBits-4            1000              3253 ns/op              93 B/op          1 allocs/op
    BenchmarkSieveWithChannel-4         1000           1212408 ns/op           13345 B/op         94 allocs/op
