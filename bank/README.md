# Bank scenario

A simple benchmark based on the original SpacetimeDB keynote. The aim is to compare transactions per second accross
databases. It is expected that SpacetimeDB will not win this scenario. For a scenario that showcases the promise of
SpacetimeDB being a database and server in one, see the [chat scenario](../chat/README.md)

## Architecture requirements

The benchmark runs against "a database system" from a separate client:

- Any database instance:
  - Must persist data to disk (within reasonable time)
  - It should have feature parity with SpacetimeDB
- One table that stores account balance information:
  - `balances` table with columns `id` and `balance`
- Out-of-the-box configuration for the database instance
    - No additional caching or optimisation

## Benchmark requirements

The benchmark can run from any client. There is no standard benchmarking client/harness - the same way SpacetimeDB did.
There are some rules for the benchmarking clients to keep it reasonable and fair:

- Requests are pre-computed and must be executed from the provided input file
- Configurable:
    - number of accounts
    - initial balance of each account
    - number of concurrent connections
    - number of maximum in-flight requests
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

## Competing systems

This is what I plan to implement here:

- Benchmarking client calling reducers on SpacetimeDB
- Benchmarking client calling stored procedures in PostgreSQL
- Benchmarking client calling API endpoints in server with an embedded database (likely SQLite)

### Production is not like this! (Goals and non-goals)

No, this is not like a production scenario. This is like the SpacetimeDB keynote: far from reality.
It lacks extra business logic, a lifelike database schema, authentication, and many more.

The goal with this benchmark is to compare apples to apples.
This benchmark also tries to be as close as possible to the original SpacetimeDB benchmark.

If you'd like a realistic production scenario, see the [chat one](../chat/README.md).
