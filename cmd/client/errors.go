package main

import "errors"

var errNoAPIResponse = errors.New("api returned zero length response")
var errAuthExpired = errors.New("authentication expired, relogin required")
var errNoAuth = errors.New("not authenticated")

var errInternalServerError = errors.New("server returned internal server error")
var errBadRequest = errors.New("server returned bad request error")
var errAPIEndpointNotFound = errors.New("api endpoint not found")
var errUnauthorized = errors.New("you are unauthorized")
