package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

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

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("\n\n%s took %s\n\n", name, elapsed)
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
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(commands))

	for _, cmd := range commands {
		go func(cmd *exec.Cmd) {
			defer wg.Done()
			output, _ := cmd.CombinedOutput()
			fmt.Println(string(output))
			// startTime := time.Now()

			// err := cmd.Run()

			// if err != nil {
			// 	fmt.Printf("%s is failing: %s\n", cmd.Args[len(cmd.Args)-1], err)
			// 	output, _ := cmd.Output()
			// 	fmt.Println(string(output))
			// }

			// elapsed := time.Since(startTime)
			// cmdName := cmd.Args[len(cmd.Args)-1]
			// fmt.Printf("%s took %s\n", cmdName, elapsed)
		}(cmd)
	}

	wg.Wait()

	fmt.Println("Packages installed successfully!!")
}
