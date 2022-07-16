# pkg/authjwt

The `jwt` package implements basic usage of JSON Web Tokens to maintain authentication for server-side verification
without explicit states or password storage. If you want to know more about how this works, [jwt.io](jwt.io) does a
much better job explaining it than I could here. A short explanation is below:

## JSON Web Token Basics

JSON Web Tokens, as implemented in this project, provide authentication by listing user permissions while protecting
against alteration. Generally, a JWT is provided by the server after the user is authenticated in some other manner.
The server notes who the user is and what permissions they have, then signs that data and includes the signature
along with it. `authjwt` tokens appear as follows:

Header:
```
{
    "alg": "sha256"
    "typ": "jwt"
}
```

Body:
```
{
    "iss": "auth"
    "sub": "<some user email or id code>"
    "access": "<some access identifier>"
    "iat": <time of token creation>
    "exp": <time of token expiration>
}
```

When `Export`ing a JWT, the header and body are base-64 encoded and signed using a server-side secret. Then, the following
output is generated:

`base64URL(Header)+"."+base64URL(Body)+"."+base64URL(Signature)`

Since the signature is generated using the information contained by the token and a secret, the token can then only be validated
if the content and signature remain unchanged. If an altered token is passed to `Verify`, a different signature will be generated
and validation will fail (return `false`)

## Usage

`authjwt` exposes:

- `NewJWT(user, access string)`: creates a new auth token with user and access information
- `Export(t JSONWebToken, secret []byte)`: exports token format after signing with the given secret
- `Verify(token string, secret []byte)`: verifies the token and and returns the resulting JSONWebToken struct and a boolean representing verification pass/fail.
The following criteria represents sucessful verification:
    1. Header/Body is not altered after creation
    2. Token was signed with the given secret
    3. Token is not expired


Example usage:

```
newToken := authjwt.NewJWT("some-user", "base-access")
fmt.Printf("Your token: %s\n", authjwt.Export(newToken, secret))

// some time later...
inputToken := <user input>
jwt, valid, _ := authjwt.Verify(inputToken, secret)
fmt.printf("Token for user %s with access %s validity: %t\n", jwt.Body.ForUser, jwt.Body.Access, valid)
```
