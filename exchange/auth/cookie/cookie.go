package cookie

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	CookieName = "exchange"
)

var key []byte = []byte(`HonestOption1234`)

func GenerateCookie(username string) (string, error) {
	buf := make([]byte, 0, 32)
	buf = append(buf, username...)
	buf = append(buf, ':')
	buf = strconv.AppendInt(buf, time.Now().Unix(), 10)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to generate cookie[1]: %s", err)
	}
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate cookie[2]: %s", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to generate cookie[3]: %s", err)
	}
	ciphertext := aesgcm.Seal(nil, nonce, buf, nil)
	ciphertext = append(ciphertext, nonce...)
	return hex.EncodeToString(ciphertext), nil
}

func CheckCookie(c *http.Cookie) (string, error) {
	val, err := hex.DecodeString(c.Value)
	if err != nil {
		return "", fmt.Errorf("failed to check cookie: failed to decode cookie from hex: %s", err)
	}
	user, err := decryptCookie(val)
	return user, err
}

func decryptCookie(val []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to check cookie: %s", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to check cookie: failed to create aesgcm: %s", err)
	}
	if len(val) < 12 {
		return "", http.ErrNoCookie
	}
	nonce := val[len(val)-12:]
	val = val[:len(val)-12]

	plain, err := aesgcm.Open(nil, nonce, val, nil)
	if err != nil {
		log.Printf("got error %q while unsealing the cookie", err)
		return "", http.ErrNoCookie
	}

	user := strings.SplitN(string(plain), ":", 2)
	if len(user) != 2 {
		return "", http.ErrNoCookie
	}
	return user[0], nil
}
