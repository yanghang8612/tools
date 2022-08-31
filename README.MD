## Awesome Tools

This repo contains some useful sub commands for ethereum and TRON

## Sub commands

| Command | Description                                 |
|:-------:|---------------------------------------------|
| `call`  | Interact with contract on TRON network      |
|  `now`  | Convert time between datetime and timestamp |
|  `abi`  | ABI related commands                        |
|  `db`   | Database related commands                   |
|  `eth`  | ETH JSON-RPC related commands               |
|  `hex`  | Hex related commands                        |
| `scan`  | TronScan related commands                   |

## Installation

```shell
go build -o /usr/local/bin/tools
ln -s /usr/local/bin/tools /bin/tt
```

### How to install `go`

You can't fucking Google or Baidu?

## Commands Usage

If you just append one arg to the command without subcommands, the program will decide what logic to execute based on the type of parameters you enter.

- arg is like eth address or TRON address

```shell
$ tt 0x000000000000000000000000000000000000dEaD
     [eth addr] - 0x000000000000000000000000000000000000dEaD
    [tron addr] - T9yD14Nj9j7xAB4dbGeiX9h8upfCg3PBbY
[32bytes in BE] - 0x000000000000000000000000000000000000000000000000000000000000dead
[32bytes in LE] - 0x000000000000000000000000000000000000dead000000000000000000000000
   [in decimal] - 57005
     [in ascii] - ޭ

$ tt T9yD14Nj9j7xAB4dbGeiX9h8upfCg3PBbY
 [eth addr] - 0x000000000000000000000000000000000000dEaD
[tron addr] - T9yD14Nj9j7xAB4dbGeiX9h8upfCg3PBbY
   [in hex] - 0x5439794431344e6a396a377841423464624765695839683875706643673350426259
```

- arg is in hexadecimal

```shell
$ tt 0x200007368616269
     [eth addr] - 0x0000000000000000000000000200007368616269
    [tron addr] - T9yD14Nj9j7xAB4dbGguV3DLSwvb8p5tdG
[32bytes in BE] - 0x0000000000000000000000000000000000000000000000000200007368616269
[32bytes in LE] - 0x0200007368616269000000000000000000000000000000000000000000000000
   [in decimal] - 144115683748307561
     [in ascii] - shabi
     
$ tt 0xa9059cbb000000000000000000000000e607f127507951682391fcc420d0b6f1bd02eb9600000000000000000000000000000000000000000000000000000000160dc080
[selector]: a9059cbb - transfer(address,uint256)
[unpack result]:
  - [arg-00]: address, 0xE607f127507951682391FcC420D0b6F1BD02Eb96 - TWwVvzy7iPVKs9oi6BdTLwjA6XJNc8h8aC
  - [arg-01]: uint256, 370000000 - 370,000,000 (9)
[in ascii] - 'PyQh# ж
```

- arg is like a timestamp

```shell
$ tt 1659929813
  [eth addr] - 0x0000000000000000000000000000000062f084d5
 [tron addr] - T9yD14Nj9j7xAB4dbGeiX9hRTdK9jdt5kW
    [in sec] - 2022-08-08 11:36:53
    [in hex] - 0x62f084d5
    [in dec] - 1,659,929,813 (10)
[origin hex] - 0x62f084d5
[padded hex] - 0x0000000000000000000000000000000000000000000000000000000062f084d5

$ tt 1659929813000
  [eth addr] - 0x000000000000000000000000000001827B86e008
 [tron addr] - T9yD14Nj9j7xAB4dbGeiXEcTq8tGAJhPmx
  [in milli] - 2022-08-08 11:36:53
    [in hex] - 0x01827b86e008
    [in dec] - 1,659,929,813,000 (13)
[origin hex] - 0x01827b86e008
[padded hex] - 0x000000000000000000000000000000000000000000000000000001827b86e008
 
$ tt 1000000000000000000
  [eth addr] - 0x00000000000000000000000000038D7eA4c68000
 [tron addr] - T9yD14Nj9j7xAB4dbGejQFKfDvTYnRiRrV
    [in hex] - 0x038d7ea4c68000
    [in dec] - 1,000,000,000,000,000 (16)
[origin hex] - 0x038d7ea4c68000
[padded hex] - 0x00000000000000000000000000000000000000000000000000038d7ea4c68000
```

- arg is like a datetime

```shell
$ tt 8-8
     [in sec] - 1659888000
   [in milli] - 1659888000000
[in datetime] - 2022-08-08 00:00:00 +0800 CST
[eval result] - 0

$ tt 8:8
     [in sec] - 1661213280
   [in milli] - 1661213280000
[in datetime] - 2022-08-23 08:08:00 +0800 CST
     [in hex] - 0x383a38

$ tt "22-08-08 08:08:08"
     [in sec] - 1659917288
   [in milli] - 1659917288000
[in datetime] - 2022-08-08 08:08:08 +0800 CST
```

- arg is just num in decimal

```shell
$ tt 1
  [eth addr] - 0x0000000000000000000000000000000000000001
 [tron addr] - T9yD14Nj9j7xAB4dbGeiX9h8unkKLxmGkn
   [max hex] - 0x01
   [max dec] - 1
[origin hex] - 0x01
[padded hex] - 0x0000000000000000000000000000000000000000000000000000000000000001

$ tt 1000
  [eth addr] - 0x00000000000000000000000000000000000003e8
 [tron addr] - T9yD14Nj9j7xAB4dbGeiX9h8unnG3nvgHP
[origin hex] - 0x03e8
[padded hex] - 0x00000000000000000000000000000000000000000000000000000000000003e8

$ tt 160
  [eth addr] - 0x00000000000000000000000000000000000000a0
 [tron addr] - T9yD14Nj9j7xAB4dbGeiX9h8unkdKNoCLK
   [max hex] - 0xffffffffffffffffffffffffffffffffffffffff
   [max dec] - 1461501637330902918203684832716283019655932542975
[origin hex] - 0xa0
[padded hex] - 0x00000000000000000000000000000000000000000000000000000000000000a0

$ tt 1000000000000000000
  [eth addr] - 0x0000000000000000000000000de0b6B3a7640000
 [tron addr] - T9yD14Nj9j7xAB4dbGuuWDPmYjeRYhJdYL
    [in hex] - 0x0de0b6b3a7640000
    [in dec] - 1,000,000,000,000,000,000 (19)
[origin hex] - 0x0de0b6b3a7640000
[padded hex] - 0x0000000000000000000000000000000000000000000000000de0b6b3a7640000
```
- arg is like a function or event

```shell
$ tt "Transfer(address from, address to, uint amt)"
 [abi valid] - Transfer(address,address,uint256)
[origin hex] - 0xddf252ad
[padded hex] - 0xddf252ad00000000000000000000000000000000000000000000000000000000

$ tt "transfer(address to, uint amt)"
 [abi valid] - transfer(address,uint256)
[origin hex] - 0xa9059cbb
[padded hex] - 0xa9059cbb0000000000000000000000000000000000000000000000000000000
```

- arg can be eval

```shell
$ tt 1e18
[eval result] - 1000000000000000000

$ tt 10e18
[eval result] - 10000000000000000000

$ tt 100+1
[eval result] - 101

$ tt 1e18/7
[eval result] - 142857142857142857

$ tt "(86470766307778918*534818468397964620*43200 + 86571716302054628*534818468397964620*561600)/1e18"
[eval result] - 27999999999999999803709
```

- arg is just string

```shell
$ tt shabi
[in hex] - 0x7368616269
```

### Command `abi`

#### Usage

```shell
$ tt abi
NAME:
   tt abi - ABI related commands

USAGE:
   tt abi command [command options] [arguments...]

COMMANDS:
   pad     Pad num(in hex or dec) to 32bytes
   split   Spilt data to each 32bytes
   unpack  Unpack data with given types
   4bytes  Get 4bytes selector for given method or event

OPTIONS:
   --help, -h  show help (default: false)
```

#### Examples

- `pad`

```shell
$ tt abi pad 0
[origin hex] 0x
[padded hex] 0x0000000000000000000000000000000000000000000000000000000000000000

$ tt abi pad 100
[origin hex] 0x64
[padded hex] 0x0000000000000000000000000000000000000000000000000000000000000064

$ tt abi pad 0xdead
[32bytes in BE] - 0x000000000000000000000000000000000000000000000000000000000000dead
[32bytes in LE] - 0xdead000000000000000000000000000000000000000000000000000000000000
```

- `split`

```shell
$ tt abi split 0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b1a2bc2ec500000000000000000000000000000000000000000000000000000f207539952d00000000000000000000000000000000000000000000000000000b1a2bc2ec500000000000000000000000000041aa6f10960ed9f7fe44aacc3aa33dd8f7da108c23
[each data word]:
0x00: 0000000000000000000000000000000000000000000000000000000000000000
0x20: 00000000000000000000000000000000000000000000000000b1a2bc2ec50000
0x40: 0000000000000000000000000000000000000000000000000f207539952d0000
0x60: 0000000000000000000000000000000000000000000000000b1a2bc2ec500000
0x80: 000000000000000000000041aa6f10960ed9f7fe44aacc3aa33dd8f7da108c23

$ tt abi split 0xa9059cbb000000000000000000000000e607f127507951682391fcc420d0b6f1bd02eb9600000000000000000000000000000000000000000000000000000000160dc080
[selector]: a9059cbb - transfer(address,uint256)
[each data word]:
0x00: 000000000000000000000000e607f127507951682391fcc420d0b6f1bd02eb96
0x20: 00000000000000000000000000000000000000000000000000000000160dc080
```

- `unpack`

```shell
$ tt abi unpack "address,address,address,address" 0x000000000000000000000000e3a2cdc25058e5dee0f4b5c1d5c7bfd5dd6836be000000000000000000000000cd420bb0b2b3fd8fa65cb53e01d47738d27949430000000000000000000000008086b67a46a54c7b42d90107e795d3e3c92f6d2100000000000000000000000065fa68800fff5a10346d1a3aa1fb2ce92f2e2971
[unpack result]:
  - [arg-00]: address, 0xe3a2CDC25058e5DEe0F4b5C1d5C7BFD5DD6836Be - TWiqRQZkKqPFV6saFw9dqHJrNG64f7QCPw
  - [arg-01]: address, 0xCd420BB0B2B3FD8Fa65Cb53e01d47738D2794943 - TUgWeKNGt22aig8YRYVAhBFnBK22A2WuZ4
  - [arg-02]: address, 0x8086B67a46a54C7b42D90107e795d3e3C92F6D21 - TMgnsd5t6516yy4xcgod251oSFWcrmWAUo
  - [arg-03]: address, 0x65fA68800FFf5A10346D1A3aA1fb2Ce92f2E2971 - TKGRE6oiU3rEzasue4MsB6sCXXSTx9BAe3

$ tt abi unpack "uint256,uint256,uint256,uint256,address" 0x00000000000000000000000000000000000000000000000000470de4df8200000000000000000000000000000000000000000000000000000429d069189e00000000000000000000000000000000000000000000000000003782dace9d9000000000000000000000000000000000000000000000000000000b1a2bc2ec50000000000000000000000000004134a0f029365c5a6af762fa14d7638a583f72595b
[unpack result]:
  - [arg-00]: uint256, 20000000000000000 - 20,000,000,000,000,000 (17)
  - [arg-01]: uint256, 300000000000000000 - 300,000,000,000,000,000 (18)
  - [arg-02]: uint256, 4000000000000000000 - 4,000,000,000,000,000,000 (19)
  - [arg-03]: uint256, 800000000000000000 - 800,000,000,000,000,000 (18)
  - [arg-04]: address, 0x34A0F029365c5a6aF762FA14D7638a583f72595b - TEmUwDDSPDF6ajToBm2QDg4p3WJUVjYZ4i

$ tt abi unpack "uint256[],uint256[],uint256[],uint256[]" 0x000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000004119011a363613ee4800000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000056bc75e2d6310000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000de1a4662fd0beb000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000de1b9f3e061d692
[unpack result]:
  - [arg-00]: uint256[]
    - [slice-00]: uint256, 1200840114937183596104 - 1,200,840,114,937,183,596,104 (22)
  - [arg-01]: uint256[]
    - [slice-00]: uint256, 100000000000000000000 - 100,000,000,000,000,000,000 (21)
  - [arg-02]: uint256[]
    - [slice-00]: uint256, 1000261351048789680 - 1,000,261,351,048,789,680 (19)
  - [arg-03]: uint256[]
    - [slice-00]: uint256, 1000285049345660562 - 1,000,285,049,345,660,562 (19)

```

- `4bytes`

```shell
$ tt abi 4bytes 'regUser(address,string)'
[origin hex] 0x5a7c7895
[padded hex] 0x5a7c789500000000000000000000000000000000000000000000000000000000

$ tt abi 4bytes 'transfer(address to,uint amt)'
 [abi valid] - transfer(address,uint256)
[origin hex] - 0xa9059cbb
[padded hex] - 0xa9059cbb00000000000000000000000000000000000000000000000000000000
```

### Command `db`

#### Usage

```shell
$ tt db
NAME:
   tt db - Database related commands

USAGE:
   tt db command [command options] [arguments...]

COMMANDS:
   count  Count the total items for given name db
   get    Get value of the given key in db
   hash   Calculate the hash for given name db
   print  Print all key-value for given name db
   diff   Diff for given db-A and db-B

OPTIONS:
   --help, -h  show help (default: false)
```

### Command `eth`

#### Usage

```shell
$ tt eth
NAME:
   tt eth - ETH JSON-RPC related commands

USAGE:
   tt eth command [command options] [arguments...]

COMMANDS:
   logs  Query eth logs with given address, from block and topics, `page` logs at a query

OPTIONS:
   --help, -h  show help (default: false)
```

#### examples

- `logs`

```shell
$ tt eth logs 0x0a3f6849f78076aefaDf113F5BED87720274dDC0 14000000 0x3c278bd500000000000000000000000000000000000000000000000000000000 100000
[████████████████████████████████████████████████████████████████████████████████████████████████████] 100%    3m 47s   1042429/1042429
{"address":"0x0a3f6849f78076aefadf113f5bed87720274ddc0","topics":["0x3c278bd500000000000000000000000000000000000000000000000000000000","0x0000000000000000000000005cab1e5286529370880776461c53a0e47d74fb63","0x000000000000000000000000dd5052bfc4d281793653b0037d46cc2d8d1fd1b5","0x0000000000000000000000000000000000000000000000000000000000000000"],"data":"0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000243c278bd5000000000000000000000000dd5052bfc4d281793653b0037d46cc2d8d1fd1b500000000000000000000000000000000000000000000000000000000","blockNumber":"0xd5bb8f","transactionHash":"0xf8ac4aa71470b2e43fc1c3dbb5a57530a1aa16eee3d16eaad491316321107f03","transactionIndex":"0x19","blockHash":"0xc67969a6d668e7aeef65d8b7c18be55dece2c8ba6f10291bc2a3ee144384fedb","logIndex":"0x19","removed":false}
{"address":"0x0a3f6849f78076aefadf113f5bed87720274ddc0","topics":["0x3c278bd500000000000000000000000000000000000000000000000000000000","0x0000000000000000000000005cab1e5286529370880776461c53a0e47d74fb63","0x000000000000000000000000e246c4ba65d95c2f902e39fbeb0047a67ab4f25a","0x0000000000000000000000000000000000000000000000000000000000000000"],"data":"0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000243c278bd5000000000000000000000000e246c4ba65d95c2f902e39fbeb0047a67ab4f25a00000000000000000000000000000000000000000000000000000000","blockNumber":"0xd6800c","transactionHash":"0x6718672e2e53b2581d94ed0e0e9580347298f9c3745220c4dca74028e309350c","transactionIndex":"0x15c","blockHash":"0x0c8c355eb711811c9f9db809a6ec3a2934a6d1bb283b268efd0b0da7b6ef5f18","logIndex":"0x244","removed":false}
```

### Command `hex`

#### Usage

```shell
$ tt hex
NAME:
   tt hex - Hex related commands

USAGE:
   tt hex command [command options] [arguments...]

COMMANDS:
   addr  Convert addr between hex, TRON-addr and eth-addr
   int   Convert num between dec and hex
   max   Get max value for the type like uint-x
   str   convert hex between str

OPTIONS:
   --help, -h  show help (default: false)
```

#### Examples

- `addr`

```shell
$ tt hex addr 0xdead
 [eth addr] - 0x000000000000000000000000000000000000dEaD
[tron addr] - T9yD14Nj9j7xAB4dbGeiX9h8upfCg3PBbY

$ tt hex addr TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t
 [eth addr] - 0xa614f803B6FD780986A42c78Ec9c7f77e6DeD13C
[tron addr] - TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t

$ tt hex addr 0xa614f803B6FD780986A42c78Ec9c7f77e6DeD13C
 [eth addr] - 0xa614f803B6FD780986A42c78Ec9c7f77e6DeD13C
[tron addr] - TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t
```

- `int`

```shell
$ tt hex int 1
[in hex] - 0x1
[in dec] - 1

$ tt hex int 1000
[in hex] - 0x03e8
[in dec] - 1000

$ tt hex int 0x1000
[in hex] - 0x1000
[in dec] - 4096
```

- `max`

```shell
$ tt hex max 8
[max hex] - 0xff
[max dec] - 255

$ tt hex max 256
[max hex] - 0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
[max dec] - 115792089237316195423570985008687907853269984665640564039457584007913129639935

$ tt hex max 160
[max hex] - 0xffffffffffffffffffffffffffffffffffffffff
[max dec] - 1461501637330902918203684832716283019655932542975
```

- `str`

```shell
$ tt hex str 0x0200007368616269
[in decimal] - 144115683748307561
  [in ascii] - shabi
```

### Command `scan`

#### Usage

```shell
$ tt scan
NAME:
   tt scan - TronScan related commands

USAGE:
   tt scan command [command options] [arguments...]

COMMANDS:
   txs       Query txs in the given net for the given account
   tx        Query the tx in the given net for the given tx hash
   speed     Query pool speed for given token
   transfer  Query all transfer records for given pool

OPTIONS:
   --help, -h  show help (default: false)
```

#### Examples

- `txs`

```shell
$ tt scan txs main TY5cfgCous8GH7rXUrNjRydTMkVHR2BNsR 0 50
[Legend]: ✅ - [Success] ⚠️  - [Revert] ⏱  - [Out_Of_Time] ⚡️ - [Out_Of_Energy] 💢 - [Other]
 1 2022-06-28 20:00:15 ab883582c7ee85d2ed1376e5c187f953c43bec984699a91555f0f8728e5e926e TAFotzexiiUJzGkBHDy9Jbn7rVHoYyWuLA ✅ getReward()
 2 2022-06-28 20:00:09 ca6c1b892c8174a9df81cef8335f904c2ef8d13bed40c36d519e8e08de743768 TWKnrGqU5dijnWT7dbPCDZV23C622T5FDv ✅ getReward()
 3 2022-06-27 20:00:09 c8973f4c4435b816ea98ca66dc7c8d8443d920ab9dde1b7f052796bd28bb40ff TAFotzexiiUJzGkBHDy9Jbn7rVHoYyWuLA ✅ getReward()
 4 2022-06-26 20:00:09 b7c43da37f5ece68c6df042bca478e908894a6f662844140f60e6bb56d7d28d3 TAFotzexiiUJzGkBHDy9Jbn7rVHoYyWuLA ✅ getReward()
 5 2022-06-25 20:00:09 1777e6b96ba06f19eafc6345694841917106e49a44747e607c1b52e3ab51396e TAFotzexiiUJzGkBHDy9Jbn7rVHoYyWuLA ✅ getReward()
 6 2022-06-24 20:00:09 bab6d8ac4ce421f61a0a58dd13f2ad1d4410598a5d68a8cd843711399e459064 TAFotzexiiUJzGkBHDy9Jbn7rVHoYyWuLA ✅ getReward()
 7 2022-06-23 20:18:00 f7f33e53df3040d710da6041dbf16e0679d21a321c2155ccee52258d4fb0213e TAFotzexiiUJzGkBHDy9Jbn7rVHoYyWuLA ✅ getReward()
 8 2022-06-22 20:00:09 dc3efc5f33741cea12898f4eb66c089dd469f08e4f580d0839fa4494dd99d7e1 TAFotzexiiUJzGkBHDy9Jbn7rVHoYyWuLA ✅ getReward()
 9 2022-06-21 20:00:09 ac6951c2c1c836311a0ebd6497fd5a27e0a8df69fd46e6a182cc11740f139fd4 TAFotzexiiUJzGkBHDy9Jbn7rVHoYyWuLA ✅ getReward()
 ...
```

- `tx`

```shell
$ tt scan tx main 4b1118a8303b23e2ef8ddd9b6b6ce5de2638f18f1d03435e692fc01e3254fd20
[Return data]:
  - In HEX: 08C379A00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000A6C6F772072657475726E00000000000000000000000000000000000000000000
  - In ASCII:y ↵low return
[From]: TGaXLDXtndZwrjrCSdaByqXubfbE1NdbqU
[To]: TP2igjG2ofLq495iz6fW93JrzeVTZRfRA6
[Method]: mine_route_u_three(address[] _tokenIn,uint256 _amountIn,uint256[] _route,uint256 _minReturn)
  - [Arg-00]: address[]
    - [slice-00]: address, 0x891cdb91d149f23B1a45D9c5Ca78a88d0cB44C18 - TNUC9Qb1rRpS5CbWLmNMxXBjyFoydXjWFR
    - [slice-01]: address, 0x0Efac3802727c5F873b887e8119fe895B5156577 - TBLQs7LqUYAgzYirNtaiX3ixnCKnhrVVCe
    - [slice-02]: address, 0xa614f803B6FD780986A42c78Ec9c7f77e6DeD13C - TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t
    - [slice-03]: address, 0x891cdb91d149f23B1a45D9c5Ca78a88d0cB44C18 - TNUC9Qb1rRpS5CbWLmNMxXBjyFoydXjWFR
  - [Arg-01]: uint256, 3209302261 - 3,209,302,261 (10)
  - [Arg-02]: uint256[]
    - [slice-00]: uint256, 30
    - [slice-01]: uint256, 14
    - [slice-02]: uint256, 12
  - [Arg-03]: uint256, 3209302261 - 3,209,302,261 (10)
```

### Command `now`

#### Examples

```shell
$ tt now
     [in sec] - 1659926195
   [in milli] - 1659926195804
[in datetime] - 2022-08-08 10:36:35.80475 +0800 CST m=+0.010608918

$ tt now 1659926195
  [if sec] - 2022-08-08 10:36:35

$ tt now 1659926195804
[if milli] - 2022-08-08 10:36:35

$ tt now "2022-08-08 10:36:35"
     [in sec] - 1659926195
   [in milli] - 1659926195000
[in datetime] - 2022-08-08 10:36:35 +0800 CST

$ tt now "2022-08-08"
     [in sec] - 1659888000
   [in milli] - 1659888000000
[in datetime] - 2022-08-08 00:00:00 +0800 CST

$ tt now "08-08"
     [in sec] - 1659888000
   [in milli] - 1659888000000
[in datetime] - 2022-08-08 00:00:00 +0800 CST

$ tt now "08-08 10:36:35"
     [in sec] - 1659926195
   [in milli] - 1659926195000
[in datetime] - 2022-08-08 10:36:35 +0800 CST

$ tt now "10:36:35"
     [in sec] - 1659926195
   [in milli] - 1659926195000
[in datetime] - 2022-08-08 10:36:35 +0800 CST
```

### Command `call`

#### Examples

```shell
$ tt call main TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t
 1. MAX_UINT()
 2. _totalSupply()
 3. addBlackList(address)
 4. allowance(address,address)
 5. approve(address,uint256)
 6. balanceOf(address)
 7. basisPointsRate()
 8. calcFee(uint256)
 9. decimals()
10. decreaseApproval(address,uint256)
11. deprecate(address)
12. deprecated()
13. destroyBlackFunds(address)
14. getBlackListStatus(address)
15. increaseApproval(address,uint256)
16. isBlackListed(address)
17. issue(uint256)
18. maximumFee()
19. name()
20. oldBalanceOf(address)
21. owner()
22. pause()
23. paused()
24. redeem(uint256)
25. removeBlackList(address)
26. setParams(uint256,uint256)
27. symbol()
28. totalSupply()
29. transfer(address,uint256)
30. transferFrom(address,address,uint256)
31. transferOwnership(address)
32. unpause()
33. upgradedAddress()

Which method you want to call: 28
You choose method: [totalSupply() view returns(uint256)]
[Trigger Result]
  - success
[Energy Used]
  - 519
[Return Data]
  - 33130268679280810

Which method you want to call: 29
You choose method: [transfer(address _to, uint256 _value) returns(bool)]
Please input arguments:
 - _to: TMuA6YqfCeX8EhbfYEg5y7S4DqzSJireY9
 - _value: 1000000
Please input from address (default zero address): TMuA6YqfCeX8EhbfYEg5y7S4DqzSJireY9
[Trigger Result]
  - success
[Energy Used]
  - 14650
[Return Data]
  - false
[Logs]
{
	address: TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t,
	data: 0x00000000000000000000000000000000000000000000000000000000000f4240,
	topics: [
		0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef,
		0x00000000000000000000000082dd6b9966724ae2fdc79b416c7588da67ff1b35,
		0x00000000000000000000000082dd6b9966724ae2fdc79b416c7588da67ff1b35,
	]
},

Which method you want to call: 29
You choose method: [transfer(address _to, uint256 _value) returns(bool)]
Please input arguments:
 - _to: TMuA6YqfCeX8EhbfYEg5y7S4DqzSJireY9
 - _value: 1000000
Please input from address (default zero address):
[Trigger Result]
  - REVERT opcode executed
[Energy Used]
  - 1984

Which method you want to call: ^C
```