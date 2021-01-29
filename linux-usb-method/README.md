# Short Explanation

This usb module only works for Linux and will enumerate plugged in USB sticks via the `lsblk` command. It will search for usb sticks based on the path for usb sticks normally being something along the lines of */media/user/*. If it finds the current user's username in any part of the output of `lsblk` (which would indicate a usb stick being plugged in), it will copy the specified file to that path, essentially copying the file to the usb stick.

# Usage

Simply run it from the command line and pass the full filepath of the file you wish to copy to the usb stick to the **-p** flag:

`./binaryname -p /home/user/Documents/malware.dat`


# Compilation

`GOOS=linux GOARCH=386 go build filename`

or 

`GOOS=linux GOARCH=amd64 go build filename`

Additional info for cross-compilation [here](https://dh1tw.de/2019/12/cross-compiling-golang-cgo-projects/)
