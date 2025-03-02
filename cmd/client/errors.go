package main

import "errors"

var errAPIEndpointNotFound = errors.New("api endpoint not found")
var errNoAPIResponse = errors.New("api returned zero length response")
var errAuthExpired = errors.New("authentication expired, relogin required")
var errNoAuth = errors.New("not authenticated")
