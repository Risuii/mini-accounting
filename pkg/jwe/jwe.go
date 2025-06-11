package jwe

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwe"
	"github.com/lestrrat-go/jwx/v2/jwk"

	Config "mini-accounting/config"
	Constants "mini-accounting/constants"
	Library "mini-accounting/library"
	CustomErrorPackage "mini-accounting/pkg/custom_error"
)

type JWE interface {
	JWEGenerateToken(claims jwt.MapClaims, secretKey string) (string, error)
	JWEValidateToken(token []byte, secretKey string) ([]byte, error)
}

type JWEImpl struct {
	config  Config.Config
	library Library.Library
}

func NewJWE(
	config Config.Config,
	library Library.Library,
) JWE {
	return &JWEImpl{
		config:  config,
		library: library,
	}
}

type CustomClaims struct {
	Credential string `json:"credential"`
	jwt.RegisteredClaims
}

func (j *JWEImpl) JWEGenerateToken(claims jwt.MapClaims, secretKey string) (string, error) {
	path := "JWTPacakge:JWEGenerateToken"

	// ADD EXPIRE TIME
	expirationTime := time.Now().Add(j.config.GetConfig().JWE.ExpiryDuration)
	claims[Constants.Expired] = expirationTime.Unix()

	// Decode the API key
	decodedKey, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		return Constants.NilString, CustomErrorPackage.New(Constants.ErrDecodeAPIKey, err, path, j.library)
	}

	// Create a JWK (JSON Web Key) from the decoded key
	key, err := jwk.FromRaw(decodedKey)
	if err != nil {
		return Constants.NilString, CustomErrorPackage.New(err, err, path, j.library)
	}

	payload, err := j.library.JsonMarshal(claims)
	if err != nil {
		return Constants.NilString, CustomErrorPackage.New(err, err, path, j.library)
	}

	// Encrypt the JWT
	encryptedJWT, err := jwe.Encrypt(payload, jwe.WithKey(jwa.DIRECT, key))
	if err != nil {
		return Constants.NilString, CustomErrorPackage.New(Constants.ErrFailedEncryptJWT, err, path, j.library)
	}

	encryptedData := string(encryptedJWT)

	return encryptedData, nil
}

func (j *JWEImpl) JWEValidateToken(token []byte, secretKey string) ([]byte, error) {
	path := "JWTPacakge:JWEValidateToken"

	// Decode the API key
	decodedKey, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		return nil, CustomErrorPackage.New(Constants.ErrDecodeAPIKey, err, path, j.library)
	}

	// Create a JWK (JSON Web Key) from the decoded key
	key, err := jwk.FromRaw(decodedKey)
	if err != nil {
		return nil, CustomErrorPackage.New(Constants.ErrFailedDecryptJWE, err, path, j.library)
	}

	// Decrypt the JWE token
	decryptedPayload, err := jwe.Decrypt([]byte(token), jwe.WithKey(jwa.DIRECT, key))
	if err != nil {
		return nil, CustomErrorPackage.New(Constants.ErrFailedDecryptJWE, err, path, j.library)
	}

	// Unmarshal the decrypted payload into jwt.MapClaims
	var claims jwt.MapClaims
	if err := json.Unmarshal(decryptedPayload, &claims); err != nil {
		return nil, CustomErrorPackage.New(Constants.ErrUnmarshalClaim, err, path, j.library)
	}

	// Validasi expired (cek jika ada klaim 'exp')
	if exp, ok := claims[Constants.Expired].(float64); ok {
		expirationTime := time.Unix(int64(exp), 0)
		if time.Now().After(expirationTime) {
			return nil, CustomErrorPackage.New(Constants.ErrInvalidJWE, Constants.ErrInvalidJWE, path, j.library)
		}
	}

	// Marshal kembali claims jika diperlukan
	claimsData, err := json.Marshal(claims)
	if err != nil {
		return nil, err
	}

	return claimsData, nil
}
