package main

import "errors"

var errApiEndpointNotFound = errors.New("api endpoint not found")
var errNoApiResponse = errors.New("api returned zero length response")
var errAuthExpired = errors.New("authentication expired, relogin required")
var errNoAuth = errors.New("not authenticated")
