# pkg/authserver

`authserver` ties the rest of the `auth` packages together as an https-accessible API for user verification.

## Basic authentication flow

For the purposes of specifying the `authserver` authentication pattern, the following key words are defined:

- AuthServer: This server running to authenticate users.
- Application: An external application using AuthServer for user authentication
- Client: A user that needs to be authenticated to access Application

### Initial Email Authentication

Initial authentication occurs when a bearer token is not provided or is not verified by AuthServer.
Application should handle fulfillment of original requests, as this process takes two requests from
Client. After this process, the client will receive a bearer token that streamlines the authentication
process.

1. Client sends a request with no authentication and an email to Application
2. Application sends a request to authenticate the email to AuthServer
3. Client sends a request with no authentication and a code to Application
4. Application sends a request to verify the code to AuthServer
    - If the code is verified, proceed.
    - Otherwise, Application notifies Client of failed authentication.
5. Application sends a request for a bearer token to AuthServer
6. Application sends the bearer token to Client
7. Application fulfills the Client's original request.

### JWT Verification

JWT verification occurs when a bearer token is provided. As this can apply to any request, it
requires no redirects and should only fail when the token expires or is altered.

1. Client sends a request with "Authorization: Bearer A.B.C" and an email to Application
2. Application sends a request to verify the bearer token to AuthServer
    - If the token is verified, proceed.
    - Otherwise, the user must re-authenticate using email.
3. Application fulfils the Client's original request.
