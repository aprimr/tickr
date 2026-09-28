package email

type VerificationData struct {
	Name string
	OTP  string
}

type PasswordResetSuccessData struct {
	Name string
}
