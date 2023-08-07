package main

import (
	"encoding/json"
	"fmt"
	"gameapp/repository/mysql"
	"gameapp/service/authservice"
	"gameapp/service/userservice"
	"io"
	"net/http"
	"time"
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
	authSvc := authservice.New("secret", "at", "rt", time.Hour*24, time.Hour*24*7)
	usvc := userservice.NewUserSvc(urepo, authSvc)
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
	authSvc := authservice.New("secret", "at", "rt", time.Hour*24, time.Hour*24*7)
	usvc := userservice.NewUserSvc(urepo, authSvc)
	response, err := usvc.Login(lr)
	if err != nil {
		w.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

		return
	}
	fmt.Println("res user login ", response)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{ "access token is" : "%s", "refresh token is" : "%s" }`, response.AccessToken, response.RefreshToken)))
}

func profile(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		fmt.Println("bad method")
		return
	}
	urepo := mysql.New()
	authSvc := authservice.New("secret", "at", "rt", time.Hour*24, time.Hour*24*7)
	usvc := userservice.NewUserSvc(urepo, authSvc)
	authToken := req.Header.Get("Authorization")
	c, err := authSvc.VerifyToken(authToken)
	if err != nil {
		w.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

		return
	}

	response, err := usvc.GetProfile(userservice.ProfileRequest{UserID: c.UserID})
	if err != nil {
		w.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))

		return
	}
	fmt.Println("res user login ", response)
	w.Write([]byte(response.Name))
}
