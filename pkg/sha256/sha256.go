package sha256

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"

	Config "mini-accounting/config"
	Library "mini-accounting/library"
)

var (
	ErrPkcs7stripDataEmpty           = errors.New("pkcs7: Data is empty")
	ErrPkcs7stripDataNotBlockAligned = errors.New("pkcs7: Data is not block-aligned")
	ErrPkcs7stripInvalidPadding      = errors.New("pkcs7: Invalid padding")
	ErrPkcs7padInvalidBlockSize      = errors.New("pkcs7: Invalid block size")
	ErrNoContentDecryption           = errors.New("decrypt: No Content")
	ErrInvalidDecryption             = errors.New("decrypt: Invalid Content")
)

type CustomSha256 interface {
	Hash(data []byte) []byte
	HashPassword(plain []byte) ([]byte, error)
	ComparePassword(plain []byte, hashedPassword []byte) bool
}
type CustomSha256Impl struct {
	library Library.Library
}

func NewCustomSha256(
	config Config.Config,
	library Library.Library,
) CustomSha256 {
	return &CustomSha256Impl{
		library: library,
	}
}

func (c *CustomSha256Impl) Hash(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}

func (c *CustomSha256Impl) HashPassword(plain []byte) ([]byte, error) {
	salt := make([]byte, 128)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	plainHash := c.Hash(plain)
	saltHash := c.Hash(salt)
	saltLen := len(saltHash)

	// sha256(pass + salt)
	combined := append(plainHash, saltHash...)
	passHash := c.Hash(combined)

	// 1/2salt[0] + passHash + 1/2salt[1]
	result := append(saltHash[:saltLen/2], passHash...)
	result = append(result, saltHash[saltLen/2:]...)

	return result, nil
}

func (c *CustomSha256Impl) ComparePassword(plain []byte, hashedPassword []byte) bool {
	hashLen := len(hashedPassword)
	saltLen := 32

	// make a copy of hashedPassword to prevent error
	h1 := make([]byte, saltLen/2)
	h2 := make([]byte, saltLen/2)
	copy(h1, hashedPassword[:saltLen/2])
	copy(h2, hashedPassword[hashLen-(saltLen/2):])

	saltHash := append(h1, h2...)
	plainHash := c.Hash(plain)

	// sha256(pass + salt)
	combined := append(plainHash, saltHash...)
	passHash := c.Hash(combined)

	// 1/2salt[0] + passHash + 1/2salt[1]
	result := append(saltHash[:saltLen/2], passHash...)
	result = append(result, saltHash[saltLen/2:]...)

	return c.library.BytesEqual(hashedPassword, result)
}
