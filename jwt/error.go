
package jwtutil

import (
	util "github.com/zyxgad/go-util/util"
)

var (
	SPLIT_ERROR = util.NewErr("JWT string split error", nil)
	MAC_NOT_SAME_ERROR = util.NewErr("The mac isn't right", nil)
	TOKEN_OUT_DATE_ERROR = util.NewErr("JWT is outdate", nil)
)
