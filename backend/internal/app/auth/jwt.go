package auth

import (
	"bufio"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Manager signs HS256 JWT and persists revoked JTI values in Redis.
type Manager struct {
	Secret       []byte
	RedisAddress string
}
type Claims struct {
	Subject   string `json:"sub"`
	ID        string `json:"jti"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

func New() (*Manager, error) {
	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 32 {
		return nil, errors.New("JWT_SECRET must contain at least 32 characters")
	}
	redis := os.Getenv("REDIS_ADDR")
	if redis == "" {
		redis = "localhost:6380"
	}
	return &Manager{Secret: []byte(secret), RedisAddress: redis}, nil
}
func (m *Manager) Sign(userID uint) (string, error) {
	random := make([]byte, 16)
	if _, err := rand.Read(random); err != nil {
		return "", err
	}
	now := time.Now().Unix()
	claims := Claims{Subject: strconv.FormatUint(uint64(userID), 10), ID: hex.EncodeToString(random), IssuedAt: now, ExpiresAt: now + 24*60*60}
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payloadRaw, _ := json.Marshal(claims)
	payload := base64.RawURLEncoding.EncodeToString(payloadRaw)
	signing := header + "." + payload
	return signing + "." + base64.RawURLEncoding.EncodeToString(m.mac(signing)), nil
}
func (m *Manager) mac(message string) []byte {
	h := hmac.New(sha256.New, m.Secret)
	_, _ = h.Write([]byte(message))
	return h.Sum(nil)
}
func (m *Manager) Verify(token string) (uint, Claims, error) {
	var claim Claims
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return 0, claim, errors.New("invalid JWT format")
	}
	headerRaw, e := base64.RawURLEncoding.DecodeString(parts[0])
	if e != nil {
		return 0, claim, e
	}
	var header struct {
		Alg  string `json:"alg"`
		Type string `json:"typ"`
	}
	if json.Unmarshal(headerRaw, &header) != nil || header.Alg != "HS256" || header.Type != "JWT" {
		return 0, claim, errors.New("invalid JWT algorithm")
	}
	signature, e := base64.RawURLEncoding.DecodeString(parts[2])
	if e != nil || !hmac.Equal(signature, m.mac(parts[0]+"."+parts[1])) {
		return 0, claim, errors.New("invalid signature")
	}
	bytes, e := base64.RawURLEncoding.DecodeString(parts[1])
	if e != nil || json.Unmarshal(bytes, &claim) != nil {
		return 0, claim, errors.New("invalid payload")
	}
	if claim.ID == "" || claim.ExpiresAt <= time.Now().Unix() || claim.IssuedAt > time.Now().Unix()+60 {
		return 0, claim, errors.New("expired or invalid token")
	}
	uid, e := strconv.ParseUint(claim.Subject, 10, 64)
	if e != nil || uid == 0 {
		return 0, claim, errors.New("invalid subject")
	}
	revoked, e := m.redis("EXISTS", "jwt:blacklist:"+claim.ID)
	if e != nil {
		return 0, claim, fmt.Errorf("redis unavailable: %w", e)
	}
	if revoked != "0" {
		return 0, claim, errors.New("token revoked")
	}
	return uint(uid), claim, nil
}
func (m *Manager) Revoke(claim Claims) error {
	ttl := claim.ExpiresAt - time.Now().Unix()
	if ttl < 1 {
		ttl = 1
	}
	_, err := m.redis("SET", "jwt:blacklist:"+claim.ID, "1", "EX", strconv.FormatInt(ttl, 10))
	return err
}

// RESP client for the two commands needed in this lab; uses no additional Go modules.
func (m *Manager) redis(args ...string) (string, error) {
	c, e := net.DialTimeout("tcp", m.RedisAddress, 3*time.Second)
	if e != nil {
		return "", e
	}
	defer c.Close()
	_ = c.SetDeadline(time.Now().Add(5 * time.Second))
	msg := fmt.Sprintf("*%d\r\n", len(args))
	for _, a := range args {
		msg += fmt.Sprintf("$%d\r\n%s\r\n", len(a), a)
	}
	if _, e = c.Write([]byte(msg)); e != nil {
		return "", e
	}
	r := bufio.NewReader(c)
	prefix, e := r.ReadByte()
	if e != nil {
		return "", e
	}
	line, e := r.ReadString('\n')
	if e != nil {
		return "", e
	}
	value := strings.TrimSpace(line)
	if prefix == '-' {
		return "", errors.New(value)
	}
	if prefix != ':' && prefix != '+' {
		return "", errors.New("unexpected Redis reply")
	}
	return value, nil
}
