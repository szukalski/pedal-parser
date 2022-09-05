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

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func readPedal(f string) {
	jsonFile, err := os.Open(f)
	check(err)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var pedal Pedal
	json.Unmarshal(byteValue, &pedal)
	fmt.Println(pedal.Id)
	for i := 0; i < len(pedal.Components); i++ {
		fmt.Println(pedal.Components[i].Category)
		fmt.Println(pedal.Components[i].SubCategory)
		fmt.Println(pedal.Components[i].Value)
		fmt.Println(pedal.Components[i].Quantity)
	}
	return
}

func main() {
	readPedal("./pedal.json")
}
