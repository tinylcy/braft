package raft

import (
	"crypto/sha256"
	"fmt"
)

func SHA256(elem interface{}) (string, error) {
	buf := GetBytes(elem)

	sum := sha256.Sum256(buf)
	ret := fmt.Sprintf("%x", sum)
	return ret, nil
}
