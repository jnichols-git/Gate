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
	"os"
	"sync"
)

/* Quick reminder of the contents of AuthServerConfig:
Config {
	Domain string
	Port int
	SMTPHost {
		Username string
		Password string
		Host string
		Port int
		Sender string
	}
	JWS {
		TokenSecret string
	}
}
*/

// An AServer holds the information needed to fulfill authentication.
type AuthServer struct {
	Config *AuthServerConfig
	// Server; not exported, used internally for controlling HTTPS server
	srv *http.Server
	// Waitgroup; needed to maintain concurrency with server
	wg sync.WaitGroup
}

// Create a new server using config.
func NewServer(cfg *AuthServerConfig) *AuthServer {
	srv := &AuthServer{
		Config: cfg,
	}
	srv.Config.SMTPHost.Username = os.Getenv(srv.Config.SMTPHost.username_ENV)
	srv.Config.SMTPHost.Password = os.Getenv(srv.Config.SMTPHost.password_ENV)
	srv.Config.JWS.TokenSecret = os.Getenv(srv.Config.JWS.tokenSecret_ENV)
	return &AuthServer{
		Config: cfg,
	}
}

func (s *AuthServer) SMTPHost() authmail.Host {
	return authmail.Host{
		Username: s.Config.SMTPHost.Username,
		Password: s.Config.SMTPHost.Password,
		Host:     s.Config.SMTPHost.Host,
		Port:     s.Config.SMTPHost.Port,
		Sender:   s.Config.SMTPHost.Sender,
	}
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
func (s *AuthServer) HandleEmailAuthRequest(w http.ResponseWriter, req *http.Request) {
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
	// Send an authentication email and write out 200=
	code := authcode.NewAuthCode(authReq.Email)
	msg := authmail.NewAuthMessage(authReq.Email, code.Code)
	authmail.SendMessage(s.SMTPHost(), authReq.Email, msg)
	succMsg := fmt.Sprintf("Authentication email sent to %s\n", authReq.Email)
	WriteResponse(w, http.StatusOK, succMsg)
}

// Handle authentication code requests
func (s *AuthServer) HandleCodeAuthRequest(w http.ResponseWriter, req *http.Request) {
	// Read in body. Send a 400 on failure
	authReq := AuthRequestBody{}
	err := ReadRequestBody(&authReq, req)
	if err != nil {
		errMsg := fmt.Sprintf("Couldn't read request body: %v\n", err)
		WriteResponse(w, http.StatusBadRequest, errMsg)
		return
	}
	if authReq.Email == "" || authReq.Code == "" {
		errMsg := fmt.Sprintf("forUser and authCode are needed for endpoint /code")
		WriteResponse(w, http.StatusBadRequest, errMsg)
		return
	}
	valid := authcode.ValidateAuthCode(authReq.Email, authReq.Code)
	if secret, okToSign := os.LookupEnv(s.Config.JWS.TokenSecret); !valid || !okToSign {
		jwt := authjwt.NewJWT(authReq.Email, map[string]interface{}{"authorized": true})
		token, _ := authjwt.Export(jwt, []byte(secret))
		WriteResponse(w, http.StatusOK, token)
	} else {
		errMsg := fmt.Sprintf("Invalid code")
		WriteResponse(w, http.StatusUnauthorized, errMsg)
	}
}

func (s *AuthServer) HandleTokenAuthRequest(w http.ResponseWriter, req *http.Request) {
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
		return
	}
	// Verify the authToken included with the request
	token, valid, err := authjwt.Verify(authReq.Token, []byte(s.Config.JWS.TokenSecret))
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

// Start the authentication server.
// Returns the dashboard used to control the server.
func (s *AuthServer) Start() {
	// Open log
	OpenLog()
	fmt.Printf("Starting server. Log file located at %s\n", LogFile)
	// Add handlers
	http.HandleFunc(fmt.Sprintf("auth.%s/mail", s.Config.Domain), s.HandleEmailAuthRequest)
	http.HandleFunc(fmt.Sprintf("auth.%s/code", s.Config.Domain), s.HandleCodeAuthRequest)
	http.HandleFunc(fmt.Sprintf("auth.%s/token", s.Config.Domain), s.HandleTokenAuthRequest)
	// Create dashboard from this AuthServer, and add its endpoint
	createDashboard(s).addEndpoints()
	// Generate address
	fulladdr := fmt.Sprintf("%s:%d", s.Config.Domain, s.Config.Port)
	crt := fmt.Sprintf("./cert/%s.crt", s.Config.Domain)
	key := fmt.Sprintf("./cert/%s.key", s.Config.Domain)
	// Fill out fields for server that aren't created by default
	s.srv = &http.Server{
		Addr:    fulladdr,
		Handler: nil,
	}
	s.wg = sync.WaitGroup{}
	err := Log("Starting auth server at https://%s.%s", "auth", fulladdr)
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		s.wg.Add(1)
		err := s.srv.ListenAndServeTLS(crt, key)
		switch err {
		case http.ErrServerClosed:
			{
				Log("Server closed successfully.")
				break
			}
		default:
			{
				Log("Server closed due to an unexpected error: %v", err)
			}
		}
		s.wg.Done()
	}()
}

// Stop the server. Waits to return until everything closes out.
func (s *AuthServer) Stop() {
	Log("Stopping server.")
	s.srv.Shutdown(context.TODO())
	s.wg.Wait()
	CloseLog()
}
