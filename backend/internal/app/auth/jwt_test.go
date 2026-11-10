package auth

import (
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

func TestJWT(t *testing.T) {
	l, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	defer l.Close()
	go func() {
		for {
			c, e := l.Accept()
			if e != nil {
				return
			}
			go func() {
				defer c.Close()
				b := make([]byte, 2048)
				n, _ := c.Read(b)
				if strings.Contains(string(b[:n]), "EXISTS") {
					c.Write([]byte(":0\r\n"))
				} else {
					c.Write([]byte("+OK\r\n"))
				}
			}()
		}
	}()
	m := &Manager{Secret: []byte("a-long-enough-signature-secret-for-testing-012345"), RedisAddress: l.Addr().String()}
	token, e := m.Sign(42)
	if e != nil {
		t.Fatal(e)
	}
	id, claims, e := m.Verify(token)
	if e != nil || id != 42 || claims.ID == "" {
		t.Fatalf("verify %v %v %v", id, claims, e)
	}
	if _, _, e = m.Verify(token + "a"); e == nil {
		t.Fatal("signature tamper accepted")
	}
	if e = m.Revoke(claims); e != nil {
		t.Fatal(e)
	}
	_ = os.Setenv("JWT_SECRET", "12345678901234567890123456789012")
	_ = time.Second
}
