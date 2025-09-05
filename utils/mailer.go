package utils

import (
	"auth-service/config"
	"fmt"
	"net/smtp"
)

func SendOTPEmail(to, otp string) error {
	auth := smtp.PlainAuth("", config.AppConfig.SmtpUser, config.AppConfig.SmtpPassword, config.AppConfig.SmtpHost)

	// Header email
	subject := "Subject: Your One-Time Password (OTP)\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	// Body email dengan HTML
	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; padding: 20px; background-color: #f4f4f4;">
			<div style="max-width: 600px; margin: auto; background-color: #fff; padding: 30px; border-radius: 8px; box-shadow: 0 0 10px rgba(0,0,0,0.1);">
				<h2 style="color: #333;">🔐 Your OTP Code</h2>
				<p style="font-size: 16px; color: #555;">Use the following One-Time Password (OTP) to proceed:</p>
				<div style="text-align: center; margin: 20px 0;">
					<span style="font-size: 32px; font-weight: bold; color: #2c3e50;">%s</span>
				</div>
				<p style="font-size: 14px; color: #888;">This OTP is valid for a limited time. Please do not share it with anyone.</p>
			</div>
		</body>
		</html>
	`, otp)

	msg := []byte(subject + mime + body)

	addr := fmt.Sprintf("%s:%s", config.AppConfig.SmtpHost, config.AppConfig.SmtpPort)
	return smtp.SendMail(addr, auth, config.AppConfig.AppEmail, []string{to}, msg)
}
