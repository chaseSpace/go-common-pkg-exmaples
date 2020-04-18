package main

import (
	"fmt"
	uuid "github.com/iris-contrib/go.uuid"
	"gopkg.in/jose.v1/crypto"
	"gopkg.in/jose.v1/jws"
	"gopkg.in/jose.v1/jwt"
	"log"
	"testing"
	"time"
)

// 	gopkg.in/jose.v1 v1.0.0-20161127122323-a941c3995164

var cl = make(jws.Claims)

func sign() []byte {
	// 对数据进行签名 得到token
	cl.SetAudience("http://gxt.com")
	cl.SetExpiration(time.Now().Add(2 * time.Second))
	//time.Sleep(time.Second)
	cl.SetNotBefore(time.Now())
	cl.SetSubject("mini_app")
	cl.SetIssuer("golang_svr")
	cl.Set("my_key", 1) // 这个值下面会用到
	cl.SetIssuedAt(time.Now())
	Uuid, _ := uuid.NewV4()
	cl.SetJWTID(Uuid.String())

	token := jws.NewJWT(cl, crypto.SigningMethodHS512)
	b, err := token.Serialize([]byte("secret_key"))
	if err != nil {
		log.Fatal(0, err)
	}
	log.Println("new token:", string(b))
	return b
}

func verify(token_b []byte) {
	// 解析、验证token

	// 按照预定结构来验证
	_Validator := jwt.Validator{
		Expected: jwt.Claims{"aud": "http://gxt.com", // 这里填的是必须要一致的值，传入的token没有就不合格
			"sub": "mini_app",
			"iss": "golang_svr"},
		EXP: 0,
		NBF: 0,
		Fn: func(claims jwt.Claims) error { // 一个自定义函数
			if !claims.Has("my_key") {
				return fmt.Errorf("token err, `my_key` not found!")
			}
			return nil
		},
	}
	token, err := jws.ParseJWT(token_b)
	if err != nil {
		log.Fatal("Wrong token err:", err)
	}

	// 调用验证方法
	err = _Validator.Validate(jws.NewJWT(jws.Claims(token.Claims()), crypto.SigningMethodHS384))
	if err != nil {
		log.Fatal("Validate token fail:", err)
	}
}

func Test_jwt(t *testing.T) {
	log.Printf("%d\n", time.Now())
	//NewToken := sign()
	//verify(NewToken)
}
