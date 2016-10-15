/* Based on https://mschoebel.info/2014/03/09/snippet-golang-webapp-login-logout */

package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Play struct {
	User   string
	PwHash [md5.Size]byte
	Count  int
}

// will be indexed by user name
var players map[string]*Play

func getUserName(r *http.Request) (userName string) {
	if cookie, err := r.Cookie("session"); err == nil {
		userName = cookie.Value
	}
	return userName
}

func clearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

// login handler

func loginHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	pass := r.FormValue("password")
	redirectTarget := "/"

	if name == "" {
		fmt.Println("no user name specified")
		http.Redirect(w, r, redirectTarget, http.StatusFound)
	}

	if pass == "" {
		fmt.Println("no password specified")
		http.Redirect(w, r, redirectTarget, http.StatusFound)
	}

	/* hash the password */
	pwHash := md5.Sum([]byte(pass))

	player, ok := players[name]
	if !ok {
		/* Create player and add to map */
		fmt.Println("create player ", name)
		players[name] = &Play{User: name, PwHash: pwHash, Count: 0}
		fmt.Println("created player", players[name])
	} else {
		if pwHash != player.PwHash {
			fmt.Println("password check failed, pwHash ", pwHash, "name.pwHash", player.PwHash)
			http.Redirect(w, r, redirectTarget, http.StatusFound)
		} else {
			fmt.Println("found player", players[name])
		}
	}

	cookie := &http.Cookie{
		Name:  "session",
		Value: name,
		Path:  "/",
	}
	http.SetCookie(w, cookie)

	redirectTarget = "/internal"

	http.Redirect(w, r, redirectTarget, http.StatusFound)
}

// logout handler

func logoutHandler(w http.ResponseWriter, r *http.Request) {
//	userName := getUserName(r)
//	delete(players, userName)
	clearSession(w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func calcHandler(w http.ResponseWriter, r *http.Request) {
	userName := getUserName(r)
	fmt.Println("calc user ", userName)
	if userName == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	player, ok := players[userName]
	if !ok {
		http.Error(w, "no player for "+userName, http.StatusInternalServerError)
		return
	}

	player.Count += 1
	fmt.Println("calc", player)

	http.Redirect(w, r, "/internal", http.StatusFound)
}

// index page

const indexPage = `
<h1>Login</h1>
<form method="post" action="/login">
    <label for="name">User name</label>
    <input type="text" id="name" name="name">
    <label for="password">Password</label>
    <input type="password" id="password" name="password">
    <button type="submit">Login</button>
</form>
`

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, indexPage)
}

// internal page

const internalPage = `
<h1>Internal</h1>
<hr>
<small>User: {{.User}}</small>
<form action="/calc" method="post">
    <p>{{.Count}}</p>
    <input type="submit" name="submitButton" value="count" />
    <input type="submit" name="submitButton" value="stop" />
</form>
<form method="post" action="/logout">
    <button type="submit">Logout</button>
</form>
`

var internalTmpl *template.Template

func internalPageHandler(w http.ResponseWriter, r *http.Request) {
	userName := getUserName(r)
	if userName == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	player, ok := players[userName]
	if !ok {
		http.Error(w, "no player for "+userName, http.StatusInternalServerError)
		return
	}

	err := internalTmpl.Execute(w, player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	var err error
	players = make(map[string]*Play)
	internalTmpl, err = template.New("internal").Parse(internalPage)

	/* handlers for GETs */
	http.HandleFunc("/", indexPageHandler)
	http.HandleFunc("/internal", internalPageHandler)

	/* handlers for POSTs */
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/calc", calcHandler)

	err = http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
