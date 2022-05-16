package transaction

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"公链系统开发/第七节课/tools"
)

/**
* @author : 哈哈
* @email : 598421227@qq.com
* @phone : 18816473550
* @DateTime : 2022/4/18 9:41
**/

//整个交易的结构体
type Transaction struct {

	//交易的唯一标识
	TXid []byte

	//多个交易输入
	Input []Input

	//多个交易输出
	Output []Output

}
//序列化
func (txs *Transaction)Serialize()([]byte,error)  {

	var result bytes.Buffer

	en := gob.NewEncoder(&result)

	err := en.Encode(txs)
	if err!=nil {
		return nil,err
	}

	return result.Bytes(),nil
}

//创建创世区块的交易
func NewCoinbase(address string) (*Transaction,error) {

	cb := Transaction{
		Output: []Output{
			{
				Value: 50,
				ScriptPutKey:[]byte(address),
			},
		},
		Input: nil,
	}

	txBytes, err := cb.Serialize()
	if err != nil {
		return nil,err
	}
	hash := tools.GetHash(txBytes)

	cb.TXid = hash

	return &cb,nil
}

//创建一个普通交易
func NewTransaction(from string,to string,amount uint,utxos []UTXO,totalAmount uint)(*Transaction,error){



	inputs := make([]Input,0)

	for _,utxo:=range utxos{

		input := NewInput(utxo.Txid, utxo.Index, []byte(from))
        inputs = append(inputs,*input)

	}

	outputs := make([]Output,0)
	//2.构建output交易输出
	if totalAmount == amount {
		for _,out := range utxos{

			output := NewOutput(out.Value,[]byte(to))
			outputs = append(outputs,output)
		}
	}else {

		lastmoney := totalAmount-amount

		for k,out := range utxos{

			if k == len(utxos)-1 {

				output :=NewOutput(out.Value-lastmoney,[]byte(to))
				backChange := NewOutput(lastmoney,[]byte(from))

				outputs = append(outputs,output,backChange)
                break
			}

			output := NewOutput(out.Value,[]byte(to))
			outputs = append(outputs,output)
		}

	}



	//3.给Txid赋值

	tx := Transaction{
		Input: inputs,
		Output: outputs,
	}

	serialize, err := tx.Serialize()
	if err != nil {
		fmt.Println(err.Error())
		return nil,err
	}

	tx.TXid = tools.GetHash(serialize)

	return &tx,nil

}
























