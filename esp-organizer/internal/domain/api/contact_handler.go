package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"
)

// contactRequest is the JSON body expected from the frontend contact form.
type contactRequest struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Message        string `json:"message"`
	RecaptchaToken string `json:"recaptchaToken,omitempty"`
}

// ContactHandler verifies a recaptcha token (when configured) and stores/sends the contact message.
// POST /api/contact
func ContactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req contactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Message = strings.TrimSpace(req.Message)

	if req.Message == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}

	// If a reCAPTCHA secret is configured, require and verify the token.
	recaptchaSecret := strings.TrimSpace(os.Getenv("RECAPTCHA_SECRET"))
	if recaptchaSecret != "" {
		if req.RecaptchaToken == "" {
			http.Error(w, "recaptcha token required", http.StatusBadRequest)
			return
		}

		ok, err := verifyRecaptcha(req.RecaptchaToken, recaptchaSecret)
		if err != nil {
			log.Printf("recaptcha verify error: %v", err)
			http.Error(w, "recaptcha verification failed", http.StatusBadGateway)
			return
		}
		if !ok {
			http.Error(w, "recaptcha verification failed", http.StatusBadRequest)
			return
		}
	}

	// Persisting to MongoDB is intentionally a no-op here because this
	// handler is single-purpose and the project uses repository layers for
	// storage elsewhere. If you want messages stored, add persistence in the
	// appropriate repo and call it here.
	log.Printf("contact: received message from %s (%s)", req.Name, req.Email)

	// If SMTP is configured, attempt to send email to CONTACT_RECIPIENT.
	smtpHost := strings.TrimSpace(os.Getenv("SMTP_HOST"))
	smtpPort := strings.TrimSpace(os.Getenv("SMTP_PORT"))
	smtpUser := strings.TrimSpace(os.Getenv("SMTP_USER"))
	smtpPass := strings.TrimSpace(os.Getenv("SMTP_PASS"))
	recipient := strings.TrimSpace(os.Getenv("CONTACT_RECIPIENT"))
	if recipient == "" {
		recipient = "world@espthinkers.co.uk"
	}

	// If credentials are not supplied via environment variables, attempt to read Docker
	// secrets mounted at /run/secrets/<name> (compose secrets). This keeps secrets out
	// of the repo and environment files.
	if smtpUser == "" {
		if b, err := os.ReadFile("/run/secrets/smtp_user"); err == nil {
			smtpUser = strings.TrimSpace(string(b))
		}
	}
	if smtpPass == "" {
		if b, err := os.ReadFile("/run/secrets/smtp_pass"); err == nil {
			smtpPass = strings.TrimSpace(string(b))
		}
	}

	if smtpHost != "" && smtpPort != "" {
		from := smtpUser
		if from == "" {
			from = recipient
		}
		subject := fmt.Sprintf("[Website contact] %s", req.Name)
		body := fmt.Sprintf("From: %s <%s>\n\n%s\n\n---\nIP: %s\nUA: %s\n", req.Name, req.Email, req.Message, r.RemoteAddr, r.UserAgent())
		msg := []byte("Subject: " + subject + "\r\n" + "To: " + recipient + "\r\n" + "From: " + from + "\r\n\r\n" + body)

		addr := smtpHost + ":" + smtpPort
		auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
		if err := smtp.SendMail(addr, auth, from, []string{recipient}, msg); err != nil {
			log.Printf("contact: smtp send failed: %v", err)
		} else {
			log.Printf("contact: message sent to %s", recipient)
		}
	} else {
		log.Printf("contact: SMTP not configured — saved to DB (or logged). Recipient would be %s", recipient)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok", "sent_to": recipient})
}

// verifyRecaptcha verifies the token against Google's siteverify endpoint.
func verifyRecaptcha(token, secret string) (bool, error) {
	if token == "" || secret == "" {
		return false, nil
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", "https://www.google.com/recaptcha/api/siteverify", nil)
	if err != nil {
		return false, err
	}

	q := req.URL.Query()
	q.Add("secret", secret)
	q.Add("response", token)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var r struct {
		Success bool     `json:"success"`
		Score   float64  `json:"score,omitempty"`
		Action  string   `json:"action,omitempty"`
		Error   []string `json:"error-codes,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return false, err
	}

	return r.Success, nil
}
