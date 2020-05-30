package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/radovskyb/watcher"

	input "github.com/tcnksm/go-input"
)

// Every release we should increase this.
const version = "0.0.2"

// Vars to the main package
var kbID = ""
var apiKey = ""
var baseURL = "https://app.knowledgeowl.com"
var snippets []Snippet

//
// Starting point of the app.
//
func main() {
	// Info version handy for debugging
	Info("Using version: " + version)
	Info("")

	// This asks the user all the API keys and such.
	bootUp()

	// This downloads all the snippets in the DB.
	getSnippets()

	// This watching the directory for changes, and then uploads them to KO
	dirWatch()
}

//
// bootUp will ask the user the questions we need to get started.
//
func bootUp() {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	// Get API key
	a, err := ui.Ask("What is your API key?", &input.Options{
		Default:  apiKey,
		Required: true,
		Loop:     true,
		Hide:     true,
	})

	if err != nil {
		panic(err)
	}

	apiKey = a

	// Get project ID
	p, err := ui.Ask("What is your knowledge base id?", &input.Options{
		Default:  kbID,
		Required: true,
		Loop:     true,
	})

	if err != nil {
		panic(err)
	}

	kbID = p

	// Server address
	b, err := ui.Ask("What is the url to your knowledge base API?", &input.Options{
		Default:  baseURL,
		Required: true,
		Loop:     true,
	})

	if err != nil {
		panic(err)
	}

	baseURL = b
}

//
// dirWatch will watch for changes of files. Once a file changes it is uploaded to the server.
//
func dirWatch() {
	w := watcher.New()

	// SetMaxEvents to 1 to allow at most 1 event's to be received
	// on the Event channel per watching cycle.
	//
	// If SetMaxEvents is not set, the default is to send all events.
	w.SetMaxEvents(1)

	// Only notify on file write.
	w.FilterOps(watcher.Write)

	// We only want files that are *.html as they are the only files are are watching.
	r := regexp.MustCompile("html$")
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	// Go routine to watch for file changes.
	go func() {
		for {
			select {
			case event := <-w.Event:
				//Info(event) // Print the event's info.
				uploadModifedFile(event.Path)
			case err := <-w.Error:
				Info(err.Error())
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch kbs recursively for changes.
	if err := w.AddRecursive("./kbs"); err != nil {
		log.Fatalln(err)
	}

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

//
// uploadModifedFile will upload any changes to a snippet to the backend server.
//
func uploadModifedFile(file string) {
	// File must containe

	// Read file contents
	body, err := ioutil.ReadFile(file)

	if err != nil {
		panic(err)
	}

	// Loop through and find this file. Then upload.
	for _, row := range snippets {
		if row.ProjectID == kbID && path.Base(file) == (row.Mergecode+".html") {
			doUpload(row.ID, string(body), file)
		}
	}
}

//
// doUpload will upload the content to the KB server
//
func doUpload(id string, bodyStr string, file string) {
	// Build Body data.
	json := []byte(`{"current_version": {"en": "` + bodyStr + `"}}`)
	body := bytes.NewBuffer(json)

	// Create client
	client := &http.Client{}

	// Create request
	url := fmt.Sprintf("%s/api/head/snippet/%s.json", baseURL, id)
	req, err := http.NewRequest("PUT", url, body)
	req.SetBasicAuth(apiKey, "X")
	req.Header.Add("Content-type", "application/json")

	// Fetch Request
	_, err = client.Do(req)

	if err != nil {
		Info("Failure : " + err.Error())
	}

	Info("Syncing: " + file)
}

//
// getSnippets will make an api call to the KB servers to get all the snippets.
//
func getSnippets() {
	// Make KB directory if it is not already there.
	_ = os.Mkdir("kbs", os.ModePerm)

	// Create client
	client := &http.Client{}

	// Create request
	url := fmt.Sprintf("%s/api/head/snippet.json?project_id=%s", baseURL, kbID)
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(apiKey, "X")
	req.Header.Add("Content-type", "application/json")

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		Info("Failure : " + err.Error())
	}

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)

	// Make sure this was a success
	if resp.StatusCode != 200 {
		Info("response Body : " + string(respBody))
		panic("Server failed to respond with a status code 200.")
	}

	// Setup response API object
	apiResponse := APISnippetResponse{}

	json.Unmarshal(respBody, &apiResponse)

	// TODO(spicer): check for error. But sometimes current_version is an {} and sometimes (likely a bug) it is []
	// if err := json.Unmarshal(respBody, &apiResponse); err != nil {
	// 	panic(err)
	// }

	//Info(String(respBody))

	// Set snippets
	snippets = apiResponse.Data

	for _, row := range snippets {
		//spew.Dump(row)
		//Info(row.CurrentVersion.En)
		// Touch the file
		dir := "kbs/" + row.ProjectID
		file := row.Mergecode + ".html"
		fullPath := dir + "/" + file
		os.MkdirAll(dir, os.ModePerm)

		// Write file out with content
		err = ioutil.WriteFile(fullPath, []byte(row.CurrentVersion.En), 0644)

		if err != nil {
			panic(err)
		}

		Info("Downloading: " + fullPath)
	}

	Info("")
}

//
// Info will Info information as the app is being used.
//
func Info(msg string) {
	fmt.Println(msg)
}

/* End File */
