package authserver

import (
	"auth/pkg/authcode"
	"auth/pkg/authjwt"
	"auth/pkg/authmail"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Server struct {
	Address string
	Port    int
	Secret  []byte
	SESHost authmail.Host
	srv     *http.Server
}

func BadRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(msg))
}

func UnauthorizedRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(msg))
}

type AuthRequestBody struct {
	Email string `json:"forUser"`
	Code  string `json:"authCode"`
}

func (s *Server) HandleEmailAuthRequest(w http.ResponseWriter, req *http.Request) {
	var body []byte = make([]byte, 0)
	bodyReader := req.Body
	body, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		errMsg := fmt.Sprintf("Couldn't read request body: %v\n", err)
		BadRequest(w, errMsg)
		return
	}
	authReq := AuthRequestBody{}
	if err := json.Unmarshal(body, &authReq); err != nil {
		errMsg := fmt.Sprintf("Request body not properly formatted: %v\n", err)
		BadRequest(w, errMsg)
		return
	}
	fmt.Printf("Received request to authenticate %s\n", authReq.Email)
	fmt.Printf("Sending authentication email to %s\n", authReq.Email)
	code := authcode.NewAuthCode(authReq.Email)
	msg := authmail.NewAuthMessage(authReq.Email, code.Code)
	authmail.SendMessage(s.SESHost, authReq.Email, msg)
	fmt.Printf("Authentication email sent to %s\n", authReq.Email)
}

func (s *Server) HandleCodeAuthRequest(w http.ResponseWriter, req *http.Request) {
	var body []byte = make([]byte, 0)
	bodyReader := req.Body
	body, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		errMsg := fmt.Sprintf("Couldn't read request body: %v\n", err)
		BadRequest(w, errMsg)
		return
	}
	authReq := AuthRequestBody{}
	if err := json.Unmarshal(body, &authReq); err != nil {
		errMsg := fmt.Sprintf("Request body not properly formatted: %v\n", err)
		BadRequest(w, errMsg)
		return
	}
	fmt.Printf("Received request to authenticate %s with code %s\n", authReq.Email, authReq.Code)
	valid := authcode.ValidateAuthCode(authReq.Email, authReq.Code)
	if valid {
		fmt.Printf("%s authenticated using authentication code\n", authReq.Email)
		w.WriteHeader(http.StatusOK)
		jwt := authjwt.NewJWT(authReq.Email, map[string]interface{}{"authorized": true})
		token, _ := authjwt.Export(jwt, s.Secret)
		w.Write([]byte(token))
	} else {
		errMsg := fmt.Sprintf("Invalid code")
		UnauthorizedRequest(w, errMsg)
	}
}

func (s *Server) HandleTokenAuthRequest(w http.ResponseWriter, req *http.Request) {
	// Get bearerToken from request Authorization header
	var bearerToken string
	if token, ok := req.Header["Authorization"]; !ok {
		errMsg := fmt.Sprintf("Request has no authorization header")
		UnauthorizedRequest(w, errMsg)
		return
	} else {
		typevalpair := strings.Split(token[0], " ")
		authType, authVal := typevalpair[0], typevalpair[1]
		if authType != "Bearer" {
			errMsg := fmt.Sprintf("Request authorization must be Bearer <token>")
			UnauthorizedRequest(w, errMsg)
			return
		}
		bearerToken = authVal
	}
	// Unmarshal and validate bearer token
	token, valid, err := authjwt.Verify(bearerToken, s.Secret)
	if err != nil {
		errMsg := fmt.Sprintf("Couldn't process bearer token: %v\n", err)
		UnauthorizedRequest(w, errMsg)
		return
	}
	if !valid {
		errMsg := fmt.Sprintf("Bearer token has is altered or expired. Re-authentication is required.")
		UnauthorizedRequest(w, errMsg)
		return
	}
	// Check bearer token against request body forUser
	var body []byte = make([]byte, 0)
	bodyReader := req.Body
	body, err = ioutil.ReadAll(bodyReader)
	if err != nil {
		errMsg := fmt.Sprintf("Couldn't read request body: %v\n", err)
		BadRequest(w, errMsg)
		return
	}
	authReq := AuthRequestBody{}
	if err := json.Unmarshal(body, &authReq); err != nil {
		errMsg := fmt.Sprintf("Request body not properly formatted: %v\n", err)
		BadRequest(w, errMsg)
		return
	}
	if authReq.Email != token.Body.ForUser {
		errMsg := fmt.Sprintf("Bearer token has been altered.")
		UnauthorizedRequest(w, errMsg)
		return
	}

	fmt.Printf("%s authenticated using bearer token\n", authReq.Email)
	w.WriteHeader(http.StatusOK)
	outToken, _ := json.Marshal(token)
	w.Write(outToken)
}

func (s *Server) Start() {
	// Test handler
	http.HandleFunc(fmt.Sprintf("auth.%s/mail", s.Address), s.HandleEmailAuthRequest)
	http.HandleFunc(fmt.Sprintf("auth.%s/code", s.Address), s.HandleCodeAuthRequest)
	http.HandleFunc(fmt.Sprintf("auth.%s/token", s.Address), s.HandleTokenAuthRequest)
	// Generate address
	fulladdr := fmt.Sprintf("%s:%d", s.Address, s.Port)
	crt := fmt.Sprintf("./cert/%s.crt", s.Address)
	key := fmt.Sprintf("./cert/%s.key", s.Address)
	// Create https.Server
	s.srv = &http.Server{
		Addr:    fulladdr,
		Handler: nil,
	}
	go func() {
		s.srv.ListenAndServeTLS(crt, key)
	}()
}

func (s *Server) Stop() {
	s.srv.Shutdown(context.TODO())
}
