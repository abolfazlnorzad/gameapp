package main

import (
	"fmt"
	"gameapp/entity"
	"gameapp/repository/mysql"
)

func main() {
	fmt.Println("find")
	my := mysql.New()

	unique, err := my.IsPhoneNumberUnique("09011216131")
	if err != nil {
		fmt.Println("unique err", err)
	}

	fmt.Println("unique res", unique)

	fmt.Println("--------------")
	fmt.Println("create new user")

	u, err := my.Create(entity.User{
		ID:          0,
		Name:        "ilia",
		PhoneNumber: "0902",
		Password:    "pass",
	})
	if err != nil {
		fmt.Println("create user err : ", err)
	}

	fmt.Println("create user res", u)

	fmt.Println("---------")

}
