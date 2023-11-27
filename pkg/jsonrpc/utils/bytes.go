package utils

import "github.com/unicornultrafoundation/go-u2u/libs/common"

func (u *Utils) LeftPadBytes(slice []byte, l int) []byte {
	return common.LeftPadBytes(slice, l)
}
