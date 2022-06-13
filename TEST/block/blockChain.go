package block

import "github.com/boltdb/bolt"

/**
* @author : 哈哈
* @email : 598421227@qq.com
* @phone : 18816473550
* @DateTime : 2022/4/11 9:46
**/

const CHAIN_DB_PATH = "./chain.db"
const BUCKET_BLOCK = "chain_blocks"
const BUCKET_STATUS = "chain_status"
const LAST_HASH = "last_hash"

type BlockChain struct {
	DB       *bolt.DB
	ListHash []byte
}

func NewBlockChain(data []byte) (*BlockChain, error) {

	db, err := bolt.Open(CHAIN_DB_PATH, 0600, nil)
	if err != nil {
		return nil, err
	}

	var lashHash []byte
	db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(BUCKET_BLOCK))

		if bucket == nil {

			createBucket, err2 := tx.CreateBucket([]byte(BUCKET_BLOCK))
			if err2 != nil {
				return err2
			}

			genesis := GenesisBlock(data)
			serialize, err2 := genesis.Serialize()
			if err2 != nil {
				return err2
			}

			err2 = createBucket.Put(genesis.Hash, serialize)
			if err2 != nil {
				return err2
			}

			createBucket2, err2 := tx.CreateBucket([]byte(BUCKET_STATUS))
			if err2 != nil {
				return err2
			}

			createBucket2.Put([]byte(LAST_HASH), genesis.Hash)

			lashHash = genesis.Hash

		} else {

			bucket2 := tx.Bucket([]byte(BUCKET_STATUS))

			lashHash = bucket2.Get([]byte(LAST_HASH))

		}

		return nil

	})

	return &BlockChain{DB: db, ListHash: lashHash}, nil

}
