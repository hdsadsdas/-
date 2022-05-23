package transaction

import "bytes"

//创建一个交易输入的结构体
type Input struct {
	//交易哈希
	Txid  []byte
	//交易输出索引位置
	Vout  int
	//解锁脚本
	ScriptSig []byte
}

func NewInput(txid []byte,vout int,scriptSig []byte)Input{
	return Input{txid,vout,scriptSig}
}

//判断input是某个人的消费
func(input *Input)IsLocked(name string)bool{
	if name ==""{
		return false
	}
	return bytes.Compare(input.ScriptSig,[]byte(name)) == 0
}