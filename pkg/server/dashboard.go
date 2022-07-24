package server

import (
	"auth/pkg/credentials"
	"auth/pkg/gatekey"
	gatemail "auth/pkg/mail"
	"crypto/tls"
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"
	"time"
)

// Persistent data for Dashboard.
type DashboardData struct {
	AppName string
	EmailOk bool
	TLSOk   bool
}

// The Dashboard is a site that allows control over the authentication server.
// It also serves as a handler for incoming requests; see ServeHTTP below
type Dashboard struct {
	Data DashboardData
	// PRIVATE
	srv            *AuthServer
	serveAddr      string
	serveDirectory string
}

// Create a dashboard from an AuthServer
func createDashboard(fromServer *AuthServer) *Dashboard {
	return &Dashboard{
		Data: DashboardData{
			AppName: "Test App",
			EmailOk: true,
			TLSOk:   false,
		},
		srv:            fromServer,
		serveAddr:      fmt.Sprintf("auth.%s/dashboard/", fromServer.Config.Domain),
		serveDirectory: "./dat/dashboardRoot",
	}
}

// Set icon in map at key string to value string in /dashboard/resource/img
// Used for templating; see tmplData below in ServeHTTP
func setIcon(in map[string]interface{}, at string, to string) {
	in[at] = fmt.Sprintf("/dashboard/resource/img/%s", to)
}

// Serve requests to dashboard
func (d *Dashboard) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Acquire and sanitize request URL
	requrl := filepath.Clean(r.URL.Path)
	reqfile := filepath.Join(d.serveDirectory, requrl)
	// Update dashboard data depending on the requested endpoint.
	tmplData := make(map[string]interface{})
	tmplData["Info"] = d.Data
	switch requrl {
	case "/dashboard":
		{
			// Check authentication token
			if authCookie, err := r.Cookie("auth-admin-jwt"); err == nil {
				key, valid, err := gatekey.Verify(authCookie.Value, []byte(d.srv.Config.JWT.TokenSecret))
				if err != nil || !valid || !key.Body.Permissions["admin"] {
					http.Redirect(w, r, "/dashboard/login", http.StatusFound)
					break
				}
			} else {
				http.Redirect(w, r, "/dashboard/login", http.StatusFound)
				break
			}
			// Get SMTP config, and send an email to the server config test email to make sure no error is returned.
			tmplData["SMTPHost"] = d.srv.SMTPHost()
			if !d.Data.EmailOk {
				if err := gatemail.SendMessage(d.srv.SMTPHost(), d.srv.Config.SMTPHost.TestEmail, []byte("Ping!")); err != nil {
					d.Data.EmailOk = true
					setIcon(tmplData, "SMTPOkIcon", "yes.svg")
				} else {
					d.Data.EmailOk = false
					setIcon(tmplData, "SMTPOkIcon", "no.svg")
				}
			} else {
				setIcon(tmplData, "SMTPOkIcon", "yes.svg")
			}
			// Database
			setIcon(tmplData, "DBOkIcon", "yes.svg")
			tmplData["DBInfo"] = "auth currently uses a local database. More support will be added in future updates here."
			// TLS
			if !d.Data.TLSOk {
				// Dial the server host using TLS.
				fulladdr := fmt.Sprintf("%s:%d", d.srv.Config.Domain, d.srv.Config.Port)
				if _, err := tls.Dial("tcp", fulladdr, &tls.Config{InsecureSkipVerify: true}); err != nil {
					setIcon(tmplData, "TLSOkIcon", "yes.svg")
					tmplData["TLSInfo"] = "TLS connection to auth successful."
				} else {
					setIcon(tmplData, "TLSOkIcon", "no.svg")
					tmplData["TLSInfo"] = `TLS connection to auth failed.
					This can occur if your certificate is self-signed or expired. Please check ./cert to ensure the correct cert and key are uploaded.
					If you are using self-signed certificates for testing, you can ignore this message.`
				}
			} else {
				setIcon(tmplData, "TLSOkIcon", "yes.svg")
				tmplData["TLSInfo"] = "TLS connection to auth successful."
			}
			// Controls
			tmplData["AuthOpen"] = d.srv.Open
			break
		}
	case "/dashboard/login":
		{

		}
	default:
		{
			WriteResponse(w, http.StatusBadRequest, "Can't go there.")
			return
		}
	}
	// Get template for this request; HTML structure is [request URL]/index.html
	templateFile := filepath.Join(reqfile, "index.html")
	tmpl := template.Must(template.ParseFiles(templateFile))
	// Execution behavior differs based on requrl, since different pages need different template
	// data (and potentially need to run requests)
	// Double-execute, once to error check and once to write
	err := tmpl.Execute(w, tmplData)
	if err != nil {
		fmt.Println(err)
	}
	//http.ServeFile(w, r, reqfile)
}

func (d *Dashboard) handleSMTP(w http.ResponseWriter, r *http.Request) {
	// TODO: Authenticate admin user(s).
	// Handle post requests through parsing form. Modify backend based on response.
	if r.Method == http.MethodPost {
		r.ParseForm()
		newHost := r.Form["smtpHost"][0]
		newSend := r.Form["sendFrom"][0]
		// TODO: These should be sanitized.
		d.srv.Config.SMTPHost.Host = newHost
		d.srv.Config.SMTPHost.Sender = newSend
	}
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func (d *Dashboard) handleControls(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		_, open := r.Form["open"]
		d.srv.Open = open
	}
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func (d *Dashboard) handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		//email := r.Form["email"][0]
		username := r.Form["username"][0]
		password := r.Form["password"][0]
		fmt.Println("Validating admin user")
		valid, user, err := credentials.ValidateUserCred(username, password)
		if err == nil && valid {
			admin, ok := user.Permissions["admin"]
			if ok && admin {
				fmt.Printf("Admin user %s logged in\n", user.Username)
				// Set cookie to admin token
				jwt := gatekey.NewGateKey("jani9652", user.Permissions, time.Duration(d.srv.Config.JWT.AdminValidTime)*time.Minute)
				token := gatekey.Export(jwt, []byte(d.srv.Config.JWT.TokenSecret))
				http.SetCookie(w, &http.Cookie{Name: "auth-admin-jwt", Value: token, Path: "/dashboard"})
				http.Redirect(w, r, "/dashboard", http.StatusFound)
			}
		}
	}
	http.Redirect(w, r, "/dashboard/login", http.StatusSeeOther)
}

func (d *Dashboard) addEndpoints() {
	http.Handle(d.serveAddr, d)
	// Use a raw FileServer for handling resources; root is a couple dirs down so non-resource content is not exposed
	resourceFS := http.FileServer(http.Dir(d.serveDirectory + "/dashboard/resource"))
	http.Handle(d.serveAddr+"resource/", http.StripPrefix("/dashboard/resource/", resourceFS))
	http.HandleFunc(d.serveAddr+"update-config-smtp", d.handleSMTP)
	http.HandleFunc(d.serveAddr+"update-config-controls", d.handleControls)
	http.HandleFunc(d.serveAddr+"login/admin-login", d.handleAdminLogin)
}
