package tools

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"os"
)

//判断文件是否存在
//返回true  代表文件存在
//返回false  代表文件不存在
func FileExist(path string)bool{
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

//将私钥序列化
func Serialize(data *ecdsa.PrivateKey)([]byte,error){

	var result bytes.Buffer
	en := gob.NewEncoder(&result)
	//将p256注册进去
	gob.Register(elliptic.P256())
	err := en.Encode(data)
	if err !=nil{
		return nil,err
	}
	return result.Bytes(),nil

}