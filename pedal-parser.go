package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Component struct {
	Category    string `json:"category"`
	SubCategory string `json:"sub-category"`
	Value       string `json:"value"`
	Quantity    int    `json:"quantity"`
}

type Pedal struct {
	Id                string      `json:"id"`
	Description       string      `json:"description"`
	BuildDoc          bool        `json:"buildDoc"`
	BuildNotes        string      `json:"buildNotes"`
	CompareToDesigner string      `json:"compareToDesigner"`
	CompareToName     string      `json:"compareToName"`
	Components        []Component `json:"components"`
	Construction      string      `json:"construction"`
	Lineage           []string    `json:"lineage"`
	Name              string      `json:"name"`
	PcbSource         string      `json:"pcbSource"`
	Schematic         bool        `json:"schematic"`
	Tags              []string    `json:"tags"`
	Type              []string    `json:"type"`
	Title             string      `json:"title"`
}

type BuildComponents struct {
	Component
	Id string `json:"id"`
}

type Build struct {
	BuildName  string            `json:"buildName"`
	Components []BuildComponents `json:"buildComponents"`
}

type Bom struct {
	Build []Build `json:"build"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readPedal(b []byte) {
	var pedal Pedal
	json.Unmarshal(b, &pedal)
	build := Build{BuildName: pedal.Id}
	for i := 0; i < len(pedal.Components); i++ {
		component := Component{
			Category:    pedal.Components[i].Category,
			SubCategory: pedal.Components[i].SubCategory,
			Value:       pedal.Components[i].Value,
			Quantity:    pedal.Components[i].Quantity,
		}
		buildComponent := BuildComponents{
			Id:        pedal.Id,
			Component: component,
		}
		build.Components = append(build.Components, buildComponent)
	}
	byteArray, err := json.Marshal(build)
	check(err)
	fmt.Println(string(byteArray))
	return
}

func main() {
	jsonFile, err := os.Open("./pedal.json")
	check(err)
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	readPedal(byteValue)
}
