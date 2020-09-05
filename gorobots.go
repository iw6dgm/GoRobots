package main

import (
	"GoRobots/count"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	// Version is the software version format v#.##-timestamp
	Version = "v1.0-20200904"
	// Separator is the OS dependent path separator
	Separator = string(os.PathSeparator)
	// RobotBinaryExt is the file extension of the compiled binary robot
	RobotBinaryExt = ".ro"
	// RobotSourceExt is the file extension of the robot source code
	RobotSourceExt = ".r"
)

var (
	// NumCPU is the number of detected CPUs/cores
	NumCPU int = runtime.NumCPU()
)

type match struct {
	Robots []string
}

type process struct {
	ID    int
	Match match
}

type processor struct {
	ID               int
	ProcessorChannel chan chan process // used to communicate between dispatcher and workers
	Channel          chan process
	End              chan bool
}

// start processor
func (p *processor) Start(crobotsExecutable string, opt string, tot int, result *sync.Map) {
	go func() {
		for {
			p.ProcessorChannel <- p.Channel // when the processor is available place channel in queue
			select {
			case process := <-p.Channel: // processor has received job
				process.Match.processCrobotsMatch(crobotsExecutable, opt, tot, result) // do work
			case <-p.End:
				return
			}
		}
	}()
}

// end processor
func (p *processor) Stop() {
	log.Printf("worker [%d] is stopping\n", p.ID)
	p.End <- true
}

// ProcessorChannel is a channel of available processors (a processor is a channel of processes)
var ProcessorChannel = make(chan chan process)

// Collector handles jobs to be sent to processors
type Collector struct {
	Process chan process // receives jobs to send to processors
	End     chan bool    // when receives bool stops processors
}

// StartDispatcher starts processors and return their Collector
func StartDispatcher(processorCount int, crobotsExecutable string, opt string, tot int, result *sync.Map) Collector {
	var i int
	var processors []processor
	input := make(chan process) // channel to recieve process
	end := make(chan bool)      // channel to spin down processors
	collector := Collector{Process: input, End: end}

	for i < processorCount {
		i++
		log.Println("Starting processor:", i)
		processor := processor{
			ID:               i,
			Channel:          make(chan process),
			ProcessorChannel: ProcessorChannel,
			End:              make(chan bool),
		}
		processor.Start(crobotsExecutable, opt, tot, result)
		processors = append(processors, processor) // stores processor
	}

	// start collector
	go func() {
		for {
			select {
			case <-end:
				for _, p := range processors {
					p.Stop() // stop processor
				}
				return
			case process := <-input:
				processor := <-ProcessorChannel // wait for available channel
				processor <- process            // dispatch process to processor
			}
		}
	}()

	return collector
}

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

func generateCombinations(list []string, size int) <-chan match {
	c := make(chan match)
	tot := len(list)
	go func(c chan match) {
		defer close(c)
		switch size {
		case 2:
			for i := 0; i < tot-2; i++ {
				for j := i + 1; j < tot-1; j++ {
					c <- match{Robots: []string{list[i], list[j]}}
				}
			}
		case 3:
			for i := 0; i < tot-3; i++ {
				for j := i + 1; j < tot-2; j++ {
					for k := j + 1; k < tot-1; k++ {
						c <- match{Robots: []string{list[i], list[j], list[k]}}
					}
				}
			}
		case 4:
			for i := 0; i < tot-4; i++ {
				for j := i + 1; j < tot-3; j++ {
					for k := j + 1; k < tot-2; k++ {
						for z := k + 1; z < tot-1; z++ {
							c <- match{Robots: []string{list[i], list[j], list[k], list[z]}}
						}
					}
				}
			}
		default:
			log.Fatal("Invalid size", size)
		}
	}(c)
	return c
}

func printRobots(tot int, result *sync.Map) {
	var bots []count.Robot

	result.Range(func(key interface{}, value interface{}) bool {
		r, _ := result.Load(key)
		robot := r.(count.Robot)
		if robot.Games > 0 {
			ties := 0
			for _, v := range robot.Ties {
				ties += v
			}
			robot.Points = robot.Wins*tot + ties
			robot.Eff = 100.0 * float32(robot.Points) / float32(tot*robot.Games)
		}
		bots = append(bots, robot)
		return true
	})
	sort.SliceStable(bots, func(i, j int) bool {
		return bots[i].Eff > bots[j].Eff
	})
	var i int = 0
	fmt.Println("#\tName\t\tGames\t\tWins\t\tTies2\t\tTies3\t\tTies4\t\tLost\t\tPoints\t\tEff%")
	for _, robot := range bots {
		i++
		fmt.Printf("%d\t%s\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\t\t%.3f\n", i, robot.Name, robot.Games, robot.Wins, robot.Ties[0], robot.Ties[1], robot.Ties[2], robot.Games-robot.Wins-(robot.Ties[0]+robot.Ties[1]+robot.Ties[2]), robot.Points, robot.Eff)
	}
}

func (m match) executeCrobotsMatch(exe string, opt string, n int) *exec.Cmd {
	switch n {
	case 2:
		return exec.Command(exe, opt, "-l200000", m.Robots[0], m.Robots[1])
	case 3:
		return exec.Command(exe, opt, "-l200000", m.Robots[0], m.Robots[1], m.Robots[2])
	case 4:
		return exec.Command(exe, opt, "-l200000", m.Robots[0], m.Robots[1], m.Robots[2], m.Robots[3])
	default:
		log.Fatal("Invalid size", n)
	}
	return nil
}

func (m match) processCrobotsMatch(crobotsExecutable string, opt string, tot int, result *sync.Map) {
	var out bytes.Buffer
	cmd := m.executeCrobotsMatch(crobotsExecutable, opt, tot)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	partial := count.ParseLogs(strings.Split(out.String(), "\n"))
	for _, robot := range partial {
		name := robot.Name
		if u, found := result.Load(name); found {
			update := u.(count.Robot)
			update.Games += robot.Games
			update.Wins += robot.Wins
			update.Ties[0] += robot.Ties[0]
			update.Ties[1] += robot.Ties[1]
			update.Ties[2] += robot.Ties[2]
			result.Store(name, count.Robot{Name: name, Games: update.Games, Wins: update.Wins, Ties: update.Ties, Points: 0, Eff: 0.0})
		} else {
			result.Store(name, count.Robot{Name: name, Games: robot.Games, Wins: robot.Wins, Ties: robot.Ties, Points: 0, Eff: 0.0})
		}
	}
}

func main() {

	log.Println("GoRobots", Version)
	log.Println("Detected CPU(s)/core(s):", NumCPU)
	runtime.GOMAXPROCS(NumCPU)
	tournamentType := flag.String("type", "", "tournament type: f2f, 3vs3 or 4vs4")
	configFile := flag.String("config", "config.yml", "YAML configuration file")
	parseLog := flag.String("parse", "", "parse log file")
	crobotsExecutable := flag.String("exe", "crobots", "Crobots executable")

	flag.Parse()

	if _, ok := schema[*tournamentType]; !ok {
		log.Fatalln("Invalid tournament type: ", *tournamentType)
	}

	tot, _ := schema[*tournamentType]

	if *parseLog != "" {
		content := logToString(*parseLog)
		result := count.ParseLogs(strings.Split(content, "\n"))
		var syncResult sync.Map
		for k, v := range result {
			syncResult.Store(k, v)
		}

		printRobots(tot, &syncResult)
		return
	}

	config := loadConfig(*configFile)
	listSize := len(config.ListRobots)

	if listSize < tot {
		log.Fatal("Robot list insufficient!")
	}

	var robots []string

	for _, r := range config.ListRobots {
		t := config.SourcePath + Separator + r + RobotBinaryExt

		if !fileExists(t) {
			log.Println("Binary robot cannot be found. Try to compile source code:", t)
			s := config.SourcePath + Separator + r + RobotSourceExt
			if fileExists(s) {
				var out bytes.Buffer
				cmd := exec.Command(*crobotsExecutable, "-c", s)
				cmd.Stdout = &out
				err := cmd.Run()
				if err != nil {
					log.Fatal(err)
				}
			} else {
				log.Fatal("Robot source code cannot be found:", s)
			}
		}

		robots = append(robots, t)
	}

	var opt string
	switch tot {
	case 2:
		opt = fmt.Sprintf("-m%d", config.MatchF2F)
	case 3:
		opt = fmt.Sprintf("-m%d", config.Match3VS3)
	case 4:
		opt = fmt.Sprintf("-m%d", config.Match4VS4)
	}
	log.Println("Start processing...")
	start := time.Now()
	var result sync.Map
	var i int
	collector := StartDispatcher(NumCPU, *crobotsExecutable, opt, tot, &result)
	for r := range generateCombinations(robots, tot) {
		i++
		collector.Process <- process{ID: i, Match: r}
	}
	duration := time.Since(start)
	log.Println("Completed in", duration)
	printRobots(tot, &result)
}
