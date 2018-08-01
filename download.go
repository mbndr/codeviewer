package codeviewer

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/urfave/cli.v1"
)

const (
	// get all file names
	apiUrl = "https://data.jsdelivr.com/v1/package/gh/highlightjs/cdn-release@9.12.0/flat"
	// download the files
	cdnUrl = "https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@9.12.0/"

	// added to cdnUrl
	cdnLanguagePath = "languages"
	cdnStylesPath   = "styles"
)

// CmdDownload downloads all js and css files for hljs
var CmdDownload = cli.Command{
	Name:   "download",
	Usage:  "Download all hljs styles and languages from jsdelivr to local cache",
	Action: download,
}

func download(c *cli.Context) error {
	var err error

	stylePath := filepath.Join(ConfigDir, StyleDir)
	langPath := filepath.Join(ConfigDir, LangDir)

	//  create subfolders
	if err = os.MkdirAll(stylePath, 0777); err != nil {
		return err
	}
	if err = os.MkdirAll(langPath, 0777); err != nil {
		return err
	}

	// get all file names
	var files apiFiles

	resp, err := http.Get(apiUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("cannot call api with status code 200")
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&files)
	if err != nil {
		return err
	}

	// iter through all files and save them in spec location
	for _, file := range files.Files {
		err = downloadFile(file.Name)
		if err != nil {
			log.Println("error downloading " + file.Name + ": " + err.Error())
		}
	}

	return nil
}

// downloadFile downloads a single file from jsdelivr
// TODO work with urls and path.Join
func downloadFile(name string) error {
	// only download build files
	if !strings.HasPrefix(name, "/build/") {
		return nil
	}

	// download file
	resp, err := http.Get(cdnUrl + name)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// create file
	f, err := os.Create(filepath.Join(ConfigDir, strings.TrimPrefix(name, "/build/")))
	if err != nil {
		return err
	}
	defer f.Close()

	// save file
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// this is gotten by one request to the jsdelivr api
type apiFiles struct {
	Files []struct { // one object is one file
		Name string `json:"name"`
	} `json:"files"`
}
