package entity

import "errors"

var ErrNotFound = errors.New("resource not found")
var ErrHasRelations = errors.New("cannot delete: entity has related records")
var ErrInvalidContent = errors.New("invalid data")
