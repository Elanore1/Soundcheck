/*package main

import "fmt"

func main() {
	fmt.Println("ouverture son")
	ffmpeg -i oli_lea.mp3 -af "volumedetected" -vn -sn -f null /dev/null

	fmt.Println("fin ouverture son")
}*/
//ffmpeg -i oli_lea.mp3 -af "volumedetected" -vn -sn -f null /dev/null

package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	file, err := os.Open("oli_lea2.mp3") // open file
	check(err)

	defer file.Close()
	buf, err := ioutil.ReadAll(file)
	check(err)

	cmd := exec.Command("ffmpeg", "-y", // Yes to all
		//"-hide_banner", "-loglevel", "panic", // Hide all logs
		"-i", "pipe:0", // take stdin as input
		"-map_metadata", "-1", // strip out all (mostly) metadata
		"-c:a", "libmp3lame", // use mp3 lame codec
		"-vsync", "2", // suppress "Frame rate very high for a muxer not efficiently supporting it"
		"-b:a", "128k", // Down sample audio birate to 128k
		"-f", "mp3", // using mp3 muxer (IMPORTANT, output data to pipe require manual muxer selecting)
		"pipe:1", // output to stdout
	)

	resultBuffer := bytes.NewBuffer(make([]byte, 5*1024*1024)) // pre allocate 5MiB buffer

	cmd.Stderr = os.Stderr    // bind log stream to stderr
	cmd.Stdout = resultBuffer // stdout result will be written here

	stdin, err := cmd.StdinPipe() // Open stdin pipe
	check(err)

	err = cmd.Start() // Start a process on another goroutine
	check(err)

	_, err = stdin.Write(buf) // pump audio data to stdin pipe
	check(err)

	err = stdin.Close() // close the stdin, or ffmpeg will wait forever
	check(err)

	err = cmd.Wait() // wait until ffmpeg finish
	check(err)

	outputFile, err := os.Create("out2.mp3") // create new file
	check(err)
	defer outputFile.Close()

	_, err = outputFile.Write(resultBuffer.Bytes()) // write result buffer to file
	check(err)
}
