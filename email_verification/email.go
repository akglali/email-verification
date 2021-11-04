package email_verification

import (
	"email-verification/helpers"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

//u can define any key
var key = "verysecurekey"

func SetupVerification(rg *gin.RouterGroup) {
	rg.POST("/sendCode", sendVerificationCode)
	rg.POST("/checkCode", checkVerificationCode)

}

func sendVerificationCode(c *gin.Context) {
	body := emailStruct{}
	data, err := c.GetRawData()
	if err != nil {
		helpers.MyAbort(c, "Input format is wrong")
		return
	}
	err = json.Unmarshal(data, &body)
	if !emailIsValid(body.Email) {
		helpers.MyAbort(c, "Check your email type!!!")
		return
	}
	currentTime := time.Now().Format("2006-01-02 3:4:5 PM")
	code, _ := generateOTP(6)
	//email is got by user
	sendEmail(code, body.Email)
	token := code + "," + currentTime + "," + body.Email
	fmt.Println(token)
	encryptedCode, err := encryptAES([]byte(token), []byte(key))
	if err != nil {
		fmt.Println(err)
		return
	}
	//sent the current time as well to save on local storage or phone storage.
	//You can save the time on ur database to check it
	c.JSON(200, gin.H{
		"encryptedCode": encryptedCode,
	})
}

func checkVerificationCode(c *gin.Context) {
	body := checkCodeStruct{}
	data, err := c.GetRawData()
	if err != nil {
		helpers.MyAbort(c, "Input format is wrong")
		return
	}
	err = json.Unmarshal(data, &body)
	if err != nil {
		helpers.MyAbort(c, "Bad Format!")
		return
	}
	fmt.Println(body)
	layout := "2006-01-02 3:4:5 PM"
	currentTime := time.Now().Format(layout)
	//Let!s get 3 element from the token,(1st one is code,2nd one is time,3rd one is the email)
	code, sentTime, email := splitString(body.EncryptedCode)

	// I am parsing the times, so we can compare two times according to same Location
	sentDate, err := time.Parse(layout, sentTime)
	currentTimeParse, err := time.Parse(layout, currentTime)

	//getting the differences
	diff := currentTimeParse.Sub(sentDate)

	//getting differences as seconds
	//diff.Seconds() gets time differences as seconds
	//diff.Hours() gets time differences as hours
	second := int(diff.Minutes())

	// if user change any of letter from the token we provided,user will be faced with an error
	if !emailIsValid(email) {
		helpers.MyAbort(c, "Check your email type!!!")
		return
		// code will be valid for 1 minute. / you can modify it however you like.
	} else if second > 1 {
		helpers.MyAbort(c, "Your code is expired")
		return

	} else {
		if code == body.Code {
			c.JSON(200, "Verification is completed!")
			return
		} else {
			helpers.MyAbort(c, "Check your code!!")
			return
		}
	}

}
