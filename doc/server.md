# pkg/server

`gateserver` ties the rest of the `gate` packages together as an https-accessible API for user verification.

---
## Gate Dashboard
---

As part of the `server` package, the Dashboard provides a place to check the health of your server,
download logs, and set permissions for certain users (banned, admin, etc.). You can access your dashboard
from `https://gate.[domain]/dashboard`.

### SMTP Configuration

You can set your SMTP host, and which email you'd like to send from, through the dashboard. These changes
are saved (see Persistent Settings), so you can be sure your server is getting authentication emails out
as intended. The TestEmail field in `dat/config/config.yaml` serves as an email to send test messages to;
when you submit new details, the Dashboard will automatically send out a blank email to that address to
make sure it can do so without error. The icon next to SMTP: indicates whether your configuration is
working.

### Database Configuration

Gate currently uses a local database through [gorm](https://github.com/go-gorm/gorm), a developer-friendly
Object-Relational Mapping library for Go. The database is implemented using sqlite, with a configuration option
for the path to the database (default `dat/database/auth.db`). More configuration is planned for future updates
regarding networked implementations and other SQL databases.

### TLS Certificate

The dashboard, on each refresh, will attempt to connect to `gate.domain` using TLS; the icon next to
TLS Certificate: indicates if this attempt is successful. The dashboard will additionally describe the error
if one occurs.

The dashboard can upload a `.crt` and `.key` file for TLS.

---
## API Specification
---

The `gate` api will be deployed as a subdomain `gate`. Initial setup will provide (TODO) provide an API key. For an application
with domain `domain.com`, requests should be main to `gate.domain.com` using header authorization `x-api-key: [key]`
as specified below.

### Authentication

Authentication endpoints take a JSON object, defined in `server.go` as AuthRequestBody. Possible arguments
are kept strictly defined in this way to avoid reading extraneous data. Possible fields are:

- `email` (string) (validated server-side)
- `username` (string)
- `password` (string)
- `newPassword` (string)
- `gateCode` (string)
- `getToken` (bool)
- `gateKey` (string): gate key from prior email or credential authentication

Each field behaves differently depending on which endpoint is being called. Any field not listed for an endpoint
will not be used; preferably they should not be included in queries.

- POST `/register`: User registration
    - Parameters
        - `email`: User email. Must be unique.
        - `username`: User username. Must be unique.
        - `password`: User password.
    - Responses
        - `200 OK`: User was registered in the auth server database.
        - `400 Bad Request`: Catch-all for registration errors; see contents for error information.
- POST `/login`: User login credential checking
    - Parameters
        - `username`: Username
        - `password`: Password
        - `getKey` (optional): Whether to return a gate key representing successful sign on.
    - Responses
        - `200 OK`: User credentials match a user in the server database. If `getToken`, body contains a bearer token.
        - `400 Bad Request`: Catch-all for login errors; see contents for error information
        - `401 Unauthorized`: User credentials are incorrect.
- POST `/resetPassword`: User password changes
    - Parameters
        - `username`: Username
        - `password`: The user's *old* password
        - `newPassword`: The user's desired *new* password
    - Responses
        - `200 OK`: User password updated successfully.
        - `400 Bad Request`: Catch-all for password reset errors; see contents for error information.
        - `401 Unauthorized`: User credentials (username/password) are incorrect.
- POST `/mail`: Sends an email with an authentication code.
    - Parameters
        - `email`: Target address
    - Responses
        - `200 OK`: Email was successfully *sent*. Golang SMTP does not throw on email bounce/complaint; response `200` does not guarantee successful delivery.
        - `400 Bad Request`: Request was poorly-formed; see contents for error information.
- POST `/code`: Validates `gateCode` for `email`
    - Parameters
        - `email`: Email address to which the authentication code was sent
        - `gateCode`: Received validation code
        - `getKey` (optional): Whether to return a gate key representing successful sign on
    - Responses
        - `200 OK`:  `gateCode` was valid. If `getToken`, body contains a bearer token.
        - `400 Bad Request`: Request was poorly-formed; see contents for error information.
        - `401 Unauthorized`: Authorization failed due to incorrect or expired `gateCode`.
- POST `/key`: Validates `gateKey`
    - Parameters
        - `gateKey`: Gate key provided with earlier authentication
    - Responses
        - `200 OK`:  `gateKey` was valid (signed by server and unmodified). Body contains the decoded contents of `gateKeyToken`.
        - `400 Bad Request`: Request was poorly-formed; see contents for error information.
        - `401 Unauthorized`: Authorization failed due to incorrect or expired `gateKey`.
