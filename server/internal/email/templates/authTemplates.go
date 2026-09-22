package templates

import "html/template"

// AccountVerificationTemplate is the raw HTTP string for the account verification OTP email
const AccountVerificationTemplate = `
	<body style="margin: 0; padding: 0; background-color: #111111; font-family: Helvetica, Arial, sans-serif; color: #e0e0e0;">
    <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%" style="background-color: #111111; background-image: linear-gradient(to right, #1a1a1a 1px, transparent 1px), linear-gradient(to bottom, #1a1a1a 1px, transparent 1px); background-size: 32px 32px; padding: 60px 20px;">
        <tr>
            <td align="center">
                <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 520px; background-color: #000000; border-radius: 0px; border: 1px solid #262626;">
                    <tr>
                        <td style="padding: 40px;">
                            
                            <!-- Brand -->
                            <div style="font-size: 12px; font-weight: bold; letter-spacing: 2px; color: #3b82f6; text-transform: uppercase; margin-bottom: 24px;">
                                Tickr
                            </div>

                            <!-- Heading -->
                            <h1 style="margin: 0 0 16px 0; color: #ffffff; font-size: 22px; font-weight: 700; letter-spacing: -0.5px;">
                                Welcome, {{.Name}}!
                            </h1>
                            
                            <!-- Body -->
                            <p style="margin: 0 0 28px 0; font-size: 15px; line-height: 1.6; color: #999999;">
                                Please use the verification code below to verify your account. This code will expire in <strong>15 minutes</strong>.
                            </p>
                            
                            <!-- OTP -->
                            <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%" style="margin-bottom: 28px;">
                                <tr>
                                    <td align="center" style="background-color: #080808; border-radius: 0px; padding: 22px; border: 1px solid #3b82f6;">
                                        <span style="font-size: 28px; font-weight: bold; letter-spacing: 12px; color: #ffffff;">{{.OTP}}</span>
                                    </td>
                                </tr>
                            </table>
                            
                            <!-- Footer -->
                            <p style="margin: 0; font-size: 13px; line-height: 1.5; color: #555555;">
                                If you didn't register for Tickr, you can safely ignore this email.
                            </p>
                            
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
	</body>
`

// Parse string to HTML template
var AccountVerification = template.Must(template.New("verification").Parse(AccountVerificationTemplate))
