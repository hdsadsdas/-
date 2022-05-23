package tools

import "os"

//判断文件是否存在
//返回true  代表文件存在
//返回false  代表文件不存在
func FileExist(path string)bool{
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}