package block

import (
	"bytes"
	"公链系统开发/第八节课/transaction"

	"errors"
	"fmt"
	"github.com/boltdb/bolt"
)

const CHAIN_DB_PATH = "./chain.db"
const BUCKET_BLOCK = "chain_block"
const BUCKET_STATUS = "chain_status"
type BlockChain struct {
	//Blocks []*Block
	DB  *bolt.DB
	LastHash []byte
}

//创建带有创世区块的区块链
func CreatChain (address string)(*BlockChain,error){
	//打开数据库，返回一个打开好的数据库对象，和err
	db, err := bolt.Open(CHAIN_DB_PATH, 0600, nil)
	if err !=nil{
		return nil,err
	}
	var lastHash []byte
	//想向数据库db中添加数据,更新包括增加、修改、删除
	err = db.Update(func(tx *bolt.Tx) error {
		//1.有桶
		// tx.Bucket 在寻找桶
		bk := tx.Bucket([]byte(BUCKET_BLOCK))
		if bk == nil{
			bk, err := tx.CreateBucket([]byte(BUCKET_BLOCK))
			if err !=nil{
				return err
			}
			//向bk中添加数据
			//创建一个coinbase交易
			err, coinbase := transaction.NewCoinBase(address)
			genesis:=CreatGenesis(*coinbase)
			serialize, err := genesis.Serialize()
			if err !=nil{
				return err
			}
			bk.Put(genesis.Hash,serialize)
			//创建桶2，用来存储最新区块的hash值
			bk2, err := tx.CreateBucket([]byte(BUCKET_STATUS))
			if err !=nil{
				return err
			}
			bk2.Put([]byte("LAST_HASH"),genesis.Hash)
			lastHash = genesis.Hash
		}else{
			bk2:=tx.Bucket([]byte(BUCKET_STATUS))
			lastHash= bk2.Get([]byte("LAST_HASH"))
		}
		return nil
	})

	bc:=BlockChain{
		DB: db,
		LastHash:lastHash,
	}




	return &bc,err
}

func (bc *BlockChain) AddBlock(tx []transaction.Transaction)error{

	newBlock:=CreatBlock(tx,bc.LastHash)
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(BUCKET_BLOCK))
		if bk == nil{
			return errors.New("没有桶1")
		}
		serialize, err := newBlock.Serialize()
		if err !=nil{
			return err
		}
		bk.Put(newBlock.Hash,serialize)

		bk2 := tx.Bucket([]byte(BUCKET_STATUS))
		if bk2 == nil{
			return errors.New("没有桶2")
		}
		bk2.Put([]byte("LAST_HASH"),newBlock.Hash)
		bc.LastHash = newBlock.Hash
		return nil
	})
	return err
}

//创建迭代器
func (bc *BlockChain)Iterator()*ChainIterator{
	//实例化迭代器
	iterator:=ChainIterator{
		DB:        bc.DB  ,
		CurrentHash: bc.LastHash,
	}
	return &iterator
}

//获取区块链中所有的区块
func(bc *BlockChain)GetAllBlock()([]*Block,error){
	iterator := bc.Iterator()
	blocks :=[]*Block{}
	for{
		if iterator.HasNext(){
			block, err := iterator.Next()
			if err !=nil{
				fmt.Println(err.Error())
				return nil,err
			}
			blocks= append(blocks,block)
		}else{
			break
		}
	}
	return blocks,nil
}

//定义一个方法，用来找某个地址的所有收入 zhang
func (bc *BlockChain)FindAllOutput(address string)[]transaction.UTXO{
	//先找所有的区块，再获取每一个区块中的所有的交易，再找是zhang的收入的
	blocks, err := bc.GetAllBlock()
	//map结构：map[key]value
	//key:收入（input）所在的交易hash  tx2
	//value:[]int 表示output的位置的下标  1，2
	//用来存储某个人的所有的收入
	allOutPuts:=make([]transaction.UTXO,0)
	if err !=nil{
		fmt.Println("寻找失败",err.Error())
		return nil
	}
	//获取每一个区块
	for _,block:=range blocks{
		//获取每一个区块中的每一个交易
		for _,tx :=range block.Txs{
			//找每一个交易中的所有的交易输出
			for outIndex,output:=range tx.Output{
				if output.IsUnlock(address) {
					utxo:=transaction.NewUTXO(tx.TXHash,outIndex,output)
					allOutPuts= append(allOutPuts,utxo)
				}
			}
		}
	}
	return allOutPuts
}

//寻找某个人的所有的消费（input）
func(bc *BlockChain)FindAllInput(name string)[]transaction.Input{
	//先找所有的区块，再获取每一个区块中的所有的交易，再找是zhang的消费的
	allInputs:=make([]transaction.Input,0)
	blocks, err := bc.GetAllBlock()
	if err !=nil{
		fmt.Println("寻找失败",err.Error())
		return nil
	}
	//找每一个区块
	for _,block:=range blocks{
		//每一个区块中的每一个交易
		for _,tx:= range block.Txs{
			//找每一个交易中的input
			for _,input :=range tx.Input{
				 if input.IsLocked(name){
				 	allInputs = append(allInputs,input)
				 }
			}
		}
	}
	return allInputs
}
// map  key  value
//      1003  1

/**
接收某个人的所有的收入和消费，并从所有的收入中去掉消费，剩下的就是可用的交易输出
 */
func(bc *BlockChain)FindSpendOutputs(alloutputs []transaction.UTXO,allintpts []transaction.Input,amount uint)([]transaction.UTXO,uint){
	//获取每一个收入
	for _,input:=range allintpts{
		for index,utxo:=range alloutputs{
			if bytes.Compare(utxo.Txid,input.Txid) ==0 || utxo.Index == input.Vout{
				//两个txid和索引下标都相等，表示这个utxo被用掉了
				//要从utxo中去掉这一笔收入，到最后剩下的utxo就是可用的余额
				alloutputs = append(alloutputs[0:index],alloutputs[index+1:]...)
				break
			}
		}
	}
	var totalAmount uint
	outputs:=make([]transaction.UTXO,0)
	//遍历找到所有的余额
	for _,output:=range alloutputs{
		totalAmount+=output.Value
		outputs= append(outputs,output)
		if totalAmount >=amount{
			//钱够用了
			break
		}
	}


	return outputs,totalAmount

}

func (bc *BlockChain)NewTransaction(from,to string,amount uint)(*transaction.Transaction,error){
	//1.构建input（交易输入）
	//a.在所有的交易中，去寻找可以使用的交易输出
	//所有可用的交易输出（余额） = 所有的收入 - 所有的消费

	allInputs := bc.FindAllInput(from)
	allOutputs := bc.FindAllOutput(from)
	//获取了from这个人的所有可用的交易输出
	//因为只需要找到满足本次交易的钱就可以了，不需要找所有的钱
	spendableOutputs,totalAmount := bc.FindSpendOutputs(allOutputs, allInputs,amount)
	//b.从所有的可用的交易输出中，取出一部分，判断是否够用
	if totalAmount  < amount{
		return nil,errors.New("余额不足")
	}
	//金额够用，可以发起交易，调用transaction包，构建交易
	tran, err := transaction.NewTransaction(from, to, amount, spendableOutputs)
	if err !=nil{
		return nil, err
	}
	return tran,nil
}

//为区块链添加新的功能：创建coinbase交易
func (bc *BlockChain)NewCoinBase(from string)(*transaction.Transaction,error){
 //核验地址是否正确
	if from == ""||len(from) == 0{
		return  nil,errors.New("地址错误")
	}
	err, cb := transaction.NewCoinBase(from)
	if err !=nil{
		return nil,err
	}
	return  cb,nil
}

//获取某个人的余额
func (bc *BlockChain) GetBalance(s string) uint {
	//获取所有的收入
	allOutputs := bc.FindAllOutput(s)
	//获取所有的消费
	allInputs := bc.FindAllInput(s)


	for _,input:=range  allInputs{
		for index,output:=range allOutputs{
			if bytes.Compare(input.Txid,output.Txid) ==0 &&input.Vout == output.Index{

				if index >= len(allOutputs){
					allOutputs = append(allOutputs[:index])
				}else{
					allOutputs = append(allOutputs[:index],allOutputs[index+1:]...)
				}

				break
			}
		}
	}


	var balance uint
	for _,output:=range allOutputs{
		balance += output.Value
	}

	return balance
}