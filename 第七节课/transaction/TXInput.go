package transaction

import "bytes"

/**
* @author : 哈哈
* @email : 598421227@qq.com
* @phone : 18816473550
* @DateTime : 2022/4/18 9:39
**/

//交易输入
type Input struct {

	//确定交易输出在哪个交易中
	TXid []byte
	//交易输出索引下标
	Vout int
	//解锁脚本
	ScriptSig []byte
}

func (input *Input) IsLocked(address string) bool {

	if address == "" {
		return false
	}

	return 0 == bytes.Compare(input.ScriptSig, []byte(address))

}

func NewInput(txid []byte,vout int,scriptSig []byte) *Input {

	return &Input{txid,vout,scriptSig}

}
