# Benchmark Results

## Environment

- OS: Ubuntu WSL
- Language: Go
- Project: Distributed Library KV
- Storage: LSM Tree

## PUT Benchmark

Command:

```bash
go test ./internal/storage -bench=BenchmarkLSMPut -benchmem
goos: linux
goarch: amd64
pkg: distributed-library-kv/internal/storage
cpu: AMD Ryzen 5 3550H with Radeon Vega Mobile Gfx
BenchmarkLSMPut-8          10000           2583442 ns/op           87976 B/op             12 allocs/op
PASS
ok      distributed-library-kv/internal/storage 25.881s
mounika@MOUNIKAHARI:~/distributed-library-kv$ go test ./...
?       distributed-library-kv/cmd/server       [no test files]
ok      distributed-library-kv/internal/kv      (cached)
?       distributed-library-kv/internal/library [no test files]
ok      distributed-library-kv/internal/storage 0.025s
mounika@MOUNIKAHARI:~/distributed-library-kv$ go test -race ./...
?       distributed-library-kv/cmd/server       [no test files]
ok      distributed-library-kv/internal/kv      (cached)
?       distributed-library-kv/internal/library [no test files]
ok      distributed-library-kv/internal/storage 1.080s
mounika@MOUNIKAHARI:~/distributed-library-kv$ go build ./...
PASS
