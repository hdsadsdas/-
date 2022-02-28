package tools

import (
	"crypto/sha256"
)

/**
* @author : 哈哈
* @email : 598421227@qq.com
* @phone : 18816473550
* @DateTime : 2022/2/21 9:25
**/

//用sha256 计算 Hash 值
func GetHash(data []byte)[]byte{

	hash := sha256.New() //创建一个 sha256 的对象

	hash.Write(data)

	return hash.Sum(nil)

}