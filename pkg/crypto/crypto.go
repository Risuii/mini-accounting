package utils

import (
	"bytes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/pbkdf2"

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

type CustomCrypto interface {
	SetPassphrase(passphrase string)
	Encrypt(plain string) (string, error)
	Decrypt(encryptedB64 string) (string, error)
}
type CustomCryptoImpl struct {
	passphrase string
	library    Library.Library
}

func NewCustomCrypto(
	config Config.Config,
	library Library.Library,
) CustomCrypto {
	return &CustomCryptoImpl{
		passphrase: config.GetConfig().App.SecretKey,
		library:    library,
	}
}

func (c *CustomCryptoImpl) SetPassphrase(passphrase string) {
	// REPLACE "passphrase" THAT WE ALREADY INITIATE ON "NewCustomCrypto"
	c.passphrase = passphrase
}

func (c *CustomCryptoImpl) Encrypt(plain string) (string, error) {
	// MARSHALING THE WORD TO BE JSON BYTE
	contentByte, err := c.library.JsonMarshal(plain)
	if err != nil {
		return "", err
	}
	// GENERATE RANDOM "iv" IN 16 BYTES
	iv := make([]byte, 16)
	_, err = c.library.RandRead(iv)
	if err != nil {
		return "", err
	}
	// GENERATE RANDOM "salt" IN 64 BYTES
	salt := make([]byte, 64)
	_, err = c.library.RandRead(salt)
	if err != nil {
		return "", err
	}
	// GENERATE AN ENCRYPTION KEY BASED ON PBKDF2 ALGORITHM
	key := pbkdf2.Key([]byte(c.passphrase), salt, 2145, 32, sha512.New)
	// CREATE A CHIPER BLOCK THAT IS NEEDED BY CBC ENCRYPTER (ENCRYPTION  PROCESS)
	block, err := c.library.AESNewCipher(key)
	if err != nil {
		return "", err
	}
	// ADD PKCS7 PADDING INTO "contentByte" THAT WE GET FROM JSON MARSHALL
	result, err := c.PKCS7PAD(contentByte)
	if err != nil {
		return "", err
	}
	// INIT A ENCRYPTION RESULT VARIABLE
	encrypted := make([]byte, len(result))
	// SETUP THE ENCRYPTER (WE USE CBC ENCRYPTER)
	mode := cipher.NewCBCEncrypter(block, iv)
	// ENCRYPTION PROCESS
	mode.CryptBlocks(encrypted, result)
	// COMBINE SALT, IV, AND ENCRYPTION RESULT INTO A SLICE OF BYTE
	combine := append(salt[:], append(iv[:], encrypted...)...)
	// ENCODE THE SLICE OF BYTE INTO BASE 64 STRING
	encodedB64 := base64.StdEncoding.EncodeToString(combine)
	// RETURN SALT, IV, AND ENCRYPTION RESULT AS BASE 64 STRING RESULT
	return encodedB64, nil
}

func (c *CustomCryptoImpl) Decrypt(encryptedB64 string) (string, error) {
	// DECODE BASE 64 STRING FROM "Encrypt" INTO A SLICE OF BYTE THAT IS CONTAINING SALT, IV, AND ENCRYPTION RESULT
	decoded64, err := c.library.Base64DecodeString(encryptedB64)
	if err != nil {
		return "", err
	}
	// CHECK WHETHER THE LENGTH OF DECODED BASE 64 == 0 OR NOT
	if c.library.GetSlicesByteLen(decoded64) == 0 {
		return "", ErrNoContentDecryption
	}
	// CHECK WHETHER DECODED BASE 64 HAS SALT, IV, AND ENCRYPTION RESULT
	if c.library.GetSlicesByteLen(decoded64) <= 80 {
		return "", ErrInvalidDecryption
	}
	// WHEN THE LENGTH OF DECODED BASE 64 > 0, THEN EXTRACT IT
	salt := decoded64[:64]      // GET SALT
	iv := decoded64[64:80]      // GET IV
	encrypted := decoded64[80:] // GET ENCRYPTION RESULT
	// GENERATE AN ENCRYPTION KEY BASED ON PBKDF2 ALGORITHM, WHICH IS USING THE COMPONENT THAT WE USE IN ENCRYPTION
	key := pbkdf2.Key([]byte(c.passphrase), salt, 2145, 32, sha512.New)
	// CREATE A CHIPER BLOCK THAT IS NEEDED BY CBC ENCRYPTER (ENCRYPTION  PROCESS)
	block, err := c.library.AESNewCipher(key)
	if err != nil {
		return "", err
	}
	// INIT A DECRYPTION RESULT VARIABLE
	plain := make([]byte, len(encrypted))
	// SETUP THE ENCRYPTER (WE USE CBC DECRYPTER)
	mode := cipher.NewCBCDecrypter(block, iv)
	// DECRYPTION PROCESS
	mode.CryptBlocks(plain, encrypted)
	// REMOVE PKCS7 PADDING THAT WE ADD IN ENCRYPTION
	finalResult, err := c.PKCS7STRIP(plain)
	if err != nil {
		return "", err
	}
	// RETURN THE DECRYPTION RESULT
	return strings.ReplaceAll(fmt.Sprintf("%s", finalResult), `"`, ""), nil
}

// pkcs7strip remove pkcs7 padding
func (c *CustomCryptoImpl) PKCS7STRIP(data []byte) ([]byte, error) {
	// GET LENGTH OF DECODED BASE64 DATA
	length := c.library.GetSlicesByteLen(data)
	// WHEN NO DATA
	if length == 0 {
		return nil, ErrPkcs7stripDataEmpty
	}
	// WHEN LENGTH AND BLOCKSIZE IS NOT VALID (LENGTH % BLOCKSIZE MUST TO BE 0)
	if length%c.library.GetAES256CBCBlockSize() != 0 {
		return nil, ErrPkcs7stripDataNotBlockAligned
	}
	// GET BYTE FROM LAST INDEX TO GET THE PADDING BECAUSE WE LOOPED PADDING AS LONG AS "c.blockSize - len(data)%c.blockSize"
	padLen := c.library.ParseInt(data[length-1])
	if padLen > c.library.GetAES256CBCBlockSize() || padLen == 0 {
		return nil, ErrPkcs7stripInvalidPadding
	}
	// GENERATE SUFFIX (THE PADDING WE GENERATE IN "PKCS7PAD")
	ref := bytes.Repeat([]byte{byte(padLen)}, padLen)
	// CHECK THE SUFFIX
	if !c.library.HasSuffix(data, ref) {
		return nil, ErrPkcs7stripInvalidPadding
	}
	// GET THE REAL DATA WHERE WE REMOVE THE SUFFIX ("length" - "padLen")
	result := data[:length-padLen]
	// RETURN THE RESULT
	return result, nil
}

// pkcs7pad add pkcs7 padding
func (c *CustomCryptoImpl) PKCS7PAD(data []byte) ([]byte, error) {
	// MAKE SURE THAT THE "blockSize" LESS THAN 256 BECAUSE WE USE AES256
	if c.library.GetAES256CBCBlockSize() <= 1 || c.library.GetAES256CBCBlockSize() >= 256 {
		return nil, ErrPkcs7padInvalidBlockSize
	}
	// GENERATE THE LENGTH OF PADDING
	/*
		-- FOR EXAMPLE:
		c.blockSize => 16
		len(data) => 9
		len(data) % c.blockSize => 9
		padLen => 16 - 9 = 7
	*/
	padLen := c.library.GetAES256CBCBlockSize() - len(data)%c.library.GetAES256CBCBlockSize()
	// GENERATE PADDING AS LONG AS THE GENERATED LENGTH OF PADDING
	/*
		-- FOR EXAMPLE:
		padLen => 16 - 9 = 7
		padding => [7 7 7 7 7 7 7]
	*/
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	// APPEND PADDING
	/*
		-- FOR EXAMPLE:
		padding => [7 7 7 7 7 7 7]
		data => [34 87 101 108 99 111 109 101 34]
		result => [34 87 101 108 99 111 109 101 34 7 7 7 7 7 7 7]
	*/
	result := append(data, padding...)
	// RETURN THE RESULT
	return result, nil
}
