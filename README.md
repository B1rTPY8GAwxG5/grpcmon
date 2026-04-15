# grpcmon

Lightweight CLI tool for monitoring and replaying gRPC traffic in development environments.

---

## Installation

```bash
go install github.com/grpcmon/grpcmon@latest
```

Or download a prebuilt binary from the [releases page](https://github.com/grpcmon/grpcmon/releases).

---

## Usage

Start monitoring gRPC traffic on a target service:

```bash
grpcmon watch --addr localhost:50051
```

Record traffic to a file and replay it later:

```bash
# Record
grpcmon record --addr localhost:50051 --out traffic.grpc

# Replay
grpcmon replay --addr localhost:50051 --in traffic.grpc
```

Filter by specific RPC methods:

```bash
grpcmon watch --addr localhost:50051 --method /mypackage.MyService/MyMethod
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--addr` | Target gRPC server address | `localhost:50051` |
| `--out` | Output file for recorded traffic | `traffic.grpc` |
| `--in` | Input file for replay | `traffic.grpc` |
| `--method` | Filter by RPC method path | all methods |
| `--tls` | Enable TLS | `false` |

---

## Requirements

- Go 1.21+
- Target service must have [gRPC server reflection](https://github.com/grpc/grpc/blob/master/doc/server-reflection.md) enabled

---

## License

MIT © [grpcmon contributors](https://github.com/grpcmon/grpcmon/blob/main/LICENSE)