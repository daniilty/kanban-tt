## Service for jwt authorization

## How to generate OpenSSL RSA Key Pair:
```bash
$ openssl genrsa -des3 -out private.pem 2048
$ openssl rsa -in private.pem -outform PEM -pubout -out public.pem
$ openssl rsa -in private.pem -out private_unencrypted.pem -outform PEM
```

## All avaliable endpoints:

* GET `/api/v1/auth/jwks` - get generated JSON Web Key Set to use it for token validation with proxy(krakend for example)
* GET `/api/v1/auth/me` - get current user info("Authorization: Bearer ..." header is required)
* POST `/api/v1/auth/login` - login user. Body:
```
{
  "email": "sample@sample.com",
  "password": "sample"
}
```
* POST `/api/v1/auth/register` - register user. Body:
```
{
  "email": "sample@sample.com",
  "password": "sample",
  "name": "Biggus Dickus",
  "userName": "sample"
}
```

## Environment variables:

* `PUBKEY` - `export PUBKEY="$(cat public.pem)"`
* `PRIVKEY` - `export PUBKEY="$(cat private_unencrypted.pem)"`
* `USERS_GRPC_ADDR` - grpc address of sharenote-users service
* `HTTP_SERVER_ADDR` - http address of this service
* `TOKEN_EXPIRY` - authorization token expiry(in seconds)

