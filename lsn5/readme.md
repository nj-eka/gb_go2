go test -bench .
goos: linux
goarch: amd64
pkg: github.com/nj-eka/gb_go2/lsn5
cpu: Intel(R) Core(TM) i7-4710HQ CPU @ 2.50GHz
BenchmarkPerf1090-4                    1        1261901215 ns/op
BenchmarkPerf5050-4                    1        2077992084 ns/op
BenchmarkPerf9010-4                    1        2874925551 ns/op
BenchmarkPerfRW1090-4                  1        1260277849 ns/op
BenchmarkPerfRW5050-4                  1        1898395314 ns/op
BenchmarkPerfRW9010-4                  1        3107647250 ns/op
PASS
ok      github.com/nj-eka/gb_go2/lsn5   14.049s
