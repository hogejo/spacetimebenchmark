# SpacetimeBenchmark

A repository to properly benchmark SpacetimeDB against other systems, following their
[performance keynote video](https://www.youtube.com/watch?v=C7gJ_UxVnSk) and my
[pull request](https://github.com/clockworklabs/SpacetimeDB/pull/4457).

Status:
- I closed the pull request
- I am moving my work here. See the [original/](original/README.md) directory.
- I am developing two scenarios to properly benchmark SpacetimeDB against other systems.

## Scenarios

### [Bank in `bank/`](bank/README.md)

This is an improved version of the original SpacetimeDB benchmark, with a full specification of the requirements.
The scenario is to test sheer DB performance like Clockwork Labs' keynote video: how many TPS can you achieve with your
database?

### [Chat in `chat/`](chat/README.md)

This is a proper benchmark of a real-world scenario. This one aims to properly show SpacetimeDB's advantage - how it is
a database and a server in one package. Because of that, the specification allows competing systems to come in a similar
package or with separate DB and backend.

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
[GitHub repo](https://github.com/clockworklabs/SpacetimeDB/blob/abbcec4ab357b956f6a2d498a666c1516edcb355/templates/keynote-2/README.md).

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

### SpacetimeDB is in a different league

I have to note that SpacetimeDB is not completely free and/or open source. I avoid copying their work in this repository
because of that. See [their licence](https://github.com/clockworklabs/SpacetimeDB/blob/master/LICENSE.txt).

This repository is GPLv3 licensed.
