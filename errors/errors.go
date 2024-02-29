package errors

import "errors"

var (
	ErrProductSKUAlreadyExist = errors.New("product sku already exist")
	ErrProductNotFound        = errors.New("product not found")
)
