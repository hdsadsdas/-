package transaction

import (
	"bytes"
	"encoding/gob"
	"time"
	"公链系统开发/考试/tools"
	"公链系统开发/考试/wallet"
)

//创建一个交易的结构体

type Transaction struct {
	//交易的唯一标识
	TXHash []byte
	//多个交易输出
	Output []Output
	//多个交易输入
	Input  []Input

	//时间戳
	TimeStamp int64
}
//创建一个普通的交易
func NewTransaction(from,to string,amount uint,spendableOutputs []UTXO)(*Transaction,error){


	//要买一个70的东西，先把余额中的每一钱进行累计，找到刚刚好够70的时候就可以了
	//不需要把所有的余额都用上  10 10 20  30  40

		//10（1张10元 ，2张5元，10张1元）  50   不够用
		//10  5   公用
		//c.根据找到的交易输出，构建input
		inputs :=make([]Input,0)
		for _,output:=range spendableOutputs{
			//构建交易输入就要引用交易输出，因为交易输入本质就是之前历史交易中的未消费的交易输出
			//不光需要公钥，签名需要私钥来完成，所有还需要私钥，需要把私钥存起来，
			//boltdb数据库  key->value
			//			地址 - > 对应地址的私钥
			input:=NewInput(output.Txid,output.Index,nil,nil)
			inputs = append(inputs,input)
	}

	//2.构建output（交易输出）
	//到此次为止，花费了多少钱
	var spandAmount uint
	outputs:=make([]Output,0)
	for _,out:=range spendableOutputs{
		spandAmount +=out.Value
		if spandAmount <=amount{
			pubHash, err := wallet.GetPubHash(to)
			if err!=nil{
				return nil, err
			}
			output:=NewOutputs(out.Value,pubHash)
			outputs = append(outputs,output)
		}else{
			//spandAmount超出了要转的金额
			spandAmount -= out.Value
			needAmount:= amount - spandAmount
			to_pubhash, err := wallet.GetPubHash(to)
			if err !=nil{
				return nil, err
			}
			output:=NewOutputs(needAmount,to_pubhash)
			outputs = append(outputs,output)

			from_pubhash, err := wallet.GetPubHash(from)
			if err !=nil{
				return nil, err
			}
			backChange :=NewOutputs(out.Value - needAmount,from_pubhash)
			outputs = append(outputs,backChange)
		}

	}
	//3.给TXHash赋值，并返回
	tx:=Transaction{
		Output: outputs,
		Input:  inputs,
	}
	tx.TimeStamp = time.Now().Unix()
	err, txBytes := tx.Serialize()
	if err !=nil{

		return nil,err
	}
	hash := tools.GetHash(txBytes)
	tx.TXHash = hash
	return &tx,nil
}

//创建一个coinbase交易  生成区块的奖励交易
func NewCoinBase(address string)(error,*Transaction){
	pubHash, err := wallet.GetPubHash(address)

	cb:=Transaction{
		Output: []Output{
			{
				Value: 50,
				ScriptPubkey:pubHash,
			},
		},
		Input:  nil,
	}
	cb.TimeStamp = time.Now().Unix()
	//计算cb的hash值
	//1、把交易变成[]byte
	err, txBytes := cb.Serialize()
	if err !=nil{
		return err,nil
	}
	hash := tools.GetHash(txBytes)
	cb.TXHash = hash
	return nil,&cb
}


func (tx *Transaction) Serialize()(error,[]byte){
	var result bytes.Buffer
	en := gob.NewEncoder(&result)
	err := en.Encode(tx)
	if err !=nil{
		return err,nil
	}
	return nil,result.Bytes()
}
