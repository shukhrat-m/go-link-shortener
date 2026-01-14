package jwt

import "github.com/golang-jwt/jwt/v5"

type JWTData struct {
	Email string
}

type JWT struct {
	SecretKey string
}

func NewJwt(secret string) *JWT {
	return &JWT{SecretKey: secret}
}

func (j *JWT) GenerateToken(data JWTData) (string, error) {
	s, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": data.Email,
	}).SignedString([]byte(j.SecretKey))

	if err != nil {
		return "", err
	}

	return s, nil
}

func (j *JWT) Parse(token string) (bool, *JWTData) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return false, nil
	}

	email := t.Claims.(jwt.MapClaims)["email"]

	return t.Valid, &JWTData{Email: email.(string)}
}
