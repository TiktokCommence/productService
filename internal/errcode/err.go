package errcode

import "errors"

var (
	ErrProductNotExist = errors.New("product not exist")
	ErrPageExceed      = errors.New("page exceeds the total page")
)
