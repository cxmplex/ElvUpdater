package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var baseUrl = "https://www.tukui.org/addons.php"
var addons map[string]string
var wowPath = "C:/Program Files (x86)/World of Warcraft/Interface/AddOns/"

func parseLinks(body *io.ReadCloser, req string, sep string) map[string]string {
	var links map[string]string
	links = make(map[string]string)

	t := html.NewTokenizer(*body)
	for {
		tokenType := t.Next()
		switch tokenType {
		case html.ErrorToken:
			return links
		case html.StartTagToken, html.EndTagToken:
			token := t.Token()
			//Every addon is a clickable link
			if "a" == token.Data {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						//dynamic selector
						if strings.Contains(attr.Val, req) {
							//dynamic separator
							id := strings.Split(attr.Val,sep)[1]

							//go to the next token (should be the a tag text
							t.Next()

							//convert []byte to string
							text := string(t.Text())

							//avoid nil assignment in map
							if len(text) == 0 {
								continue
							}
							links[id] = string(text)
						}
					}
				}
			}
		}
	}
	return links
}

func getAddonInfo(id *string) *string {
	//all addons follow this var scheme
	resp, err := http.Get(baseUrl + "?id=" + *id)
	checkErr(err)

	defer resp.Body.Close()

	//compile regex to parse version
	re, err := regexp.Compile(`The latest version of this addon is <b class="VIP">([0-9\.]+)`)
	checkErr(err)

	text, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	res := re.FindAllStringSubmatch(string(text), -1)
	return &res[0][1]
}

func init() {
	fmt.Println("ElvUI Updater v1.0")
	fmt.Println("Github.com/complexitydev\n")
	//fetch web page contents
	resp, err := http.Get(baseUrl)
	checkErr(err)

	defer resp.Body.Close()

	//parse every link that contains addons.php?id=\d+
	addons = parseLinks(&resp.Body, "addons.php?", "?id=")

	fmt.Println("Available addons:\n")
	fmt.Println("ElvUI - Latest")
	for k, v := range addons {
		fmt.Print("ID: " + k + " | ")
		fmt.Println(v)
	}
	fmt.Println()
}