package block

import (
	"bytes"
	"公链系统开发/第八节课/pow"
	"公链系统开发/第八节课/transaction"

	"encoding/gob"
	"time"
)

type Block struct {
	TimeStamp int64  //创建区块时的时间戳
	PrevHash  []byte
	//Data  []byte
	Txs    []transaction.Transaction
	Nonce int64 //随机数
	Hash []byte
}

func(block *Block)GetTimeStamp()int64{
	return block.TimeStamp
}
func(block *Block)GetPrevHash()[]byte{
	return block.PrevHash
}
func(block *Block)GetData()[]transaction.Transaction{
	return block.Txs
}

/*
创建区块
 */
func CreatBlock(txs []transaction.Transaction,prevHash []byte)*Block{
	block:=Block{
		TimeStamp: time.Now().Unix(),
		PrevHash:  prevHash,
		Txs:txs,

	}
	//实现了blockinterface的一个实例
	pow:=pow.NewPow(&block)
	nonce,hash:=pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return &block
}

//func (block *Block)SetHash()[]byte{
//
//	time:=[]byte(strconv.FormatInt(block.TimeStamp,10))
//	nonce:=[]byte(strconv.FormatInt(block.Nonce,10))
//	hashByte:=bytes.Join([][]byte{time,block.Data,block.PrevHash,nonce},[]byte{})
//	return tools.GetHash(hashByte)
//}

//序列化：把block结构体类型转成[]byte
func(block *Block)Serialize()([]byte,error){
	var result bytes.Buffer
	en := gob.NewEncoder(&result)
	err := en.Encode(block)
	if err !=nil{
		return nil,err
	}
	return result.Bytes(), nil
}
//反序列化：把[]byte转成block结构体
func DeSerialize(data []byte)(*Block,error){
	//var result bytes.Buffer
	reader := bytes.NewReader(data)
	de := gob.NewDecoder(reader)
	var block *Block
	err := de.Decode(&block)
	if err !=nil{
		return nil, err
	}
	return block, nil
}

func CreatGenesis(tx transaction.Transaction)*Block{
	return CreatBlock([]transaction.Transaction{tx},nil)
}