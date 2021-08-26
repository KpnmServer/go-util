
package jwtutil

import (
	time   "time"
	hmac   "crypto/hmac"
	sha256 "crypto/sha256"
	crand "crypto/rand"
	base64 "encoding/base64"
)


func timeNow()(time.Time){
	return time.Now().UTC()
}

func timeNowUnix()(int64){
	return time.Now().UTC().Unix()
}

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

func encodeB64Url(bytes []byte)(b64 string){
	return base64.URLEncoding.EncodeToString(bytes)
}

func decodeB64Url(b64 string)(bytes []byte, err error){
	return base64.URLEncoding.DecodeString(Base64AddTail(b64))
}

func b64RmTail(b64_ string)(b64 string){
	n := len(b64_)
	for b64_[n - 1] == '=' { n -= 1 }
	return b64_[:n]
}

func b64AddTail(b64 string)(b64_ string){
	b64_ = b64
	for len(b64_) % 4 != 0 { b64_ += "=" }
	return b64_
}

func makeCRandBytes(leng int)(bytes []byte){
	if leng < 0 {
		return nil
	}
	if leng == 0 {
		return make([]byte, 0)
	}
	bytes = make([]byte, leng)
	i := 0
	for i < leng {
		n, err := crand.Read(bytes[i:])
		if err != nil || n == 0 {
			return []byte{}
		}
		i += n
	}
	return bytes
}

const (
	DATE_FORMAT = "2006-01-02 15:04:05 -0700"
)
