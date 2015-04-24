# relevel
REPL for LevelDB

## Installation

```bash
go get github.com/akiomik/relevel
```

## Usage

```bash
# create new database
$ relevel new awesome-data

# launch REPL
$ relevel console awesome-data
relevel - 0.0.1

relevel> keys

Found: 0 keys

relevel> put scala play
Created: scala

relevel> put haskell yesod
Created: haskell

relevel> get scala
play

Got: scala

relevel> keys
haskell
scala

Found: 2 keys

relevel>
```

## Other solutions
* [lev](https://github.com/hij1nx/lev)
* [RockDB Ldb tool](https://github.com/facebook/rocksdb/wiki/Ldb-Tool)
