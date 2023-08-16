package middleware

import (
	"fmt"
	"gameapp/service/authservice"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"log"
)

func Auth(s authservice.Service, config authservice.Config) echo.MiddlewareFunc {
	log.Println("middleware is  here")
	return echojwt.WithConfig(echojwt.Config{
		ContextKey:    "user",
		SigningKey:    []byte(config.SignKey),
		SigningMethod: "HS256",
		//ErrorHandler: func(c echo.Context, err error) error {
		//	if err != nil {
		//		_, ok := err.(*echojwt.TokenExtractionError)
		//		if ok {
		//			fmt.Println("hahahah")
		//		}
		//	}
		//	return nil
		//},
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			claim, err := s.VerifyToken(auth)
			if err != nil {
				fmt.Println("err in auth middleware", err)
				return nil, err
			}
			return claim, nil
		},
	})
}
