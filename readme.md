# Bootdev_Server
This is a project from [boot.dev](www.boot.dev) to learn building HTTP web server in Go. The project builds a JSON API, incorporates webhooks and JWTs and more.


## API
Here's an explanation to all API endpoints, how to send your requests and what the responses will look alike. 

If your requests fail or there is an internal error, there will be corresponding responses ``4xx``/``5xx`` with an JSON message. The message contains information about the error and why it occured.
```
{
    "error":    "[string]"
}
```


### /api/users
#### POST /api/users
With this HTTP.Request you can create new users. The request needs JSON data in the following format. Every email can only be used once for an account and the password isn't allowed to be empty.
```
{
    "email":    "[string]",
    "password": "[string]"
}
```
If all is correct the response is ``201`` with JSON data in the following format:
```
{
    "id":               "[uuid.string]",
    "created_at":       "[time]",
    "updated_at":       "[time]",
    "email":            "[string]",
    "is_chirpy_red":    "[bool]"
}
```



#### PUT /api/users
With this HTTP.Request you can update an already existing user. The request needs an authorization token in the header and JSON data in the same format as the [``POST /api/users``](#post-apiusers), but in addition the user needs to be already logged in and to have the authorization token in the request header, replace ``tokenString`` with the users token.
```
HTTP Header
Authorization: Bearer tokenString
```
If all is correct the response is ``200`` with JSON data like the response of [``POST /api/users``](#post-apiusers). If something went off, there will be corresponding responses ``4xx``/``5xx`` with an JSON message.
```
{
    "error":    "[string]"
}
```


### /api/login
#### POST /api/login
This HTTP.Request is to login to existing user. The request needs JSON data in the following format which correspond to the user data.
```
{
    "email":    "[string]",
    "password": "[string]"
}
```
If all is correct the response is ``200`` with JSON data in the following format:
```
{
    "id":               "[uuid.string]",
    "created_at":       "[time]",
    "updated_at":       "[time]",
    "email":            "[string]",
    "token":            "[string]",
    "refresh_token":    "[string]",
    "is_chirpy_red":    "[bool]"
}
```
The ``token`` is used to authorize user specific commands like [``POST /api/chirps``](#post-apichirps) and is viable for an hour. The ``refresh_token`` is used to refresh the ``token`` and is viable for 60 days.

If something went off, there will be corresponding responses ``4xx``/``5xx`` with an JSON message.
```
{
    "error":    "[string]"
}
```


### /api/refresh
#### POST /api/refresh
Request to this endpoint to refresh an already logged in user. If the user is authorized you'll get a new authorization token. The request needs the authorization token in the request header, replace ``refresh_token_string`` with the users refresh_token.
```
HTTP Header
Authorization: Bearer refresh_token_string
```
If all is correct the response is ``200`` with JSON data in the following format, serving a new ``token`` (not ``refresh_token``):
```
{
    "token":    "[string]"
}
```
If something went off, there will be corresponding responses ``4xx``/``5xx`` with an JSON message.
```
{
    "error":    "[string]"
}
```


### /api/revoke
#### POST /api/revoke

### /api/chirps
#### POST /api/chirps

#### GET /api/chirps

#### GET /api/chirps/{chirpID}

#### DELETE /api/chirps/{chirpID}

### /api/polka/webhooks
#### POST /api/polka/webhooks

