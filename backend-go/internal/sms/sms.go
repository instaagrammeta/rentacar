// Package sms integrates with the zudsms.tj SMS gateway.
package sms

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"sync/atomic"
	"time"
)

// Sender sends SMS messages through the zudsms.tj HTTP API.
type Sender struct {
	Enabled bool
	URL     string
	Login   string
	From    string
	Secret  string
	client  *http.Client
}

// New builds an SMS sender.
func New(enabled bool, url, login, from, secret string) *Sender {
	return &Sender{
		Enabled: enabled,
		URL:     url,
		Login:   login,
		From:    from,
		Secret:  secret,
		client:  &http.Client{Timeout: 15 * time.Second},
	}
}

var (
	nonDigit = regexp.MustCompile(`\D`)
	counter  int64
)

// normalizePhone strips formatting and ensures a Tajik country code (992).
func normalizePhone(phone string) string {
	d := nonDigit.ReplaceAllString(phone, "")
	switch {
	case len(d) == 9: // local 9-digit number
		d = "992" + d
	case len(d) == 10 && d[0] == '0':
		d = "992" + d[1:]
	}
	return d
}

// Configured reports whether the gateway credentials are present.
func (s *Sender) Configured() bool {
	return s.Enabled && s.Login != "" && s.From != "" && s.Secret != ""
}

// Send delivers a single SMS. The signature is sha256(txn;login;sender;phone;secret).
func (s *Sender) Send(phone, message string) error {
	if !s.Enabled {
		return fmt.Errorf("SMS-уведомления отключены (RENTACAR_SMS_ENABLED=false)")
	}
	if s.Login == "" || s.From == "" || s.Secret == "" {
		return fmt.Errorf("SMS-шлюз не настроен (login/sender/secret)")
	}

	txnID := strconv.FormatInt(time.Now().Unix(), 10) + fmt.Sprintf("%04d", atomic.AddInt64(&counter, 1)%10000)
	p := normalizePhone(phone)
	raw := txnID + ";" + s.Login + ";" + s.From + ";" + p + ";" + s.Secret
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(raw)))

	payload, _ := json.Marshal(map[string]string{
		"from":         s.From,
		"phone_number": p,
		"msg":          message,
		"login":        s.Login,
		"txn_id":       txnID,
		"str_hash":     hash,
	})

	resp, err := s.client.Post(s.URL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("не удалось отправить SMS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("шлюз SMS вернул ошибку: %s %s", resp.Status, string(body))
	}
	return nil
}
