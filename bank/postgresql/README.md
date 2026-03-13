# Plain PostgreSQL implementation of the Bank scenario

## Architecture

Client -> PostgreSQL

## Set up PostgreSQL

Create an empty database:

```shell
initdb -D <data-directory>
```

Start PostgreSQL:

```shell
postgres -D <data-directory> -h 127.0.0.1 -p 5432 -k ""
```

Connect and create the `benchmark` database:

```shell
psql -h 127.0.0.1 postgres -c "CREATE DATABASE benchmark;"
```

## Generating the requests

Use the [generate-requests](../generate-requests/README.md) project to generate the requests.

## Running the benchmark

Build and run the benchmark:

```shell
go build .
./benchmark -input <requests-file>
```

### Verification only

If needed (as specified by the benchmark), you can verify the database contents:

```shell
./benchmark -input <requests-file> -verify-only
```

This will check the total balance of all accounts and the number of accounts.
