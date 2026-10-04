# LSM Tree Benchmark Results

## 1. Project

**Project:** Distributed Library Key-Value Store

**Storage Engine:** LSM Tree

**Language:** Go

**Operating System:** Ubuntu 26.04.1 LTS (WSL)

**Go Version:** Go 1.26.0

**CPU:** AMD Ryzen 5 3550H with Radeon Vega Mobile Gfx

**Architecture:** amd64

---

## 2. LSM Tree PUT Benchmark

### Command

```bash
go test ./internal/storage -bench=BenchmarkLSMPut -benchmem
```

### Result

```text
BenchmarkLSMPut-8
10000
2583442 ns/op
87976 B/op
12 allocs/op
PASS
```

The LSM Tree successfully completed the PUT benchmark.

---

## 3. LSM Tree GET Benchmark

### Command

```bash
go test ./internal/storage -bench=BenchmarkLSMGet -benchmem
```


---

## 4. RocksDB Sequential PUT Benchmark

**RocksDB Version:** 9.11.2

### Command

```bash
db_bench --benchmarks=fillseq --num=10000 --value_size=100 --use_existing_db=0
```

### Result

```text
fillseq : 76.120 micros/op
13134 ops/sec
0.761 seconds
10000 operations
```

**Throughput:** 13,134 operations/second.

---

## 5. RocksDB Sequential GET Benchmark

### Command

```bash
db_bench --benchmarks=readseq --num=10000 --reads=10000 --use_existing_db=1
```

### Result

```text
readseq : 2.602 micros/op
382350 ops/sec
0.026 seconds
10000 operations
```

**Throughput:** 382,350 operations/second.

---

## 6. RocksDB Random GET Benchmark

### Command

```bash
db_bench --benchmarks=readrandom --num=10000 --reads=10000 --use_existing_db=1
```

### Result

```text
readrandom : 11.505 micros/op
86795 ops/sec
0.115 seconds
10000 operations
10000 of 10000 found
```

**Throughput:** 86,795 operations/second.

All 10,000 requested random reads successfully found their records.

---

## 7. Correctness Testing

### Test command

```bash
go test ./...
```

### Result

```text
PASS
```

The complete Go project passed the test suite.

---

## 8. Race Detection

### Command

```bash
go test -race ./...
```

### Result

```text
PASS
```

No race detector failures were reported.

---

## 9. Build Verification

### Command

```bash
go build ./...
```

### Result

```text
PASS
```

The complete project builds successfully.

---

## 10. Benchmark Summary

| System   | Operation      |                    Throughput |
| -------- | -------------- | ----------------------------: |
| LSM Tree | PUT            | See Go benchmark result above |
| RocksDB  | Sequential PUT |                13,134 ops/sec |
| RocksDB  | Sequential GET |               382,350 ops/sec |
| RocksDB  | Random GET     |                86,795 ops/sec |

### Notes

The LSM Tree and RocksDB benchmarks were run on the same Ubuntu WSL environment.

The RocksDB benchmarks used 10,000 records and a 100-byte value size.

The benchmark results are intended to provide a basic performance comparison. The workloads and benchmarking tools are not completely identical, so the numbers should be interpreted as experimental measurements rather than a strict apples-to-apples performance comparison.

---

## 11. Milestone Result

The LSM storage milestone has been substantially completed:

* MemTable implemented
* WAL implemented
* SSTable implemented
* LSM Tree implemented
* WAL recovery implemented
* SSTable flushing implemented
* Compaction implemented
* Tombstone-based deletion implemented
* Unit tests implemented
* Race detection completed
* Go build verified
* LSM  benchmark completed
* RocksDB benchmark completed

