# Rebuilding Redis in Go

A Redis server built from scratch in Go, written by hand as a learning project for the CodeCrafters "Build Your Own Redis" challenge. No framework, no Redis client libraries. Just a TCP server, a RESP protocol parser, and a command registry I refactored as the command set grew.

Building it was the point. Letting an AI autocomplete it would have been faster, but would have taught me nothing, so the implementation code is mine. The design notes and the mistakes are written up in a six part series on my blog:

- Part 1, Bytes, handlers, and resisting autocomplete: https://ryanflores.dev/blog/rebuilding-redis-in-go-part-1
- Part 2, Lists and concurrency: https://ryanflores.dev/blog/rebuilding-redis-in-go-part-2
- Part 3, Streams and the art of the naive implementation: https://ryanflores.dev/blog/rebuilding-redis-in-go-part-3
- Part 4, Transactions and the big refactor: https://ryanflores.dev/blog/rebuilding-redis-in-go-part-4
- Part 5, Optimistic locking with WATCH: https://ryanflores.dev/blog/rebuilding-redis-in-go-part-5
- Part 6, Replication and the WAIT command: https://ryanflores.dev/blog/rebuilding-redis-in-go-part-6

## What is implemented

- **Strings and basics:** PING, ECHO, SET (with expiry), GET, INCR, TYPE
- **Lists:** RPUSH, LPUSH, LRANGE, LLEN, LPOP, BLPOP (blocking)
- **Streams:** XADD, XRANGE, XREAD, including blocking reads
- **Transactions:** MULTI, EXEC, DISCARD, with a command result model that let every command return errors the same way
- **Optimistic locking:** WATCH, UNWATCH, using dirty flags on watched keys
- **Replication:** REPLCONF, PSYNC, WAIT, INFO replication, propagation of writes from a master to replicas, and the empty RDB transfer on sync
- **RDB configuration:** CONFIG GET for the persistence directory and filename

## How it is organized

- `resp/`: encoding and decoding of the RESP protocol.
- `commands/`: one package per command family (basic, lists, streams, transactions, optimistic locking, replication). Commands register themselves in `commands/registry.go` as a typed `Command` value, so a command can carry more than a function later (description, help text) without another refactor.
- `app/`: the TCP server and per connection handling.
- `cli/`: server setup commands that do not belong to a client connection, such as port configuration and replica setup.

Concurrency is handled with mutexes around shared state. That was a deliberate order for learning: understand the locking first, then judge where channels would actually help.

## Run it

```sh
./your_program.sh
```

This starts the server on port 6379 by default, the entry point the CodeCrafters runner expects. From another terminal:

```sh
redis-cli PING
```

## Status

The challenge work is parked with replication complete. RDB persistence (loading a real RDB file) is the unfinished piece.
