package main

import (
	"io/ioutil"
	"strings"

	"./remote"
)

func extractVersion(tagline string) string {
	var result string
	var tagOpen = "<strong>"
	var tagClose = "</strong>"
	posStart := strings.Index(tagline, tagOpen)
	if posStart != -1 {
		posStart += len(tagOpen)
		posEnd := strings.Index(tagline, tagClose)
		if posEnd != -1 {
			result = tagline[posStart:posEnd]
		}
	}
	return result
}

func getArchiveLink(httpText *string, version string) string {
	posEnd := strings.Index(*httpText, version+".tar.xz")
	var posStart int
	for posStart = posEnd; (*httpText)[posStart:posStart+4] != "http"; posStart-- {
		if posStart == 0 {
			return ""
		}
	}
	return (*httpText)[posStart : posEnd+len(version)+7]
}

func parseConfigFile(configFile string) map[string]string {
	var conf map[string]string
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		println(err.Error())
		return conf
	}
	lines := strings.Split(string(bytes), "\n")
	var pos int
	for _, line := range lines {
		pos = strings.Index(line, "=")
		if pos != -1 {
			conf[line[:pos]] = line[pos+1:]
		}
	}
	return conf
}

func main() {
	httpSource, err := remote.GetHTTPText(remote.KernelsSourceHost)
	if err != nil {
		println(err)
		return
	}
	var ver, link string
	lines := strings.Split(httpSource, "\n")
	for idx, line := range lines {
		if strings.Index(line, "stable:") != -1 {
			ver = extractVersion(lines[idx+1])
			link = getArchiveLink(&httpSource, ver)
			println(ver, link)
			break
		}
	}
	if len(link) > 0 {
		//filename := link[strings.LastIndex(link, "/")+1 : len(link)]
		//remote.DownloadFile(link, filename)
		remote.DownloadFile("https://www.farmanager.com/files/Far30b5454.x86.20190823.msi", "FarManager.msi")
	}
	config := parseConfigFile("./default.conf")
	configPath, exists := config["config_path"]
	if exists {
		println(configPath)
	}
}

