# relevel
REPL for LevelDB

## Requirements

* LevelDB 1.7 or above

## Installation

```bash
go get github.com/akiomik/relevel
```

## Usage

### Create new leveldb data

```bash
# relevel new <name>
relevel new awesome-data
```

### Launch leveldb REPL

```bash
# relevel console <name>
$ relevel console awesome-data
```

### Use command-line interface

```bash
# relevel exec <name> <query>
$ relevel exec awesome-data "put banana monkey"
```

## Queries

* keys

```
relevel> keys
```

* get <key>

```
relevel> get banana
```

* put <key> <value>

```
relevel> put banana monkey
```

* delete <key>

```
relevel> delete banana
```

## Other solutions
* [ldb](https://github.com/hij1nx/ldb)
* [lev](https://github.com/hij1nx/lev)
* [RocksDB Ldb tool](https://github.com/facebook/rocksdb/wiki/Ldb-Tool)
* [level-repl](https://github.com/lapwinglabs/leveldb-repl)
