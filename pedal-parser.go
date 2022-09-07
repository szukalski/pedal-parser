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

type BomBuildComponents struct {
	Component
	Id string `json:"id"`
}

type BomBuild struct {
	BuildName  string               `json:"buildName"`
	Components []BomBuildComponents `json:"buildComponents"`
}

type Bom struct {
	Build []BomBuild `json:"build"`
}

type Build struct {
	BuildName string  `json:"buildName"`
	Pedals    []Pedal `json:"pedals"`
}

type BuildList struct {
	Build []Build `json:"build"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readPedal(b []byte) []BomBuildComponents {
	var bomBuildComponents []BomBuildComponents
	var pedal Pedal
	json.Unmarshal(b, &pedal)
	for i := 0; i < len(pedal.Components); i++ {
		component := Component{
			Category:    pedal.Components[i].Category,
			SubCategory: pedal.Components[i].SubCategory,
			Value:       pedal.Components[i].Value,
			Quantity:    pedal.Components[i].Quantity,
		}
		buildComponent := BomBuildComponents{
			Id:        pedal.Id,
			Component: component,
		}
		bomBuildComponents = append(bomBuildComponents, buildComponent)
	}
	return bomBuildComponents
}

func readBuildList(b []byte) Bom {
	var buildList BuildList
	json.Unmarshal(b, &buildList)
	var bom Bom
	for i := 0; i < len(buildList.Build); i++ {
		build := BomBuild{BuildName: buildList.Build[i].BuildName}
		for j := 0; j < len(buildList.Build[i].Pedals); j++ {
			pedalJson, err := os.Open("../pedals/" + buildList.Build[i].Pedals[j].Id + ".pedal.json")
			check(err)
			byteValue, _ := ioutil.ReadAll(pedalJson)
			var pedal Pedal
			json.Unmarshal(byteValue, &pedal)
			pedalComponents := readPedal(byteValue)
			fmt.Println(pedalComponents)
		}
		bom.Build = append(bom.Build, build)
	}
	return bom
}

func checkpoint() {
	fmt.Println("checkpoint")
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

	build := readBuildList(byteValue)
	byteArray, err := json.Marshal(build)
	check(err)
	fmt.Println(string(byteArray))
}
