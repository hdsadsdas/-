package tools

import "crypto/sha256"

/**
* @author : 哈哈
* @email : 598421227@qq.com
* @phone : 18816473550
* @DateTime : 2022/4/11 9:39
**/

func GetHash(data []byte)[]byte{
	hash:=sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}
