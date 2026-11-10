package storage

// Minimal S3-compatible PUT (AWS Signature V4) for MinIO uploads.
// Uses standard library only, and uploads images and videos into a private-write/public-read bucket.
import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

type Client struct {
	Endpoint, AccessKey, SecretKey, Bucket string
	HTTP                                   *http.Client
}

func New() *Client {
	env := func(k, v string) string {
		if s := os.Getenv(k); s != "" {
			return s
		}
		return v
	}
	return &Client{Endpoint: env("MINIO_ENDPOINT", "http://localhost:9000"), AccessKey: env("MINIO_ROOT_USER", "root"), SecretKey: env("MINIO_ROOT_PASSWORD", "rootpassword"), Bucket: env("MINIO_BUCKET", "hypertension-media"), HTTP: &http.Client{Timeout: 40 * time.Second}}
}
func digest(b []byte) string { s := sha256.Sum256(b); return hex.EncodeToString(s[:]) }
func sign(key []byte, message string) []byte {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte(message))
	return mac.Sum(nil)
}
func (s *Client) Upload(file multipart.File, original string, max int64) (string, error) {
	if file == nil {
		return "", errors.New("empty upload")
	}
	// LimitReader reads max+1, allowing size-limit enforcement even with a lying Content-Length.
	content, e := io.ReadAll(io.LimitReader(file, max+1))
	if e != nil {
		return "", e
	}
	if len(content) == 0 || int64(len(content)) > max {
		return "", fmt.Errorf("file must contain 1 to %d bytes", max)
	}
	mimeType := http.DetectContentType(content)
	if !strings.HasPrefix(mimeType, "image/") && !strings.HasPrefix(mimeType, "video/") && !strings.HasPrefix(mimeType, "application/octet-stream") {
		return "", errors.New("only image/video uploads are supported")
	}
	ext := strings.ToLower(path.Ext(original))
	if !map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".mp4": true, ".mov": true, ".webm": true, ".ogg": true}[ext] {
		return "", errors.New("unsupported file extension")
	}
	randBytes := make([]byte, 16)
	if _, e = rand.Read(randBytes); e != nil {
		return "", e
	}
	name := hex.EncodeToString(randBytes) + ext
	endpoint, e := url.Parse(strings.TrimRight(s.Endpoint, "/"))
	if e != nil {
		return "", e
	}
	uri := "/" + s.Bucket + "/" + name
	urlString := strings.TrimRight(s.Endpoint, "/") + uri
	req, e := http.NewRequest(http.MethodPut, urlString, bytes.NewReader(content))
	if e != nil {
		return "", e
	}
	now := time.Now().UTC()
	stamp := now.Format("20060102T150405Z")
	date := now.Format("20060102")
	sum := digest(content)
	req.Header.Set("x-amz-date", stamp)
	req.Header.Set("x-amz-content-sha256", sum)
	req.Header.Set("Content-Type", mimeType)
	host := endpoint.Host
	headers := "host:" + host + "\n" + "x-amz-content-sha256:" + sum + "\n" + "x-amz-date:" + stamp + "\n"
	canonical := "PUT\n" + uri + "\n\n" + headers + "\nhost;x-amz-content-sha256;x-amz-date\n" + sum
	scope := date + "/us-east-1/s3/aws4_request"
	stringToSign := "AWS4-HMAC-SHA256\n" + stamp + "\n" + scope + "\n" + digest([]byte(canonical))
	key := sign([]byte("AWS4"+s.SecretKey), date)
	key = sign(key, "us-east-1")
	key = sign(key, "s3")
	key = sign(key, "aws4_request")
	req.Header.Set("Authorization", fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=host;x-amz-content-sha256;x-amz-date, Signature=%s", s.AccessKey, scope, hex.EncodeToString(sign(key, stringToSign))))
	res, e := s.HTTP.Do(req)
	if e != nil {
		return "", e
	}
	defer res.Body.Close()
	if res.StatusCode/100 != 2 {
		body, _ := io.ReadAll(io.LimitReader(res.Body, 2048))
		return "", fmt.Errorf("MinIO returned %s: %s", res.Status, strings.TrimSpace(string(body)))
	}
	return name, nil
}
