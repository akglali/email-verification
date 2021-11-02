package email_verification

import (
	"email-verification/helpers"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/smtp"
	"time"
)

//u can define any key with 32 bits
var key = "alisultaniseviyor"

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
	encryptedCode, err := encryptAES([]byte(code), []byte(key))
	if err != nil {
		fmt.Println(err)
		return
	}
	//sent the current time as well to save on local storage or phone storage.
	//You can save the time on ur database to check it
	c.JSON(200, gin.H{
		"encryptedCode": encryptedCode,
		"sent_date":     currentTime,
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
	layout := "2006-01-02 3:4:5 PM"
	currentTime := time.Now().Format(layout)

	// I am parsing the times, so we can compare two times according to same Location
	sentDate, err := time.Parse(layout, body.SentDate)
	currentTimeParse, err := time.Parse(layout, currentTime)

	//getting the differences
	diff := currentTimeParse.Sub(sentDate)

	//getting differences as seconds
	second := int(diff.Seconds())

	encryptCode, _ := base64.StdEncoding.DecodeString(body.EncryptedCode)

	decryptedCode, err := decryptAES(encryptCode, []byte(key))
	if err != nil {
		fmt.Println(err)
		return
	}
	code := body.DecryptedCode

	if second > 30 {
		c.JSON(400, "Your code is expired")
		return
	} else {
		if code == string(decryptedCode) {
			c.JSON(200, "Verification is completed!")
		} else {
			c.JSON(400, "Check your code !!")
		}
	}

}

func sendEmail(code, mail string) {
	//put ur e-mail address that you want to sent e-mail by.
	from := "exampleMail@gmail.com"
	pass := "examplePass"

	to := []string{
		mail,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	message := []byte("To: " + mail + "\r\n" +
		"Subject: Verification Code\r\n" +
		"\r\n" +
		"Hello dear,\r\n" + "Your code is\n" +
		code)

	auth := smtp.PlainAuth("", from, pass, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Is Successfully sent.")

}
