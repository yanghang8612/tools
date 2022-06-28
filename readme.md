## Awesome Tools

This repo contains some useful sub commands for ethereum and TRON

### Sub commands

| Command | Description                   |
|:-------:|-------------------------------|
|  `db`   | Database related commands     |
| `addr`  | Address related commands      |
|  `vm`   | EVM related commands          |
| `scan`  | TronScan related commands     |
|  `eth`  | ETH JSON-RPC related commands |

### Installation

```shell
protoc --go_out=./proto ./proto/*.proto
go build -o /usr/local/bin/tools
ln -s /usr/local/bin/tools /bin/tt
```

#### How to install `protoc` or `go`

Can't you fucking Google it?