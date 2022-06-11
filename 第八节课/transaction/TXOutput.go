package transaction

import (
	"bytes"
	"公链系统开发/第八节课/wallet"
)

//创建一个交易输出结构体

type Output struct {
	//币的金额
	Value uint
	//锁定脚本
	ScriptPubkey []byte
}


func NewOutputs(value uint,scriptpubkey []byte)Output{
	return Output{value,scriptpubkey}
}

//判断某个人是否能解开交易输出（判断这笔钱是不是某个人的）
func(output *Output)IsUnlock(name string)bool{

	pubHash, err := wallet.GetPubHash(name)
	if err != nil {
		return false
	}

	if name ==""{
		return false
	}

	return bytes.Compare(output.ScriptPubkey,pubHash) == 0
}