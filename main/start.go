package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

//helper function to avoid err!=nil repetition
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func encodeVersionHistory(versions map[string]string) {
	file, err := os.Create("updateHistory.txt")
	checkErr(err)

	//save version history as gob serialized map
	e := gob.NewEncoder(file)
	e.Encode(versions)
	file.Close()
}

func decodeVersionHistory(file *os.File) map[string]string {
	var versions = make(map[string]string)

	//decode gob map
	d := gob.NewDecoder(file)
	err := d.Decode(&versions)
	checkErr(err)

	file.Close()
	return versions
}

func getLatestVersions(selections *[]string) map[string]string {
	var versions = make(map[string]string)
	for _, id := range *selections {
		fmt.Println("Processing " + addons[id])
		version := *getAddonInfo(&id)
		if len(version) == 0 {
			fmt.Println("Error processing version, skipping")
			continue
		}
		fmt.Println("Latest Version: " + version)
		fmt.Println("-----------------------------------")
		versions[id] = version
	}
	return versions
}

func processSelections(selections *[]string) {
	var versions = make(map[string]string)

	file, err := os.Open("updateHistory.txt")
	if err != nil {
		fmt.Println("I noticed this may be your first time using the app.")
		fmt.Println("In order to set the version history, we must first update the existing addons.")

		versions = getLatestVersions(selections)

		//save version history as gob serialized map
		encodeVersionHistory(versions)

		//download each addon
		downloadSelections(selections)
	} else {
		var queue []string
		//retrieve download history
		versions = decodeVersionHistory(file)

		//check for updates
		latest := getLatestVersions(selections)
		for id, v := range latest {
			if versions[id] != v {
				queue = append(queue, id)
			}
		}
		//save version history
		encodeVersionHistory(versions)

		//download each addon
		downloadSelections(&queue)
	}
}

func printSelections(selections *[]string) {
	for _, id := range *selections {
		fmt.Println(addons[id])
	}
}

func takeUserInput(message string) *string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(message)
	text, _ := reader.ReadString('\n')
	return &text
}


func main() {
	b, err := ioutil.ReadFile("config.txt")
	selectionMessage := "Please enter the IDs of the addons you want updated.\nEx: 18,9,72"
	var selections []string
	if err != nil {
		fmt.Println("Config file was not found!")
		//Direct the user to enter in a new selection
		selections = strings.Split(*takeUserInput(selectionMessage), ",")
	} else {
		//Load the users selections and ask them to confirm
		fmt.Println("Config file found, loading selections!\n")
		//fmt.Println("ElvUI - Latest")
		selections = strings.Split(string(b), ",")
		printSelections(&selections)
		fmt.Println("\nWould you like to modify these selections?")
		res := *takeUserInput("Y/n")
		if res == "Y" {
			selections = strings.Split(*takeUserInput(selectionMessage), ",")
		}
	}
	processSelections(&selections)
}