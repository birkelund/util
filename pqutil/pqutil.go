package pqutil

import "github.com/lib/pq"

func IsUndefinedObjectError(err *pq.Error) bool {
	return err.Code == "42704"
}

func IsUndefinedTableError(err *pq.Error) bool {
	return err.Code == "42P01"
}

func IsInvalidPasswordError(err *pq.Error) bool {
	return err.Code == "28P01"
}
