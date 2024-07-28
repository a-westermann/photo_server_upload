package main

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"os"
	"time"
)

// Sleep 3 min until find a file in the local directory
// Copy it over to the server
// Delete the local file if it is successful

var localDir = "C:\\Users\\Andrew\\Code Projects\\Pi_Photo_Saver\\photo_saver_server\\UploadPhotos\\"
var host = "192.168.1.86" // := -> declare and assign a variable at the same time
var user = "andweste"
var port = 22
var privateKeyPath = "C:\\Users\\Andrew\\Code Projects\\Pi_Photo_Saver\\photo_saver_server\\private_key.ppk"
var uploadPath = "\\home\\andweste\\AshleyUploadPhotos\\"
var logFilePath = "logfile.txt"
var logFile *os.File

func main() {
	// Try to get Stat on the file. If error, it means it doesn't exist. So create it
	if _, err := os.Stat(logFilePath); err != nil {
		_, err := os.Create(logFilePath)
		if err != nil {
			panic(err)
		}
	}

	// Open the log file w/ appending privileges
	logFile, _ = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	defer logFile.Close()
	logFile.Write([]byte(fmt.Sprintf("\nbeginning program: %s", time.DateTime)))

	pollFolder()
}

func pollFolder() {
	// Check if there are any files to copy to the server
	files := check_for_files(localDir)
	if files != nil && len(files) > 0 {
		logFile.Write([]byte(fmt.Sprintf("\nFound photos to upload: %s", time.DateTime)))
		sendFiles(files)
	}

	logFile.Write([]byte(fmt.Sprintf("\nNothing to upload...: %s", time.DateTime)))
	time.Sleep(3 * time.Minute)

	pollFolder()
}

func sendFiles(files []os.DirEntry) {
	pemBytes, err := os.ReadFile(privateKeyPath)
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

	for _, file := range files {
		upload(sftpClient, fmt.Sprintf("%s%s", localDir, file.Name()),
			fmt.Sprintf("%s%s", uploadPath, file.Name()))

		//download(sftpClient, fmt.Sprintf("/home/andweste/AshleyPictures/%s", filename),
		//	filename)
	}
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

func upload(sftpClient *sftp.Client, srcPath string, destinationPath string) {
	sourceFile, err := os.Open(srcPath)
	if err != nil {
		logFile.Write([]byte(fmt.Sprintf("\nerror getting source file: %s", err)))
	}

	destinationFile, err := sftpClient.Create(fmt.Sprintf("%s%s", uploadPath, sourceFile.Name()))
	if err != nil {
		logFile.Write([]byte(fmt.Sprintf("\nerror creating destination file: %s", err)))
	}
	defer destinationFile.Close()

	destinationFile.ReadFrom(sourceFile)
}

func check_for_files(path string) []os.DirEntry {
	files, _ := os.ReadDir(path)
	return files
}
