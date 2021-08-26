
package jwtutil

type auto_encoder struct{
	default_encoder

	change_interval int64
	key_length int
}

func NewAutoEncoder(interval int64, keylen int, outtime_ ...int64)(*auto_encoder){
	return &auto_encoder{
		default_encoder: *NewEncoder(makeCRandBytes(keylen), outtime_...),
		change_interval: interval,
		key_length: keylen,
	}
}

func (cdr *auto_encoder)checkKey(){
	if cdr.last_change_time + cdr.change_interval >= timeNowUnix() {
		return
	}
	cdr.RandKey(cdr.key_length)
}

func (cdr *auto_encoder)ChangeKey(key []byte){
	cdr.default_encoder.ChangeKey(key)
}

func (cdr *auto_encoder)RandKey(leng int){
	cdr.ChangeKey(makeCRandBytes(leng))
}

func (cdr *auto_encoder)Encode(json Json)(token string){
	cdr.checkKey()
	return cdr.default_encoder.Encode(json)
}

func (cdr *auto_encoder)Decode(token string)(json Json, isout bool, err error){
	cdr.checkKey()
	return cdr.default_encoder.Decode(token)
}
