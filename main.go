package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
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
	OutputFile string `json:"output"`
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

func done(config config) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	netClient.Post(config.ServerAddr+"/done", "text/plain", bytes.NewBufferString(config.JobName))
}

func main() {
	configFile := flag.String("config", "", "specify the config file location")
	command := flag.String("command", "", "the command to run")
	flag.Parse()
	config := readFileContent(configFile)
	fmt.Println(config)
	go runAlive(config)
	runCommand(command)
	done(config)
}
