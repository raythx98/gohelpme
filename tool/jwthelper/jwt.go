package jwthelper

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
)

type IJwt interface {
	// NewAccessToken creates a new access token.
	NewAccessToken(subject string) (string, error)
	// NewRefreshToken creates a new refresh token.
	NewRefreshToken(subject string) (string, error)
	// Authenticate authenticates the request using the token type.
	//
	// It returns an error if the request is not authenticated.
	Authenticate(request *http.Request, tokenType TokenType) error
	// GetSubject extracts the subject from the JWT token.
	//
	// It returns the subject and an error if the subject cannot be extracted.
	GetSubject(request *http.Request) (string, error)
}

// New creates a new Jwt instance.
func New(config Config, configProvider ConfigProvider) *Jwt {
	return &Jwt{
		hmacSecret:           configProvider.GetHmacSecret(),
		issuer:               config.Issuer,
		audiences:            config.Audiences,
		accessTokenValidity:  config.AccessTokenValidity,
		refreshTokenValidity: config.RefreshTokenValidity,
		SigningMethod:        jwt.SigningMethodHS512,
	}
}

func (j *Jwt) NewAccessToken(subject string) (string, error) {
	return j.CreateToken(subject, AccessToken, j.accessTokenValidity)
}

func (j *Jwt) NewRefreshToken(subject string) (string, error) {
	return j.CreateToken(subject, RefreshToken, j.refreshTokenValidity)
}

func (j *Jwt) CreateToken(subject string, tokenType TokenType, tokenValidity time.Duration) (string, error) {
	claims := &CustomClaims{
		TokenType: string(tokenType),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   subject,
			Audience:  j.audiences,
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(tokenValidity)},
			NotBefore: &jwt.NumericDate{Time: time.Now()},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		},
	}

	token := jwt.NewWithClaims(j.SigningMethod, claims)
	return token.SignedString(j.hmacSecret)
}

func (j *Jwt) Authenticate(r *http.Request, tokenType TokenType) error {
	bearerToken, err := request.BearerExtractor{}.ExtractToken(r)
	if err != nil {
		return err
	}

	token, err := j.Parse(bearerToken)
	if err != nil {
		//s.BaseService.Logger.Debug("filterValidBearerAuthTokens: invalid token",
		//	zap.Error(err), zap.String("token", bearerAuthToken))
		return fmt.Errorf("cannot parse token: %v", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return fmt.Errorf("cannot parse claims")
	}

	if !token.Valid || claims.TokenType != string(tokenType) {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func (j *Jwt) GetSubject(r *http.Request) (string, error) {
	bearerToken, err := request.BearerExtractor{}.ExtractToken(r)
	if err != nil {
		return "", err
	}

	token, err := j.Parse(bearerToken)
	if err != nil {
		return "", fmt.Errorf("cannot parse token: %v", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return "", fmt.Errorf("cannot parse claims")
	}

	return claims.RegisteredClaims.GetSubject()
}

func (j *Jwt) Parse(bearerAuthToken string) (*jwt.Token, error) {
	var token *jwt.Token
	token, err := jwt.ParseWithClaims(
		bearerAuthToken,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return j.hmacSecret, nil
		},
		jwt.WithValidMethods([]string{j.SigningMethod.Alg()}),
		jwt.WithIssuer(j.issuer))
	return token, err
}
