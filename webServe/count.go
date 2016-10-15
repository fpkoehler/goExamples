/* based on lessions from
 * https://astaxie.gitbooks.io/build-web-application-with-golang
 * and https://golang.org/doc/articles/wiki
 * Use a hidden token to keep track of web session.  No login.
 */

package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const startHtml = `
<html>
<head>
<title></title>
</head>
<body>
<form action="/calc/{{.}}" method="post">
    <input type="submit" value="start">
</form>
</body>
</html>`

const playHtml = `
<html>
<head>
<title></title>
</head>
<body>
<form action="/calc/{{.Token}}" method="post">
    <p>{{.Count}}</p>
    <input type="submit" name="submitButton" value="count" />
    <input type="submit" name="submitButton" value="stop" />
</form>
</body>
</html>
`

/* The token can be embedded into the page, but decided to use the URL instead */
//    <input type="hidden" name="token" value="{{.Token}}">

//var templates = template.Must(template.ParseFiles("login.html", "play.html"))

var startTmpl, playTmpl *template.Template

type Play struct {
	Token string
	Count int
}

// will be indexed by token
var players map[string]*Play

/* for validating/restricting the URL path */
var validPath = regexp.MustCompile("^/(start|play|calc)/([a-zA-Z0-9]+)$")

/* Used to extract the token from the URL.  See regex above. */
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil // The title is the second subexpression.
}

/* function closure technique used to extract token from URL then call
 * the handler */
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handler:", r.URL.Path)
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func start(w http.ResponseWriter, r *http.Request) {
	defer fmt.Println("---")
	fmt.Println("start method:", r.Method)

	// create token that will be stored on page
	crutime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(crutime, 10))
	token := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println("token:", token)

	// initialize player session
	if _, ok := players[token]; ok {
		http.Error(w, "player with token already in use: "+token, http.StatusInternalServerError)
		return
	}

	/* Create player and add to map */
	players[token] = &Play{Token: token, Count: -1}

	// render
	err := startTmpl.Execute(w, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/* handles the POST */
func calc(w http.ResponseWriter, r *http.Request, token string) {
	defer fmt.Println("---")
	fmt.Println("calc method:", r.Method)

	r.ParseForm()

	button := r.Form.Get("submitButton")
	if button == "stop" {
		fmt.Println("stop pressed")
		http.Redirect(w, r, "/start", http.StatusFound)
		return
	}

	fmt.Println("calc token:", token)
	if token == "" {
		http.Error(w, "no token in calc", http.StatusInternalServerError)
		return
	}

	if _, ok := players[token]; !ok {
		http.Error(w, "no player for token "+token, http.StatusInternalServerError)
		return
	}

	players[token].Count += 1

	http.Redirect(w, r, "/play/"+token, http.StatusFound)
}

/* handles get and returns the html page */
func play(w http.ResponseWriter, r *http.Request, token string) {
	defer fmt.Println("---")

	fmt.Println("play method:", r.Method) // get request method

	fmt.Println("play token:", token)
	if token == "" {
		http.Error(w, "no token in play", http.StatusInternalServerError)
		return
	}

	player, ok := players[token]
	if !ok {
		http.Error(w, "no player for token "+token, http.StatusInternalServerError)
		return
	}

	fmt.Println(player)
	err := playTmpl.Execute(w, player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	var err error
	players = make(map[string]*Play)
	startTmpl, err = template.New("start").Parse(startHtml)
	playTmpl, err = template.New("play").Parse(playHtml)

	http.HandleFunc("/start", start)
	http.HandleFunc("/play/", makeHandler(play))
	http.HandleFunc("/calc/", makeHandler(calc))

	err = http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
