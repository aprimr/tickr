package domain

type OTPType string

const (
	OTPTypeAccountVerification OTPType = "account_verification"
	OTPTypeForgotPassword      OTPType = "forgot_password"
)
