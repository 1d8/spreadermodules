package main
import (
	"fmt"
	"strings"
	"os/exec"
	"flag"
)

//only works if Linux system has lsblk installed
func findUSBDevices(filepath string) {
	if len(filepath) == 0 {
		fmt.Println("[!] You must specify a filepath using the -p flag for the file to copy to the usb stick!")
		return
	}
	allUSBS, _ := exec.Command("lsblk").Output()
	username, _ := exec.Command("whoami").Output()
	usernameSplit := strings.Split(string(username), "\n")[0]
	if len(allUSBS) < 1 {
		fmt.Println("USB devices not found")
		return
	} else {
		allUSBSSplit := strings.Split(string(allUSBS), "\n")
		for i:=0; i<=len(allUSBSSplit)-1; i++ {
			//searches for username in lsblk path output to find usb sticks
			if strings.Index(strings.Split(allUSBSSplit[i], " ")[len(strings.Split(allUSBSSplit[i], " "))-1], usernameSplit) != -1 {
				fmt.Println("[+] USB path found!")
				fmt.Println("[+] Attempting to copy file from argument to USB path...")
				_, err := exec.Command("cp", filepath, strings.Split(allUSBSSplit[i], " ")[len(strings.Split(allUSBSSplit[i], " "))-1]).Output()
				if err != nil {
					fmt.Println(err)
					return
				} else {
					fmt.Printf("[+] Successfully copied %s to available usb stick!", filepath)
					return
				}
			}
			//fmt.Println(strings.Split(allUSBSSplit[i], " ")[len(strings.Split(allUSBSSplit[i], " "))-1])
		}
	}
	return
}


func main() {
	var filepath string
	flag.StringVar(&filepath, "p", "", "File path of file to copy to usb stick")
	flag.Parse()
	findUSBDevices(filepath)

}
