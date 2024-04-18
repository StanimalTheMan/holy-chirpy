# holy-chirpy

## Summary

This project aims to build a web server using Go and the standard library. Twitter at its core just allows short posts. This project is a very light version for learning purposes without frameworks that aims to moderate inappropriate sentiment by leveraging OpenAI moderation.

## How to use

1.  Use Postman or some other http client

POST /api/users - create a new user by supplying ("email": "YOUR EMAIL", "password": "PASSWORD") in JSON body.

POST /api/login - using the same email and password, make request. Use access token from "token" in http response.

POST /api/chirps - test that only post that isn't flagged by OpenAI moderation is created by supplying {"body": "CONTENT"} in your JSON body. Set Authorization header using token from previous POST /api/login http response.
