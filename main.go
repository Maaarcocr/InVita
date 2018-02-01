package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func runCommand(command *string) {
	commandSplitted := strings.Split(*command, " ")
	cmd := exec.Command(commandSplitted[0], commandSplitted[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		//TODO
		log.Fatal(err)
	}
}

type config struct {
	JobName    string `json:"name"`
	ServerAddr string `json:"addr"`
}

func readFileContent(configFile *string) config {
	raw, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	var config config
	json.Unmarshal(raw, &config)
	return config
}

func runAlive(config config) {
	for {
		time.Sleep(time.Second * 5)
		var netClient = &http.Client{
			Timeout: time.Second * 10,
		}
		netClient.Post(config.ServerAddr+"/alive", "text/plain", bytes.NewBufferString(config.JobName))
	}
}

func copyFile(src, dest string) {
	from, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
}

func done(config config, workDir, persistenDir *string) {
	if _, err := os.Stat(path.Join(*workDir, "output.json")); !os.IsNotExist(err) {
		copyFile(path.Join(*workDir, "output.json"), path.Join(*persistenDir, config.JobName+".output.json"))
	}

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	netClient.Post(config.ServerAddr+"/done", "text/plain", bytes.NewBufferString(config.JobName))
}

func initialise(wordDir *string, configFile *string) {
	copyFile(*configFile, path.Join(*wordDir, "input.json"))
}

func main() {
	configFile := flag.String("config", "", "specify the config file location")
	command := flag.String("command", "", "the command to run")
	workDir := flag.String("workdir", "", "The working directory")
	persistentDir := flag.String("persistentDir", "", "The persisten directory")
	flag.Parse()
	initialise(workDir, configFile)
	config := readFileContent(configFile)
	go runAlive(config)
	runCommand(command)
	done(config, workDir, persistentDir)
}
