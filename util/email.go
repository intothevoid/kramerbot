package util

import (
	"fmt"
	"html"
	"log"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/intothevoid/kramerbot/models"
)

// EmailService sends transactional emails via SMTP (STARTTLS, port 587).
// If Host or Username is empty, Send is a no-op (graceful degradation for
// self-hosted deployments without an SMTP server configured).
type EmailService struct {
	cfg SMTPConfig
}

// NewEmailService creates an EmailService from the provided config.
func NewEmailService(cfg SMTPConfig) *EmailService {
	return &EmailService{cfg: cfg}
}

// Enabled reports whether SMTP is configured (Host is sufficient; auth is optional).
func (s *EmailService) Enabled() bool {
	return s.cfg.Host != ""
}

// loginAuth implements smtp.Auth using the AUTH LOGIN mechanism required by
// Outlook/Hotmail. Go's stdlib only provides AUTH PLAIN which Outlook rejects.
type loginAuth struct {
	username, password string
}

func (a *loginAuth) Start(_ *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", nil, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch strings.ToLower(strings.TrimSpace(string(fromServer))) {
		case "username:":
			return []byte(a.username), nil
		case "password:":
			return []byte(a.password), nil
		default:
			return nil, fmt.Errorf("unexpected server challenge: %q", fromServer)
		}
	}
	return nil, nil
}

// envelopeFrom extracts the bare email address from the From field for use as
// the SMTP envelope sender. smtp.SendMail requires a plain address with no
// display name (e.g. "noreply@example.com", not "Name <noreply@example.com>").
func (s *EmailService) envelopeFrom() string {
	addr, err := mail.ParseAddress(s.cfg.From)
	if err != nil {
		return s.cfg.From // fallback: already a plain address
	}
	return addr.Address
}

// Send delivers an HTML email. Returns nil without sending if SMTP is not configured.
// When Username is empty, authentication is skipped (suitable for unauthenticated local SMTP relays).
func (s *EmailService) Send(to, subject, htmlBody string) error {
	if !s.Enabled() {
		return nil
	}

	var auth smtp.Auth
	if s.cfg.Username != "" {
		auth = &loginAuth{s.cfg.Username, s.cfg.Password}
	}
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		s.cfg.From, to, subject, htmlBody,
	)
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	log.Printf("[smtp] sending email to=%s via=%s from=%s", to, addr, s.envelopeFrom())
	return smtp.SendMail(addr, auth, s.envelopeFrom(), []string{to}, []byte(msg))
}

// SendVerificationEmail sends the account verification email.
func (s *EmailService) SendVerificationEmail(to, verifyLink string) error {
	subject := "Verify your KramerBot account"
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family:sans-serif;max-width:480px;margin:0 auto;padding:32px 16px;background:#fdf8f0">
  <div style="background:#c0392b;border-radius:12px 12px 0 0;padding:24px;text-align:center">
    <span style="color:#fff;font-size:22px;font-weight:bold">KramerBot - Aussie Deals</span>
  </div>
  <div style="background:#fff;border-radius:0 0 12px 12px;padding:32px;border:1px solid #e5e7eb;border-top:none">
    <h2 style="color:#1a1a1a;margin-top:0">Verify your email</h2>
    <p style="color:#555">Thanks for signing up! Click the button below to verify your email address and activate your account.</p>
    <a href="%s" style="display:inline-block;background:#c0392b;color:#fff;padding:12px 28px;border-radius:8px;text-decoration:none;font-weight:bold;margin:16px 0">
      Verify Email Address
    </a>
    <p style="color:#888;font-size:13px">This link expires in 24 hours. If you didn't create an account, you can ignore this email.</p>
    <hr style="border:none;border-top:1px solid #e5e7eb;margin:24px 0">
    <p style="color:#aaa;font-size:12px">Or copy this link into your browser:<br>%s</p>
  </div>
</body>
</html>`, verifyLink, verifyLink)
	return s.Send(to, subject, body)
}

// SendPasswordResetEmail sends the password reset email.
func (s *EmailService) SendPasswordResetEmail(to, resetLink string) error {
	subject := "Reset your KramerBot password"
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family:sans-serif;max-width:480px;margin:0 auto;padding:32px 16px;background:#fdf8f0">
  <div style="background:#c0392b;border-radius:12px 12px 0 0;padding:24px;text-align:center">
    <span style="color:#fff;font-size:22px;font-weight:bold">KramerBot - Aussie Deals</span>
  </div>
  <div style="background:#fff;border-radius:0 0 12px 12px;padding:32px;border:1px solid #e5e7eb;border-top:none">
    <h2 style="color:#1a1a1a;margin-top:0">Reset your password</h2>
    <p style="color:#555">We received a request to reset the password for your account. Click the button below to choose a new password.</p>
    <a href="%s" style="display:inline-block;background:#c0392b;color:#fff;padding:12px 28px;border-radius:8px;text-decoration:none;font-weight:bold;margin:16px 0">
      Reset Password
    </a>
    <p style="color:#888;font-size:13px">This link expires in 1 hour. If you didn't request a password reset, you can ignore this email.</p>
    <hr style="border:none;border-top:1px solid #e5e7eb;margin:24px 0">
    <p style="color:#aaa;font-size:12px">Or copy this link into your browser:<br>%s</p>
  </div>
</body>
</html>`, resetLink, resetLink)
	return s.Send(to, subject, body)
}

// SendDailySummary sends the nightly deal digest email.
// ozbDeals should be pre-filtered to OZB_SUPER type, sorted by votes descending.
// amzDeals should be pre-filtered to AMZ_DAILY type.
func (s *EmailService) SendDailySummary(to string, ozbDeals []models.OzBargainDeal, amzDeals []models.CamCamCamDeal) error {
	subject := "KramerBot Daily Deal Summary 🔥"

	const limit = 10
	var ozbRows strings.Builder
	for i, d := range ozbDeals {
		if i >= limit {
			break
		}
		ozbRows.WriteString(fmt.Sprintf(`
      <tr>
        <td style="padding:10px 0;border-bottom:1px solid #f0ede4">
          <a href="%s" style="color:#c0392b;font-weight:bold;text-decoration:none;font-size:14px">%s</a>
          <div style="margin-top:4px">
            <span style="background:#F5C518;color:#1a1a1a;border-radius:4px;padding:2px 7px;font-size:12px;font-weight:bold">🔺 %s votes</span>
            <span style="color:#aaa;font-size:12px;margin-left:8px">%s</span>
          </div>
        </td>
      </tr>`, html.EscapeString(d.Url), html.EscapeString(d.Title), html.EscapeString(d.Upvotes), html.EscapeString(d.PostedOn)))
	}
	if ozbRows.Len() == 0 {
		ozbRows.WriteString(`<tr><td style="padding:10px 0;color:#aaa;font-size:13px">No top deals today.</td></tr>`)
	}

	var amzRows strings.Builder
	for i, d := range amzDeals {
		if i >= limit {
			break
		}
		amzRows.WriteString(fmt.Sprintf(`
      <tr>
        <td style="padding:10px 0;border-bottom:1px solid #f0ede4">
          <a href="%s" style="color:#c0392b;font-weight:bold;text-decoration:none;font-size:14px">%s</a>
          <div style="margin-top:4px">
            <span style="background:#e8f4f8;color:#1a1a1a;border-radius:4px;padding:2px 7px;font-size:12px">📦 Amazon Daily</span>
          </div>
        </td>
      </tr>`, html.EscapeString(d.Url), html.EscapeString(d.Title)))
	}
	if amzRows.Len() == 0 {
		amzRows.WriteString(`<tr><td style="padding:10px 0;color:#aaa;font-size:13px">No Amazon daily deals today.</td></tr>`)
	}

	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family:sans-serif;max-width:560px;margin:0 auto;padding:32px 16px;background:#FFFEF7">
  <div style="background:#c0392b;border-radius:12px 12px 0 0;padding:24px;text-align:center">
    <span style="color:#fff;font-size:22px;font-weight:bold">KramerBot — Daily Summary</span>
  </div>
  <div style="background:#fff;border-radius:0 0 12px 12px;padding:32px;border:1px solid #e5e7eb;border-top:none">
    <h2 style="color:#c0392b;margin-top:0;font-size:17px">🔥 Top OzBargain Deals</h2>
    <table style="width:100%%;border-collapse:collapse">%s</table>
    <h2 style="color:#c0392b;margin-top:28px;font-size:17px">📦 Amazon Daily Deals</h2>
    <table style="width:100%%;border-collapse:collapse">%s</table>
    <hr style="border:none;border-top:1px solid #e5e7eb;margin:28px 0">
    <p style="color:#aaa;font-size:12px;text-align:center">
      KramerBot Daily Summary · Manage your preferences in the Dashboard
    </p>
  </div>
</body>
</html>`, ozbRows.String(), amzRows.String())

	return s.Send(to, subject, body)
}
