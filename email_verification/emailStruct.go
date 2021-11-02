package email_verification

type emailStruct struct {
	Email string
}

//after we send the code to the email this struct is going to be used
type checkCodeStruct struct {
	EncryptedCode string
	DecryptedCode string
	SentDate      string
}
