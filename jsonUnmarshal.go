package main

import (
	"encoding/json"
	"fmt"
)

type Link struct {
	PeerId   int
	PeerName string
}

type Box struct {
	MaskList []int
	LinkInfo map[string]Link
}

type CommandOutput struct {
	Boxes map[string]Box
}

func main() {
	var jsonBlob = []byte(`{
		"boxes": {
			"box1": {
				"maskList" : [7,8,9],
				"linkInfo" : {
					"1" : {
						"count" : 0,
						"peerId" : 99,
						"peerName" : "square1"
					},
					"2" : {
						"count" : 0,
						"peerId" : 98,
						"peerName" : "square2"
					}
				}
			},
			"box2": {
				"maskList" : [7,6,5],
				"linkInfo" : {
					"5" : {
						"count" : 0,
						"peerId" : 80,
						"peerName" : "square5"
					},
					"7" : {
						"count" : 0,
						"peerId" : 70,
						"peerName" : "square4"
					}
				}
			}
		} 
	}`)

	var commandOutput CommandOutput
	err := json.Unmarshal(jsonBlob, &commandOutput)
	if err != nil {
		fmt.Println("error:", err)
	}
	//fmt.Println(commandOutput)

	for boxName, box := range commandOutput.Boxes {
		fmt.Println(boxName)
		for link, linkInfo := range box.LinkInfo {
			fmt.Println(link, linkInfo)
		}
	}

}
