# Generate requests for the Bank scenario

## Format

The first line of the output contains the expected number of accounts and the initial balance. For example:

```text
1000000 100000
```

All following lines contain requests in the format:

```text
<request-id> <expected-success> <from-account> <to-account> <amount>
```

Where:
- `<request-id>` is a unique identifier for the request (starting from 0)
- `<expected-success>` is either `1` for regular requests and `0` if the request is expected to fail
- `<from-account>` is the account ID to transfer money from
- `<to-account>` is the account ID to transfer money to
- `<amount>` is the amount of money that should be transferred

## Probability

Since the implementations of Zipf distributions are not deterministic across different languages - the generator is
using uniform distribution by default. In that mode, `coverage` controls what percentage (specified as a ratio between 0
and 1) of account IDs should be used for the requests. A `0.1` coverage means that only the first 10% of the account IDs
will show up in the requests.

If `zipf` is set, `coverage` controls the alpha/exponent parameter of the Zipf distribution.
