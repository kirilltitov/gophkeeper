package gophkeeper

import "errors"

var ErrEmptyLogin = errors.New("login string is empty")
var ErrEmptyPassword = errors.New("password string is empty")
var ErrAuthFailed = errors.New("wrong login or password")
var ErrNoAuth = errors.New("user not authorized for this action")
var ErrAuthExpired = errors.New("your authorization has expired")
