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
)

// A Server holds the information needed to fulfill authentication.
type Server struct {
	// Server port/address
	Address string
	Port    int
	// JWT signing secret (TODO: move this somewhere more sensible)
	Secret []byte
	// Host; see doc/mail
	SESHost authmail.Host
	// Server; not exported, used internally for controlling HTTPS server
	srv *http.Server
}

// Write out an HTTP response with status code 400 (bad request)
// 400 denotes malformed requests, i.e. JSON improperly formatted, text encoding wrong
// 400 should *not* be used for failed authentication.
func BadRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(msg))
}

// Write out an HTTP response with status code 401 (unauthorized)
// 401 denotes lack of, or failed, authentication for a given resource.
// 401 should *not* be used for users that are authenticated but lacking permissions; see 403
func UnauthorizedRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(msg))
}

// Write out a response with code and msg
// code should be an http library constant. see doc/server for code usage
func WriteResponse(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

// Request body format for all authentication requests.
type AuthRequestBody struct {
	Email string `json:"forUser"`
	Code  string `json:"authCode"`
	Token string `json:"authToken"`
}

// Read request body.
// An error here usually indicates a malformed request and should return a 400.
func ReadRequestBody(out *AuthRequestBody, req *http.Request) error {
	bodyReader := req.Body
	body, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, out); err != nil {
		return err
	}
	return nil
}

// Handle email authentication requests
func (s *Server) HandleEmailAuthRequest(w http.ResponseWriter, req *http.Request) {
	authReq := AuthRequestBody{}
	err := ReadRequestBody(&authReq, req)
	if err != nil {
		errMsg := fmt.Sprintf("Couldn't read request body: %v\n", err)
		WriteResponse(w, http.StatusBadRequest, errMsg)
		return
	}
	// Also throw a 400 if there's no included email
	if authReq.Email == "" {
		errMsg := fmt.Sprintf("forUser is needed for endpoint /mail")
		WriteResponse(w, http.StatusBadRequest, errMsg)
	}
	// Send an authentication email and write out 200
	fmt.Printf("Received request to authenticate %s\n", authReq.Email)
	fmt.Printf("Sending authentication email to %s\n", authReq.Email)
	code := authcode.NewAuthCode(authReq.Email)
	msg := authmail.NewAuthMessage(authReq.Email, code.Code)
	authmail.SendMessage(s.SESHost, authReq.Email, msg)
	succMsg := fmt.Sprintf("Authentication email sent to %s\n", authReq.Email)
	WriteResponse(w, http.StatusOK, succMsg)
}

// Handle authentication code requests
func (s *Server) HandleCodeAuthRequest(w http.ResponseWriter, req *http.Request) {
	// Read in body. Send a 400 on failure
	authReq := AuthRequestBody{}
	err := ReadRequestBody(&authReq, req)
	if err != nil {
		errMsg := fmt.Sprintf("Couldn't read request body: %v\n", err)
		WriteResponse(w, http.StatusBadRequest, errMsg)
		return
	}
	fmt.Printf("Received request to authenticate %s with code %s\n", authReq.Email, authReq.Code)
	if authReq.Email == "" || authReq.Code == "" {
		errMsg := fmt.Sprintf("forUser and authCode are needed for endpoint /code")
		WriteResponse(w, http.StatusBadRequest, errMsg)
	}
	valid := authcode.ValidateAuthCode(authReq.Email, authReq.Code)
	if valid {
		fmt.Printf("%s authenticated using authentication code\n", authReq.Email)
		jwt := authjwt.NewJWT(authReq.Email, map[string]interface{}{"authorized": true})
		token, _ := authjwt.Export(jwt, s.Secret)
		WriteResponse(w, http.StatusOK, token)
	} else {
		errMsg := fmt.Sprintf("Invalid code")
		WriteResponse(w, http.StatusUnauthorized, errMsg)
	}
}

func (s *Server) HandleTokenAuthRequest(w http.ResponseWriter, req *http.Request) {
	// Read in body. Send a 400 on failure
	authReq := AuthRequestBody{}
	err := ReadRequestBody(&authReq, req)
	if err != nil {
		errMsg := fmt.Sprintf("Couldn't read request body: %v\n", err)
		WriteResponse(w, http.StatusBadRequest, errMsg)
		return
	}
	if authReq.Token == "" {
		errMsg := fmt.Sprintf("authToken is needed for endpoint /token")
		WriteResponse(w, http.StatusBadRequest, errMsg)
	}
	// Verify the authToken included with the request
	token, valid, err := authjwt.Verify(authReq.Token, s.Secret)
	if err != nil {
		errMsg := fmt.Sprintf("Couldn't process bearer token: %v\n", err)
		WriteResponse(w, http.StatusUnauthorized, errMsg)
		return
	}
	if !valid {
		errMsg := fmt.Sprintf("Bearer token has is altered or expired. Re-authentication is required.")
		WriteResponse(w, http.StatusUnauthorized, errMsg)
		return
	}

	outToken, _ := json.Marshal(token)
	WriteResponse(w, http.StatusOK, string(outToken))
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
