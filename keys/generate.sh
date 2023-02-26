#!/bin/sh

openssl genpkey -algorithm ed25519 -out jwt_access_private.pem &&
openssl pkey -in jwt_access_private.pem -pubout -out jwt_access_public.pem

openssl genpkey -algorithm ed25519 -out jwt_refresh_private.pem &&
openssl pkey -in jwt_refresh_private.pem -pubout -out jwt_refresh_public.pem
