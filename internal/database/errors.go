package database

import "errors"

type DBErrors struct {
	UserAlreadyExists error
}

var Errors = DBErrors{
	UserAlreadyExists: errors.New("User already exists."),
}
