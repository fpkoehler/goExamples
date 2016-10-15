package main

import (
	"os"
    "fmt"
    "bufio"
    "io/ioutil"
    "net/http"
)

var fileName string

func getInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter file path: ")
		fileName, _ = reader.ReadString('\n')
		fileName = fileName[:len(fileName)-1] // chop off '\n'
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadFile(fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

    fmt.Fprint(w, string(body))
}

func main() {
	fileName = "x.html"
	go getInput()
	fmt.Println("http://localhost:9090")
    http.HandleFunc("/", handler)
    http.ListenAndServe(":9090", nil)
}