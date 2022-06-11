package tools

import (
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
)


//进行sha256
func GetHash(data []byte)[]byte{
	hash:=sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

//进行ripemd160计算
func Ripemd160(data []byte)[]byte  {

	hash := ripemd160.New()
	hash.Write(data)
	return hash.Sum(nil)
}
