package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"os"
	"strconv"
)

func main() {

	Run()

}

func Run() {

	args := os.Args

	switch args[1] {

	case "login":
		login()

	case "lookMoney":
		lookMoney()

	case "saveMoney":
		saveMoney()

	case "withdrawMoney":
		withdrawMoney()

	case "help":
		fmt.Println("login")
		fmt.Println("--username   用户名   曾凯")
		fmt.Println("--password   用户密码  123456")
		fmt.Println("lookMoney")
		fmt.Println("saveMoney")
		fmt.Println("--money   取款金额")
		fmt.Println("withdrawMoney")
		fmt.Println("--money    存款金额")

	default:
		fmt.Println("未知选项  ")

	}

}

func login() {
	login := flag.NewFlagSet("login", flag.ExitOnError)

	username := login.String("username", "", "用户名")
	password := login.String("password", "", "用户密码")

	login.Parse(os.Args[2:])

	if *username == "曾凯" && *password == "123456" {

		var db, _ = bolt.Open("./chain.db", 0600, nil)
		defer db.Close()

		db.Update(func(tx *bolt.Tx) error {

			bucket := tx.Bucket([]byte("login"))

			if bucket == nil {

				createBucket, err := tx.CreateBucket([]byte("login"))
				if err != nil {
					return err
				}

				createBucket.Put([]byte("money"), []byte("100"))

			} else {

				fmt.Println("你已经登录")
				return nil
			}

			return nil

		})

	} else {

		fmt.Println("账号或密码错误")
		return

	}

}

func lookMoney() {

	var db, _ = bolt.Open("./chain.db", 0600, nil)

	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte("login"))
		if bucket == nil {
			fmt.Println("请先登录")
			return errors.New("未登录")
		} else {

			money := bucket.Get([]byte("money"))

			atoi,_ := strconv.Atoi(string(money))


			fmt.Println("取钱成功当前账号还剩", atoi, "元")

			return nil
		}

	})

}

func saveMoney() {

	var db, _ = bolt.Open("./chain.db", 0600, nil)

	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte("login"))

		if bucket == nil {
			fmt.Println("请先登录")
			return errors.New("未登录")
		} else {

			saveMoney := flag.NewFlagSet("saveMoney", flag.ExitOnError)

			sMoney := saveMoney.Int("money", 0, "存款金额")

			saveMoney.Parse(os.Args[2:])

			money := bucket.Get([]byte("money"))

			atoi, _ := strconv.Atoi(string(money))

			if *sMoney > atoi || *sMoney < 0 {
				fmt.Println("当前账号余额不足或输入数额有误")
				return errors.New("当前账号余额不足或输入数额有误")
			} else {

				atoi = atoi - *sMoney



				bucket.Put([]byte("money"), []byte(strconv.Itoa(atoi)))
				fmt.Println("取钱成功当前账号还剩", atoi, "元")
			}

		}

		return nil

	})

}

func withdrawMoney() {

	var db, _ = bolt.Open("./chain.db", 0600, nil)

	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte("login"))

		if bucket == nil {
			fmt.Println("请先登录")
			return errors.New("未登录")
		} else {

			withdrawMoney := flag.NewFlagSet("withdrawMoney", flag.ExitOnError)

			wMoney := withdrawMoney.Int("money", 0, "存款金额")

			withdrawMoney.Parse(os.Args[2:])

			money := bucket.Get([]byte("money"))

			atoi, _ := strconv.Atoi(string(money))

			if *wMoney < 0 {
				fmt.Println("输入数额有误")
				return errors.New("当前账号余额不足或输入数额有误")
			} else {

				atoi = atoi + *wMoney

				fmt.Println(atoi)

				bucket.Put([]byte("money"), []byte(strconv.Itoa(atoi)))
				fmt.Println("存钱成功当前账号还剩", atoi, "元")
			}

		}

		return nil

	})

}
