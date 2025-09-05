package utils

import (
	"auth-service/config"
	"fmt"
	"net/smtp"
)

// SendOTPEmail mengirim email berisi OTP ke pengguna.
// Sekarang menerima cfg *config.Config sebagai parameter.
func SendOTPEmail(to, otp string, cfg *config.Config) error {
	auth := smtp.PlainAuth("", cfg.SmtpUser, cfg.SmtpPassword, cfg.SmtpHost)

	// Header email
	subject := "Subject: Your One-Time Password (OTP)\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	// Body email dengan HTML
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Your OTP Code</h2>
			<p>Use the following One-Time Password (OTP) to proceed:</p>
			<h3>%s</h3>
			<p>This OTP is valid for %d minutes.</p>
		</body>
		</html>
	`, otp, int(cfg.OTPDuration.Minutes()))

	msg := []byte(subject + mime + body)

	addr := fmt.Sprintf("%s:%s", cfg.SmtpHost, cfg.SmtpPort)
	return smtp.SendMail(addr, auth, cfg.AppEmail, []string{to}, msg)
}

// SendResetPasswordEmail mengirim email berisi token untuk reset password.
func SendResetPasswordEmail(to, token string, cfg *config.Config) error {
	auth := smtp.PlainAuth("", cfg.SmtpUser, cfg.SmtpPassword, cfg.SmtpHost)

	subject := "Subject: Your Password Reset Token\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Password Reset Request</h2>
			<p>Use the following token to reset your password:</p>
			<h3>%s</h3>
			<p>This token is valid for %d minutes. If you did not request a password reset, please ignore this email.</p>
		</body>
		</html>
	`, token, int(cfg.ResetPasswordTokenDuration.Minutes()))

	msg := []byte(subject + mime + body)
	addr := fmt.Sprintf("%s:%s", cfg.SmtpHost, cfg.SmtpPort)
	return smtp.SendMail(addr, auth, cfg.AppEmail, []string{to}, msg)
}
