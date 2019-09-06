package main

import (
	"dcx.com/auth/provider/password"
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

var (
	flagSet    = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	userFlag   = flagSet.String("u", "", "新账号用户名")
	passFlag   = flagSet.String("p", "", "新账号密码")
	roleFlag   = flagSet.String("r", "admin", "新账号角色")
	nameFlag   = flagSet.String("n", "", "新账号显示名字")
	dbFlag     = flagSet.String("db", "", "数据库名")
	dbUserFlag = flagSet.String("duser", "root", "数据库用户名")
	dbPassFlag = flagSet.String("dpass", "", "数据库密码")
	dbHostFlag = flagSet.String("dhost", "127.0.0.1:3306", "数据库地址,eg:127.0.0.1:3306")
)

func main() {

	flagSet.Parse(os.Args[1:])

	var (
		loginProvider = password.New()
		err           error
	)

	if *userFlag == "" {
		panic("请使用-u=输入新账号用户名")
	}

	if *passFlag == "" {
		panic("请使用-p=输入新账号密码")
	}

	if *nameFlag == "" {
		*nameFlag = *userFlag
	}

	loginProvider.DB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=true&loc=Asia%%2FShanghai", *dbUserFlag, *dbPassFlag, *dbHostFlag, *dbFlag))
	if err != nil {
		panic(err)
	}

	if user, err := loginProvider.AddUser(*userFlag, *passFlag, *roleFlag, *nameFlag); err != nil {
		panic(err)
	} else {
		fmt.Printf("创建用户成功，ID：%v，用户名：%v，密码：%v，角色：%v，显示名字：%v", user.ID, user.LoginName, *passFlag, user.Role, user.Name)
	}

}
