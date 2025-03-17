package gophkeeper

import "errors"

// ErrEmptyLogin is an error indicating missing login field in login request.
var ErrEmptyLogin = errors.New("login string is empty")

// ErrEmptyPassword is an error indicating missing password field in login request.
var ErrEmptyPassword = errors.New("password string is empty")

// ErrAuthFailed is an error indicating missing login in DB or password mismatch.
var ErrAuthFailed = errors.New("wrong login or password")

// ErrNoAuth is an error indicating missing authorization for a certain action.
var ErrNoAuth = errors.New("user not authorized for this action")
