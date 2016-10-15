package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type Game struct {
	TeamV  string
	TeamH  string
	Status string
	ScoreV string
	ScoreH string
}

type Week struct {
	Num   int `xml:"num,attr"`
	Games []Game
}

type Season struct {
	Year int `xml:"Year,attr"`
	Week []Week
}

func main() {
	/* http://stackoverflow.com/questions/1821811/how-to-read-write-from-to-file 
	 * xml.Unmarshal()'s first argument is a []byte.  A little scary that the only
	 * way to get a []byte from the file is to read the entire file.  
	*/
	b, err := ioutil.ReadFile("season2016.xml")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var season Season
	err = xml.Unmarshal(b, &season)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Println(season)
}
