package block

import (
	"bytes"
	"encoding/gob"
	"time"
	"公链系统开发/第三节课/pow"
)

/**
* @author : 哈哈
* @email : 598421227@qq.com
* @phone : 18816473550
* @DateTime : 2022/2/21 9:13
**/

type Block struct {

	PrevHash []byte //上个区块的Hash值

	TimeStamp int64 //时间戳

	Data []byte //交易信息

	Hash []byte //当前区块的哈希值

	Nonce int64 //随机数

}

/**
初始化区块(创建区块)
*/
func NewBlock(data []byte, prevHash []byte) *Block {

	block := Block{
		PrevHash:  prevHash,
		TimeStamp: time.Now().Unix(),
		Data:      data,
	}

	Pow := pow.NewPow(block.TimeStamp,block.PrevHash,block.Data)

	hash, nonce := Pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return &block

}

/**
获取当前区块的hash值
*/
//func (b *Block) SetHash() []byte {
//
//	//区块的hash ： 时间戳 + 上一个区块的hash值 + 交易信息 + 随机数
//
//	//将时间戳转换为[]byte    strconv.FormatInt 这个方法是将int 类型 转为 字符串
//	time := []byte(strconv.FormatInt(b.TimeStamp, 10))
//
//	//随机数
//	nonce := []byte(strconv.FormatInt(b.Nonce, 10))
//
//	//拼接 bytes.Join() 第一个参数 存放拼接的内容  第二个参数 以什么字符拼接
//		join := bytes.Join([][]byte{b.PrevHash, b.Data, time,nonce}, []byte{})
//
//	return tools.GetHash(join)
//
//}


/**
创建创始区块
 */
func GenesisBlock(data []byte)*Block {

	genesis := NewBlock(data, nil)

	return genesis
}


/**
序列化
 */
func (block *Block)Serialize()([]byte ,error)  {

	var result bytes.Buffer

	//gob
	en := gob.NewEncoder(&result)

	err := en.Encode(&block)

	if err != nil {

		return nil , err

	}

	return result.Bytes(),nil

}

/**
反序列化
 */
func DeSerialize(data []byte) (*Block,error) {

	reader := bytes.NewReader(data)

	de := gob.NewDecoder(reader)

	var block *Block

	err := de.Decode(&block)

	if err != nil {
		return nil, err
	}

	return block,nil

}