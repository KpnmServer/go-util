
package jwtutil

type auto_encoder struct{
	Encoder

	change_interval int64
	key_length int
}

func NewAutoEncoder(encoder Encoder, keylen int, interval int64)(ae *auto_encoder){
	ae = &auto_encoder{
		Encoder: encoder,
		change_interval: interval,
		key_length: keylen,
	}
	if ae.getkey() == nil {
		ae.RandKey(ae.key_length)
	}
	return ae
}

func (cdr *auto_encoder)checkKey(){
	if cdr.lastChangeTime() + cdr.change_interval >= timeNowUnix() {
		return
	}
	cdr.RandKey(cdr.key_length)
}

func (cdr *auto_encoder)RandKey(leng int){
	cdr.ChangeKey(makeCRandBytes(leng))
}

func (cdr *auto_encoder)Encode(json Json)(token string){
	cdr.checkKey()
	return cdr.Encoder.Encode(json)
}

func (cdr *auto_encoder)Decode(token string)(json Json, err error){
	cdr.checkKey()
	return cdr.Encoder.Decode(token)
}
