# JWT service using Ethereum / ENS identities

A goland API that implements a challenge/response mechanism in order to 
issue JWTs using Ethereum identities.

The Openhouse team has modified this API to issue JWTs signed with HS256 so that
they can be used in authentication with a Jitsi server.

## Install

```bash
$ go get github.com/openhouse-project/eth-auth
```

## Starting the server

The server needs the following ENV variables to be set before running:

* `ORIGINS` - list of allowed Origins separated by space (your client app URL)
* `ETH_PRIVKEY` - a private key that is used to sign the JWTs. This should be shared with your Jitsi server.
* `LOGGING` - whether to log requests to stdout

Example of how to start the server:

```bash
$ export ORIGINS="http://localhost:8888 https://example.org"
$ export ETH_PRIVKEY="your-exported-eth-key-in-hexa"
$ export LOGGING="true"
$ go run server.go
```

## Jitsi server configuration

If you configure Jitsi using the [Docker instructions](https://jitsi.github.io/handbook/docs/devops-guide/devops-guide-docker#lets-encrypt-configuration), there are a few variables in the
.env file you will need to modify. This is a checklist so you don't miss anything;
many of these are unrelated to this auth server:

* ENABLE_LETSENCRYPT if you want to call Jitsi with https
* LETSENCRYPT_DOMAIN
* LETSENCRYPT_EMAIL
* ENABLE_HTTP_REDIRECT is also for LetsEncrypt
* ENABLE_AUTH to enable authentication at all
* ENABLE_GUESTS - leave uncommented to require JWTs
* AUTH_TYPE=jwt
* JWT_APP_ID=openhouse_client (or change to your app ID configured in auth.go)
* JWT_APP_SECRET={put your shared secret here, the same as `ETH_PRIVKEY` above}


## How to use the API

### Obtaining the challenge
Basically, the way the API works is that a client will send a `GET` request to
`/login/{ethAddress}` in order to obtain a unique challenge from the server
that will then be presented to the user in order to be signed with the user's
Ethereum key.

Replace `http://api.example.org` with your own domain.

Request:
```bash
curl 'http://api.example.org/login/0x91ff16a5ffb07e2f58600afc6ff9c1c32ded1f81'
```

Response:
```js
{
  address: "0x91ff16a5ffb07e2f58600afc6ff9c1c32ded1f81",
  challenge: "JiqPLBbLBdCfWZoS"
}
```

### Obtaining the JWT

Next, the client must send a `POST` request to `/login/{ethAddress}`, containing the
signed challenge. The API then will validate the user's signature and issue the JWT
if the signature is good, together with the token's expiration time.

Request:
```bash
curl 'http://api.example.org/login/0x91ff16a5ffb07e2f58600afc6ff9c1c32ded1f81' \
  -X POST \
  -H 'Content-Type: application/json' \
  --data-binary '{"signature": "0x5114fb7...33f5c031c"}'
```

Response:
```js
{
  expires: "2020-11-06T15:06:38.602022706Z",
  token: "eyJleHAiOjE....fyHd7kPlg",
  user: "0x350F72a69D....67C2EBE98dA"
}
```

### Refreshing a JWT before expiration

Clients can obtain a new token as they get closer to the expiration time,
by sending a `GET` request to `/refresh` using the (still valid) JWT as a `Bearer`
token within an `Authorization` header. The response is similar to the one above,
containing a new expiration date and token:

Request:
```bash
curl 'http://api.example.org/refresh' \
  -H 'Authorization: Bearer eyJleHAiOjE....fyHd7kPlg'
```

Response:
```js
{
  expires: "2020-11-06T15:06:38.602022706Z",
  token: "eyJleHAiOjE2MD....fY1qv8Oxjw",
  user: "0x350F72a69D....67C2EBE98dA"
}
```



