package auth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JWTClaims struct {
	UserID primitive.ObjectID `json:"user_id"`
	Role   string             `json:"role"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func NewJWTManager() (*JWTManager, error) {
	// Read private key
	privateKeyBytes, err := os.ReadFile("keys/private.pem")
	if err != nil {
		return nil, fmt.Errorf("could not read private key: %v", err)
	}

	// Parse private key
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing private key")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	// Read public key
	publicKeyBytes, err := os.ReadFile("keys/public.pem")
	if err != nil {
		return nil, fmt.Errorf("could not read public key: %v", err)
	}

	// Parse public key
	block, _ = pem.Decode(publicKeyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	ecdsaPublicKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not an ECDSA key")
	}

	return &JWTManager{
		privateKey: privateKey,
		publicKey:  ecdsaPublicKey,
	}, nil
}

func (m *JWTManager) GenerateToken(userID primitive.ObjectID, role string) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(m.privateKey)
}

func (m *JWTManager) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}
