package errorlist

import "errors"

var (
	ErrCacheNotFound = errors.New("value not found in cache")
)
