
package jwtutil

import (
	errors "errors"
)

var (
	NULL_POINT_ERR = errors.New("Null point error")
	SPLIT_ERROR = errors.New("JWT string split error")
	MAC_NOT_SAME_ERROR = errors.New("The mac isn't right")
	TOKEN_OUT_DATE_ERROR = errors.New("JWT is outdate")
)
