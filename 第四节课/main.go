package main

import (
	"fmt"
	"公链系统开发/第四节课/block"
)

func main() {

	bc, err := block.NewChain([]byte("创始区块"))
	defer bc.DB.Close()
	if err != nil {
		fmt.Println("hahha",err.Error())
		return
	}

	err = bc.AddBlock([]byte("111"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = bc.AddBlock([]byte("222"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	iterator := bc.Iterator()
	for{
		if iterator.HasNext(){
			bk, err := iterator.Next()
			if err !=nil{
				fmt.Println(err.Error())
				return
			}
			fmt.Println(string(bk.Data))
		}else{
			break
		}
	}

}
