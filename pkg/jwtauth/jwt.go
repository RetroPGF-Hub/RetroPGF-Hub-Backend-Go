package jwtauth

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/config"
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
)

type (
	AuthFactory interface {
		SignToken() (string, error)
	}

	Claims struct {
		UserId   string `json:"userId"`
		Email    string `json:"email"`
		Username string `json:"username"`
		Source   string `json:"source"`
		Profile  string `json:"profile"`
	}

	AuthMapClaims struct {
		*Claims
		jwt.RegisteredClaims
		isErr error
	}

	authConcrete struct {
		Secret     []byte
		PrivateKey *rsa.PrivateKey
		Claims     *AuthMapClaims `json:"claims"`
	}
	apiKey struct{ *authConcrete }
)

func NewAccessToken(cfg *config.Jwt, claims *Claims, expiredAt int64, subject string) *authConcrete {

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.PrivateKeyPem))
	// privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(prvKey))
	if err != nil {
		return &authConcrete{
			PrivateKey: nil,
			Claims: &AuthMapClaims{
				Claims:           nil,
				RegisteredClaims: jwt.RegisteredClaims{},
				isErr:            fmt.Errorf("reading private key errors : %s", err.Error()),
			},
		}
	}
	return &authConcrete{
		PrivateKey: privateKey,
		Claims: &AuthMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "mix.com",
				Subject:   subject,
				Audience:  []string{"mix.com"},
				ExpiresAt: jwtTimeDurationCal(expiredAt),
				NotBefore: jwt.NewNumericDate(now()),
				IssuedAt:  jwt.NewNumericDate(now()),
			},
			isErr: nil,
		},
	}
}

func (a *authConcrete) SignToken() (string, error) {
	if a.Claims.isErr != nil {
		return "", a.Claims.isErr
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, a.Claims)
	token, err := jwtToken.SignedString(a.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("signing token error: %s", err.Error())
	}
	return token, nil
}

func ParseToken(tokenString string, cfg *config.Jwt) (*AuthMapClaims, error) {
	// publicKeyPem is a string

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cfg.PublicKeyPem))
	if err != nil {
		return nil, fmt.Errorf("reading public key errors: %s", err.Error())
	}

	token, err := jwt.ParseWithClaims(tokenString, &AuthMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("parsing token error: %s", err.Error())
	}

	if claims, ok := token.Claims.(*AuthMapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func now() time.Time {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return time.Now().In(loc)
}

func jwtTimeDurationCal(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(now().Add(time.Duration(t * int64(math.Pow10(9)))))
}

func NewApiKey(privateKeyPem string, secret string) AuthFactory {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPem))
	// privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(prvKey))
	if err != nil {
		return &apiKey{
			authConcrete: &authConcrete{
				Secret:     nil,
				PrivateKey: nil,
				Claims: &AuthMapClaims{
					Claims:           &Claims{},
					RegisteredClaims: jwt.RegisteredClaims{},
				},
			},
		}
	}

	return &apiKey{
		authConcrete: &authConcrete{
			Secret:     []byte(secret),
			PrivateKey: privateKey,
			Claims: &AuthMapClaims{
				Claims: &Claims{},
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "mix.com",
					Subject:   "api-key",
					Audience:  []string{"mix.com"},
					ExpiresAt: jwtTimeDurationCal(31560000),
					NotBefore: jwt.NewNumericDate(now()),
					IssuedAt:  jwt.NewNumericDate(now()),
				},
			},
		},
	}
}

func ParseTokenGrpc(secret string, tokenString string) (*AuthMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error: unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errors.New("error: token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("error: token is expired")
		} else {
			return nil, errors.New("error: token is invalid")
		}
	}

	if claims, ok := token.Claims.(*AuthMapClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New("error: claims type is invalid")
	}
}

// Apikey  generator
var apiKeyInstant string

// Work only once
var once sync.Once

func SetApiKey(cfg *config.Jwt) {
	once.Do(func() {
		apiKeyInstant, _ = NewApiKey(cfg.PrivateKeyPem, cfg.ApiSecretKey).SignToken()
		// apiKeyInstant := NewApiKey(secret)
		// _ = apiKeyInstant
		// log.Printf("%+v", apiKeyInstant)
	})
}

func SetApiKeyInContext(pctx *context.Context) {
	*pctx = metadata.NewOutgoingContext(*pctx, metadata.Pairs("auth", apiKeyInstant))
}
