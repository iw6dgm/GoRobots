package main

import (
	"GoRobots/count"
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
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
	Version = "v1.22-20200907"
	// Separator is the OS dependent path separator
	Separator = string(os.PathSeparator)
	// RobotBinaryExt is the file extension of the compiled binary robot
	RobotBinaryExt = ".ro"
	// RobotSourceExt is the file extension of the robot source code
	RobotSourceExt = ".r"
	// Header is the output header
	Header = "#\tName\t\tGames\t\tWins\t\tTies2\t\tTies3\t\tTies4\t\tLost\t\tPoints\t\tEff%"
	// OutputFormat is a single row output format
	OutputFormat = "%d\t%s\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\t\t%.3f\n"
)

var (
	// NumCPU is the number of detected CPUs/cores
	NumCPU int = runtime.NumCPU()
)

type match struct {
	Robots []string
}

type result struct {
	m sync.Mutex
	r map[string]count.Robot
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

// commandExists checks if an executable exists
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func generateCombinations(list []string, size int, c chan<- match) {
	tot := len(list)
	switch size {
	case 2:
		for i := 0; i < tot-1; i++ {
			for j := i + 1; j < tot; j++ {
				c <- match{Robots: []string{list[i], list[j]}}
			}
		}
	case 3:
		for i := 0; i < tot-2; i++ {
			for j := i + 1; j < tot-1; j++ {
				for k := j + 1; k < tot; k++ {
					c <- match{Robots: []string{list[i], list[j], list[k]}}
				}
			}
		}
	case 4:
		for i := 0; i < tot-3; i++ {
			for j := i + 1; j < tot-2; j++ {
				for k := j + 1; k < tot-1; k++ {
					for z := k + 1; z < tot; z++ {
						c <- match{Robots: []string{list[i], list[j], list[k], list[z]}}
					}
				}
			}
		}
	default:
		log.Fatal("Invalid size", size)
	}
}

func generateCombinationsForBenchmark(robot string, list []string, size int, c chan<- match) {
	tot := len(list)
	switch size {
	case 2:
		for i := 0; i < tot-1; i++ {
			c <- match{Robots: []string{list[i], robot}}
		}
	case 3:
		for i := 0; i < tot-2; i++ {
			for j := i + 1; j < tot-1; j++ {
				c <- match{Robots: []string{list[i], list[j], robot}}
			}
		}
	case 4:
		for i := 0; i < tot-3; i++ {
			for j := i + 1; j < tot-2; j++ {
				for k := j + 1; k < tot-1; k++ {
					c <- match{Robots: []string{list[i], list[j], list[k], robot}}
				}
			}
		}
	default:
		log.Fatal("Invalid size", size)
	}
}

func check(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}

func printRobots(out *string, tot int, result map[string]*count.Robot) {
	var bots []count.Robot

	for _, robot := range result {
		if robot.Games > 0 {
			ties := 0
			for _, v := range robot.Ties {
				ties += v
			}
			robot.Points = robot.Wins*tot + ties
			robot.Eff = 100.0 * float32(robot.Points) / float32(tot*robot.Games)
		}
		bots = append(bots, *robot)
	}
	sort.SliceStable(bots, func(i, j int) bool {
		return bots[i].Eff > bots[j].Eff
	})
	printToStd := func() {
		var i int = 0
		fmt.Println(Header)
		for _, robot := range bots {
			i++
			fmt.Printf(OutputFormat, i, robot.Name, robot.Games, robot.Wins, robot.Ties[0], robot.Ties[1], robot.Ties[2], robot.Games-robot.Wins-(robot.Ties[0]+robot.Ties[1]+robot.Ties[2]), robot.Points, robot.Eff)
		}
	}
	if *out != "" {
		outputFile := *out
		f, err := os.Create(outputFile)
		if check(err) {
			printToStd()
			return
		}
		defer f.Close()
		w := bufio.NewWriter(f)
		var i int = 0
		_, err = fmt.Fprintln(w, Header)
		if check(err) {
			printToStd()
			return
		}
		for _, robot := range bots {
			i++
			_, err = fmt.Fprintf(w, OutputFormat, i, robot.Name, robot.Games, robot.Wins, robot.Ties[0], robot.Ties[1], robot.Ties[2], robot.Games-robot.Wins-(robot.Ties[0]+robot.Ties[1]+robot.Ties[2]), robot.Points, robot.Eff)
			if check(err) {
				printToStd()
				return
			}
		}
		w.Flush()
	} else {
		printToStd()
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

func (m match) processCrobotsMatch(crobotsExecutable string, opt string, tot int, result map[string]*count.Robot, mutex *sync.Mutex) {
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
		// sync.Map doesn't seem to work here
		mutex.Lock()
		if _, found := result[name]; found {
			result[name].Games += robot.Games
			result[name].Wins += robot.Wins
			result[name].Ties[0] += robot.Ties[0]
			result[name].Ties[1] += robot.Ties[1]
			result[name].Ties[2] += robot.Ties[2]
		} else {
			result[name] = &count.Robot{Name: name, Games: robot.Games, Wins: robot.Wins, Ties: robot.Ties, Points: 0, Eff: 0.0}
		}
		mutex.Unlock()
	}
}

func worker(id int, matches <-chan match, crobotsExecutable string, opt string, tot int, result map[string]*count.Robot, mutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	for m := range matches {
		m.processCrobotsMatch(crobotsExecutable, opt, tot, result, mutex)
	}
}

func checkAndCompile(r string, crobotsExecutable string, path func(r string) string) string {
	t := path(r) + RobotBinaryExt
	if !fileExists(t) {
		log.Println("Binary robot cannot be found:", t, "Trying to compile source code")
		s := path(r) + RobotSourceExt
		if fileExists(s) {
			var out bytes.Buffer
			cmd := exec.Command(crobotsExecutable, "-c", s)
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("Robot source code cannot be found: ", s)
		}
	}
	return t
}

func main() {

	log.Println("GoRobots", Version)
	log.Println("Detected CPU(s)/core(s):", NumCPU)
	runtime.GOMAXPROCS(NumCPU)
	tournamentType := flag.String("type", "", "tournament type: f2f, 3vs3 or 4vs4")
	configFile := flag.String("config", "config.yml", "YAML configuration file")
	parseLog := flag.String("parse", "", "parse log file only (no tournament)")
	crobotsExecutable := flag.String("exe", "crobots", "Crobots executable")
	benchRobot := flag.String("bench", "", "robot (full path, no extension) to create a benchmark tournament for")
	testMode := flag.Bool("test", false, "test mode, check configuration and exit")
	randomMode := flag.Bool("random", false, "generate random matches (4vs4 only)")
	limit := flag.Int("limit", 0, "limit random numer of matches (random mode only)")
	out := flag.String("out", "", "output report to file")

	flag.Parse()

	if _, ok := schema[*tournamentType]; !ok {
		log.Fatalln("Invalid tournament type: ", *tournamentType)
	}

	tot, _ := schema[*tournamentType]

	if *parseLog != "" {
		content := logToString(*parseLog)
		result := count.ParseLogs(strings.Split(content, "\n"))
		printRobots(out, tot, result)
		return
	}

	if *randomMode {
		if tot != 4 {
			log.Fatal("Random mode supported for 4vs4 only")
		}

		if *limit <= 0 {
			log.Fatal("Limit missing or invalid in random mode: ", *limit)
		}
	} else {
		if *limit != 0 {
			log.Println("Limit ignored in non-random mode")
		}
	}

	config := loadConfig(*configFile)
	listSize := len(config.ListRobots)
	if listSize < tot {
		log.Fatal("Robot list insufficient!")
	}

	var robots []string

	for _, r := range config.ListRobots {
		t := checkAndCompile(r, *crobotsExecutable, func(r string) string {
			return config.SourcePath + Separator + r
		})

		robots = append(robots, t)
	}

	if *benchRobot != "" {
		checkAndCompile(*benchRobot, *crobotsExecutable, func(r string) string {
			return r
		})
	}

	if *testMode {
		log.Println("Test mode completed. Exit")
		return
	}

	crobots := *crobotsExecutable

	if !commandExists(crobots) {
		log.Fatal("Crobots executable not found: ", crobots)
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
	result := make(map[string]*count.Robot)
	jobs := make(chan match, NumCPU)
	var wg sync.WaitGroup
	var mutex sync.Mutex
	for w := 1; w <= NumCPU; w++ {
		wg.Add(1)
		go worker(w, jobs, crobots, opt, tot, result, &mutex, &wg)
	}

	var br string = ""
	if *benchRobot != "" {
		br = *benchRobot + RobotBinaryExt
		log.Println("Benchmark tournament for", br)
	}

	if *randomMode {
		l := *limit
		log.Println("Random mode enable. Limit", l)

		for i := 0; i < l; i++ {
			rand.Seed(time.Now().UnixNano())
			perm := rand.Perm(listSize)
			if br != "" {
				a, b, c := perm[0], perm[1], perm[2]
				jobs <- match{Robots: []string{robots[a], robots[b], robots[c], br}}
			} else {
				a, b, c, d := perm[0], perm[1], perm[2], perm[3]
				jobs <- match{Robots: []string{robots[a], robots[b], robots[c], robots[d]}}
			}
		}
	} else if br != "" {
		generateCombinationsForBenchmark(br, robots, tot, jobs)
	} else {
		generateCombinations(robots, tot, jobs)
	}
	close(jobs)
	wg.Wait()
	duration := time.Since(start)
	log.Println("Completed in", duration)
	printRobots(out, tot, result)
}
