package main

import (
	"encoding/json"
	"fmt"
	"gameapp/repository/mysql"
	"gameapp/service/userservice"
	"io"
	"net/http"
)

func main() {

	fmt.Println("welcome to game app")
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/profile", profile)
	err := http.ListenAndServe("localhost:7777", nil)
	if err != nil {
		fmt.Println("err in listenandserver : ", err)
	}
}

func register(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		fmt.Println("bad method")
		return
	}
	data, err := io.ReadAll(req.Body)
	if err != nil {
		w.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

		return
	}
	var rr userservice.RegisterRequest
	uErr := json.Unmarshal(data, &rr)
	if uErr != nil {
		fmt.Println("uErr", uErr)
		w.Write([]byte(
			fmt.Sprintf(`{"error": "%v"}`, uErr.Error()),
		))

		return
	}
	urepo := mysql.New()
	usvc := userservice.NewUserSvc(urepo)
	response, err := usvc.Register(rr)
	if err != nil {
		w.Write([]byte(
			fmt.Sprintf(`{"error": "%v"}`, err.Error()),
		))
		return
	}
	fmt.Println("res user register ", response)
	w.Write([]byte(`{ "message" : "done ! registred"}`))
}

func login(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		fmt.Println("bad method")
		return
	}
	data, err := io.ReadAll(req.Body)
	if err != nil {
		w.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

		return
	}
	var lr userservice.LoginRequest
	uErr := json.Unmarshal(data, &lr)
	if uErr != nil {
		w.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, uErr.Error()),
		))

		return
	}
	urepo := mysql.New()
	usvc := userservice.NewUserSvc(urepo)
	response, err := usvc.Login(lr)
	if err != nil {
		w.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

		return
	}
	fmt.Println("res user login ", response)
	w.Write([]byte(fmt.Sprintf(`{ "token" : %s }`, response.Token)))
}

func profile(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		fmt.Println("bad method")
		return
	}
	data, err := io.ReadAll(req.Body)
	if err != nil {
		w.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

		return
	}
	var pr userservice.ProfileRequest
	uErr := json.Unmarshal(data, &pr)
	if uErr != nil {
		w.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, uErr.Error()),
		))

		return
	}
	urepo := mysql.New()
	usvc := userservice.NewUserSvc(urepo)
	response, err := usvc.GetProfile(pr)
	if err != nil {
		w.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

		return
	}
	fmt.Println("res user login ", response)
	w.Write([]byte(response.Name))
}
