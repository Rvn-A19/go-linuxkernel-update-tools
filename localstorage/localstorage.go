package localstorage

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// DefaultConfigName contains default path to config for update tool.
const (
	DefaultConfigName string = "./default.conf"
)

// ParseConfigFile gets all needed vars from config.
func ParseConfigFile(configFile string) map[string]string {
	var conf = make(map[string]string)
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

// CompareVersions compares version1 and version2 as string.
// E.g. "a.b.c" >= "d.e.f"?
func CompareVersions(v1 string, v2 string) (int, error) {
	var slice1, slice2 = strings.Split(v1, "."), strings.Split(v2, ".")
	var idMax = len(slice1)
	if len(slice1) > len(slice2) {
		idMax = len(slice2)
	}
	var err error
	var id1, id2 int
	for i := 0; i < idMax; i++ {
		id1, err = strconv.Atoi(slice1[i])
		if err != nil {
			return 0, err
		}
		id2, err = strconv.Atoi(slice2[i])
		if err != nil {
			return 0, err
		}
		if id1 > id2 {
			return 1, nil
		}
		if id2 < id1 {
			return -1, nil
		}
	}
	return 0, nil
}

// ShouldUpdate checks if new version should be downloaded.
func ShouldUpdate(newVersion, kernelDir string) (bool, error) {
	var err error
	var fi []os.FileInfo
	fi, err = ioutil.ReadDir(kernelDir)
	if err != nil {
		return false, err
	}
	if len(fi) == 0 {
		return true, nil
	}
	var verFlag int
	var i = 0
	var existingLastVersion = fi[i].Name()
	i = 1
	for ; i < len(fi); i++ {
		verFlag, err = CompareVersions(fi[i].Name(), existingLastVersion)
		if err == nil && verFlag == 1 {
			existingLastVersion = fi[i].Name()
		}
	}
	verFlag, err = CompareVersions(newVersion, existingLastVersion)
	var result = err == nil && verFlag == 1
	return result, err
}
