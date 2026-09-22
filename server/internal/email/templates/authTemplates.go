package templates

import "fmt"

// AccountVerification generates the HTML body for the email verification OTP
func AccountVerification(name, otp string) string {
	return fmt.Sprintf(`
		<div>
			<h2>Welcome to Tickr, %s!</h2>
			<p>Thanks for signing up. Please use the verification code below to verify your account:</p>
			<div>%s</div>
		</div>
	`, name, otp)
}
