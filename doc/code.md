# pkg/gatecode

The `gatecode` package provides support for generating and verifying short-term authentication codes.

## Usage

`gatecode` primarily exposes 2 functions: `NewGateCode` and `ValidateGateCode`.

- `NewGateCode` takes an email and generates an `gateCode` that expires 1 minute after its creation.
This code is stored in memory until validation is attempted on it.
- `ValidateGateCode` takes an email and string code and checks if it is valid; the email must be mapped to that
specific string code, and the expiration must not have passed. The authorization code is then immediately removed
from memory.
