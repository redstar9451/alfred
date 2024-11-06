package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sahilm/fuzzy"
)

type AlfredConf struct {
	Valid        bool   `json:"valid"`
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle"`
	Arg          string `json:"arg"`
	Autocomplete string `json:"autocomplete"`
}

type AlfredNode struct {
	Name string     `json:"name"`
	Conf AlfredConf `json:"conf"`
}

type AlfredNodes []AlfredNode

type AlfredConfig struct {
	Nodes AlfredNodes `json:"nodes"`
}

func (n AlfredNodes) String(i int) string {
	return n[i].Name + "," + n[i].Conf.Title + "," + n[i].Conf.Subtitle + "," + n[i].Conf.Arg
}

func (n AlfredNodes) Len() int {
	return len(n)
}

type AlfredJson struct {
	Items []AlfredConf `json:"items"`
}

func LoadConfiguration(file string) AlfredConfig {
	var config AlfredConfig
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		panic("load configuration error")
	}
	return config
}

func errorNode(errMsg string) *AlfredNode {
	errNode := AlfredNode{
		Name: "Error Msg",
	}
	errNode.Conf.Valid = true
	errNode.Conf.Title = errMsg
	errNode.Conf.Autocomplete = errMsg
	return &errNode
}

func displayError(errMsg string) {
	errNode := errorNode(errMsg)
	alfredJson := AlfredJson{
		Items: []AlfredConf{errNode.Conf},
	}
	jsonStr, _ := json.Marshal(alfredJson)
	fmt.Println(string(jsonStr))
	return
}

func main() {
	helpString := "usage: alfred <config filename>"
	if len(os.Args) < 2 {
		displayError("invalid arguments, " + helpString)
		return
	}

	fileName := os.Args[1]
	// cwd, _ := os.Getwd()
	// path := filepath.Join(cwd, "conf", fileName)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		displayError("invalid config file " + fileName)
		return
	}

	config := LoadConfiguration(fileName)

	nameMap := make(map[string]int)
	for _, conf := range config.Nodes {
		if nameMap[conf.Name] != 0 {
			displayError("duplicated conf file " + conf.Name)
			return
		}
		nameMap[conf.Name] = 1
	}

	var matchNodes AlfredNodes
	if len(os.Args) > 2 {
		searchStr := os.Args[2]
		results := fuzzy.FindFrom(searchStr, config.Nodes)
		for _, r := range results {
			matchNodes = append(matchNodes, config.Nodes[r.Index])
		}
	} else {
		matchNodes = config.Nodes
	}

	res := AlfredJson{}
	for _, conf := range matchNodes {
		res.Items = append(res.Items, conf.Conf)
	}
	jsonStr, _ := json.Marshal(res)
	fmt.Println(string(jsonStr))
	return
}
