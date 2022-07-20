package authserver

import (
	"auth/pkg/authmail"
	"crypto/tls"
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"
)

// Persistent data for Dashboard.
// TODO: this needs to be save-able
type DashboardData struct {
	AppName   string
	EmailOk   bool
	DBTesting bool
	TLSOk     bool
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
			AppName:   "Test App",
			EmailOk:   true,
			DBTesting: true,
			TLSOk:     false,
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
			// Get SMTP config, and send an email to the server config test email to make sure no error is returned.
			tmplData["SMTPHost"] = d.srv.SMTPHost()
			if !d.Data.EmailOk {
				if err := authmail.SendMessage(d.srv.SMTPHost(), d.srv.Config.SMTPHost.TestEmail, []byte("Ping!")); err != nil {
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
			tmplData["DBInfo"] = "Database support incoming."
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
	fmt.Printf("Received %s", r.Method)
	if r.Method == http.MethodPost {
		r.ParseForm()
		newHost := r.Form["smtpHost"][0]
		newSend := r.Form["sendFrom"][0]
		// TODO: These should be sanitized.
		d.srv.Config.SMTPHost.Host = newHost
		d.srv.Config.SMTPHost.Sender = newSend
	}
	http.Redirect(w, r, "./", http.StatusSeeOther)
}

func (d *Dashboard) handleDB(w http.ResponseWriter, r *http.Request) {
	// TODO: Authenticate admin user(s).
	// Handle post requests through parsing form. Modify backend based on response.
	fmt.Printf("Received %s", r.Method)
	if r.Method == http.MethodPost {
		r.ParseForm()
		testing := r.Form["Testing Mode"][0]
		// TODO: These should be sanitized.
		d.Data.DBTesting = testing == "true"
	}
	http.Redirect(w, r, "./", http.StatusSeeOther)
}

func (d *Dashboard) addEndpoints() {
	http.Handle(d.serveAddr, d)
	// Use a raw FileServer for handling resources; root is a couple dirs down so non-resource content is not exposed
	resourceFS := http.FileServer(http.Dir(d.serveDirectory + "/dashboard/resource"))
	http.Handle(d.serveAddr+"resource/", http.StripPrefix("/dashboard/resource/", resourceFS))
	http.HandleFunc(d.serveAddr+"update-config-smtp", d.handleSMTP)
}
