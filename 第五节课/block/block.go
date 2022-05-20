package block

import (
	"bytes"
	"encoding/gob"
	"time"
	pow2 "公链系统开发/第五节课/pow"
)

type Block struct {
	PrevHash []byte //上个区块的Hash值

	TimeStamp int64 //时间戳

	Data []byte //交易信息

	Hash []byte //当前区块的哈希值

	Nonce int64 //随机数
}

func NewBlock(prevHash []byte,data []byte)*Block{

	block := Block{
		PrevHash: prevHash,
		TimeStamp: time.Now().Unix(),
		Data: data,
	}

	pow:=pow2.NewPow(block.TimeStamp,block.PrevHash,block.Data)

	hash,nonce:=pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return &block

}
//序列化：把结构体block转成[]byte
func (block *Block)Serialize()([]byte,error){

	var result bytes.Buffer

	en := gob.NewEncoder(&result)

	err := en.Encode(block)

	if err !=nil{
		return nil,err
	}

	return result.Bytes(),nil

}

//把[]byte转成block
func DeSerialize(data []byte)(*Block,error){
	//var result bytes.Buffer
	reader := bytes.NewReader(data)

	de := gob.NewDecoder(reader)

	var block *Block

	err := de.Decode(&block)

	if err !=nil{
		return nil,err
	}

	return block, nil
}

func GenesisBlock(data []byte)*Block{

	return NewBlock(nil,data)

}