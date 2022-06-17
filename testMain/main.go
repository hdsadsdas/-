package main

import (
	"flag"
	"fmt"
	"os"
)

/**
* @author : 哈哈
* @email : 598421227@qq.com
* @phone : 18816473550
* @DateTime : 2022/3/28 9:15
**/

func main() {

	sum := flag.NewFlagSet("sum", flag.ExitOnError)

	if os.Args[1] == "sum" {
		num1 := sum.Int("num1", 0, "num1")
		num2 := sum.Int("num2", 0, "num2")
		sum.Parse(os.Args[2:])
		fmt.Println("zk")
		fmt.Println(*num1+*num2)
	}

	sub := flag.NewFlagSet("sub", flag.ExitOnError)

	if os.Args[1] == "sub" {
		num1 := sub.Int("num1", 0, "num1")
		num2 := sub.Int("num2", 0, "num2")
		sub.Parse(os.Args[2:])
		fmt.Println("zk")
		fmt.Println(*num1-*num2)
	}

	mul := flag.NewFlagSet("mul", flag.ExitOnError)

	if os.Args[1] == "mul" {
		num1 := mul.Int("num1", 0, "num1")
		num2 := mul.Int("num2", 0, "num2")
		mul.Parse(os.Args[2:])
		fmt.Println("zk")
		fmt.Println(*num1/(*num2))
	}

	div := flag.NewFlagSet("div", flag.ExitOnError)

	if os.Args[1] == "div" {
		num1 := div.Int("num1", 0, "num1")
		num2 := div.Int("num2", 0, "num2")
		div.Parse(os.Args[2:])
		fmt.Println("zk")
		fmt.Println(*num1*(*num2))
	}







	return

	today := flag.NewFlagSet("today", flag.ExitOnError)

	if os.Args[1] == "today" {

		day := today.String("day", "星期天", "星期")
		today.Parse(os.Args[2:])
		fmt.Println(*day)

	}

	return

	args := os.Args

	for key,value := range args{

		fmt.Println(key,value)

	}

	age := flag.Int("age", 18, "年龄")

	flag.Parse()

	fmt.Println(*age)

}

func sum()  {

}