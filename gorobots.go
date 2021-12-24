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
	"time"

	"gopkg.in/yaml.v3"
)

const (
	// Version is the software version format v#.#.#-timestamp
	Version = "v1.5.1-20211224"
	// Separator is the OS dependent path separator
	Separator = string(os.PathSeparator)
	// RobotSourceExt is the file extension of the robot source code
	RobotSourceExt = ".r"
	// RobotBinaryExt is the file extension of the compiled binary robot
	RobotBinaryExt = RobotSourceExt + "o"
	// Header is the output header
	Header = "#\tName\t\tGames\t\tWins\t\tTies2\t\tTies3\t\tTies4\t\tLost\t\tPoints\t\tEff%"
	// OutputFormat is a single row output format
	OutputFormat = "%d\t%s\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\t\t%d\t\t%.3f\n"
	// SQLOutputFormat template for SQL updates
	SQLOutputFormat = "UPDATE `results_%s` SET `games`=`games`+%d, `wins`=`wins`+%d, `ties`=`ties`+%d, `points`=`points`+%d WHERE `robot`='%s';\n"
	// StdMatchLimitCycles is the standard Crobots limit as number of cycles for a single match
	StdMatchLimitCycles = "-l200000"
)

var (
	// NumCPU is the number of detected CPUs/cores
	NumCPU int = runtime.NumCPU()
	// Crobots is the crobots executable
	Crobots string
	// EOF End-of-line
	EOF = []byte("\n")
)

type signal struct{}

// Match holds a list of robots for a single crobots match
type Match struct {
	Robots []string
}

// Result holds a single crobots results
type Result struct {
	Robots map[string]*count.Robot
}

// tournament config from YAML file
type tournamentConfig struct {
	Label      string   `yaml:"label"`
	MatchF2F   int      `yaml:"matchF2F"`
	Match3VS3  int      `yaml:"match3VS3"`
	Match4VS4  int      `yaml:"match4VS4"`
	SourcePath string   `yaml:"sourcePath"`
	ListRobots []string `yaml:"listRobots"`
}

// tournament types (modes)
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

func logToString(file string) []byte {
	// Read entire file content, giving us little control but
	// making it very simple. No need to close the file.
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	return content
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

// Verbose helper to display tournament progression
type Verbose struct {
	Enabled  bool
	TotComb  int
	Counter  int
	LastPerc int
	Inc      func(i int) int
}

// Print progression
func (v *Verbose) Print(i int) {
	v.Counter += v.Inc(i)
	perc := 100 * v.Counter / v.TotComb
	if perc > v.LastPerc {
		log.Printf("%d%% completed...\n", perc)
		v.LastPerc = perc
	}
}

// generate standard, ordered tournament match combinations
func generateCombinations(list []string, size int, c chan<- *Match, verbose bool) {
	tot := len(list)
	switch size {
	case 2:
		v := Verbose{
			Enabled:  verbose,
			LastPerc: 0,
			TotComb:  tot * (tot - 1) / 2,
			Inc: func(c int) int {
				return tot - c - 1
			},
		}
		for i := 0; i < tot-1; i++ {
			for j := i + 1; j < tot; j++ {
				c <- &Match{Robots: []string{list[i], list[j]}}
			}
			if v.Enabled {
				v.Print(i)
			}
		}
	case 3:
		v := Verbose{
			Enabled:  verbose,
			LastPerc: 0,
			TotComb:  tot * (tot - 1) * (tot - 2) / 6,
			Inc: func(c int) int {
				return (tot - c - 1) * (tot - c - 2) / 2
			},
		}
		for i := 0; i < tot-2; i++ {
			for j := i + 1; j < tot-1; j++ {
				for k := j + 1; k < tot; k++ {
					c <- &Match{Robots: []string{list[i], list[j], list[k]}}
				}
			}
			if v.Enabled {
				v.Print(i)
			}
		}
	case 4:
		v := Verbose{
			Enabled:  verbose,
			LastPerc: 0,
			TotComb:  tot * (tot - 1) * (tot - 2) * (tot - 3) / 24,
			Inc: func(c int) int {
				return (tot - c - 1) * (tot - c - 2) * (tot - c - 3) / 6
			},
		}
		for i := 0; i < tot-3; i++ {
			for j := i + 1; j < tot-2; j++ {
				for k := j + 1; k < tot-1; k++ {
					for z := k + 1; z < tot; z++ {
						c <- &Match{Robots: []string{list[i], list[j], list[k], list[z]}}
					}
				}
			}
			if v.Enabled {
				v.Print(i)
			}
		}
	default:
		log.Fatal("Error: invalid size", size)
	}
}

// generate standard, ordered sub-combinations for a benchmark tournament
func generateCombinationsForBenchmark(robot string, list []string, size int, c chan<- *Match, verbose bool) {
	tot := len(list)
	switch size {
	case 2:
		v := Verbose{
			Enabled:  verbose,
			LastPerc: 0,
			TotComb:  tot,
			Inc: func(c int) int {
				return 1
			},
		}
		for i := 0; i < tot; i++ {
			c <- &Match{Robots: []string{list[i], robot}}
			if v.Enabled {
				v.Print(i)
			}
		}
	case 3:
		v := Verbose{
			Enabled:  verbose,
			LastPerc: 0,
			TotComb:  tot * (tot - 1) / 2,
			Inc: func(c int) int {
				return tot - c - 1
			},
		}
		for i := 0; i < tot-1; i++ {
			for j := i + 1; j < tot; j++ {
				c <- &Match{Robots: []string{list[i], list[j], robot}}
			}
			if v.Enabled {
				v.Print(i)
			}
		}
	case 4:
		v := Verbose{
			Enabled:  verbose,
			LastPerc: 0,
			TotComb:  tot * (tot - 1) * (tot - 2) / 6,
			Inc: func(c int) int {
				return (tot - c - 1) * (tot - c - 2) / 2
			},
		}
		for i := 0; i < tot-2; i++ {
			for j := i + 1; j < tot-1; j++ {
				for k := j + 1; k < tot; k++ {
					c <- &Match{Robots: []string{list[i], list[j], list[k], robot}}
				}
			}
			if v.Enabled {
				v.Print(i)
			}
		}
	default:
		log.Fatal("Error: invalid size", size)
	}
}

func check(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}

// print out tournament results to stdout or file or SQL updates
// if errors occur always print out to stdout
func printRobots(out, sql, tournamentType *string, tot int, result *Result) {
	var bots []count.Robot

	for _, robot := range result.Robots {
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
	// local function to print out to stdout if no output file is specified
	// or as at last resort should errors occur
	printToStd := func() {
		var i int = 0
		fmt.Println(Header)
		for _, robot := range bots {
			i++
			fmt.Printf(OutputFormat, i, robot.Name, robot.Games, robot.Wins, robot.Ties[0], robot.Ties[1], robot.Ties[2], robot.Games-robot.Wins-(robot.Ties[0]+robot.Ties[1]+robot.Ties[2]), robot.Points, robot.Eff)
		}
	}
	// Output to SQL updates
	if *sql != "" {
		outputFile := *sql
		t := *tournamentType
		f, err := os.Create(outputFile)
		if check(err) {
			printToStd()
			return
		}
		defer f.Close()
		w := bufio.NewWriter(f)
		for _, robot := range bots {
			_, err = fmt.Fprintf(w, SQLOutputFormat, t, robot.Games, robot.Wins, robot.Ties[0]+robot.Ties[1]+robot.Ties[2], robot.Points, robot.Name)
			if check(err) {
				printToStd()
				return
			}
		}
		w.Flush()
	}
	// Output to tab separated file
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

// given a match returns a Crobots command ready to be executed
func (m *Match) executeCrobotsMatch(opt string, n int) ([]byte, error) {
	switch n {
	case 2:
		return exec.Command(Crobots, opt, StdMatchLimitCycles, m.Robots[0], m.Robots[1]).Output()
	case 3:
		return exec.Command(Crobots, opt, StdMatchLimitCycles, m.Robots[0], m.Robots[1], m.Robots[2]).Output()
	case 4:
		return exec.Command(Crobots, opt, StdMatchLimitCycles, m.Robots[0], m.Robots[1], m.Robots[2], m.Robots[3]).Output()
	default:
		log.Fatal("Error: invalid size", n)
	}
	return nil, fmt.Errorf("something went horribly wrong")
}

// given a match execute Crobots command and parse output to update partial results
func (m *Match) processCrobotsMatch(opt string, tot int, result chan<- *Result) {
	out, err := m.executeCrobotsMatch(opt, tot)
	if err != nil {
		log.Fatal(err)
	}
	if len(out) == 0 {
		log.Fatal("no output from Crobots match")
	}
	partial := count.ParseLogs(bytes.Split(out, EOF))
	result <- &Result{Robots: partial}
}

// goroutine to handle a batch of matches
func worker(id int, matches <-chan *Match, opt string, tot int, result chan<- *Result, s chan<- signal) {
	for m := range matches {
		m.processCrobotsMatch(opt, tot, result)
	}
	var end signal
	s <- end
}

// check binary robot exists or try to compile its source code
func checkAndCompile(r string, path func(r string) string) string {
	t := path(r) + RobotBinaryExt
	if !fileExists(t) {
		log.Println("Warning: binary robot cannot be found:", t, "Compiling source code")
		s := path(r) + RobotSourceExt
		if fileExists(s) {
			var out bytes.Buffer
			cmd := exec.Command(Crobots, "-c", s)
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("Error: robot source code cannot be found: ", s)
		}
	}
	return t
}

// generate matches for benchmark using current random permutation and queue them all
func benchMatches(limit, slices int, perm []int, robots []string, jobs chan *Match, br string) int {
	t := min(limit, slices)
	j := 0
	for i := 0; i < t; i++ {
		current := perm[j : j+3]
		a, b, c := current[0], current[1], current[2]
		jobs <- &Match{Robots: []string{robots[a], robots[b], robots[c], br}}
		j += 3
	}
	return t
}

// generate matches using current random permutation and queue them all
func randomMatches(limit, slices int, perm []int, robots []string, jobs chan *Match) int {
	t := min(limit, slices)
	j := 0
	for i := 0; i < t; i++ {
		current := perm[j : j+4]
		a, b, c, d := current[0], current[1], current[2], current[3]
		jobs <- &Match{Robots: []string{robots[a], robots[b], robots[c], robots[d]}}
		j += 4
	}
	return t
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {

	log.Println("GoRobots", Version)
	log.Println("Detected CPU(s)/core(s):", NumCPU)
	tournamentType := flag.String("type", "", "tournament type: f2f, 3vs3 or 4vs4")
	configFile := flag.String("config", "config.yml", "YAML configuration file")
	parseLog := flag.String("parse", "", "parse log file only (no tournament)")
	crobotsExecutable := flag.String("exe", "crobots", "Crobots executable")
	benchRobot := flag.String("bench", "", "robot (full path, no extension) to create a benchmark tournament for")
	testMode := flag.Bool("test", false, "test mode, check configuration and exit")
	randomMode := flag.Bool("random", false, "random mode: generate random matches for 4vs4 only")
	limit := flag.Int("limit", 0, "limit random number of matches (random mode only)")
	out := flag.String("out", "", "output results to file")
	sql := flag.String("sql", "", "output results as SQL updates to file")
	cpu := flag.Int("cpu", NumCPU, "number of threads (CPUs/cores)")
	verbose := flag.Bool("verbose", false, "verbose mode: print tournament progression percentage")

	flag.Parse()

	c := *cpu

	if c < 1 || c > NumCPU {
		log.Println("Invalid parameter cpu", c, ". Using default", NumCPU)
	} else {
		NumCPU = c
	}

	log.Println("Using", NumCPU, "CPU(s)/core(s)")
	runtime.GOMAXPROCS(NumCPU)

	if _, ok := schema[*tournamentType]; !ok {
		log.Fatalln("Error: invalid tournament type: ", *tournamentType)
	}

	tot, _ := schema[*tournamentType]

	if *parseLog != "" {
		content := logToString(*parseLog)
		result := &Result{Robots: count.ParseLogs(bytes.Split(content, EOF))}
		printRobots(out, sql, tournamentType, tot, result)
		return
	}

	if *randomMode {
		if tot != 4 {
			log.Fatal("Error: random mode supported for 4vs4 only")
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
		log.Fatal("Error: robot list too small!")
	}

	var robots []string

	Crobots = *crobotsExecutable

	if !commandExists(Crobots) {
		log.Fatal("Error: Crobots executable not found ", Crobots)
	}

	for _, r := range config.ListRobots {
		t := checkAndCompile(r, func(r string) string {
			return config.SourcePath + Separator + r
		})

		robots = append(robots, t)
	}

	if *benchRobot != "" {
		checkAndCompile(*benchRobot, func(r string) string {
			return r
		})
		// sanity check
		base := count.GetName(*benchRobot)
		for _, r := range robots {
			b := count.GetName(r)
			if base == b {
				log.Fatal("Error: duplicate name detected as configuration already contains ", base)
			}
		}
	}

	var br string = ""
	if *benchRobot != "" {
		br = *benchRobot + RobotBinaryExt
		log.Println("Benchmark", *tournamentType, "tournament for", *benchRobot)
	}

	if *testMode {
		log.Println("Test mode completed. Exit")
		return
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
	log.Println("Start processing", *tournamentType, "...")
	start := time.Now()
	result := make(chan *Result)
	total := &Result{Robots: make(map[string]*count.Robot)}
	jobs := make(chan *Match, NumCPU)
	sig := make(chan signal)

	for w := 1; w <= NumCPU; w++ {
		go worker(w, jobs, opt, tot, result, sig)
	}

	go func() {
		for p := range result {
			for _, robot := range p.Robots {
				name := robot.Name
				if update, found := total.Robots[name]; found {
					update.Games += robot.Games
					update.Wins += robot.Wins
					update.Ties[0] += robot.Ties[0]
					update.Ties[1] += robot.Ties[1]
					update.Ties[2] += robot.Ties[2]
				} else {
					total.Robots[name] = &count.Robot{Name: name, Games: robot.Games, Wins: robot.Wins, Ties: robot.Ties, Points: 0, Eff: 0.0}
				}
			}
		}
		var end signal
		sig <- end
	}()

	if *randomMode {
		l := *limit

		var slices int
		if br != "" {
			slices = listSize / 3
		} else {
			slices = listSize / 4
		}

		v := Verbose{
			Enabled:  *verbose,
			LastPerc: 0,
			TotComb:  l,
			Inc: func(c int) int {
				return c
			},
		}
		log.Println("Random mode enabled. Limit", l)

		for l > 0 {
			var t int
			rand.Seed(time.Now().UnixNano())
			perm := rand.Perm(listSize)
			if br != "" {
				t = benchMatches(l, slices, perm, robots, jobs, br)
			} else {
				t = randomMatches(l, slices, perm, robots, jobs)
			}
			l -= t
			if v.Enabled {
				v.Print(t)
			}
		}
	} else if br != "" {
		generateCombinationsForBenchmark(br, robots, tot, jobs, *verbose)
	} else {
		generateCombinations(robots, tot, jobs, *verbose)
	}
	close(jobs)
	for i := 0; i < NumCPU; i += 1 {
		<-sig
	}
	close(result)
	<-sig

	duration := time.Since(start)
	log.Println("Completed in", duration)
	printRobots(out, sql, tournamentType, tot, total)
}
