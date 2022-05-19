package block

import (
	"bytes"
	"errors"
	"github.com/boltdb/bolt"
)

/**
* @author : 哈哈
* @email : 598421227@qq.com
* @phone : 18816473550
* @DateTime : 2022/3/14 9:40
**/

type ChainIterator struct {
	DB *bolt.DB
	//指定当前位置 也是值当前区块的hash值
	currentHash []byte
}

func (iterator *ChainIterator) Next() (*Block,error) {

	var block *Block
	var err error

	//View  同一时间允许多个人进行查看数据
	err = iterator.DB.View(func(tx *bolt.Tx) error {

		bk := tx.Bucket([]byte(BUCKET_BLOCK))
		if bk == nil {
			return errors.New("没有桶")
		}

		//获取得到最后一个区块的信息
		blockBytes := bk.Get(iterator.currentHash)

		//反序列化
		block , err = DeSerialize(blockBytes)
		if err != nil {
			return err
		}

		iterator.currentHash = block.PrevHash

		return nil
	})

	if err != nil {
		return nil,err
	}


	return block,nil

}


func (iterator *ChainIterator)HasNext()bool{

	int := bytes.Compare(iterator.currentHash, nil)

	return int != 0

}
