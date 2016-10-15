package main

import (
	"encoding/xml"
	"fmt"
	"os"
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
	Week []Week // xml unmarshal does not like arrays, so use a slice
}

func main() {
	season := Season{Year: 2016,Week:make([]Week,0,3)}

	week := 0
	season.Week = append(season.Week,Week{Num:week+1})
	season.Week[week].Games = make([]Game, 0, 16)

	season.Week[week].Games = append(season.Week[week].Games,
		Game{TeamV: "jags",
			TeamH:  "cards",
			Status: "final",
			ScoreV: "3",
			ScoreH: "7"})

	season.Week[week].Games = append(season.Week[week].Games,
		Game{TeamV: "broncos",
			TeamH:  "panthers",
			Status: "final",
			ScoreV: "24",
			ScoreH: "10"})

	week = 1
	season.Week = append(season.Week,Week{Num:week+1})
	season.Week[week].Games = make([]Game, 0, 16)

	season.Week[week].Games = append(season.Week[week].Games,
		Game{TeamV: "seahawks",
			TeamH:  "49ers",
			Status: "final",
			ScoreV: "0",
			ScoreH: "53"})

	season.Week[week].Games = append(season.Week[week].Games,
		Game{TeamV: "redskins",
			TeamH:  "eagles",
			Status: "final",
			ScoreV: "10",
			ScoreH: "13"})

	seasonXMLFile, err := os.Create("season2016.xml")
	if err != nil {
		fmt.Println("can not open season2016.xml")
		return
	}
	defer seasonXMLFile.Close()


	enc := xml.NewEncoder(seasonXMLFile)
	enc.Indent("", "    ")
	if err := enc.Encode(season); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
