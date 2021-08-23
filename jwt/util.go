
package jwtutil

import (
	hmac   "crypto/hmac"
	sha256 "crypto/sha256"
	// unsafe "unsafe"
	
	util "github.com/zyxgad/go-util/util"
)

func hmacSha256(data []byte, key []byte)(code []byte){
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

func equalMac(mac1, mac2 []byte)(bool){
	if mac1 == nil || mac2 == nil {
		return false
	}
	return hmac.Equal(mac1, mac2)
}


func strToBytes(str string)([]byte){
	return ([]byte)(str)
	// return *((*[]byte)(unsafe.Pointer(&str)))
}

func bytesToStr(bytes []byte)(string){
	return (string)(bytes)
	// return *((*string)(unsafe.Pointer(&bytes)))
}

type Json util.JsonType


const (
	DATE_FORMAT = "2006-01-02 15:04:05 -0700"
)
