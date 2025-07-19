package utils

import (
	"errors"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
)

type Credential struct {
	logger *slog.Logger
	secret string
}

type CredFunctionOptions func(*Credential)

func NewCredential(secret string, opt ...CredFunctionOptions) *Credential {
	cred := &Credential{
		secret: secret,
		logger: slog.Default(),
	}
	for _, o := range opt {
		o(cred)
	}
	return cred
}

func (c *Credential) ParseToken(jwtToken string) (map[string]any, error) {
	// This is a placeholder for actual JWT parsing logic.
	// In a real application, you would use a JWT library to parse and validate the token.
	c.logger.Debug("Parsing token", "token", jwtToken)
	// Simulate a parsed token
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		c.logger.Error("Failed to parse token", "error", err)
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.logger.Debug("Token parsed successfully", "claims", claims)
		return claims, nil
	} else {
		c.logger.Error("Invalid token claims")
		return nil, errors.New("invalid token claims")
	}
}

func CredWithLogger(logger *slog.Logger) CredFunctionOptions {
	return func(c *Credential) {
		c.logger = logger
	}
}
