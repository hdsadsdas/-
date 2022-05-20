package main

import (
	"fmt"
	"公链系统开发/第五节课/block"
)

func main() {

		bc,err:=block.NewChain([]byte("创世区块"))
		defer  bc.DB.Close()
		if err !=nil{
			fmt.Println(err.Error())
			return
		}
	bc.AddBlock([]byte("2223333"))
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


	//bk, err := iterator.Next()
	//if err !=nil{
	//	fmt.Println(err.Error())
	//	return
	//}
	//fmt.Println(string(bk.Data))
	//
	//bk1, err := iterator.Next()
	//fmt.Println(string(bk1.Data))
	//
	//bk2, err := iterator.Next()
	//fmt.Println(string(bk2.Data))

	//fmt.Println(bc.LastHash)
	//   if err !=nil{
	//   	fmt.Println("失败")
	//   }else{
	//   	fmt.Println("成功")
	//   }


}
