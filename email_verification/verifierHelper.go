package email_verification

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	emailverifier "github.com/AfterShip/email-verifier"
	"io"
	"net/smtp"
	"strings"
)

var (
	verifier = emailverifier.NewVerifier()
	otpChars = "1234567890"
)

// checking is mail valid or not!
func emailIsValid(email string) bool {

	ret, err := verifier.Verify(email)
	if err != nil {
		fmt.Println("verify email address failed, error is: ", err)
		return false
	}
	if !ret.Syntax.Valid {
		fmt.Println("email address syntax is invalid")
		return false
	}

	fmt.Println("email validation result", ret)

	return true

}

//encryption part is provided here
func encryptAES(plaintext []byte, key []byte) (ciphertext []byte, err error) {
	k := sha256.Sum256(key)
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

//decryption part is provided here
func decryptAES(ciphertext []byte, key []byte) (plaintext []byte, err error) {
	k := sha256.Sum256(key)
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}

//6 digits code is generated!
func generateOTP(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}

func sendEmail(code, mail string) {
	//put ur e-mail address that you want to sent e-mail by.
	from := "example@gmail.com"
	//put your email' password!!!
	pass := "examplePassword"

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

func splitString(token string) (string, string, string) {
	encryptCodeBase64, _ := base64.StdEncoding.DecodeString(token)
	decryptedCode, err := decryptAES(encryptCodeBase64, []byte(key))
	if err != nil {
		fmt.Println(err)
		return "", "", ""
	}
	splitToken := strings.Split(string(decryptedCode), ",")
	code := splitToken[0]
	sentTime := splitToken[1]
	email := splitToken[2]
	return code, sentTime, email
}
