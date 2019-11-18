package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"syscall"

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
		println("Incorrect directory format: it should contain at least one kernel source dir, or be empty")
		println(err.Error())
		return
	}
	if !isUpdate {
		println("You have latest sources.")
		return
	}
	var cwd string
	cwd, err = os.Getwd()
	var latestLocalVersion = localstorage.GetLatestLocalVersion(config["kernels_dir"])
	println("New kernel source available")
	if err = os.Chdir(config["kernels_dir"]); err != nil {
		println(err.Error())
		return
	}
	if err = os.Mkdir(kInfo.Version, 0755); err != nil {
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
	// E.g. (/path/to/script.sh <local old version> </path/to/kernel_config> <path to downloaded archive>).
	var postGetScripts, exist = config["post_get_scripts"]
	if exist {
		println("Executing", postGetScripts)
		var binary string
		binary, err = exec.LookPath("bash")
		if err != nil {
			println(err.Error())
			return
		}

		var curPath string
		curPath, err = os.Getwd()
		if err != nil {
			println(err.Error())
			return
		}

		var args = []string{binary, postGetScripts, curPath, config["config_path"], latestLocalVersion}
		var env = os.Environ()

		err = syscall.Exec(binary, args, env)
		if err != nil {
			println(err.Error())
		}
	}
	os.Chdir(cwd)

}
