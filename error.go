package aster

import "errors"

var (
	errModeInvalid   = errors.New("error mode is invalid")
	errMaskInvalid   = errors.New("error mask is invalid")
	errEmptyByte     = errors.New("error receive empty data")
	errHeaderInvalid = errors.New("header invalid")
)
