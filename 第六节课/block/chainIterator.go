package block

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
)

type ChainIterator struct {
	DB *bolt.DB
	//标志位，表示当前迭代器所迭代到的位置
	currentHash []byte
}

//使用迭代器找上一个区块，因为是从后向前找区块
func (iterator *ChainIterator) Next() (*Block, error) {
	var block *Block
	var err error
	//同一时间允许多个人进行查看数据
	err = iterator.DB.View(func(tx *bolt.Tx) error {

		bk := tx.Bucket([]byte(BUCKET_BLOCK))

		if bk == nil {
			return errors.New("没有桶")
		}

		//最后一个区块的信息 []byte类型
		blockBytes := bk.Get(iterator.currentHash)

		//想获取最后一个区块的prevhash
		block, err = DeSerialize(blockBytes)

		iterator.currentHash = block.PrevHash
		return nil
	})
	return block, err
}

//判断是否还有下一个区块
func (iterator *ChainIterator) HasNext() bool {

	int := bytes.Compare(iterator.currentHash, nil)

	fmt.Println("AAAA",int)

	fmt.Println(int != 0)

	return int != 0
}
