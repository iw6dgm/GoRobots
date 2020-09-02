package main

import (
	"GoRobots/count"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	Version        = "v1.0a 02/09/2020"
	Separator      = string(os.PathSeparator)
	RobotBinaryExt = ".ro"
	RobotSourceExt = ".r"
)

type tournamentConfig struct {
	Label      string   `yaml:"label"`
	MatchF2F   int      `yaml:"matchF2F"`
	Match3VS3  int      `yaml:"match3VS3"`
	Match4VS4  int      `yaml:"match4VS4"`
	SourcePath string   `yaml:"sourcePath"`
	ListRobots []string `yaml:"listRobots"`
}

var schema = map[string]int{"f2f": 2, "3vs3": 3, "4vs4": 4}

func loadConfig(config string) tournamentConfig {
	f, err := os.Open(config)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var cfg tournamentConfig
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func logToString(file string) string {
	// Read entire file content, giving us little control but
	// making it very simple. No need to close the file.
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string
	return string(content)
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {

	log.Println("GoRobots", Version)

	tournamentType := flag.String("type", "all", "tournament type: all, f2f, 3vs3 or 4vs4")
	configFile := flag.String("config", "config.yml", "YAML configuration file")
	parseLog := flag.String("parse", "", "parse log file")

	flag.Parse()

	if _, ok := schema[*tournamentType]; !ok && *tournamentType != "all" {
		log.Fatalln("Invalid tournament type: ", *tournamentType)
	}

	if *parseLog != "" {

		if *tournamentType == "all" {
			log.Fatalln("Cannot parse log for all tournament types")
			return
		}
		tot := schema[*tournamentType]

		content := logToString(*parseLog)
		result := count.ParseLogs(strings.Split(content, "\n"))
		var robots []count.Robot

		for _, robot := range result {

			if robot.Games > 0 {
				ties := 0
				for _, v := range robot.Ties {
					ties += v
				}
				robot.Points = robot.Wins*tot + ties
				robot.Eff = 100.0 * float32(robot.Points) / float32(tot*robot.Games)
			}
			robots = append(robots, robot)
		}

		fmt.Printf("result: %v\n", robots)
		return
	}

	config := loadConfig(*configFile)

	listSize := len(config.ListRobots)

	if listSize < 1 {
		log.Fatal("Robot list empty!")
	}

	var robots []string

	for _, r := range config.ListRobots {
		t := config.SourcePath + Separator + r + RobotBinaryExt
		/*
			if !fileExists(t) {
				log.Fatal(fmt.Sprintf("Robot %s cannot be found", t))
			}*/

		robots = append(robots, t)
	}

	log.Println("Robots", robots)

	// cmd := exec.Command("crobots", "-m10", "-l200000", "/home/joshua/crobots/bench.r", "/home/joshua/crobots/bench.r")
	// var out bytes.Buffer
	// cmd.Stdout = &out
	// err := cmd.Run()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// result := count.ParseLogs(strings.Split(out.String(), "\n"))

	// fmt.Printf("result: %v\n", result)
}
