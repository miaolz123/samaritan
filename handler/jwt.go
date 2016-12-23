package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	tokenKey = []byte("XXXXXXXXXXXXXXXX")
)

func makeToken(sub string) (token string) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		Subject:   sub,
	})
	token, _ = t.SignedString(tokenKey)
	return
}

func parseToken(token string) (sub string) {
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		return
	}
	t, _ := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return tokenKey, nil
	})
	if t != nil {
		if claims, ok := t.Claims.(*jwt.StandardClaims); ok && t.Valid {
			return claims.Subject
		}
	}
	return
}
