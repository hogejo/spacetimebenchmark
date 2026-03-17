# Bank scenario

A simple benchmark based on the original SpacetimeDB keynote. The aim is to compare transactions per second between
databases. It is expected that SpacetimeDB will not win this scenario. For a scenario that showcases the promise of
SpacetimeDB being a database and server in one, see the [chat scenario](../chat/README.md).

## Architecture requirements

The benchmark runs against "a database system" from a separate client. Both the database system and the client can be
custom implemented. The requirements are:

- Any database instance:
  - Must persist data to disk (within reasonable time)
- One table that stores account balance information:
  - `balances` table with columns `id` and `balance`
- Out-of-the-box configuration for the database instance
  - No additional caching or optimisation
  - Configuration changes and optimisation to match SpacetimeDB (inner) workings are allowed

## Benchmark requirements

The benchmark can run from any client. There is no standard benchmarking client/harness - the same way SpacetimeDB did.
There are some rules for the benchmarking clients to keep it reasonable and fair:

- Requests are pre-computed and must be executed from the provided input file
- Configurable:
  - number of accounts
  - distribution of accounts among the requests
    - uniform distribution on a portion of the accounts OR
    - Zipf distribution
  - initial balance of each account
  - number of concurrent connections
  - number of (total) maximum in-flight requests
- Each request's result must be verified against the expected result (success or failure)
- Some (configurable number of) requests should fail to check the database's robustness:
  - trying to transfer from the same account to the same account
  - trying to transfer from an account that does not exist
  - trying to transfer from an account that does not have enough funds
  - trying to transfer to an account that does not exist
- No account balance can go below zero
  - The database system must prevent this and fail any request that would result in a negative balance
- Warmup is allowed as part of the benchmark run
- At the end of the benchmark the total sum of money must not change in the database, which the client must verify
  - This verification must be repeated after a database shutdown and restart

NOTE: Uniform distribution on a portion of the accounts is preferred, as Zipf implementations can vary between languages
and libraries.

### Production is not like this! (Goals and non-goals)

No, this is not like a production scenario. This is like the SpacetimeDB keynote: far from reality.
It lacks extra business logic, a lifelike database schema, authentication, and many more.

The goal with this benchmark is to compare apples to apples.
This benchmark also tries to be as close as possible to the original SpacetimeDB benchmark.

If you'd like a realistic production scenario, see the [chat one](../chat/README.md).

## Competing systems

### SpacetimeDB

Implementation is pending. Maybe Clockwork Labs will join the effort.

### PostgreSQL

The architecture is similar to the original version I wrote in response to the SpacetimeDB keynote. Limiting the global
maximum in-flight requests is not implemented yet.

A PostgreSQL pool (via [pgx](https://github.com/jackc/pgx)) is used to manage connections. Connections are used by
goroutine workers. Both the number of connections and the number of goroutines are configurable.

### SQLite

A simple Go server accepts TCP connections and responds to RPCs:

- `transfer <from> <to> <amount>`
- `get <account>`

## Initial results

I've only run small tests on small cloud instances (4 CPUs with 16GB RAM). The results are **scaled up** to match the
original SpacetimeDB benchmark.

| System      | Results                                  |
|-------------|------------------------------------------|
| SpacetimeDB | *expected 100k TPS<br/>based on keynote* |
| PostgreSQL  | **~ 88.000 TPS**                         |
| SQLite      | *< 2.000 TPS*                            |

### Infrastructure

Client and server on separate machines:

- Hetzner Cloud
- 4 dedicated CPUs
- 16GB RAM
- SSD storage

Multiplier to match SpacetimeDB benchmark: x4
