# SpacetimeBenchmark

A repository to properly benchmark SpacetimeDB against other systems, following their
[performance keynote video](https://www.youtube.com/watch?v=C7gJ_UxVnSk) and my
[pull request](https://github.com/clockworklabs/SpacetimeDB/pull/4457).

Status: I am currently living my regular life and waiting for SpacetimeDB's response.

## An improved scenario

THIS IS A WORK IN PROGRESS DEFINITION OF A BENCHMARK SCENARIO.

- Let's keep it simple and close to the original SpacetimeDB benchmark
- Let's not create another database benchmarking suite or comparison
  - Unless SpacetimeDB and others put effort into this
- Let's keep it fun but respectful

Most important question:
- Must there be an RPC server between the client and the database?
- If so it must have a function and that must be benchmarked too.

### WIP: Scenario requirements

This is what I plan to implement in coming commits:

- One table that stores account balance information
- One database instance
- Out-of-the-box configuration for the database instance
    - No caching or optimisation
- Requests can be pre-computed and cached, or generated on the fly
- Configurable:
    - number of accounts
    - initial balance of each account
    - number of concurrent connections
    - number of maximum in-flight requests
- Each request's result must be verified against the expected result
- Some (configurable number of) requests should fail to check the database's robustness:
    - trying to transfer from the same account to the same account
    - trying to transfer from an account that does not exist
    - trying to transfer from an account that does not have enough funds
    - trying to transfer to an account that does not exist
- No account balance can go below zero
  - The system must prevent this and fail any request that would result in a negative balance

### Planned systems:

Yeah, just moving the code I wrote for the PR into a separate repository:

- Benchmarking client calling reducers on SpacetimeDB
- Benchmarking client calling stored procedures in PostgreSQL

### Production is not like this! (Goals and non-goals)

No, this is not like a production scenario.
It lacks extra business logic, a lifelike database schema, authentication, etc.

The goal with this benchmark is to compare apples to apples.
This benchmark also tries to be as close as possible to the original SpacetimeDB benchmark. See below.

## Background

SpacetimeDB 2.0 was released with the aforementioned video on YouTube. In the video, they claimed that SpacetimeDB was
"1000x faster" (according to the title of the video) and provided the following performance numbers:

| System                            | TPS (~0% Contention) | TPS (~80% Contention) |
|-----------------------------------|----------------------|-----------------------|
| **SpacetimeDB**                   | **107,850**          | **103,590**           |
| SQLite + Node HTTP + Drizzle      | 7,845                | 7,652                 |
| Bun + Drizzle + Postgres          | 7,115                | 2,074                 |
| Postgres + Node HTTP + Drizzle    | 6,429                | 2,798                 |
| Supabase + Node HTTP + Drizzle    | 6,310                | 1,268                 |
| CockroachDB + Node HTTP + Drizzle | 5,129                | 197                   |
| PlanetScale + Node HTTP + Drizzle | 477                  | 30                    |
| Convex                            | 438                  | 58                    |

This shows more than a 10x performance improvement over other systems.

They also provided the source code for the benchmark at their
[GitHub repo](https://github.com/clockworklabs/SpacetimeDB/blob/master/templates/keynote-2/README.md).

### Details of the SpacetimeDB benchmark

As they mention in their README (but not in the video):

> When running the benchmark for SpacetimeDB on higher-end hardware we found out that we were actually bottlnecked on
> our test TypeScript client. To get the absolute most out of the performance of SpacetimeDB we wrote a custom Rust
> client that allows us to send a much larger number of requests then we could otherwise. We didn't do this for the
> other backends/databases as they maxed out before the client.

Thus, the comparison is done with **very** different architectures:

| System                            | Architecture                                            |
|-----------------------------------|---------------------------------------------------------|
| SpacetimeDB                       | Integrated platform (Rust)                              |
| SQLite + Node HTTP + Drizzle      | Node.js HTTP server → Drizzle ORM → SQLite              |
| Bun + Drizzle + Postgres          | Bun HTTP server → Drizzle ORM → PostgreSQL              |
| Postgres + Node HTTP + Drizzle    | Node.js HTTP server → Drizzle ORM → PostgreSQL          |
| Supabase + Node HTTP + Drizzle    | Node.js HTTP server → Drizzle ORM → Supabase (Postgres) |
| CockroachDB + Node HTTP + Drizzle | Node.js HTTP server → Drizzle ORM → CockroachDB         |
| PlanetScale + Node HTTP + Drizzle | Node.js HTTP server → Drizzle ORM → PlanetScale (Cloud) |
| Convex                            | Integrated platform                                     |

### The promise of SpacetimeDB

> SpacetimeDB is a database that is also a server.

The promise of SpacetimeDB - that is highlighted very well in their video - is that you can say goodbye to the usual
overhead that comes with deploying a "traditional" architecture consisting of a database and a server.

### Alternative architectures

To achieve good results and performance in any scenario, one has to do some work and properly plan the architecture.
At the end of that process, SpacetimeDB could come out as a good alternative. Radically different architectures can
emerge too, like they did before. Think DuckDB + Parquet files for example.

Throwing that aside and trying to stay within the scope of the benchmark, there are a couple of architectures I can
think of that can similarly reduce the overhead:

- An application server with an embedded database
  - Java + H2
  - Go + BoltDB
  - Rust + Sled
  - Your favourite language + SQLite
- A database server with stored procedures, notification channels with or without an application server in front of it:
  - PostgreSQL + PL/pgSQL ([or other](https://wiki.postgresql.org/wiki/PL_Matrix))

Feel free to help me extend this list with your favourite combination. Let's keep it simple and close to the original
SpacetimeDB benchmark.

### SpacetimeDB is in a different league

I have to note that SpacetimeDB is not completely free and/or open source.
See [their licence](https://github.com/clockworklabs/SpacetimeDB/blob/master/LICENSE.txt).
