package main

import (
	"fmt"
	"公链系统开发/第二节课/block"
)

func main() {

	bc := block.NewChain([]byte("创始区块"))
	//把一个信息存到第二个区块
	bc.AddBlock([]byte("第一个区块"))

	fmt.Println()
    fmt.Println("创始区块的信息",string(bc.Blocks[0].Data))
    fmt.Println("创世区块nonce值为",bc.Blocks[0].Nonce)
	fmt.Println("第一个区块nonce值为",bc.Blocks[1].Nonce)


}
