package main

import (
	"email-verification/email_verification"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(func(context *gin.Context) {
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "*")
		context.Header("Access-Control-Allow-Methods", "*")
		if context.Request.Method == "OPTIONS" {
			context.Status(200)
			context.Abort()
		}
	})
	sendVerifyCode := r.Group("/verify")
	email_verification.SetupVerification(sendVerifyCode)

	err := r.Run(":8000")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}
