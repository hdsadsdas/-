package block

import (
	"bytes"
	"strconv"
	"time"
	"公链系统开发/第一节课/tools"
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

	block.Hash = block.SetHash()

	return &block

}

/**
获取当前区块的hash值
*/
func (b *Block) SetHash() []byte {

	//当前区块的Hash值 是由时间戳 + 上一个区块的hash + 交易信息 形成一个字符串，计算这个字符串的hash值

	//将时间戳转换为[]byte    strconv.FormatInt 这个方法是将int 类型 转为 字符串
	time := []byte(strconv.FormatInt(b.TimeStamp, 10))

	//拼接 bytes.Join() 第一个参数 存放拼接的内容  第二个参数 以什么字符拼接
	join := bytes.Join([][]byte{b.PrevHash, b.Data, time}, []byte{})

	return tools.GetHash(join)

}


/**
创建创始区块
 */
func GenesisBlock(data []byte)*Block {

	genesis := NewBlock(data, nil)

	return genesis
}
