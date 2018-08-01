# Codeviewer
This is just a fun project to look at your source code syntax highlighted in the browser.

It uses [hljs](https://highlightjs.org/) and is written in Go and JavaScript.

![Screenshot](./scrot.png)

(The theme chooser can be toggled with `CTRL`)

## Features
* All hljs styles and languages
* Realtime switching the styles
* Once build you have no external dependencies
* Working offline once downloaded the highlightjs assets (`codeviewer download`)

## Usage
Currently you have to build your executable yourself.
```bash
git clone https://github.com/mbndr/codeviewer
go build -o codeviewer codeviewer/cmd/main.go
```
Done that you can move the executable to your path and use it.
```bash
# Download hljs assets (to $HOME/.config)
codeviewer download
# Start server als visit http://localhost:8080
codeviewer serve -f myFile.go
# For further options use
codeviewer -h
```