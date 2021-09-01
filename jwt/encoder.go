
package jwtutil

import (
	strings "strings"
	time    "time"

	ujson "github.com/KpnmServer/go-util/json"
)

type Json = map[string]interface{}

var DEFAULT_OUTTIME int64 = 60 * 60 * 24 * 7

func SetOutdate(json Json, outtime time.Duration)(Json){
	json["iat"] = timeNow().Add(outtime).Unix()
	json["exp"] = timeNowUnix()
	return json
}

type Encoder interface{
	ChangeKey(key []byte)
	Encode(json Json)(token string)
	Decode(token string)(json Json, err error)
	getkey() []byte
	lastChangeTime() int64
}

type default_encoder struct{
	lastkey []byte
	key []byte
	last_change_time int64
	outtime int64
}

func NewEncoder(key []byte, outtime_ ...int64)(*default_encoder){
	outtime := DEFAULT_OUTTIME
	if len(outtime_) > 0 {
		outtime = outtime_[0]
	}
	return &default_encoder{
		lastkey: nil,
		key: key,
		last_change_time: timeNowUnix(),
		outtime: outtime,
	}
}

func (cdr *default_encoder)ChangeKey(key []byte){
	if key == nil {
		return
	}
	cdr.lastkey = cdr.key
	cdr.key = make([]byte, len(key))
	copy(cdr.key, key)
	cdr.last_change_time = timeNowUnix()
}

func (cdr *default_encoder)Encode(json Json)(token string){
	if json == nil || cdr == nil || cdr.key == nil {
		return ""
	}
	head := b64RmTail(encodeB64Url(([]byte)(`{"alg":"HS256","typ":"JWT"}`)))
	fdata := head + "." + b64RmTail(encodeB64Url(ujson.EncodeJson(json)))
	code := b64RmTail(encodeB64Url(hmacSha256(([]byte)(fdata), cdr.key)))
	token += fdata + "." + code
	return token
}

func (cdr *default_encoder)Decode(token string)(json Json, err error){
	if cdr.lastkey != nil && cdr.last_change_time + 60 * 60 * 24 * 7 < timeNowUnix() {
		cdr.lastkey = nil
	}

	if cdr == nil || cdr.key == nil {
		return nil, NULL_POINT_ERR
	}
	arr := strings.Split(token, ".")
	if len(arr) != 3 {
		return nil, SPLIT_ERROR
	}
	mac1 := hmacSha256(([]byte)(token)[0:len(arr[0]) + 1 + len(arr[1])], cdr.key)
	mac2, _ := decodeB64Url(arr[2])
	if !equalMac(mac1, mac2) {
		if cdr.lastkey == nil {
			return nil, MAC_NOT_SAME_ERROR
		}else if cdr.last_change_time + cdr.outtime < timeNowUnix(){
			cdr.lastkey = nil
			return nil, MAC_NOT_SAME_ERROR
		}else if !equalMac(hmacSha256(([]byte)(token)[0:len(arr[0]) + 1 + len(arr[1])], cdr.lastkey), mac2) {
			return nil, MAC_NOT_SAME_ERROR
		}
	}
	var data []byte
	data, err = decodeB64Url(arr[1])
	if err != nil {
		return nil, err
	}
	err = ujson.DecodeJson(data, &json)
	if err != nil {
		return nil, err
	}
	if outdate0, ok := json["iat"]; ok && outdate0 != nil && (int64)(outdate0.(float64)) <= timeNowUnix() {
		return json, TOKEN_OUT_DATE_ERROR
	}
	return json, nil
}

func (cdr *default_encoder)getkey()([]byte){
	return cdr.key
}

func (cdr *default_encoder)lastChangeTime()(int64){
	return cdr.last_change_time
}
