# Distributed Library Key-Value Store — Project Progress Report

Based on our work so far, your project has progressed from a basic Go key-value store to a distributed application with storage-engine components, replication, simplified Raft leader election, Docker, Kubernetes, Helm, and monitoring infrastructure.

Your application and monitoring pods were running successfully in Kubernetes at our last check. The remaining issue is the Windows–WSL connection, which interrupted access to your development environment. Grafana's dashboard and metrics visualization still need final verification.

## 1. What we have completed

Milestone 1 — Basic key-value store

Implemented basic PUT, GET, and DELETE operations, plus a library use case and Go unit tests.

Observed output: `Book: Computer Networks` and `B103 deleted successfully`.

Milestone 2 — Storage engine

Implemented the Write-Ahead Log (WAL), MemTable, SSTable, and LSM-tree engine with tests.

Observed: `go test ./internal/storage` and the full Go test suite passed at the recorded checkpoints.

Milestone 3 — Distributed replication and Raft

Added replication between three nodes and simplified Raft-style leader election.

Observed: node2 became leader after node1 was stopped; a PUT through node2 was retrievable from node3.

Milestone 4 — Docker and Docker Compose

Built container images and configured a three-node setup with container networking.

Observed: node service-name resolution and container networking worked. Older Compose containers are currently stopped.

Milestone 5 — Kubernetes

Created a Kind cluster and deployed the three application nodes.

Last verified: the control-plane node was `Ready`, and all three application pods were `Running`.

Milestone 6 — Helm

Created and validated a Helm chart, then deployed the application through Helm.

Observed: Helm release `distributed-kv` was deployed, and `HELM_TEST` written through node1 was readable from node3.

Milestone 7 — Prometheus and Grafana installation

Installed the `kube-prometheus-stack` in the `monitoring` namespace.

Last verified: Prometheus was `2/2 Running`, Grafana was `3/3 Running`, and the other listed monitoring pods were also running.

These checks establish that the components were working at the recorded checkpoints. We still need a final end-to-end run after WSL is restored.

## 2. Important outputs obtained in the project

These are the actual outputs recorded during our earlier tests.

### A. Basic library operations

```
Book: Computer Networks
B103 deleted successfully
```

This demonstrates retrieving a book record and deleting a record.

### B. Go tests and build

Commands that passed at earlier checkpoints:

```
go test ./...
go test -race ./...
go build ./...
```

These verified the code and concurrency checks at those checkpoints. They should be rerun for final confirmation.

### C. Distributed fault-tolerance test

After stopping node1, node2 became leader. A write through node2 succeeded, and the value was retrieved from node3.

```
Node node2 became LEADER
PUT success: true
GET found: true
GET value: Distributed Library KV Working
```

This demonstrates the tested failover and cross-node read path. It does not by itself prove complete production-grade Raft consensus.

### D. Kubernetes and Helm test

The last Kubernetes check returned:

```
NAME                           STATUS   ROLES
distributed-kv-control-plane   Ready    control-plane
```

Application pods:

```
node1   1/1 Running
node2   1/1 Running
node3   1/1 Running
```

The actual pod names included generated suffixes; the shortened names above are for readability.

The Helm deployment also passed a write/read test:

```
PUT success: true
GET found: true
GET value: Helm Distributed KV Working
```

### E. Benchmark result

Our recorded Go LSM PUT benchmark was:

| Metric                    | Result          |
| ------------------------- | --------------- |
| Operations                | 10,000          |
| Time per operation        | 2,583,442 ns/op |
| Memory per operation      | 87,976 B/op     |
| Allocations per operation | 12 allocs/op    |

This is a preliminary result, not the final benchmark report. We still need to rerun it under a consistent setup and complete the comparison with RocksDB.

## 3. What remains to finish

1\. Restore WSL access

Fix `Wsl/Service/0x8007274c` so Ubuntu terminals open reliably.

2\. Finish Grafana verification

Log in, confirm Prometheus is available as a data source, and display useful metrics. Installing Grafana is complete; dashboard verification is not.

3\. Final reliability tests

Test recovery after a pod failure and verify PUT/GET behavior after failover. Confirm recovery and compaction behavior in the storage engine.

4\. Final benchmarks

Run the Go benchmarks, collect comparable RocksDB results, and prepare a performance table.

5\. Documentation and GitHub

Complete the README, architecture diagram, installation instructions, test evidence, and benchmark report. Commit and push the final changes.

## 4. Overall status

You have completed the major implementation and deployment work. The project is not yet fully finalized because the final monitoring verification, complete recovery/benchmark evidence, documentation, and final GitHub update remain.

The immediate priority is to fix WSL. Once Ubuntu opens reliably, we can resume from the existing setup rather than rebuilding it.

Progress summary based on the project work and terminal outputs recorded in our conversation.&#x20;
