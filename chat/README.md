# Chat scenario

The benchmark version of the [SpacetimeDB chat tutorial](https://spacetimedb.com/docs/tutorials/chat-app). The aim is to
properly showcase the performance and features of SpacetimeDB while comparing free alternatives with existing
open-source products. Based on their promises, it is expected that SpacetimeDB will win this scenario.

## Changes to the original tutorial

- Private messages are supported
- Users are authenticated
- TLS is mandatory for all communication

## Architecture requirements

See below for client/benchmark requirements. The server-side requirements are:

- The implementation must use a proper database engine
  - Can be embedded into the server or running separately
  - Must persist data to disk (within reasonable time)
- The suggested schema is two tables:
  - `users` for storing user identities
  - `messages` for storing all messages
- Custom configuration and tuning is allowed on the server side
- All requests must be authenticated
- Client communication must be implemented via either:
  - HTTPS API requests and SEE (Server-Sent Events) OR
  - Websockets with TLS

### Example database schema

```sql
CREATE TABLE "users"
(
    identity UUID    PRIMARY KEY,
    name     TEXT    NULL,
    online   BOOLEAN NOT NULL
);

CREATE TABLE "messages"
(
    sender    UUID      NOT NULL,
    recipient UUID      NULL,
    sent      TIMESTAMP NOT NULL,
    text      TEXT      NOT NULL,
    FOREIGN KEY (sender) REFERENCES "users" (identity),
    FOREIGN KEY (recipient) REFERENCES "users" (identity)
);
```

## Benchmark requirements

The benchmark is run from **one standard client** against all server-side implementations. The requirements are:

- Requests are pre-computed and must be executed from the provided input file
- Configurable:
  - number of total users
  - distribution of users among the requests
    - uniform distribution on a portion of the accounts (first N% of users sampled only)
  - number of concurrent connections (users connected)
  - ratio of private vs. public messages
- All connected (online) users must receive notifications about newly connected and disconnected users
- Each request's result must be verified:
  - message notifications must be sent to all connected users (or recipient, if connected)
  - messages must be available to all newly connected users (or recipient, when connected)
- All expected messages must be verified at the end of the benchmark
  - This verification must be repeated after a database shutdown and restart

### Pending requirements

The failure modes of the different architectures are unknown. Because of that, the benchmark is using unlimited
throughput for now: provided that all checks pass, each client/user should send as many messages as possible.

## Client implementation

TODO: reference the client implementation once done
