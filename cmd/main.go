package main

import (
	"log"
	"os"

	"github.com/mbndr/codeviewer"
	"gopkg.in/urfave/cli.v1"
)



func main() {
	// config dir must be specified
	codeviewer.ConfigDir = codeviewer.GetConfigDir()
	if codeviewer.ConfigDir == "" {
		log.Fatal("set $XDG_CONFIG_HOME or $HOME to specify config dir")
	}

	app := cli.NewApp()
	app.Name = "code-viewer"
	app.Usage = "view highlighted source code in the browser"

	app.Commands = []cli.Command{
		codeviewer.CmdDownload, // Download all hljs styles
		codeviewer.CmdServe, // Run webserver
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}