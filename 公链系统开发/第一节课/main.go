package main

import (
	"fmt"
	"公链系统开发/第一节课/block"
)

func main() {

	block1 := block.NewBlock([]byte("hello"),nil)

	fmt.Println(string(block1.Data))

	block2 := block.NewBlock([]byte("lala"),block1.Hash)

	fmt.Println(string(block2.Data))

	bc := block.NewChain([]byte("创始区块"))

    fmt.Println(string(bc.Blocks[0].Data))

	//把一个信息存到第二个区块
	bc.AddBlock([]byte("第一个区块"))
	fmt.Println(string(bc.Blocks[1].Data))

	//循环打印
	for _,b := range bc.Blocks{
		fmt.Printf("上个区块的Hash值%x\n",b.PrevHash)
		fmt.Printf("时间戳%d\n",b.TimeStamp)
		fmt.Println("信息:  "+string(b.Data))
		fmt.Printf("当前区块的Hash值%x\n",b.Hash)
	}

}
