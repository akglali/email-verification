package email_verification

type emailStruct struct {
	Email string
}

type checkCodeStruct struct {
	EncryptedCode string
	DecryptedCode string
	SentDate      string
}
