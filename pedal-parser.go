package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
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

type BomComponents struct {
	Component
	Id string `json:"id"`
}

type Bom struct {
	BuildName  string          `json:"buildName"`
	Components []BomComponents `json:"buildComponents"`
}

type BuildList struct {
	BuildName string  `json:"buildName"`
	Pedals    []Pedal `json:"pedals"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getPedalComponents(pedalId string) []BomComponents {
	pedalJson, err := os.Open("../pedals/" + pedalId + ".pedal.json")
	check(err)
	byteValue, _ := ioutil.ReadAll(pedalJson)
	var bomComponents []BomComponents
	var pedal Pedal
	json.Unmarshal(byteValue, &pedal)
	for i := 0; i < len(pedal.Components); i++ {
		component := Component{
			Category:    pedal.Components[i].Category,
			SubCategory: pedal.Components[i].SubCategory,
			Value:       pedal.Components[i].Value,
			Quantity:    pedal.Components[i].Quantity,
		}
		var id string
		switch pedal.Components[i].Category {
		case "Enclosure", "Knob":
			id = pedal.Id
		default:
			id = "All"
		}
		buildComponent := BomComponents{
			Id:        id,
			Component: component,
		}
		bomComponents = append(bomComponents, buildComponent)
	}
	return bomComponents
}

func buildBom(b []byte) Bom {
	var buildList BuildList
	json.Unmarshal(b, &buildList)
	bom := Bom{
		BuildName: buildList.BuildName,
	}
	for i := 0; i < len(buildList.Pedals); i++ {
		pedalComponents := getPedalComponents(buildList.Pedals[i].Id)
		bom.Components = append(bom.Components, pedalComponents...)
	}
	return bom
}

func checkpoint() {
	fmt.Println("checkpoint")
	return
}

func printBomToCsv(b []byte) {
	var bom Bom
	json.Unmarshal(b, &bom)
	sort.Slice(bom.Components, func(i, j int) bool {
		return bom.Components[i].Category < bom.Components[j].Category
	})
	for j := 0; j < len(bom.Components); j++ {
		fmt.Println(bom.BuildName + "," + bom.Components[j].Id + "," + bom.Components[j].Category + "," + bom.Components[j].SubCategory + "," + bom.Components[j].Value + "," + fmt.Sprint(bom.Components[j].Quantity))
	}
	return
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: pedal-parser 'pedal-buildlist.json'")
		os.Exit(1)
	}
	buildList, err := os.Open(os.Args[1])
	check(err)
	defer buildList.Close()
	byteValue, _ := ioutil.ReadAll(buildList)

	build := buildBom(byteValue)
	byteArray, err := json.Marshal(build)
	check(err)
	printBomToCsv(byteArray)
	//fmt.Println(string(byteArray))
}
