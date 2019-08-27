package remote

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

const (
	// KernelsSourceHost - kernel.org url.
	KernelsSourceHost = "https://www.kernel.org/"
)

func humanRepr(bytesCount int) string {
	var Kilo, Mega, Giga = 1024, 1024 * 1024, 1024 * 1024 * 1024
	if bytesCount >= Giga {
		return strconv.Itoa(bytesCount/Giga) + "G"
	}
	if bytesCount >= Mega {
		return strconv.Itoa(bytesCount/Mega) + "M"
	}
	if bytesCount >= Kilo {
		return strconv.Itoa(bytesCount/Kilo) + "K"
	}
	return strconv.Itoa(bytesCount) + " byte(s)"
}

// GetHTTPText reads http body from kernelsSourceHost.
func GetHTTPText(host string) (string, error) {
	var resp *http.Response
	var err error

	resp, err = http.Get(host)

	if err != nil {
		println(err.Error())
		return "", err
	}

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil

}

// DownloadFile saves file from remote link.
func DownloadFile(link string, filename string) error {
	var err error
	var resp *http.Response
	resp, err = http.Head(link)

	if err != nil {
		println(err.Error())
		return err
	}

	var contentLength int
	var cl string
	cl = resp.Header.Get("Content-Length")
	if len(cl) == 0 {
		println("No content length header")
		return http.ErrContentLength
	}

	contentLength, err = strconv.Atoi(cl)

	if err != nil {
		println(err.Error())
		return err
	}
	resp, err = http.Get(link)
	if err != nil {
		println(err.Error())
		return err
	}

	var MB4 = 1024 * 1024 * 4
	if contentLength < MB4 {
		println("saving", contentLength, "bytes to", filename)
		var bytes []byte
		bytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			println(err.Error())
			return err
		}
		ioutil.WriteFile(filename, bytes, 0600)
		return nil
	}
	var f *os.File
	f, err = os.Create(filename)
	if err != nil {
		println(err.Error())
		return err
	}
	defer f.Close()

	// Read by 512-K chunks.
	buf := make([]byte, 512*1024)
	var n int
	var total = 0
	var spaces = "               "
	for {
		n, err = io.ReadFull(resp.Body, buf)
		f.Write(buf[:n])
		total += n
		fmt.Printf("\rreceived %s of %s (%d %%)%s", humanRepr(total), humanRepr(contentLength), (int)(total/(contentLength/100)), spaces)
		if err == io.EOF || (err == io.ErrUnexpectedEOF && contentLength == total) {
			println("\nDone")
			return nil
		}
		if err != nil {
			println("\n", err.Error())
			return err
		}
	}
}
