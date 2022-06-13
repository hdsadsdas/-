package pow

import (
	"bytes"
	"math/big"
	"strconv"
	"公链系统开发/TEST/tools"
)

/**
* @author : 哈哈
* @email : 598421227@qq.com
* @phone : 18816473550
* @DateTime : 2022/4/11 9:13
**/

const BITS = 1

type POW struct {
	PrevHash  []byte
	TimeStamp int64
	Data      []byte
	Target    *big.Int
}

func NewPOW(prevHash []byte, timeStamp int64, data []byte) *POW {

	target := big.NewInt(0)

	target = target.Lsh(target, 255-BITS)

	pow := POW{
		PrevHash:  prevHash,
		TimeStamp: timeStamp,
		Data:      data,
		Target:    target,
	}

	return &pow

}

func (pow *POW) Run() (Hash []byte, Nonce int64) {

	num := big.NewInt(0)

	timeStamp := []byte(strconv.FormatInt(pow.TimeStamp, 10))

	for {
		NonceBytes := []byte(strconv.FormatInt(Nonce, 10))
		hashByets := bytes.Join([][]byte{pow.Data, pow.PrevHash, timeStamp, NonceBytes}, []byte{})
		Hash = tools.GetHash(hashByets)
		num := num.SetBytes(Hash)

		if num.Cmp(pow.Target) == -1 {

			return

		}

		Nonce++
	}

	return nil, 0

}
