package services

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/crypto/argon2"
)

type param struct {
	memory, iterations, saltLength, keyLength uint32
	parallelism                               uint8
}

type Username string

func generateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return salt, err
	}
	return salt, nil
}

func decodeHash(encodedHash string) (salt, hash []byte, p *param, err error) {
	s := strings.Split(encodedHash, "$")
	if len(s) != 6 {
		return nil, nil, nil, errors.New("not an encoded hash")
	}
	var version int
	_, err = fmt.Sscanf(s[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errors.New("versions don't match")
	}
	p = &param{}
	_, err = fmt.Sscanf(s[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}
	salt, err = base64.RawStdEncoding.Strict().DecodeString(s[4])
	if err != nil {
		return nil, nil, nil, err
	}
	hash, err = base64.RawStdEncoding.Strict().DecodeString(s[5])
	if err != nil {
		return nil, nil, nil, err
	}
	if len(hash) != 32 {
		return nil, nil, nil, errors.New("length doesn't match")
	}
	p.keyLength = uint32(len(hash))
	return salt, hash, p, nil
}

func GenerateHash(password string) (string, error) {
	salt, err := generateSalt()
	if err != nil {
		return "", err
	}
	p := param{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}
	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)
	b64encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	b64encodedHash := base64.RawStdEncoding.EncodeToString(hash)
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.memory, p.iterations, p.parallelism, b64encodedSalt, b64encodedHash)
	return encodedHash, nil
}

func GenerateJwtToken(username string) (string, error) {
	privateKey := os.Getenv("PRIVATE_KEY")
	key, err := jwk.FromRaw([]byte(privateKey))
	if err != nil {
		return "", err
	}
	builder := jwt.
		NewBuilder().
		Issuer(os.Getenv("ISSUER")).
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(7 * 24 * time.Hour)).
		Audience([]string{os.Getenv("AUDIENCE")}).
		Subject(username)
	token, err := builder.Build()
	if err != nil {
		return "", err
	}
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, key))
	if err != nil {
		return "", err
	}
	return string(signed), nil
}

func VerifyToken(token string) (jwt.Token, error) {
	privateKey := os.Getenv("PRIVATE_KEY")
	key, err := jwk.FromRaw([]byte(privateKey))
	if err != nil {
		return nil, err
	}
	t, err := jwt.Parse(
		[]byte(token),
		jwt.WithKey(jwa.HS256, key),
		jwt.WithValidate(true),
		jwt.WithAudience(os.Getenv("AUDIENCE")),
		jwt.WithIssuer(os.Getenv("ISSUER")),
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func VerifyPassword(password, encodedHash string) (bool, error) {
	salt, hash, p, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}
	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)
	if subtle.ConstantTimeCompare(hash, otherHash) == 0 {
		return false, nil
	}
	return true, nil
}

func IsAuthenticated(r *http.Request) (string, error) {
	var n Username = "username"
	name := r.Context().Value(n)
	if name == nil {
		return "", errors.New("not authenticated")
	}
	return name.(string), nil
}
