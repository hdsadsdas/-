package tools

import "crypto/sha256"

func GetHash(data []byte)[]byte{
	hash:=sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}
