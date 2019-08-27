package kernelorgparser

import (
	"strings"
)

// VersionInfo provides information about latest kernel source archives.
type VersionInfo struct {
	Version     string
	ArchiveLink string
}

// ExtractVersion finds latest kernel version in http source.
func ExtractVersion(tagline string) string {
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

// GetArchiveLink founds link to archive by version in text.
func GetArchiveLink(httpText *string, version string) string {
	posEnd := strings.Index(*httpText, version+".tar.xz")
	var posStart int
	for posStart = posEnd; (*httpText)[posStart:posStart+4] != "http"; posStart-- {
		if posStart == 0 {
			return ""
		}
	}
	return (*httpText)[posStart : posEnd+len(version)+7]
}

// GetInformation extracts data from http source.
func GetInformation(httpText *string) VersionInfo {
	var xInfo VersionInfo
	lines := strings.Split(*httpText, "\n")
	for idx, line := range lines {
		if strings.Index(line, "stable:") != -1 {
			xInfo.Version = ExtractVersion(lines[idx+1])
			xInfo.ArchiveLink = GetArchiveLink(httpText, xInfo.Version)
			break
		}
	}
	return xInfo
}
