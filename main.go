package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

type Package struct {
	Name    string `yaml:"name"`
	Install string `yaml:"install"`
}

type Config struct {
	Packages []Package `yaml:"packages"`
}

func parseConfig(path string) (config *Config, err error) {
	contents, err := ioutil.ReadFile(path)

	if err != nil {
		return
	}

	err = yaml.Unmarshal(contents, &config)

	if err != nil {
		return
	}

	return
}

func main() {
	config, err := parseConfig(os.Args[1])

	if err != nil {
		fmt.Println(err)
	}

	var commands []*exec.Cmd

	for _, pkg := range config.Packages {
		cmdArgs := strings.Split(pkg.Install, " ")
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		commands = append(commands, cmd)

		fmt.Println("Install package: ", pkg.Name)
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(commands))

	for _, cmd := range commands {
		go func(cmd *exec.Cmd) {
			defer wg.Done()
			results, _ := cmd.CombinedOutput()
			fmt.Println(string(results))
		}(cmd)
	}

	wg.Wait()

	fmt.Println("Packages installed successfully!!")
}
