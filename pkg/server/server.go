package server

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/mail"
	"os"
	"sync"
	"time"

	"github.com/jakenichols2719/gate/pkg/credentials"
	"github.com/jakenichols2719/gate/pkg/gatecode"
	"github.com/jakenichols2719/gate/pkg/gatekey"
	gatemail "github.com/jakenichols2719/gate/pkg/mail"
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
	Config *AuthServerConfig // Configuration settings
	Open   bool              // Is server open to API calls?
	// Server; not exported, used internally for controlling HTTPS server
	srv *http.Server
	// Waitgroup; needed to maintain concurrency with server
	wg sync.WaitGroup
}

// Create a new AuthServer (authentication server) using an AuthServerConfig.
//
// Input:
//   - cfg *AuthServerConfig: Target configuration. Should be non-nil.
// Output:
//   - *AuthServer: A new server object containing configuration, an *http.Server, and a waitgroup
func NewServer(cfg *AuthServerConfig) *AuthServer {
	return &AuthServer{
		Config: cfg,
		Open:   true,
	}
}

// Create a mail host from the calling server's configuration.
//
// Calling:
//   - srv *AuthServer: Server to use with this SMTP configuration
// Output:
//   - (gate)mail.Host: SMTP host for use with pkg/mail.
func (srv *AuthServer) SMTPHost() gatemail.Host {
	return gatemail.Host{
		Username: srv.Config.SMTPHost.Username,
		Password: srv.Config.SMTPHost.Password,
		Host:     srv.Config.SMTPHost.Host,
		Port:     srv.Config.SMTPHost.Port,
		Sender:   srv.Config.SMTPHost.Sender,
	}
}

// Write out an HTTP response with a code/message.
//
// Input:
//   - w http.ResponseWriter: These are always given to http-response-capable functions; just pass that here.
//   - code int: HTTP response code.
//   - msg string: Message to send with the code. Technically optional, though very good form to include an informative body, especially with codes >399.
func WriteResponse(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

// Request body format for all authentication requests.
type AuthRequestBody struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
	Code        string `json:"authCode"`
	GetKey      bool   `json:"getKey"`
	Key         string `json:"gateKey"`
}

// Read the body of an http request with AuthRequestBody params.
//
// Input:
//   - out *AuthRequestBody: Pointer to an AuthRequestBody object to read data into
//   - req *http.Request: Request to read from. This uses ioutil.ReadAll, which means it depletes the buffer; trying to call
//   any other read on the request after ReadRequestBody will make the body appear to be empty.
// Output:
//   - Error, if one occurs. Non-POST requests and invalid JSON will cause this.
func ReadRequestBody(out *AuthRequestBody, req *http.Request) error {
	if req.Method != http.MethodPost {
		return errors.New("gate requests MUST be POST requests.")
	}
	if apikey := req.Header.Get("x-api-key"); apikey != "" {
		// Check API key "user"
		if validKey, _, err := credentials.ValidateUserCred("api", apikey); err != nil || !validKey {
			return errors.New("couldn't validate the x-api-key header field")
		}
	} else {
		return errors.New("gate requests require the x-api-key header with a valid API key")
	}
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

// Credential registration
func (s *AuthServer) handleCredRegiRequest(w http.ResponseWriter, req *http.Request) {
	if !s.Open {
		WriteResponse(w, http.StatusInternalServerError, "Server is currently disabled")
		return
	}
	authReq := AuthRequestBody{}
	err := ReadRequestBody(&authReq, req)
	if err != nil {
		errMsg := fmt.Sprintf("Couldn't read request body: %v\n", err)
		WriteResponse(w, http.StatusBadRequest, errMsg)
		return
	}
	if authReq.Email == "" || authReq.Username == "" || authReq.Password == "" {
		errMsg := fmt.Sprintf("email, username, and password are needed for endpoint /register\n")
		WriteResponse(w, http.StatusBadRequest, errMsg)
		return
	}
	// Make sure email is valid
	if _, err := mail.ParseAddress(authReq.Email); err != nil {
		errMsg := fmt.Sprintf("invalid email\n")
		WriteResponse(w, http.StatusBadRequest, errMsg)
		return
	}
	// Register user.
	if err := credentials.RegisterUser(authReq.Email, authReq.Username, authReq.Password, nil); err != nil {
		errMsg := fmt.Sprintf("Registration failed: %v\n", err)
		WriteResponse(w, http.StatusBadRequest, errMsg)
	} else {
		succMsg := fmt.Sprintf("User %s registered successfully under email %s\n", authReq.Username, authReq.Email)
		WriteResponse(w, http.StatusOK, succMsg)
	}
}

// Credential authorization
func (s *AuthServer) handleCredAuthRequest(w http.ResponseWriter, req *http.Request) {
	if !s.Open {
		WriteResponse(w, http.StatusInternalServerError, "Server is currently disabled")
		return
	}
	authReq := AuthRequestBody{}
	err := ReadRequestBody(&authReq, req)
	if err != nil {
		errMsg := fmt.Sprintf("Couldn't read request body: %v\n", err)
		WriteResponse(w, http.StatusBadRequest, errMsg)
		return
	}
	if authReq.Username == "" || authReq.Password == "" {
		errMsg := fmt.Sprintf("username and password are needed for endpoint /register\n")
		WriteResponse(w, http.StatusBadRequest, errMsg)
		return
	}
	valid, entry, err := credentials.ValidateUserCred(authReq.Username, authReq.Password)
	if !valid {
		errMsg := fmt.Sprintf("Invalid credentials\n")
		WriteResponse(w, http.StatusUnauthorized, errMsg)
		return
	}
	jwt := gatekey.NewGateKey(authReq.Username, entry.Permissions, time.Duration(s.Config.GateKey.UserValidTime)*time.Minute)
	token := gatekey.Export(jwt, []byte(s.Config.GateKey.GatekeySecret))
	if authReq.GetKey {
		WriteResponse(w, http.StatusOK, token)
	} else {
		WriteResponse(w, http.StatusOK, "no token requested; set getToken=true in request body for an auth token\n")
	}
}

// Credential change
func (s *AuthServer) handlePwdChangeRequest(w http.ResponseWriter, req *http.Request) {
	if !s.Open {
		WriteResponse(w, http.StatusInternalServerError, "Server is currently disabled")
		return
	}
	authReq := AuthRequestBody{}
	err := ReadRequestBody(&authReq, req)
	if err != nil {
		errMsg := fmt.Sprintf("Couldn't read request body: %v\n", err)
		WriteResponse(w, http.StatusBadRequest, errMsg)
		return
	}
	if authReq.Username == "" || authReq.Password == "" || authReq.NewPassword == "" {
		errMsg := fmt.Sprintf("username, password, and new password are needed for endpoint /changePassword\n")
		WriteResponse(w, http.StatusBadRequest, errMsg)
		return
	}
	err = credentials.ChangeUserPassword(authReq.Username, authReq.Password, authReq.NewPassword)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to reset password: %v\n", err)
		WriteResponse(w, http.StatusUnauthorized, errMsg)
	} else {
		succMsg := fmt.Sprintf("Password changed successfully. Please log back in.\n")
		WriteResponse(w, http.StatusOK, succMsg)
	}
}

// Handle email authentication requests
func (s *AuthServer) HandleEmailAuthRequest(w http.ResponseWriter, req *http.Request) {
	if !s.Open {
		WriteResponse(w, http.StatusInternalServerError, "Server is currently disabled")
		return
	}
	authReq := AuthRequestBody{}
	err := ReadRequestBody(&authReq, req)
	if err != nil {
		errMsg := fmt.Sprintf("Couldn't read request body: %v\n", err)
		WriteResponse(w, http.StatusBadRequest, errMsg)
		return
	}
	// Also throw a 400 if there's no included email
	if authReq.Email == "" {
		errMsg := fmt.Sprintf("email is needed for endpoint /mail\n")
		WriteResponse(w, http.StatusBadRequest, errMsg)
	}
	// Send an authentication email and write out 200=
	code := gatecode.NewGateCode(authReq.Email)
	msg := gatemail.NewAuthMessage(authReq.Email, code)
	gatemail.SendMessage(s.SMTPHost(), authReq.Email, msg)
	succMsg := fmt.Sprintf("Authentication email sent to %s\n", authReq.Email)
	WriteResponse(w, http.StatusOK, succMsg)
}

// Handle authentication code requests
func (s *AuthServer) HandleCodeAuthRequest(w http.ResponseWriter, req *http.Request) {
	if !s.Open {
		WriteResponse(w, http.StatusInternalServerError, "Server is currently disabled")
		return
	}
	// Read in body. Send a 400 on failure
	authReq := AuthRequestBody{}
	err := ReadRequestBody(&authReq, req)
	if err != nil {
		errMsg := fmt.Sprintf("Couldn't read request body: %v\n", err)
		WriteResponse(w, http.StatusBadRequest, errMsg)
		return
	}
	if authReq.Email == "" || authReq.Code == "" {
		errMsg := fmt.Sprintf("email and authCode are needed for endpoint /code\n")
		WriteResponse(w, http.StatusBadRequest, errMsg)
		return
	}
	valid := gatecode.ValidateGateCode(authReq.Email, authReq.Code)
	if !valid {
		errMsg := fmt.Sprintf("Invalid code")
		WriteResponse(w, http.StatusUnauthorized, errMsg)
	}
	jwt := gatekey.NewGateKey(authReq.Email, map[string]bool{"authorized": true}, time.Duration(s.Config.GateKey.UserValidTime)*time.Minute)
	token := gatekey.Export(jwt, []byte(s.Config.GateKey.GatekeySecret))
	if authReq.GetKey {
		WriteResponse(w, http.StatusOK, token)
	} else {
		WriteResponse(w, http.StatusOK, "no token requested; set getToken=true in request body for an auth token\n")
	}
}

func (s *AuthServer) HandleKeyAuthRequest(w http.ResponseWriter, req *http.Request) {
	if !s.Open {
		WriteResponse(w, http.StatusInternalServerError, "Server is currently disabled")
		return
	}
	// Read in body. Send a 400 on failure
	authReq := AuthRequestBody{}
	err := ReadRequestBody(&authReq, req)
	if err != nil {
		errMsg := fmt.Sprintf("Couldn't read request body: %v\n", err)
		WriteResponse(w, http.StatusBadRequest, errMsg)
		return
	}
	if authReq.Key == "" {
		errMsg := fmt.Sprintf("authToken is needed for endpoint /token\n")
		WriteResponse(w, http.StatusBadRequest, errMsg)
		return
	}
	// Verify the authToken included with the request
	token, valid, err := gatekey.Verify(authReq.Key, []byte(s.Config.GateKey.GatekeySecret))
	if err != nil {
		errMsg := fmt.Sprintf("Couldn't process bearer token: %v\n", err)
		WriteResponse(w, http.StatusUnauthorized, errMsg)
		return
	}
	if !valid {
		errMsg := fmt.Sprintf("Bearer token has is altered or expired. Re-authentication is required.\n")
		WriteResponse(w, http.StatusUnauthorized, errMsg)
		return
	}

	outToken, _ := json.Marshal(token)
	WriteResponse(w, http.StatusOK, string(outToken))
}

// Start an authentication server.
// This has two different behaviors; a new server with an empty database and any other server.
// New users will have to create an admin account, be provided with a randomly-generated password and an API key,
// and then asked to restart the server.
//
// Calling:
//   - s *AuthServer: Server to run. Should be initialized with NewServer before calling Start.
func (s *AuthServer) Start() error {
	// Open log
	OpenLog()
	fmt.Printf("Starting server. Log file located at %s\n", LogFile)
	// Open database
	credentials.OpenDB("./dat/database/auth.db")
	// Check entries. Count as first run if empty; otherwise, attempt to log in as admin. Admin credentials must be provided to run Gate.
	if credentials.Entries() == 0 {
		Log("Welcome to Gate. Your database is empty; your provided admin credentials will be used to create the admin account.")
		var password string
		Log("Your password is below. It will never be output again--save it somewhere secure.")
		pwd := make([]byte, 32)
		rand.Read(pwd[:])
		password = base64.RawURLEncoding.EncodeToString(pwd)
		fmt.Println(password)
		err := credentials.RegisterUser(s.Config.Admin.Email, s.Config.Admin.Username, password, map[string]bool{"admin": true})
		if err != nil {
			return err
		}
		Log("All Gate API calls require an API key. Your API key is below. It will never be output again--save it somewhere secure.")
		ak := make([]byte, 32)
		rand.Read(ak[:])
		apikey := base64.RawURLEncoding.EncodeToString(ak)
		fmt.Println(apikey)
		err = credentials.RegisterUser("nil", "api", apikey, map[string]bool{"apikey": true})
		if err != nil {
			return err
		}
		Log("Gate will now restart. Please set the gate-admin-password secret to log in.")
		os.Exit(0)
	} else {
		username, password := s.Config.Admin.Username, s.Config.Admin.Password
		if valid, user, err := credentials.ValidateUserCred(username, password); valid && user.Permissions["admin"] {
			Log("Logged in as admin user '%s'", username)
		} else if err != nil {
			return err
		} else {
			return fmt.Errorf("Couldn't start server with admin credentials for unknown reason.")
		}
	}
	// Add handlers
	http.HandleFunc("/register", s.handleCredRegiRequest)
	http.HandleFunc("/login", s.handleCredAuthRequest)
	http.HandleFunc("/resetPassword", s.handlePwdChangeRequest)
	http.HandleFunc("/mail", s.HandleEmailAuthRequest)
	http.HandleFunc("/code", s.HandleCodeAuthRequest)
	http.HandleFunc("/key", s.HandleKeyAuthRequest)
	// Create dashboard from this AuthServer, and add its endpoint
	createDashboard(s).addEndpoints()
	// Generate address
	// fulladdr := fmt.Sprintf("%s:%d", s.Config.Domain, s.Config.Port)
	crt := "/run/secrets/gate-ssl-crt"
	key := "/run/secrets/gate-ssl-key"
	fmt.Printf("%s, %s\n", crt, key)
	// Fill out fields for server that aren't created by default
	//fmt.Println(s.Config.Address)
	s.srv = &http.Server{
		Addr: s.Config.Address,
	}
	s.wg = sync.WaitGroup{}
	err := Log("Starting auth server at https://%s", s.srv.Addr)
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
	return nil
}

// Stop the server. Waits to return until everything closes out.
func (s *AuthServer) Stop() {
	Log("Stopping server.")
	s.srv.Shutdown(context.TODO())
	s.wg.Wait()
	CloseLog()
}
