package service

import (
	"fmt"
	"github.com/nikogura/guestbook/state"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// NewVisitorContent is the initial landing page
const NewVisitorContent = `
<html>
	<body>
		<p>
				Howdy stranger.  Ain't seen you around these parts. </br>

				Sign your name below if you like: </br>
		</p>
		<form action="/guestbook/sign" method="post">
			<div>
				Name: <input type="text" name="visitor" lab>
			</div>
			<div>
				<input type="submit" value="Sign Guestbook">
			</div>
		</form>
	</body>
</html>
`

// page displayed after someone signs the guest book
const postSubmitTemplateContent = `
<html>
	<body>
		<p>Good to meet you {{.Name}}.  Your IP is: {{.IP}}  <br/><br/>If I see you again, I'll greet you by name.</p>
	</body>
</html>
`

const returningVisitorTemplateContent = `
<html>
	<body>
		<p>Good to see you again {{.Name}}!</p>
	</body>
</html>`

var gm *state.GORMStateManager

// template object processing the above  doing this all inline just for simplicity's sake
var postSubmitTemplate = template.Must(template.New("postSubmit").Parse(postSubmitTemplateContent))

var returningVisitorTemplate = template.Must(template.New("returning").Parse(returningVisitorTemplateContent))

// RootHandler is the default handler
func RootHandler(w http.ResponseWriter, r *http.Request) {
	ip := ParseIP(r)

	visitor, err := gm.GetVisitor(ip)
	if err != nil {
		log.Printf("Error searching db for visitor: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if visitor.Name == "" {
		fmt.Fprint(w, NewVisitorContent)
	} else {
		err = returningVisitorTemplate.Execute(w, visitor)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

// SignHandler is for handling the submitted form
func SignHandler(w http.ResponseWriter, r *http.Request) {
	ip := ParseIP(r)

	visitor := state.Visitor{
		Name: r.FormValue("visitor"),
		IP:   ip,
	}

	visitor, err := gm.NewVisitor(visitor)
	if err != nil {
		log.Printf("failed to save visitor: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = postSubmitTemplate.Execute(w, visitor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Run actually runs the server
func Run(address string, manager *state.GORMStateManager) error {
	gm = manager

	http.HandleFunc("/guestbook", RootHandler)
	http.HandleFunc("/guestbook/sign", SignHandler)

	log.Printf("Server starting on address %s", address)

	return http.ListenAndServe(address, nil)
}

// ParseIP returns the visitor's IP address
func ParseIP(r *http.Request) (ip string) {
	log.Printf("Visitor Remote Addr: %s", r.RemoteAddr)
	log.Printf("Visitor Forwarded Addr: %s", r.Header.Get("X-Forwarded-For"))

	if r.Header.Get("X-Forwarded-For") != "" {
		// this could actually be several IPs, we're just gonna play simple and grab the first one
		ips := strings.Split(r.Header.Get("X-Forwarded-For"), ", ")
		if len(ips) > 0 {
			ip = ips[0]
		}

	} else {
		parts := strings.Split(r.RemoteAddr, ":")
		if len(parts) > 0 {
			ip = parts[0]
		}
	}

	return ip
}
