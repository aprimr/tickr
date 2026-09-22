package templates

import "html/template"

// AccountVerificationTemplate is the raw HTTP string for the account verification OTP email
const AccountVerificationTemplate = `
<div>
	<h2>Welcome to Tickr, {{.Name}}!</h2>
	<p>Thanks for signing up. Please use the verification code below to verify your account:</p>
	<div style="font-size: 24px; font-weight: bold;">{{.OTP}}</div>
</div>
`

// Parse string to HTML template
var AccountVerification = template.Must(template.New("verification").Parse(AccountVerificationTemplate))
