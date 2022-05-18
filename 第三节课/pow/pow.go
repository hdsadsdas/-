package pow

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"
	"公链系统开发/第三节课/tools"
)

/**
* @author : 哈哈
* @email : 598421227@qq.com
* @phone : 18816473550
* @DateTime : 2022/2/28 9:50
**/

//在256维hash里面前面至少要有16个零
const BITS = 16

type ProofOfWork struct {

	//Block *block.Block   //当前要验证的区块

	TimeStamp int64 //时间戳
	PrevHash []byte //上个区块的Hash值
	Data []byte //交易信息

	Target *big.Int // 二进制hash值

}


func NewPow(TimeStamp int64 ,PrevHash []byte,Data []byte )*ProofOfWork  {

	//1.创建一个初始值为1的target

	target := big.NewInt(1)

	//2.左移256 - targetBit - 1 得到要移动的位数
	target = target.Lsh(target,256 - BITS - 1)

    return &ProofOfWork{TimeStamp,PrevHash,Data,target}

}


//用来寻找随机数
func (pow *ProofOfWork)Run()([]byte,int64){

	//初始化随机数
   var nonce int64 = 0

   //block := pow.Block
   //
   //block.Nonce = nonce

	//将时间戳转换为[]byte    strconv.FormatInt 这个方法是将int 类型 转为 字符串
	time := []byte(strconv.FormatInt(pow.TimeStamp, 10))

	//初始化大整形
	num := big.NewInt(1)


	for  {

		//随机数
		nonceByte := []byte(strconv.FormatInt(nonce, 10))

		//拼接切片
		hashByets := bytes.Join([][]byte{pow.Data, pow.PrevHash, time, nonceByte}, []byte{})

		//得到hash值
		hash := tools.GetHash(hashByets)

		//将hash值转换为二进制hash值
		hashBigInt := num.SetBytes(hash)

		//fmt.Println("当前的nonce为",nonce)
		fmt.Printf("\r%X",hash)

		//大整形比较
		if (hashBigInt.Cmp( pow.Target) == -1) {
			return hash,nonce
		}

		nonce ++

	}


}














