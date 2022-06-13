package main

import (
	"fmt"
	"公链系统开发/TEST/block"
)

func main() {

	chain, err := block.NewBlockChain([]byte("哈哈哈"))

	fmt.Println(err)

	fmt.Println(chain.ListHash)

}
