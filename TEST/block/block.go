package block

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
	"公链系统开发/TEST/pow"
)

/**
* @author : 哈哈
* @email : 598421227@qq.com
* @phone : 18816473550
* @DateTime : 2022/4/11 8:49
**/

type Block struct {

	PrevHash []byte //上个区块的Hash值

	TimeStamp int64 //时间戳

	Data []byte //交易信息

	Hash []byte //当前区块的哈希值

	Nonce int64 //随机数

}

//创建区块
func NewBlock(prevHash []byte,data []byte)*Block{

	block := &Block{
		PrevHash: prevHash,
		TimeStamp: time.Now().Unix(),
		Data: data,
	}

	p := pow.NewPOW(block.PrevHash,block.TimeStamp,block.Data)

	hash, nonce := p.Run()

	block.Hash = hash
	block.Nonce = nonce

	fmt.Println("AAAA",hash)

	return block

}

//创建创世区块
func GenesisBlock(data []byte)*Block{
	return NewBlock(nil,data)
}

//序列化
func (block *Block)Serialize()([]byte,error)  {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil {
		return nil, err
	}

	return result.Bytes() ,nil

}