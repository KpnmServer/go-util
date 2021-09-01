
package jwtutil

import (
	os "os"
	ioutil "io/ioutil"
)

type file_encoder struct{
	Encoder

	file string
}

func NewFileEncoder(encoder Encoder, file string)(fe *file_encoder){
	fe = &file_encoder{
		Encoder: encoder,
		file: file,
	}
	if file, err := os.Open(fe.file); err == nil {
		defer file.Close()
		if bytes, err := ioutil.ReadAll(file); err == nil {
			if bytes, err = decodeB64(bytes); err == nil {
				fe.Encoder.ChangeKey(bytes)
			}
		}
	}
	return fe
}

func (cdr *file_encoder)ChangeKey(key []byte){
	cdr.Encoder.ChangeKey(key)
	if file, err := os.Create(cdr.file); err == nil {
		defer file.Close()
		file.Write(encodeB64(key))
	}
}
