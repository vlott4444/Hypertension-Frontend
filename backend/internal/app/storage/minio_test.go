package storage

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type rfile struct{ *bytes.Reader }

func (r rfile) Close() error { return nil }
func TestMinioUploadSignature(t *testing.T) {
	called := false
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if !strings.HasPrefix(r.Header.Get("Authorization"), "AWS4-HMAC-SHA256 Credential=") {
			t.Error("missing S3 signature")
		}
		if r.Method != "PUT" || !strings.HasPrefix(r.URL.Path, "/hypertension-media/") {
			t.Errorf("wrong request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(200)
	}))
	defer s.Close()
	c := &Client{Endpoint: s.URL, AccessKey: "root", SecretKey: "rootpassword", Bucket: "hypertension-media", HTTP: s.Client()}
	data := []byte("\x89PNG\r\n\x1a\n" + strings.Repeat("x", 40))
	f := rfile{bytes.NewReader(data)}
	name, e := c.Upload(f, "stage.png", 1024)
	if e != nil {
		t.Fatal(e)
	}
	if !called || !strings.HasSuffix(name, ".png") {
		t.Fatalf("upload returned %q", name)
	}
	_, e = c.Upload(rfile{bytes.NewReader(data)}, "bad.exe", 1024)
	if e == nil {
		t.Fatal("accepted forbidden extension")
	}
	_ = multipart.ErrMessageTooLarge
}
