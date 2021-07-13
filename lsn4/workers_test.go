package main

import "testing"

const(
	attempts = 10
	jobs     = 1000
	workers  = 20
)

func TestAtomicWorkers(t *testing.T) {
	for i := 0; i < attempts; i++ {
		result := runAtomicWorkers(jobs, workers)
		if result != jobs{
			t.Errorf("atomic: Result = %d | Expected %d;\n", result, jobs)
		}
	}
}

func TestMutexWorkers(t *testing.T) {
	for i := 0; i < attempts; i++ {
		result := runMutexWorkers(jobs, workers)
		if result != jobs{
			t.Errorf("mutex: Result = %d | Expected %d;\n", result, jobs)
		}
	}
}

//goos: linux
//goarch: amd64
//pkg: github.com/nj-eka/gb_go2/lsn4
//cpu: Intel(R) Core(TM) i7-4710HQ CPU @ 2.50GHz
//BenchmarkAtomicWorkers-4   	1000000000	         0.01537 ns/op
func BenchmarkAtomicWorkers(b *testing.B) {
	runAtomicWorkers(jobs, workers)
}

//goos: linux
//goarch: amd64
//pkg: github.com/nj-eka/gb_go2/lsn4
//cpu: Intel(R) Core(TM) i7-4710HQ CPU @ 2.50GHz
//BenchmarkMutexWorkers-4   	1000000000	         0.01392 ns/op
func BenchmarkMutexWorkers(b *testing.B) {
	runMutexWorkers(jobs, workers)
}
