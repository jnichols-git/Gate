# pkg/authserver

`authserver` ties the rest of the `auth` packages together as an https-accessible API for user verification.

## API Endpoints

The `auth` api will be deployed as a subdomain `auth`. Initial setup will provide an API key. For an application
with domain `domain.com`, requests should be main to `auth.domain.com` using header authorization `x-api-key: [key]`
as specified below.

### Authentication

Authentication endpoints take a JSON object, defined in `authserver.go` as AuthRequestBody, that has 3 fields:

- `forUser` (string): email being used to authenticate
- `authCode` (string): email authentication code
- `authToken` (string): authentication token from prior email authentication

Endpoints will specify which fields are needed below, and will ignore the other fields; they can safely
be left empty if they aren't being used.

- POST `/mail`: Sends an authentication code to `forUser`
    - `200 OK`: Authentication email sent to `forUser`. Does not guarantee delivery.
    - `400 Bad Request`: Request was poorly-formed; body should provide more information on the error.
- POST `/code`: Validates `authCode` for `forUser`
    - `200 OK`:  `authCode` was valid. Body contains a bearer token.
    - `400 Bad Request`: Request was poorly-formed; body should provide more information on the error.
    - `401 Unauthorized`: Authorization failed due to incorrect or expired `authCode`.
- POST `/token`: Validates `authToken`
    - `200 OK`: `authToken` was valid. Body contains the decoded contents of `authToken`.
    - `400 Bad Request`: Request was poorly-formed; body should provide more information on the error.
    - `401 Unauthorized`: Authorization failed due to incorrect or expired `authToken`.
