# ldbc
Command-line client for LevelDB

## Installation

```bash
go get github.com/akiomik/ldbc
```

## Usage

```bash
# create new database
$ ldbc new awesome-data

# launch REPL
$ ldbc console awesome-data
ldbc - 0.0.1

ldbc> keys

Found: 0 keys

ldbc> put scala play
Created: scala

ldbc> put haskell yesod
Created: haskell

ldbc> get scala
play

Got: scala

ldbc> keys
haskell
scala

Found: 2 keys

ldbc>
```

## Other solutions
* [lev](https://github.com/hij1nx/lev)
* [RockDB Ldb tool](https://github.com/facebook/rocksdb/wiki/Ldb-Tool)
