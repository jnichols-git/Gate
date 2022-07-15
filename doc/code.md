# pkg/authcode

The `authcode` package provides support for generating and verifying authentication codes.

## Usage

`authcode` primarily exposes 2 functions: `NewAuthCode` and `ValidateAuthCode`.

- `NewAuthCode` takes an email and generates an `authorizationCode` that expires 1 minute after its creation.
This code is stored in memory until validation is attempted on it.
- `ValidateAuthCode` takes an email and string code and checks if it is valid; the email must be mapped to that
specific string code, and the expiration must not have passed. The authorization code is then immediately removed
from memory.
