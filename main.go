package main

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"os"
)

func main() {
	host := "192.168.1.86" // := -> declare and assign a variable at the same time
	user := "andweste"
	port := 22
	// Private key
	pemBytes, err := os.ReadFile("H:\\Coding\\Python Projects\\RaspberryPi\\private_key_ftp")
	signer, err := ssh.ParsePrivateKey(pemBytes)
	auths := []ssh.AuthMethod{ssh.PublicKeys(signer)} // normal Array declaration

	config := ssh.ClientConfig{
		User:            user,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	print(fmt.Printf("Connecting to: %s:%d\n", host, port))

	addr := fmt.Sprintf("%s:%d", host, port) // string formatting

	conn, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to [%s]: %v\n", addr, err)
		os.Exit(1)
	}
	defer conn.Close() // after the function finishes, this will be called. Defer moves it to the end
	// Defer is often used to ensure freeing up resources after things finish

	sftpClient, err := sftp.NewClient(conn) // create the sftp client
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create SFTP client: %v\n", err)
		os.Exit(1)
	}
	defer sftpClient.Close()

	filename := "ashleyw-3.jpg"
	download(sftpClient, fmt.Sprintf("/home/andweste/AshleyPictures/%s", filename),
		filename)
}

func download(sftpClient *sftp.Client, srcPath string, destinationPath string) {
	sourceFile, err := sftpClient.Open(srcPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open source file: %v\n", err)
		os.Exit(1)
	}
	defer sourceFile.Close()

	// Create the destination file to write onto
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open destination file: %v\n", err)
	}
	defer destinationFile.Close()

	destinationFile.ReadFrom(sourceFile)

}
