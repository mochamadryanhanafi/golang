package utils

import (
	"auth-service/config"
	"fmt"
	"net/smtp"
	"time"
)

// SendOTPEmail mengirim email berisi OTP ke pengguna.
func SendOTPEmail(to, otp string, cfg *config.Config) error {
	auth := smtp.PlainAuth("", cfg.SmtpUser, cfg.SmtpPassword, cfg.SmtpHost)

	// Header email
	subject := "Subject: Your One-Time Password (OTP)\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	// Body email dengan HTML
	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="id">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f4f4f4;
					margin: 0;
					padding: 0;
				}
				.container {
					max-width: 600px;
					margin: 40px auto;
					background-color: #ffffff;
					border-radius: 8px;
					overflow: hidden;
					box-shadow: 0 4px 15px rgba(0,0,0,0.1);
				}
				.header {
					background-color: #007bff;
					color: #ffffff;
					padding: 20px;
					text-align: center;
				}
				.header h1 {
					margin: 0;
					font-size: 24px;
				}
				.content {
					padding: 30px;
					text-align: center;
					color: #333333;
				}
				.content p {
					font-size: 16px;
					line-height: 1.5;
				}
				.otp-code {
					display: inline-block;
					background-color: #e9ecef;
					color: #000000;
					font-size: 36px;
					font-weight: bold;
					letter-spacing: 4px;
					padding: 15px 25px;
					border-radius: 6px;
					margin: 20px 0;
				}
				.footer {
					background-color: #f8f9fa;
					color: #6c757d;
					font-size: 12px;
					text-align: center;
					padding: 20px;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>Kode Verifikasi Anda</h1>
				</div>
				<div class="content">
					<p>Silakan gunakan kode berikut untuk menyelesaikan proses verifikasi Anda. 
					Kode ini hanya berlaku selama <strong>%d menit</strong>.</p>
					<div class="otp-code">%s</div>
					<p>Jika Anda tidak meminta kode ini, mohon abaikan email ini demi keamanan akun Anda.</p>
				</div>
				<div class="footer">
					<p>&copy; %d Nama Perusahaan Anda. Semua Hak Cipta Dilindungi.</p>
				</div>
			</div>
		</body>
		</html>
	`, int(cfg.OTPDuration.Minutes()), otp, time.Now().Year())

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
