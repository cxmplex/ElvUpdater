package main

import (
	"fmt"
	"github.com/mholt/archiver"
	"io"
	"net/http"
	"os"
)

func downloadFileAndExtract(file string, url string) {
	fmt.Println("Downloading: " + addons[file])
	//create file
	w, err := os.Create(addons[file] + ".zip")
	checkErr(err)

	defer w.Close()

	resp, err := http.Get(url)
	checkErr(err)

	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)
	checkErr(err)

	err = archiver.Zip.Open(addons[file] + ".zip", wowPath)
	checkErr(err)
}

func downloadSelections(selections *[]string) {
	downloadUrl := "https://www.tukui.org/addons.php?download="
	for _, id := range *selections {
		downloadFileAndExtract(id, downloadUrl+id)
	}
}