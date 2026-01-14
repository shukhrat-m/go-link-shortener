package auth

import (
	"go/adv-demo/configs"
	"go/adv-demo/pkg/jwt"
	"go/adv-demo/pkg/request"
	"go/adv-demo/pkg/response"
	"net/http"
)

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}

type AuthHandler struct {
	*configs.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{Config: deps.Config, AuthService: deps.AuthService}

	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[LoginRequest](&w, r)
		if err != nil {
			return
		}

		email, err := handler.AuthService.AuthenticateUser(body.Email, body.Password)

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.NewJwt(handler.Auth.Secret).GenerateToken(jwt.JWTData{Email: email})

		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		loginResponse := LoginResponse{Token: token}
		response.JSON(w, 200, loginResponse)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[RegisterRequest](&w, r)

		if err != nil {
			return
		}

		email, err := handler.AuthService.RegisterUser(body.Email, body.Password, body.Name)

		if err != nil {
			http.Error(w, "Failed to register user: "+err.Error(), http.StatusBadRequest)
			return
		}

		token, err := jwt.NewJwt(handler.Auth.Secret).GenerateToken(jwt.JWTData{Email: email})

		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		registerResponse := RegisterResponse{Token: token}
		response.JSON(w, 200, registerResponse)
	}
}
