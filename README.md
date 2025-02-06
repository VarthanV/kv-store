# kv-store

A [Redis serialization protocol specification(RESP)](https://redis.io/docs/latest/develop/reference/protocol-spec/) compliant key-value store written in Golang. An attempt to decode how Redis works under the hood and understand the design decisions.

## Features

- RESP Protocol Implementation
- In-memory key-value store
- Thread-safe operations
- Supports multiple data structures (Strings, Lists, Hashes)

## Supported Commands

### String Operations
- `PING [message]` - Test server connectivity
- `SET key value` - Store a key-value pair
- `GET key` - Retrieve value by key
- `DEL key` - Delete a key
- `INCR key` - Increment the integer value of a key by one
- `DECR key` - Decrement the integer value of a key by one
- `APPEND key value` - Append a value to an existing string

### Hash Operations
- `HSET key field value [field value ...]` - Store hash field-value pairs
- `HGET key field` - Retrieve hash field value
- `HGETALL key` - Retrieve all field-value pairs of a hash

### List Operations
- `LPUSH key value [value ...]` - Insert elements at the head of a list
- `RPUSH key value [value ...]` - Insert elements at the tail of a list
- `LPOP key` - Remove and return the first element of a list
- `RPOP key` - Remove and return the last element of a list

## Getting Started

### Prerequisites

- Go 1.16 or higher

### Installation

```bash
git clone https://github.com/VarthanV/kv-store.git
cd kv-store
go build
```

### Running the Server
```bash
./kv-store
```
The server starts on port 6363 by default.


### Usage
You can interact with the server using any Redis client. For example, using ``redis-cli ``

```bash
redis-cli -p 6363
```
### Example Commands
```bash
127.0.0.1:6363> PING
PONG

127.0.0.1:6363> SET counter 10
OK

127.0.0.1:6363> INCR counter
11

127.0.0.1:6363> HSET user:1 name "John" age "30"
OK

127.0.0.1:6363> HGETALL user:1
1) "name"
2) "John"
3) "age"
4) "30"
```

## Project Structure
- ``command/`` - Command handlers and implementation
- ``pkg/``- Core packages and utilities
- ``resp/`` - RESP protocol implementation

## Implementation Details
- Thread-safe operations using mutex locks
- In-memory storage with map data structures
- RESP protocol compliance for Redis client compatibility

## Pipeline

- [] Implement registry pattern for dynamic command registration and cleaner code base
- [] Append only file log persistence - WIP