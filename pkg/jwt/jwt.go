package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type (
	AuthFactory interface {
		SignToken() (string, error)
	}

	Claims struct {
		UserId   string `json:"userId"`
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	AuthMapClaims struct {
		*Claims
		jwt.RegisteredClaims
	}

	authConcrete struct {
		PrivateKey *rsa.PrivateKey
		PublicKey  *rsa.PublicKey
		Secret     []byte
		Claims     *AuthMapClaims `json:"claims"`
	}

	accessToken  struct{ *authConcrete }
	refreshToken struct{ *authConcrete }
	apiKey       struct{ *authConcrete }
)

func (a *authConcrete) SignToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, a.Claims)
	ss, err := token.SignedString(a.Secret)
	if err != nil {
		return "", fmt.Errorf("Singtoken error: %s", err.Error())
	}
	return ss, nil
}

func now() time.Time {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return time.Now().In(loc)
}

// Note that: t is a second unit
// base on current time
func jwtTimeDurationCal(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(now().Add(time.Duration(t * int64(math.Pow10(9)))))
}

// base on previous time in cookie
func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func NewAccessToken(expiredAt int64, claims *Claims) AuthFactory {

	// example key na if you want to hack it
	privatePemString := `-----BEGIN PRIVATE KEY-----
	hehe
	/Qby4r8dKcOsOjD0U5t2tFyD77iuO6JSvxTvUCBw5aIGeikAaogm9F3hnXOJ+k+J
	gEd4l7ED/qDnkvYu7E1qjwwd3KbyEy4cSJIx0csEKcYSTSyIDVUP3/0o40KD8y3i
	n1unMlXyGaDjpp+REUr3I+fE
	-----END PRIVATE KEY-----`

	block, _ := pem.Decode([]byte(privatePemString))
	parseResult, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
	privateKey := parseResult.(*rsa.PrivateKey)

	publicPemString := `-----BEGIN PUBLIC KEY-----
	hehe
	fF1bvjCz0jnPHRZtPDdQaX5MuUsX4/t263csggxJoFVkrMlAGHafra1DVZwaXK5w
	xT35XhbWzZdVm26NqDe+JNk2pG2CqVDxNf2dCmjCaGGQORfr5mgL07RDoYnexNFd
	TwIDAQAB
	-----END PUBLIC KEY-----`

	blockPub, _ := pem.Decode([]byte(publicPemString))
	parseResultPub, _ := x509.ParsePKCS8PrivateKey(blockPub.Bytes)
	publicKey := parseResultPub.(*rsa.PublicKey)

	return &accessToken{
		authConcrete: &authConcrete{
			PrivateKey: privateKey,
			PublicKey:  publicKey,
			Claims: &AuthMapClaims{
				Claims: claims,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "mix.com",
					Subject:   "access-token",
					Audience:  []string{"mix.com"},
					ExpiresAt: jwtTimeDurationCal(expiredAt),
					NotBefore: jwt.NewNumericDate(now()),
					IssuedAt:  jwt.NewNumericDate(now()),
				},
			},
		},
	}
}

func NewRefreshToken(secret string, expiredAt int64, claims *Claims) AuthFactory {
	return &refreshToken{
		authConcrete: &authConcrete{
			Secret: []byte(secret),
			Claims: &AuthMapClaims{
				Claims: claims,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "mix.com",
					Subject:   "refresh-token",
					Audience:  []string{"mix.com"},
					ExpiresAt: jwtTimeDurationCal(expiredAt),
					NotBefore: jwt.NewNumericDate(now()),
					IssuedAt:  jwt.NewNumericDate(now()),
				},
			},
		},
	}
}

func ReloadToken(secret string, expiredAt int64, claims *Claims) (string, error) {
	obj := &refreshToken{
		authConcrete: &authConcrete{
			Secret: []byte(secret),
			Claims: &AuthMapClaims{
				Claims: claims,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "mix.com",
					Subject:   "refresh-token",
					Audience:  []string{"mix.com"},
					ExpiresAt: jwtTimeRepeatAdapter(expiredAt),
					NotBefore: jwt.NewNumericDate(now()),
					IssuedAt:  jwt.NewNumericDate(now()),
				},
			},
		},
	}

	return obj.SignToken()
}

func NewApiKey(secret string) AuthFactory {
	return &apiKey{
		authConcrete: &authConcrete{
			Secret: []byte(secret),
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

func ParseToken(secret string, tokenString string) (*AuthMapClaims, error) {
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

// // Apikey  generator
// var apiKeyInstant string

// // Work only once
// var once sync.Once

// func SetApiKey(secret string) {
// 	once.Do(func() {
// 		apiKeyInstant, _ = NewApiKey(secret).SignToken()
// 	})
// }

// func SetApiKeyInContext(pctx *context.Context) {
// 	*pctx = metadata.NewOutgoingContext(*pctx, metadata.Pairs("auth", apiKeyInstant))
// }
