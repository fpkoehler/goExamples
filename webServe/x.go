/* based on lessions from 
 * https://astaxie.gitbooks.io/build-web-application-with-golang
 */

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"crypto/md5"
	"time"
	"io"
	"strconv"
)

var templates = template.Must(template.ParseFiles("login.html", "play.html"))

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //Parse url parameters passed, then parse the response packet for the POST body (request body)
	// attention: If you do not call ParseForm method, the following data can not be obtained form
	fmt.Println(r.Form) // print information on server side.
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	fmt.Println("---")
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") // write data to response
}

type Play struct {
	Token string
	Count int
}

var gPlay Play

func login(w http.ResponseWriter, r *http.Request) {
    fmt.Println("login method:", r.Method) // get request method
    if r.Method == "GET" {
        crutime := time.Now().Unix()
        h := md5.New()
        io.WriteString(h, strconv.FormatInt(crutime, 10))
        token := fmt.Sprintf("%x", h.Sum(nil))

        err := templates.ExecuteTemplate(w, "login.html", token)
        	if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
        fmt.Println("---")

    } else {
        // log in request
        r.ParseForm()
        token := r.Form.Get("token")
        if token != "" {
            // check token validity
            fmt.Println("token:", token)
            gPlay.Token = token
            gPlay.Count = 0
        } else {
            // give error if no token
        }
        // chapter 4.3 cross site scripting
        fmt.Println("username length:", len(r.Form["username"][0]))
        fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username"))) // print in server side
        fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
        fmt.Println("---")
//        template.HTMLEscape(w, []byte(r.Form.Get("username"))) // respond to client
//        err := templates.ExecuteTemplate(w, "play.html", &play)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//		}
		http.Redirect(w, r, "/play", http.StatusFound)
    }
}

func play(w http.ResponseWriter, r *http.Request) {
    fmt.Println("play method:", r.Method) // get request method
    if r.Method == "GET" {
        err := templates.ExecuteTemplate(w, "play.html", &gPlay)
        	if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
        fmt.Println("---")
    } else {
        r.ParseForm()
        token := r.Form.Get("token")
        if token != "" {
            // check token validity
            fmt.Println("token:", token)
//            gPlay.Token = token
            gPlay.Count += 1
        } else {
        	fmt.Println("no token in play.html")
            // give error if no token
        }
		http.Redirect(w, r, "/play", http.StatusFound)
    }
}

func main() {
	http.HandleFunc("/", sayhelloName) // setting router rule
	http.HandleFunc("/login", login)
	http.HandleFunc("/play", play)
	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
