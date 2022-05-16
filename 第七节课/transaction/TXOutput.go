package transaction

import "bytes"

/**
* @author : 哈哈
* @email : 598421227@qq.com
* @phone : 18816473550
* @DateTime : 2022/4/18 9:29
**/

//交易输出
type Output struct {

	//描述交易输出的金额
	Value uint

	//锁定脚本
	ScriptPutKey []byte

}

//判断某个人是否能解开交易输出（判断这笔钱是否是某个人的）
func (Output *Output)IsUnlock(address string)bool  {

	if address == "" {
		return false
	}

	return 0 == bytes.Compare(Output.ScriptPutKey, []byte(address))


}

func NewOutput(value uint,scriptPutKey []byte)Output{
	return Output{value,scriptPutKey}
}
