package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func getLinksUnformatted(httpBodyUnformatted string) []string {
	var prefix = "href=\"http"
	var prefixLen = len(prefix)
	var xvars = make([]string, 0, 40)
	var length = len(httpBodyUnformatted)
	var i2 = 0
	for i := 0; i < length-prefixLen; {
		if httpBodyUnformatted[i:i+prefixLen] == prefix {
			i2 = i + prefixLen
			for ; httpBodyUnformatted[i2] != '"'; i2++ {
				if i2 == length-1 {
					return xvars
				}
			}
			// Cut the href=" part.
			xvars = append(xvars, httpBodyUnformatted[i+(prefixLen-4):i2])
			i = i2 + 1
		} else {
			i++
		}
	}
	return xvars
}

func extractAllLinksFromBody(httpBody string) []string {
	var allLinks = make([]string, 0, 80)
	var lines = strings.Split(httpBody, "\n")
	for _, l := range lines {
		if strings.Index(l, "<a") != -1 {
			if httpLink := extractLink(l); len(httpLink) > 0 {
				allLinks = append(allLinks, httpLink)
			}
		}
	}
	return allLinks
}

func extractLink(httpStr string) string {
	res := ""
	i1 := strings.Index(httpStr, "\"http")
	var i2 int
	if i1 != -1 {
		for i2 = i1 + 1; httpStr[i2] != '"'; i2++ {
			if i2 == len(httpStr) {
				break
			}
		}
	}
	res = httpStr[i1+1 : i2]
	return res
}

/// includeQuery checks if string is in strings slice.
func includeQuery(pile []string, needle string) bool {
	for _, elem := range pile {
		if elem == needle {
			return true
		}
	}
	return false
}

func getVersions(links []string) []string {
	var versions = make([]string, 0, 30)
	prefix := "linux-"
	postfix := ".tar.xz"
	var i1 int
	for _, text := range links {
		i1 = strings.Index(text, prefix)
		if i1 != -1 && strings.Index(text, postfix) != -1 {
			versions = append(versions, text[i1+len(prefix):len(text)-len(postfix)])
		}
	}
	return versions
}

func isErr(err error) bool {
	return err != nil
}

// ShowUsage prints short usage message.
func ShowUsage() {
	println("Usage:\n\thref_links <html source file>\n\tRead content from stdin: href_links -\ne.g. curl https://example.com/ | href_links -\nBy default utility tries read ./home.html")

}

func main() {
	document := "home.html"
	if len(os.Args) > 1 {
		document = os.Args[1]
	}

	var bytes []byte
	var err error

	if document == "-" {
		bytes, err = ioutil.ReadAll(os.Stdin)
	} else {
		bytes, err = ioutil.ReadFile(document)
	}

	if isErr(err) {
		fmt.Fprintln(os.Stderr, err.Error())
		ShowUsage()
		return
	}
	var links = make([]string, 0, 30)
	var check string

	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		if strings.Index(line, "tar.xz") != -1 {
			check = extractLink(line)
			if !includeQuery(links, check) {
				links = append(links, check)
			}
		}
	}

	var versions = getVersions(links)
	for i, link := range links {
		println(link, ":", versions[i])
	}
	for _, ahref := range getLinksUnformatted(string(bytes)) {
		println(ahref)
	}
}
