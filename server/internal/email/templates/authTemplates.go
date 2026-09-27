package templates

import "html/template"

// AccountVerificationTemplate is the raw HTTP string for the account verification OTP email
const AccountVerificationTemplate = `
	<body style="margin: 0; padding: 0; background-color: #111111; font-family: Helvetica, Arial, sans-serif; color: #e0e0e0;">

    <!-- Hidden Text for Inbox Preview -->
    <div style="display: none; font-size: 1px; color: #111111; line-height: 1px; max-height: 0px; max-width: 0px; opacity: 0; overflow: hidden; mso-hide: all;">
        Your Tickr verification code is {{.OTP}}. It expires in 15 minutes.
    </div>
    <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%" style="background-color: #111111; background-image: linear-gradient(to right, #053061 1px, transparent 1px), linear-gradient(to bottom, #053061 1px, transparent 1px); background-size: 32px 32px; padding: 60px 20px;">
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
                                Hi, {{.Name}}!
                            </h1>
                            
                            <!-- Body -->
                            <p style="margin: 0 0 28px 0; font-size: 15px; line-height: 1.6; color: #999999;">
                                Please use the verification code below to complete your verification and access your Tickr account. This code will expire in <strong>15 minutes</strong>.
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
                                If you did not request this, please ignore this message.
                            </p>

                            <!-- Copyright & link Footer -->
                            <div style="margin-top: 24px; padding-top: 20px; border-top: 1px solid #222222; font-size: 11px; color: #444444; text-align: left;">
                                &copy; 2026 Tickr. All rights reserved.
                            </div>                  
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
