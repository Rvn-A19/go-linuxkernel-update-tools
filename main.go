package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"./kernelorgparser"
	"./localstorage"
	"./remote"
)

func main() {
	var err error

	config := localstorage.ParseConfigFile("./default.conf")
	if len(config) == 0 {
		fmt.Fprintln(os.Stderr, "No config file.")
		return
	}

	var html string
	html, err = remote.GetHTTPText(remote.KernelsSourceHost)
	if err != nil {
		println(err.Error())
		return
	}
	var kInfo = kernelorgparser.GetInformation(&html)
	var isUpdate bool
	isUpdate, err = localstorage.ShouldUpdate(kInfo.Version, config["kernels_dir"])
	if err != nil {
		println(err.Error())
		return
	}
	if isUpdate {
		var cwd string
		cwd, err = os.Getwd()
		println("New kernel source available")
		if err = os.Chdir(config["kernels_dir"]); err != nil {
			println(err.Error())
			return
		}
		if err = os.Mkdir(kInfo.Version, 0660); err != nil {
			println(err.Error())
			return
		}
		if err = os.Chdir(kInfo.Version); err != nil {
			println(err.Error())
			return
		}
		var filename = path.Base(kInfo.ArchiveLink)
		remote.DownloadFile(kInfo.ArchiveLink, filename)
		// We're done - run post-get scripts :
		// E.g. (/path/to/script.sh <version> </path/to/kernel_config> <path to downloaded archive>).
		var postGetScripts, exist = config["post_get_scripts"]
		if exist {
			println("Executing", postGetScripts)
			var path, _ = os.Getwd()
			var bashRun = exec.Command("bash", postGetScripts, kInfo.Version, config["config_path"], path)
			var binOut []byte
			binOut, err = bashRun.Output()
			if err != nil {
				println(err.Error())
			} else {
				println(string(binOut))
			}
		}
		os.Chdir(cwd)
	} else {
		println("You have latest sources")
	}
}
