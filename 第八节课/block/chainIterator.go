package block

import (
	"bytes"
	"errors"
	"github.com/boltdb/bolt"
)

type ChainIterator struct {
	//数据库，区块在数据中
	DB *bolt.DB
	//标志位：标志当前迭代到哪一个区块  Current现代
	CurrentHash []byte
}

//功能：寻找区块信息,并且把标志位向前移动一位
func (iterator *ChainIterator) Next()(*Block,error){
	var block *Block
	var err error
	//读操作
	err=iterator.DB.View(func(tx *bolt.Tx) error {
		//找桶1
		bk := tx.Bucket([]byte(BUCKET_BLOCK))
		if bk == nil{
			return errors.New("没用桶1")
		}
		//根据标志位来找对应的block
		blockBytes := bk.Get(iterator.CurrentHash)
		//反序列化
		block, err = DeSerialize(blockBytes)
		if err !=nil{
			return err
		}
		iterator.CurrentHash = block.PrevHash
		return nil
	})
	return block,err
}

//功能：判断是否还有下一个区块
func (iterator *ChainIterator)HasNext()bool{
	//判断当前区块的上一个hash值是不是创世区块的上一个hash值

	//为什么不判断当前区块的hash值是不是创世区块的hash值
	//因为创世区块的prevhash是特殊的，nil
	int := bytes.Compare(iterator.CurrentHash, nil)
	//if int == 0{
	//	return false
	//}
	return int != 0 //还有上一个区块 返回true
}