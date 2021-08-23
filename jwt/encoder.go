
package jwtutil

import (
	strings "strings"
	time    "time"

	util "github.com/zyxgad/go-util/util"
)


func SetOutdate(json Json, outtime time.Duration)(Json){
	json["iat"] = util.GetTimeAfter(outtime).Unix()
	json["exp"] = util.GetTimeNow().Unix()
	return json
}

type Encoder interface{
	ChangeKey(key []byte)
	Encode(json Json)(token string)
	Decode(token string)(json Json, isout bool, err error)
}

type default_encoder struct{
	Encoder

	lastkey []byte
	key []byte
	_last_change_time int64
}

func NewEncoder(key []byte)(*default_encoder){
	return &default_encoder{
		lastkey: nil,
		key: key,
		_last_change_time: util.GetTimeNow().Unix(),
	}
}

func (cdr *default_encoder)ChangeKey(key []byte){
	if key == nil {
		return
	}
	cdr.lastkey = cdr.key
	cdr.key = make([]byte, len(key))
	copy(cdr.key, key)
	cdr._last_change_time = util.GetTimeNow().Unix()
}

func (cdr *default_encoder)Encode(json Json)(token string){
	if json == nil || cdr == nil || cdr.key == nil {
		return ""
	}
	head := util.Base64RmTail(util.EncodeBase64Url(([]byte)(`{"alg":"HS256","typ":"JWT"}`)))
	fdata := head + "." + util.Base64RmTail(util.EncodeBase64Url(strToBytes(util.EncodeJson((util.JsonType)(json)))))
	code := util.Base64RmTail(util.EncodeBase64Url(hmacSha256(strToBytes(fdata), cdr.key)))
	token += fdata + "." + code
	return token
}

func (cdr *default_encoder)Decode(token string)(json Json, isout bool, err error){
	isout = false

	if cdr.lastkey != nil && cdr._last_change_time + 60 * 60 * 24 * 7 < util.GetTimeNow().Unix() {
		cdr.lastkey = nil
	}

	if cdr == nil || cdr.key == nil {
		return nil, false, util.NULL_POINT_ERR
	}
	arr := strings.Split(token, ".")
	if len(arr) != 3 {
		return nil, false, SPLIT_ERROR
	}
	mac1 := hmacSha256(strToBytes(token)[0:len(arr[0]) + 1 + len(arr[1])], cdr.key)
	mac2, _ := util.DecodeBase64Url(arr[2])
	if !equalMac(mac1, mac2) {
		if cdr.lastkey == nil || !equalMac(hmacSha256(strToBytes(token)[0:len(arr[0]) + 1 + len(arr[1])], cdr.lastkey), mac2) {
			return nil, false, MAC_NOT_SAME_ERROR
		}
		isout = true
	}
	var data []byte
	data, err = util.DecodeBase64Url(arr[1])
	if err != nil {
		return nil, false, err
	}
	json = (Json)(util.DecodeJson(bytesToStr(data)))
	if outdate0, ok := json["iat"]; ok && outdate0 != nil && util.JsonToInt64(outdate0) <= util.GetTimeNow().Unix() {
		return json, false, TOKEN_OUT_DATE_ERROR
	}
	return json, isout, nil
}

type auto_encoder struct{
	default_encoder

	_change_interval int64
	_key_length int
}

func NewAutoEncoder(interval int64, keylen int)(*auto_encoder){
	return &auto_encoder{
		default_encoder: default_encoder{
			lastkey: nil,
			key: util.MakeCRandBytes(keylen),
			_last_change_time: util.GetTimeNow().Unix(),
		},
		_change_interval: interval,
		_key_length: keylen,
	}
}

func (cdr *auto_encoder)checkKey(){
	if cdr._last_change_time + cdr._change_interval >= util.GetTimeNow().Unix() {
		return
	}
	cdr.RandKey(cdr._key_length)
}

func (cdr *auto_encoder)ChangeKey(key []byte){
	cdr.default_encoder.ChangeKey(key)
}

func (cdr *auto_encoder)RandKey(leng int){
	cdr.ChangeKey(util.MakeCRandBytes(leng))
}

func (cdr *auto_encoder)Encode(json Json)(token string){
	cdr.checkKey()
	return cdr.default_encoder.Encode(json)
}

func (cdr *auto_encoder)Decode(token string)(json Json, isout bool, err error){
	cdr.checkKey()
	return cdr.default_encoder.Decode(token)
}

