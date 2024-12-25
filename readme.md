# Bootdev_Server
This is a guided project from [boot.dev](https://www.boot.dev) to learn building HTTP web server in Go. The project builds a local server application with JSON API, incorporates webhooks, JWTs/Authorization and more.

## Install
The application is build for learning purposes, so there is no inherited installation purpose. If you still want to you can install it with:\
``go install github.com/christopherhanke/bootdev_server``


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
    "id":               "[uuid]",
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
    "id":               "[uuid]",
    "created_at":       "[time]",
    "updated_at":       "[time]",
    "email":            "[string]",
    "token":            "[string]",
    "refresh_token":    "[string]",
    "is_chirpy_red":    "[bool]"
}
```
The ``token`` is used to authorize user specific commands like [``POST /api/chirps``](#post-apichirps) and is viable for an hour. The ``refresh_token`` is used to refresh the ``token`` and is viable for 60 days.


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


### /api/revoke
#### POST /api/revoke
This endpoint is to revoke a given refresh_token. This way a user can not refresh their authorization token. This does not revoke an existing authorization time if still valid. To revoke a refresh token it only needs a HTTP request with the token which should be revoked given in the request header, replace ``refresh_token_string`` with the users refresh_token.
```
HTTP Header
Authorization: Bearer refresh_token_string
```
If all is correct the response is ``204``. 


### /api/chirps
#### POST /api/chirps
Creates a chirp in the database. Needs a authorization token in the request header, replace ``tokenString`` with the authorization token from the user posting the chirp. The chirp needs to come via JSON data in the request, the chirp ``text`` in a body field. The ``text`` is only allowed to have up to 140 characters. If the ``text`` is longer than that, the request will fail with ``400`` response. 
```
HTTP Header
Authorization: Bearer tokenString
```
```
{
    "body": "text"
}
```
If all is well you'll get a ``201`` response with the chirp data in JSON format.
```
{
    "id":           "[uuid]",
    "created_at":   "[time]",
    "updated_at":   "[time]",
    "body":         "[string]",
    "user_id":      "[uuid]"
}
```


#### GET /api/chirps
Get all chirps from database in a JSON response. All chirps will be in a list in descending order by their time of creation.
```
[
    {
        "id":           "[uuid]",
        "created_at":   "[time]",
        "updated_at":   "[time]",
        "body":         "[string]",
        "user_id":      "[uuid]"
    },
    ...
]
```
You can add a query to the request like ``?author_id=IDstring`` to filter the chirps you want to get. At the moment this queries are supported.
* ``authord_id``:   needs the ID_string of the author to filter for their chirps
* ``sort``:         accepts ``asc`` or ``desc`` to sort the chirps accordingly by their time of creation. Default is ``asc``.

With queries added to the request there are possible errors:
* ``404`` given ID_string for an author could not be resolved

#### GET /api/chirps/{chirpID}
Get a specific chirp by its ID replace ``{chirpID}`` with the ID_string of the chirp to get. If the ID_string is invalid there will be a ``404`` response. Otherwise the response will return ``200`` the chirp data in JSON.
```
{
    "id":           "[uuid]",
	"created_at":   "[time]",
	"updated_at":   "[time]",
	"body":         "[string]",
	"user_id":      "[uuid]"
}
```


#### DELETE /api/chirps/{chirpID}
Deleting a specific chirp is only allowed for the author of that chirp. You need the authorization token of that user in the request header and the ID_string of the chirp in question. Replace ``{chirpID}`` with the ID_string of the chirp.
```
HTTP Header
Authorization: Bearer tokenString
```
If no error occurs the response is ``204`` without a response body. 

Possible errors:
* ``401`` authorization token invalid
* ``403`` wrong user authorization
* ``404`` given ID_string could not be resolved to a chirp


### /api/polka/webhooks
#### POST /api/polka/webhooks
The polka/webhooks endpoint realizes the event listener for the polka API. In this example the webhook is used to realize the upgrade of the user to ChirpyRedAccount. The request header needs the authorization key. The request body is in JSON format.
```
{
    "event":    "[string]",
    "data":     {
        "user_id":  "[uuid]"
    }
}
```
Possible errors:
* ``400`` request body did not fit JSON format
* ``401`` authorization failed
* ``404`` given user_id is not valid

