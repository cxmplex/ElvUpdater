package main

import (
	"fmt"
	"github.com/mholt/archiver"
	"io"
	"net/http"
	"os"
)

var wowPath = "C:/Program Files (x86)/World of Warcraft/_retail_/Interface/AddOns/"

func downloadFileAndExtract(file string, url string) {
	fmt.Println("Downloading: " + addons[file])
	//create file
	w, err := os.Create(wowPath + addons[file] + ".zip")
	checkErr(err)

	defer w.Close()

	resp, err := http.Get(url)
	checkErr(err)

	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)
	checkErr(err)

	fmt.Println(addons[file] + ".zip")
	fmt.Println(wowPath)
	err = archiver.Zip.Open(addons[file]+".zip", wowPath)
	checkErr(err)
}

func downloadSelections(selections *[]string) {
	downloadUrl := "https://www.tukui.org/addons.php?download="
	for _, id := range *selections {
		downloadFileAndExtract(id, downloadUrl+id)
	}
}
