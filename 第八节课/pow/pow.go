package pow

import (
	"bytes"
	"公链系统开发/第八节课/tools"
	"公链系统开发/第八节课/transaction"

	"math/big"
	"strconv"
)
const BITS = 10 //目标值前面有多少个0
//0001111  前面的0越多，越麻烦
type ProofOfWork struct {
	//Block *block.Block //给谁工作的区块

	//TimeStamp int64
	//PrevHash []byte
	//Data []byte

	Block BlockInterface
	Target  *big.Int  //要判断的系统给定的hash
}

type BlockInterface interface {
	GetTimeStamp() int64
	GetPrevHash() []byte
	GetData() []transaction.Transaction
}

/**
创建一个pow的实例，并且返回
 */
func NewPow(block BlockInterface)*ProofOfWork{

//var a int = 1
	target:= big.NewInt(1) // 000000...000001
	//00000000000..00
	//前面有20个0，1移动多少位可以变成这样
	target=target.Lsh(target,255-BITS)
	pow:=ProofOfWork{
		Block: block,
		Target: target,
	}
	return &pow
}

/**
寻找满足条件的随机数的
 */
func (pow *ProofOfWork)Run()(int64,[]byte){
	var nonce int64
	nonce = 0
	//先得到区块的hash
	block:=pow.Block
	timeByte:=[]byte(strconv.FormatInt(block.GetTimeStamp(),10))

	num:=big.NewInt(0)
	//统一类型 []byte -> 转成 大整数类型
	for{
		nonceByte:=[]byte(strconv.FormatInt(nonce,10))
		//把block.GetData()的类型转为[]byte
		txsBytes := []byte{}
		for _,value :=range block.GetData(){
			_, txBytes := value.Serialize()
			txsBytes = append(txsBytes,txBytes...)
		}
		hashByte:=bytes.Join([][]byte{txsBytes,block.GetPrevHash(),nonceByte,timeByte},[]byte{})
		hash:=tools.GetHash(hashByte)
		//fmt.Println("正在寻找nonce，当前的nonce为",nonce)
		num = num.SetBytes(hash)
		if(num.Cmp(pow.Target)==-1){
			return nonce,hash
		}
		nonce++
	}
	return 0, nil
	/*
	if(a< target){

	}
	 */

}