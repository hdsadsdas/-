package tools

import (
	"bytes"
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

func Serialize(data interface{})([]byte,error)  {

	var result bytes.Buffer
	en := gob.NewEncoder(&result)

	err := en.Encode(data)
	if err != nil {
		return nil,err
	}

	return result.Bytes(),nil
}