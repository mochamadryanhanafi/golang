package utils

import (
	"auth-service/config"
	"fmt"
	"net/smtp"
)

func SendOTPEmail(to, otp string) error {
	auth := smtp.PlainAuth("", config.AppConfig.SmtpUser, config.AppConfig.SmtpPassword, config.AppConfig.SmtpHost)
	msg := []byte(fmt.Sprintf("Subject: Your OTP Code\n\nYour OTP is: %s", otp))

	addr := fmt.Sprintf("%s:%s", config.AppConfig.SmtpHost, config.AppConfig.SmtpPort)
	return smtp.SendMail(addr, auth, config.AppConfig.AppEmail, []string{to}, msg)
}
