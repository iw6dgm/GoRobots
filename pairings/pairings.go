package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Constants
const (
	Separator  = string(os.PathSeparator)
	GROUP_SIZE = 64 // Desired group size
	matchF2F   = 1000
	match3vs3  = 15
	match4vs4  = 1
	label      = "group"
)

// Tournament lists
var (
	tournament1990 = []string{}
	tournament1991 = []string{}
	tournament1992 = []string{}
	tournament1993 = []string{}
	tournament1994 = []string{}
	tournament1995 = []string{}
	tournament1996 = []string{}
	tournament1997 = []string{}
	tournament1998 = []string{}
	tournament1999 = []string{}
	tournament2000 = []string{}
	tournament2001 = []string{}
	tournament2002 = []string{}
	tournament2003 = []string{}
	tournament2004 = []string{}
	tournament2007 = []string{}
	tournament2010 = []string{}
	tournament2011 = []string{}
	tournament2012 = []string{}
	tournament2013 = []string{}
	tournament2015 = []string{}
	tournament2020 = []string{}
	tournament2025 = []string{}
	micro          = []string{}
	crobs          = []string{}
	others         = []string{}

	tournaments = [][]string{}
	TYPES       = []string{"f2f", "3vs3", "4vs4"}

	rounds [][]string
	robots int
)

func main() {
	setup()
	countRobots()

	attempts := 0

	for {
		attempts++
		fmt.Printf("ATTEMPT : %d\n", attempts)
		shuffle()
		collect()
		if !alternativePairing() {
			break
		}
	}

	show()
	buildConfigFileYAML()
	//buildSQLInserts() // optional - not needed if using `tournament` scripts
}

/* uncomment this and comment the main above if you want to use a single list of robots */
// func main() {
// 	setupFromSingleList()
// 	show()
// 	buildConfigFileYAML()
// }

// Show pairings (plain text)
func show() {
	n := 1
	for _, round := range rounds {
		if len(round) > 0 {
			fmt.Printf("------- Group %d (size %d) ------\n", n, len(round))
			for _, s := range round {
				fmt.Println(s)
			}
			n++
		}
	}
}

// Save pairings into YAML files
func buildConfigFileYAML() {
	n := 1
	for _, round := range rounds {
		if len(round) > 0 {
			//fmt.Printf("------- CFG group%d ------\n", n)
			f, err := os.Create(fmt.Sprintf("%s%d.yml", label, n))
			check(err)
			defer f.Close()
			w := bufio.NewWriter(f)
			_, err = fmt.Fprintf(w, "matchF2F: %d\nmatch3VS3: %d\nmatch4VS4: %d\nsourcePath: '.'\n", matchF2F, match3vs3, match4vs4)
			check(err)
			_, err = fmt.Fprintf(w, "label: '%s%d'\n", label, n)
			check(err)

			_, err = w.WriteString("listRobots: [\n")
			check(err)
			for _, s := range round {
				_, err = fmt.Fprintf(w, "'%s',\n", s)
				check(err)
			}
			_, err = w.WriteString("\n]")
			check(err)
			w.Flush()
			n++
		}
	}
}

// Prints SQL Insert statements to initialise reports tables
func buildSQLInserts() {
	n := 1
	for _, round := range rounds {
		if len(round) > 0 {
			fmt.Printf("------- SQL group%d ------\n", n)

			var values strings.Builder
			for i, s := range round {
				values.WriteString(fmt.Sprintf("('%s')", filepath.Base(s)))
				if i != len(round)-1 {
					values.WriteString(",\n")
				}
			}
			values.WriteString(";")

			sql := values.String()
			for _, table := range TYPES {
				fmt.Printf("------- %s -------\n", table)
				fmt.Printf("DELETE FROM `results_%s`;\n", table)
				fmt.Printf("INSERT INTO `results_%s` (robot) VALUES\n", table)
				fmt.Println(sql)
			}
			n++
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func collect() {
	tournaments = [][]string{
		tournament1990, tournament1991, tournament1992, tournament1993, tournament1994,
		tournament1995, tournament1996, tournament1997, tournament1998, tournament1999,
		tournament2000, tournament2001, tournament2002, tournament2003, tournament2004,
		tournament2007, tournament2010, tournament2011, tournament2012, tournament2013,
		tournament2015, tournament2020, tournament2025, crobs, micro, others,
	}

	// Shuffle tournaments
	rand.Shuffle(len(tournaments), func(i, j int) {
		tournaments[i], tournaments[j] = tournaments[j], tournaments[i]
	})
}

func countRobots() {
	robots = len(micro) + len(crobs) + len(others) +
		len(tournament1990) + len(tournament1991) + len(tournament1992) + len(tournament1993) +
		len(tournament1994) + len(tournament1995) + len(tournament1996) + len(tournament1997) +
		len(tournament1998) + len(tournament1999) + len(tournament2000) + len(tournament2001) +
		len(tournament2002) + len(tournament2003) + len(tournament2004) + len(tournament2007) +
		len(tournament2010) + len(tournament2011) + len(tournament2012) + len(tournament2013) +
		len(tournament2015) + len(tournament2020) + len(tournament2025)

	fmt.Printf("TOTAL Robots :%d\n", robots)
}

func alternativePairing() bool {
	groupCount := robots / GROUP_SIZE
	if robots%GROUP_SIZE > 0 {
		groupCount++
	}

	groupIndex := 0
	rounds = make([][]string, groupCount)

	for i := 0; i < groupCount; i++ {
		rounds[i] = []string{}
	}

	for _, tournament := range tournaments {
		for _, r := range tournament {
			// Check if robot already exists in the current group by comparing basenames
			if slices.ContainsFunc(rounds[groupIndex], func(existing string) bool {
				return filepath.Base(existing) == filepath.Base(r)
			}) {
				fmt.Printf("Robot %s generated a conflict\n", filepath.Base(r))
				return true // has conflicts
			}
			rounds[groupIndex] = append(rounds[groupIndex], r)

			groupIndex++
			if groupIndex == groupCount {
				groupIndex = 0
			}
		}
	}

	return false // no conflicts
}

func shuffle() {
	shuffleSlice(tournament1990)
	shuffleSlice(tournament1991)
	shuffleSlice(tournament1992)
	shuffleSlice(tournament1993)
	shuffleSlice(tournament1994)
	shuffleSlice(tournament1995)
	shuffleSlice(tournament1996)
	shuffleSlice(tournament1997)
	shuffleSlice(tournament1998)
	shuffleSlice(tournament1999)
	shuffleSlice(tournament2000)
	shuffleSlice(tournament2001)
	shuffleSlice(tournament2002)
	shuffleSlice(tournament2003)
	shuffleSlice(tournament2004)
	shuffleSlice(tournament2007)
	shuffleSlice(tournament2010)
	shuffleSlice(tournament2011)
	shuffleSlice(tournament2012)
	shuffleSlice(tournament2013)
	shuffleSlice(tournament2015)
	shuffleSlice(tournament2020)
	shuffleSlice(crobs)
	shuffleSlice(micro)
	shuffleSlice(others)
}

// shuffleSlice shuffles a slice in place
func shuffleSlice(slice []string) {
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

// use this from a single list of robots
func setupFromSingleList() {
	var list = []string{ /* add robots here */ }

	robots = len(list)
	groupCount := robots / GROUP_SIZE
	if robots%GROUP_SIZE > 0 {
		groupCount++
	}
	pairing := func() bool {

		groupIndex := 0
		rounds = make([][]string, groupCount)

		for i := 0; i < groupCount; i++ {
			rounds[i] = []string{}
		}
		for _, r := range list {
			// Check if robot already exists in the current group by comparing basenames
			if slices.ContainsFunc(rounds[groupIndex], func(existing string) bool {
				return filepath.Base(existing) == filepath.Base(r)
			}) {
				fmt.Printf("Robot %s generated a conflict\n", filepath.Base(r))
				return true // has conflicts
			}
			rounds[groupIndex] = append(rounds[groupIndex], r)

			groupIndex++
			if groupIndex == groupCount {
				groupIndex = 0
			}
		}
		return false
	}

	attempts := 0

	for {
		attempts++
		fmt.Printf("ATTEMPT : %d\n", attempts)
		shuffleSlice(list)
		if !pairing() {
			break
		}
	}
}

func setupMidi() {
	var path string
	fmt.Print("Loading others... ")
	others = []string{
		fmt.Sprintf("aminet%santiclock", Separator),
		fmt.Sprintf("aminet%sbeaver", Separator),
		fmt.Sprintf("aminet%sblindschl", Separator),
		fmt.Sprintf("aminet%sblindschl2", Separator),
		fmt.Sprintf("aminet%smirobot", Separator),
		fmt.Sprintf("aminet%sopfer", Separator),
		fmt.Sprintf("aminet%sschwan", Separator),
		fmt.Sprintf("aminet%stron", Separator),
		fmt.Sprintf("cplusplus%sselvaggio", Separator),
		fmt.Sprintf("cplusplus%svikingo", Separator),
	}
	fmt.Printf("%d robot(s)\n", len(others))

	fmt.Print("Loading 1990... ")
	path = fmt.Sprintf("1990%s", Separator)
	tournament1990 = []string{
		path + "et_1",
		path + "et_2",
		path + "hunter",
		path + "killer",
		path + "nexus_1",
		path + "rob1",
		path + "scanner",
		path + "york",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1990))

	fmt.Print("Loading 1991... ")
	path = fmt.Sprintf("1991%s", Separator)
	tournament1991 = []string{
		path + "blade3",
		path + "casimiro",
		path + "ccyber",
		path + "clover",
		path + "diagonal",
		path + "et_3",
		path + "f1",
		path + "fdig",
		path + "geltrude",
		path + "genius_j",
		path + "gira",
		path + "gunner",
		path + "jazz",
		path + "nexus_2",
		path + "paolo101",
		path + "paolo77",
		path + "poor",
		path + "qibo",
		path + "robocop",
		path + "runner",
		path + "sara_6",
		path + "seeker",
		path + "warrior2",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1991))

	fmt.Print("Loading 1992... ")
	path = fmt.Sprintf("1992%s", Separator)
	tournament1992 = []string{
		path + "666",
		path + "ap_1",
		path + "assassin",
		path + "baeos",
		path + "banzel",
		path + "bronx-00",
		path + "bry_bry",
		path + "crazy",
		path + "cube",
		path + "cw",
		path + "d47",
		path + "daitan3",
		path + "dancer",
		path + "deluxe",
		path + "dorsai",
		path + "et_4",
		path + "et_5",
		path + "flash",
		path + "genesis",
		path + "hunter",
		path + "ice",
		path + "jack",
		path + "jager",
		path + "johnny",
		path + "lead1",
		path + "marika",
		path + "mimo6new",
		path + "mrcc",
		path + "mut",
		path + "ninus6",
		path + "nl_1a",
		path + "nl_1b",
		path + "ola",
		path + "paolo",
		path + "pavido",
		path + "phobos_1",
		path + "pippo92",
		path + "pippo",
		path + "raid",
		path + "random",
		path + "revenge3",
		path + "robbie",
		path + "robocop2",
		path + "robocop",
		path + "sassy",
		path + "spider",
		path + "sp",
		path + "superv",
		path + "t1000",
		path + "thunder",
		path + "triangol",
		path + "trio",
		path + "uanino",
		path + "warrior3",
		path + "xdraw2",
		path + "zorro",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1992))

	fmt.Print("Loading 1993... ")
	path = fmt.Sprintf("1993%s", Separator)
	tournament1993 = []string{
		path + "am_174",
		path + "ap_2",
		path + "ares",
		path + "argon",
		path + "aspide",
		path + "beast",
		path + "biro",
		path + "blade8",
		path + "boom",
		path + "brain",
		path + "cantor",
		path + "castore",
		path + "casual",
		path + "corner1d",
		path + "corner3",
		path + "courage",
		path + "(c)",
		path + "crob1",
		path + "deluxe_2",
		path + "deluxe_3",
		path + "didimo",
		path + "duke",
		path + "elija",
		path + "fermo",
		path + "flash2",
		path + "food5",
		path + "godel",
		path + "gunnyboy",
		path + "hamp1",
		path + "hamp2",
		path + "hell",
		path + "horse",
		path + "isaac",
		path + "kami",
		path + "lazy",
		path + "mimo13",
		path + "mister2",
		path + "mister3",
		path + "mohawk",
		path + "mutation",
		path + "ninus17",
		path + "nl_2a",
		path + "nl_2b",
		path + "p68",
		path + "p69",
		path + "penta",
		path + "phobos_2",
		path + "pippo93",
		path + "pognant",
		path + "poirot",
		path + "polluce",
		path + "premana",
		path + "puyopuyo",
		path + "raid2",
		path + "rapper",
		path + "r_cyborg",
		path + "r_daneel",
		path + "robocop3",
		path + "spartaco",
		path + "target",
		path + "tm",
		path + "torneo",
		path + "vannina",
		path + "vocus",
		path + "warrior4",
		path + "wassilij",
		path + "wolfgang",
		path + "zulu",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1993))

	fmt.Print("Loading 1994... ")
	path = fmt.Sprintf("1994%s", Separator)
	tournament1994 = []string{
		path + "8bismark",
		path + "anglek2",
		path + "apache",
		path + "bachopin",
		path + "baubau",
		path + "biro",
		path + "blob",
		path + "circlek1",
		path + "corner3b",
		path + "corner4",
		path + "deluxe_4",
		path + "deluxe_5",
		path + "didimo",
		path + "dima10",
		path + "dima9",
		path + "emanuela",
		path + "ematico",
		path + "fastfood",
		path + "flash3",
		path + "funky",
		path + "giali1",
		path + "hal9000",
		path + "heavens",
		path + "horse2",
		path + "iching",
		path + "jet",
		path + "ken",
		path + "lazyii",
		path + "matrox",
		path + "maverick",
		path + "miaomiao",
		path + "nemesi",
		path + "ninus75",
		path + "patcioca",
		path + "pioppo",
		path + "pippo94a",
		path + "pippo94b",
		path + "polipo",
		path + "randwall",
		path + "robot1",
		path + "robot2",
		path + "sdix3",
		path + "sgnaus",
		path + "shadow",
		path + "superfly",
		path + "the_dam",
		path + "t-rex",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1994))

	fmt.Print("Loading 1995... ")
	path = fmt.Sprintf("1995%s", Separator)
	tournament1995 = []string{
		path + "andrea",
		path + "animal",
		path + "apache95",
		path + "archer",
		path + "b115e2",
		path + "b52",
		path + "biro",
		path + "boss",
		path + "camillo",
		path + "carlo",
		path + "circle",
		path + "cri95",
		path + "diablo",
		path + "flash4",
		path + "hal9000",
		path + "heavens",
		path + "horse3",
		path + "kenii",
		path + "losendos",
		path + "mikezhar",
		path + "ninus99",
		path + "paccu",
		path + "passion",
		path + "peribolo",
		path + "pippo95",
		path + "rambo",
		path + "rocco",
		path + "saxy",
		path + "sel",
		path + "skizzo",
		path + "star",
		path + "stinger",
		path + "tabori-1",
		path + "tabori-2",
		path + "tequila",
		path + "tmii",
		path + "tox",
		path + "t-rex",
		path + "tricky",
		path + "twins",
		path + "upv-9596",
		path + "xenon",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1995))

	fmt.Print("Loading 1996... ")
	path = fmt.Sprintf("1996%s", Separator)
	tournament1996 = []string{
		path + "aleph",
		path + "andrea96",
		path + "ap_4",
		path + "carlo96",
		path + "diablo2",
		path + "drago5",
		path + "d_ray",
		path + "fb3",
		path + "gevbass",
		path + "golem",
		path + "gpo2",
		path + "hal9000",
		path + "heavnew",
		path + "hider2",
		path + "infinity",
		path + "jaja",
		path + "memories",
		path + "murdoc",
		path + "natas",
		path + "newb52",
		path + "pacio",
		path + "pippo96a",
		path + "pippo96b",
		path + "!",
		path + "risk",
		path + "robot1",
		path + "robot2",
		path + "rudolf",
		path + "second3",
		path + "s-seven",
		path + "tatank_3",
		path + "tronco",
		path + "uht",
		path + "xabaras",
		path + "yuri",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1996))

	fmt.Print("Loading 1997... ")
	path = fmt.Sprintf("1997%s", Separator)
	tournament1997 = []string{
		path + "1&1",
		path + "abyss",
		path + "ai1",
		path + "andrea97",
		path + "arale",
		path + "belva",
		path + "carlo97",
		path + "ciccio",
		path + "colossus",
		path + "diablo3",
		path + "diabolik",
		path + "drago6",
		path + "erica",
		path + "fable",
		path + "flash5",
		path + "fya",
		path + "gevbass2",
		path + "golem2",
		path + "gundam",
		path + "hal9000",
		path + "jedi",
		path + "kill!",
		path + "me-110c",
		path + "ncmplt",
		path + "paperone",
		path + "pippo97",
		path + "raid3",
		path + "robivinf",
		path + "rudolf_2",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1997))

	fmt.Print("Loading 1998... ")
	path = fmt.Sprintf("1998%s", Separator)
	tournament1998 = []string{
		path + "ai2",
		path + "bartali",
		path + "carla",
		path + "coppi",
		path + "dia",
		path + "dicin",
		path + "eva00",
		path + "eva01",
		path + "freedom",
		path + "fscan",
		path + "goblin",
		path + "goldrake",
		path + "hal9000",
		path + "heavnew",
		path + "maxheav",
		path + "ninja",
		path + "paranoid",
		path + "pippo98",
		path + "plump",
		path + "quarto",
		path + "rattolo",
		path + "rudolf_3",
		path + "son-goku",
		path + "sottolin",
		path + "stay",
		path + "stighy98",
		path + "themicro",
		path + "titania",
		path + "tornado",
		path + "traker1",
		path + "traker2",
		path + "vision",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1998))

	fmt.Print("Loading 1999... ")
	path = fmt.Sprintf("1999%s", Separator)
	tournament1999 = []string{
		path + "11",
		path + "aeris",
		path + "akira",
		path + "alezai17",
		path + "alfa99",
		path + "alien",
		path + "ap_5",
		path + "bastrd!!",
		path + "cancer",
		path + "carlo99",
		path + "#cimice#",
		path + "cortez",
		path + "cyborg",
		path + "dario",
		path + "dav46",
		path + "defender",
		path + "elisir",
		path + "flash6",
		path + "hal9000",
		path + "ilbestio",
		path + "jedi2",
		path + "ka_aroth",
		path + "kakakatz",
		path + "lukather",
		path + "mancino",
		path + "marko",
		path + "mcenrobo",
		path + "m_hingis",
		path + "minatela",
		path + "new",
		path + "nexus_2",
		path + "nl_3a",
		path + "nl_3b",
		path + "obiwan",
		path + "omega99",
		path + "panduro",
		path + "panic",
		path + "pippo99",
		path + "pizarro",
		path + "quarto",
		path + "quingon",
		path + "rudolf_4",
		path + "satana",
		path + "shock",
		path + "songohan",
		path + "stealth",
		path + "storm",
		path + "surrende",
		path + "t1001",
		path + "themicro",
		path + "titania2",
		path + "vibrsper",
		path + "zero",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1999))

	fmt.Print("Loading 2000... ")
	path = fmt.Sprintf("2000%s", Separator)
	tournament2000 = []string{
		path + "bach_2k",
		path + "defender",
		path + "doppia_g",
		path + "flash7",
		path + "jedi3",
		path + "mancino",
		path + "marine",
		path + "m_hingis",
		path + "navaho",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2000))

	fmt.Print("Loading 2001... ")
	path = fmt.Sprintf("2001%s", Separator)
	tournament2001 = []string{
		path + "burrfoot",
		path + "charles",
		path + "cisc",
		path + "cobra",
		path + "copter",
		path + "defender",
		path + "fizban",
		path + "gers",
		path + "grezbot",
		path + "hammer",
		path + "heavnew",
		path + "homer",
		path + "klr2",
		path + "kyashan",
		path + "max10",
		path + "mflash2",
		path + "microdna",
		path + "midi_zai",
		path + "mnl_1a",
		path + "mnl_1b",
		path + "murray",
		path + "neo0",
		path + "nl_5b",
		path + "pippo1a",
		path + "pippo1b",
		path + "raistlin",
		path + "ridicol",
		path + "risc",
		path + "rudy_xp",
		path + "sdc2",
		path + "staticii",
		path + "thunder",
		path + "vampire",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2001))

	fmt.Print("Loading 2002... ")
	path = fmt.Sprintf("2002%s", Separator)
	tournament2002 = []string{
		path + "01",
		path + "adsl",
		path + "anakin",
		path + "asterix",
		path + "bruenor",
		path + "colera",
		path + "colosseum",
		path + "copter_2",
		path + "corner5",
		path + "doom2099",
		path + "dynamite",
		path + "enigma",
		path + "groucho",
		path + "halman",
		path + "harpo",
		path + "idefix",
		path + "kyash_2",
		path + "marco",
		path + "mazinga",
		path + "medioman",
		path + "mg_one",
		path + "mind",
		path + "neo_sifr",
		path + "ollio",
		path + "padawan",
		path + "peste",
		path + "pippo2a",
		path + "pippo2b",
		path + "regis",
		path + "scsi",
		path + "serse",
		path + "ska",
		path + "stanlio",
		path + "staticxp",
		path + "supernov",
		path + "tifo",
		path + "tigre",
		path + "todos",
		path + "tomahawk",
		path + "vaiolo",
		path + "vauban",
		path + "yoyo",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2002))

	fmt.Print("Loading 2003... ")
	path = fmt.Sprintf("2003%s", Separator)
	tournament2003 = []string{
		path + "730",
		path + "adrian",
		path + "ares",
		path + "barbarian",
		path + "blitz",
		path + "briscolo",
		path + "bruce",
		path + "cadderly",
		path + "cariddi",
		path + "cvirus2",
		path + "cvirus",
		path + "danica",
		path + "dynacond",
		path + "falco",
		path + "foursquare",
		path + "frame",
		path + "herpes",
		path + "ici",
		path + "instict",
		path + "irpef",
		path + "janu",
		path + "kyash_3c",
		path + "kyash_3m",
		path + "lbr1",
		path + "lbr",
		path + "lebbra",
		path + "mg_two",
		path + "minicond",
		path + "morituro",
		path + "nautilus",
		path + "nemo",
		path + "neo_sel",
		path + "piiico",
		path + "pippo3b",
		path + "pippo3",
		path + "red_wolf",
		path + "scanner",
		path + "scilla",
		path + "sirio",
		path + "sith",
		path + "sky",
		path + "spaceman",
		path + "tartaruga",
		path + "valevan",
		path + "virus2",
		path + "virus",
		path + "yoda",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2003))

	fmt.Print("Loading 2004... ")
	path = fmt.Sprintf("2004%s", Separator)
	tournament2004 = []string{
		path + "adam",
		path + "b_selim",
		path + "!caos",
		path + "ciclope",
		path + "coyote",
		path + "diodo",
		path + "fisco",
		path + "gostar",
		path + "gotar2",
		path + "gotar",
		path + "irap",
		path + "ires",
		path + "magneto",
		path + "mg_three",
		path + "mystica",
		path + "n3g4_jr",
		path + "n3g4tivo",
		path + "new_mini",
		path + "pippo04a",
		path + "pippo04b",
		path + "poldo",
		path + "puma",
		path + "rat-man",
		path + "ravatto",
		path + "rotar",
		path + "selim_b",
		path + "unlimited",
		path + "wgdi",
		path + "zener",
		path + "!zeus",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2004))

	fmt.Print("Loading 2007... ")
	path = fmt.Sprintf("2007%s", Separator)
	tournament2007 = []string{
		path + "angel",
		path + "back",
		path + "brontolo",
		path + "electron",
		path + "e",
		path + "gongolo",
		path + "iceman",
		path + "mammolo",
		path + "microbo1",
		path + "microbo2",
		path + "midi1",
		path + "neutron",
		path + "pippo07a",
		path + "pippo07b",
		path + "pisolo",
		path + "pyro",
		path + "rythm",
		path + "tobey",
		path + "t",
		path + "zigozago",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2007))

	fmt.Print("Loading 2010... ")
	path = fmt.Sprintf("2010%s", Separator)
	tournament2010 = []string{
		path + "buffy",
		path + "cancella",
		path + "copia",
		path + "enkidu",
		path + "eurialo",
		path + "hal9010",
		path + "macchia",
		path + "niso",
		path + "party",
		path + "pippo10a",
		path + "reuben",
		path + "stitch",
		path + "sweat",
		path + "taglia",
		path + "toppa",
		path + "wall-e",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2010))

	fmt.Print("Loading 2011... ")
	path = fmt.Sprintf("2011%s", Separator)
	tournament2011 = []string{
		path + "ataman",
		path + "coeurl",
		path + "gerty",
		path + "hal9011",
		path + "jeeg",
		path + "minion",
		path + "nikita",
		path + "origano",
		path + "ortica",
		path + "pain",
		path + "piperita",
		path + "pippo11a",
		path + "pippo11b",
		path + "tannhause",
		path + "unmaldestr",
		path + "vector",
		path + "wall-e_ii",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2011))

	fmt.Print("Loading 2012... ")
	path = fmt.Sprintf("2012%s", Separator)
	tournament2012 = []string{
		path + "avoider",
		path + "beat",
		path + "british",
		path + "camille",
		path + "china",
		path + "dampyr",
		path + "easyjet",
		path + "flash8c",
		path + "flash8e",
		path + "gerty2",
		path + "grezbot2",
		path + "gunnyb29",
		path + "hal9012",
		path + "lycan",
		path + "mister2b",
		path + "mister3b",
		path + "pippo12a",
		path + "pippo12b",
		path + "power",
		path + "puffomic",
		path + "puffomid",
		path + "q",
		path + "ryanair",
		path + "silversurf",
		path + "torchio",
		path + "wall-e_iii",
		path + "yeti",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2012))

	fmt.Print("Loading 2013... ")
	path = fmt.Sprintf("2013%s", Separator)
	tournament2013 = []string{
		path + "axolotl",
		path + "destro",
		path + "eternity",
		path + "frisa_13",
		path + "gerty3",
		path + "ghostrider",
		path + "guanaco",
		path + "gunnyb13",
		path + "hal9013",
		path + "jarvis",
		path + "lamela",
		path + "leopon",
		path + "ncc-1701",
		path + "osvaldo",
		path + "pippo13a",
		path + "pippo13b",
		path + "pray",
		path + "ug2k",
		path + "wall-e_iv",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2013))

	fmt.Print("Loading 2015... ")
	path = fmt.Sprintf("2015%s", Separator)
	tournament2015 = []string{
		path + "antman",
		path + "aswhup",
		path + "avoider", // same as 2012
		path + "babadook",
		path + "circles15",
		path + "colour",
		path + "coppi15mc1",
		path + "coppi15md1",
		path + "flash9",
		path + "frank15",
		path + "gerty4",
		path + "hal9015",
		path + "hulk",
		path + "linabo15",
		path + "lluke",
		path + "mcfly",
		path + "mike3",
		path + "pippo15a",
		path + "pippo15b",
		path + "puppet",
		path + "randguard",
		path + "salippo",
		path + "sidewalk",
		path + "thor",
		path + "tux",
		path + "tyrion",
		path + "wall-e_v",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2015))

	fmt.Print("Loading 2020... ")
	path = fmt.Sprintf("2020%s", Separator)
	tournament2020 = []string{
		path + "antman_20",
		path + "b4b",
		path + "brexit",
		path + "coppi20mc1",
		path + "coppi20md1",
		path + "discotek",
		path + "flash10",
		path + "gerty5",
		path + "hal9020",
		path + "hulk_20",
		path + "jarvis2",
		path + "leavy2",
		path + "loneliness",
		path + "pippo20a",
		path + "pippo20b",
		path + "wall-e_vi",
		path + "wizard2",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2020))

	fmt.Print("Loading 2025... ")
	path = fmt.Sprintf("2025%s", Separator)
	tournament2025 = []string{
		path + "blue", // late entry
		path + "extrasmall",
		path + "flash11",
		path + "gerty6",
		path + "hal9025",
		path + "hulk_25",
		path + "hydra",
		path + "kerberos",
		path + "meeseeks1",
		path + "meeseeks2",
		path + "nova",
		path + "pippo25a",
		path + "rabbitc",
		path + "rotaprinc8",
		path + "sentry2",
		path + "sentry3",
		path + "sgorbio",
		path + "slant6",
		path + "supremo",
		path + "trouble3",
		path + "ultron_25",
		path + "wall-e_vii",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2025))

	fmt.Print("Loading crobs... ")
	path = fmt.Sprintf("crobs%s", Separator)
	crobs = []string{
		path + "adversar",
		path + "agressor",
		path + "antru",
		path + "assassin",
		path + "b4",
		path + "bishop",
		path + "bouncer",
		path + "boxer",
		path + "cassius",
		path + "catfish3",
		path + "chase",
		path + "chaser",
		path + "cooper1",
		path + "cooper2",
		path + "cornerkl",
		path + "counter",
		path + "counter2",
		path + "cruiser",
		path + "cspotrun",
		path + "danimal",
		path + "dave",
		path + "di",
		path + "dirtyh",
		path + "duck",
		path + "dumbname",
		path + "etf_kid",
		path + "flyby",
		path + "fred",
		path + "friendly",
		path + "grunt",
		path + "gsmr2",
		path + "h-k",
		path + "hac_atak",
		path + "hak3",
		path + "hitnrun",
		path + "hunter",
		path + "huntlead",
		path + "intrcptr",
		path + "jagger",
		path + "jason100",
		path + "kamikaze",
		path + "killer",
		path + "leader",
		path + "leavy",
		path + "lethal",
		path + "maniac",
		path + "marvin",
		path + "mini",
		path + "ninja",
		path + "nord",
		path + "nord2",
		path + "ogre",
		path + "ogre2",
		path + "ogre3",
		path + "perizoom",
		path + "pest",
		path + "phantom",
		path + "pingpong",
		path + "politik",
		path + "pzk",
		path + "pzkmin",
		path + "quack",
		path + "quikshot",
		path + "rabbit10",
		path + "rambo3",
		path + "rapest",
		path + "reflex",
		path + "robbie",
		path + "rook",
		path + "rungun",
		path + "samurai",
		path + "scan",
		path + "scanlock",
		path + "scanner",
		path + "secro",
		path + "sentry",
		path + "shark3",
		path + "shark4",
		path + "silly",
		path + "slead",
		path + "sniper",
		path + "spinner",
		path + "spot",
		path + "squirrel",
		path + "stalker",
		path + "stush-1",
		path + "topgun",
		path + "tracker",
		path + "trial4",
		path + "twedlede",
		path + "twedledm",
		path + "venom",
		path + "watchdog",
		path + "wizard",
		path + "xecutner",
		path + "xhatch",
		path + "yal",
	}
	fmt.Printf("%d robot(s)\n", len(crobs))

	fmt.Print("Loading micro... ")
	path = fmt.Sprintf("micro%s", Separator)
	micro = []string{
		path + "caccola",
		path + "carletto",
		path + "chobin",
		path + "dream",
		path + "ld",
		path + "lucifer",
		path + "marlene",
		path + "md8",
		path + "md9",
		path + "mflash",
		path + "minizai",
		path + "pacoon",
		path + "pikachu",
		path + "pippo00a",
		path + "pippo00",
		path + "pirla",
		path + "p",
		path + "rudy",
		path + "static",
		path + "tanzen",
		path + "uhm",
		path + "zioalfa",
		path + "zzz",
	}
	fmt.Printf("%d robot(s)\n", len(micro))
}

func setupMicro() {
	fmt.Print("Loading others... ")
	others = []string{
		fmt.Sprintf("aminet%santiclock", Separator),
		fmt.Sprintf("aminet%smirobot", Separator),
		fmt.Sprintf("aminet%sschwan", Separator),
		fmt.Sprintf("aminet%stron", Separator),
		fmt.Sprintf("cplusplus%sselvaggio", Separator),
		fmt.Sprintf("cplusplus%svikingo", Separator),
	}
	fmt.Printf("%d robot(s)\n", len(others))

	var path string

	fmt.Print("Loading 1990... ")
	path = fmt.Sprintf("1990%s", Separator)
	tournament1990 = []string{
		path + "et_1",
		path + "et_2",
		path + "hunter",
		path + "nexus_1",
		path + "scanner",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1990))

	fmt.Print("Loading 1991... ")
	path = fmt.Sprintf("1991%s", Separator)
	tournament1991 = []string{
		path + "blade3",
		path + "ccyber",
		path + "diagonal",
		path + "et_3",
		path + "fdig",
		path + "genius_j",
		path + "gira",
		path + "gunner",
		path + "jazz",
		path + "nexus_2",
		path + "paolo101",
		path + "paolo77",
		path + "poor",
		path + "robocop",
		path + "runner",
		path + "seeker",
		path + "warrior2",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1991))

	fmt.Print("Loading 1992... ")
	path = fmt.Sprintf("1992%s", Separator)
	tournament1992 = []string{
		path + "ap_1",
		path + "assassin",
		path + "baeos",
		path + "banzel",
		path + "bronx-00",
		path + "bry_bry",
		path + "crazy",
		path + "d47",
		path + "daitan3",
		path + "dancer",
		path + "deluxe",
		path + "et_4",
		path + "et_5",
		path + "flash",
		path + "genesis",
		path + "hunter",
		path + "ice",
		path + "johnny",
		path + "mimo6new",
		path + "mut",
		path + "ninus6",
		path + "nl_1a",
		path + "nl_1b",
		path + "ola",
		path + "paolo",
		path + "pavido",
		path + "phobos_1",
		path + "pippo",
		path + "raid",
		path + "random",
		path + "revenge3",
		path + "robbie",
		path + "robocop2",
		path + "robocop",
		path + "superv",
		path + "t1000",
		path + "thunder",
		path + "trio",
		path + "uanino",
		path + "warrior3",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1992))

	fmt.Print("Loading 1993... ")
	path = fmt.Sprintf("1993%s", Separator)
	tournament1993 = []string{
		path + "ares",
		path + "argon",
		path + "aspide",
		path + "beast",
		path + "biro",
		path + "boom",
		path + "casual",
		path + "corner1d",
		path + "corner3",
		path + "courage",
		path + "(c)",
		path + "crob1",
		path + "deluxe_2",
		path + "didimo",
		path + "elija",
		path + "fermo",
		path + "flash2",
		path + "gunnyboy",
		path + "hell",
		path + "horse",
		path + "isaac",
		path + "kami",
		path + "lazy",
		path + "mimo13",
		path + "mohawk",
		path + "ninus17",
		path + "nl_2a",
		path + "nl_2b",
		path + "phobos_2",
		path + "pippo93",
		path + "pognant",
		path + "premana",
		path + "raid2",
		path + "rapper",
		path + "r_cyborg",
		path + "r_daneel",
		path + "robocop3",
		path + "spartaco",
		path + "target",
		path + "torneo",
		path + "vannina",
		path + "wassilij",
		path + "wolfgang",
		path + "zulu",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1993))

	fmt.Print("Loading 1994... ")
	path = fmt.Sprintf("1994%s", Separator)
	tournament1994 = []string{
		path + "anglek2",
		path + "baubau",
		path + "biro",
		path + "circlek1",
		path + "corner3b",
		path + "didimo",
		path + "dima10",
		path + "dima9",
		path + "emanuela",
		path + "ematico",
		path + "heavens",
		path + "iching",
		path + "jet",
		path + "nemesi",
		path + "ninus75",
		path + "pioppo",
		path + "pippo94b",
		path + "robot1",
		path + "robot2",
		path + "superfly",
		path + "the_dam",
		path + "t-rex",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1994))

	fmt.Print("Loading 1995... ")
	path = fmt.Sprintf("1995%s", Separator)
	tournament1995 = []string{
		path + "andrea",
		path + "b115e2",
		path + "carlo",
		path + "circle",
		path + "diablo",
		path + "flash4",
		path + "heavens",
		path + "mikezhar",
		path + "ninus99",
		path + "rocco",
		path + "sel",
		path + "skizzo",
		path + "tmii",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1995))

	fmt.Print("Loading 1996... ")
	path = fmt.Sprintf("1996%s", Separator)
	tournament1996 = []string{
		path + "andrea96",
		path + "carlo96",
		path + "drago5",
		path + "d_ray",
		path + "gpo2",
		path + "murdoc",
		path + "natas",
		path + "risk",
		path + "tronco",
		path + "yuri",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1996))

	fmt.Print("Loading 1997... ")
	path = fmt.Sprintf("1997%s", Separator)
	tournament1997 = []string{
		path + "ciccio",
		path + "drago6",
		path + "erica",
		path + "fya",
		path + "pippo97",
		path + "raid3",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1997))

	fmt.Print("Loading 1998... ")
	path = fmt.Sprintf("1998%s", Separator)
	tournament1998 = []string{
		path + "carla",
		path + "fscan",
		path + "maxheav",
		path + "pippo98",
		path + "plump",
		path + "themicro",
		path + "traker1",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1998))

	fmt.Print("Loading 1999... ")
	path = fmt.Sprintf("1999%s", Separator)
	tournament1999 = []string{
		path + "ap_5",
		path + "flash6",
		path + "mcenrobo",
		path + "nexus_2",
		path + "surrende",
		path + "themicro",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1999))

	fmt.Print("Loading 2000... ")
	// Empty for micro setup
	fmt.Printf("%d robot(s)\n", len(tournament2000))

	fmt.Print("Loading 2001... ")
	path = fmt.Sprintf("2001%s", Separator)
	tournament2001 = []string{
		path + "burrfoot",
		path + "charles",
		path + "cisc",
		path + "cobra",
		path + "copter",
		path + "gers",
		path + "grezbot",
		path + "hammer",
		path + "homer",
		path + "klr2",
		path + "kyashan",
		path + "max10",
		path + "mflash2",
		path + "microdna",
		path + "midi_zai",
		path + "mnl_1a",
		path + "mnl_1b",
		path + "murray",
		path + "neo0",
		path + "pippo1a",
		path + "pippo1b",
		path + "raistlin",
		path + "ridicol",
		path + "risc",
		path + "rudy_xp",
		path + "sdc2",
		path + "staticii",
		path + "thunder",
		path + "vampire",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2001))

	fmt.Print("Loading 2002... ")
	path = fmt.Sprintf("2002%s", Separator)
	tournament2002 = []string{
		path + "01",
		path + "adsl",
		path + "anakin",
		path + "copter_2",
		path + "corner5",
		path + "doom2099",
		path + "groucho",
		path + "idefix",
		path + "kyash_2",
		path + "marco",
		path + "mazinga",
		path + "mind",
		path + "neo_sifr",
		path + "pippo2a",
		path + "pippo2b",
		path + "regis",
		path + "scsi",
		path + "ska",
		path + "stanlio",
		path + "staticxp",
		path + "supernov",
		path + "tigre",
		path + "vaiolo",
		path + "vauban",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2002))

	fmt.Print("Loading 2003... ")
	path = fmt.Sprintf("2003%s", Separator)
	tournament2003 = []string{
		path + "730",
		path + "barbarian",
		path + "blitz",
		path + "briscolo",
		path + "bruce",
		path + "cvirus",
		path + "danica",
		path + "falco",
		path + "foursquare",
		path + "frame",
		path + "herpes",
		path + "ici",
		path + "instict",
		path + "janu",
		path + "kyash_3m",
		path + "lbr1",
		path + "lbr",
		path + "lebbra",
		path + "minicond",
		path + "morituro",
		path + "nemo",
		path + "neo_sel",
		path + "piiico",
		path + "pippo3b",
		path + "pippo3",
		path + "red_wolf",
		path + "scilla",
		path + "sirio",
		path + "tartaruga",
		path + "valevan",
		path + "virus",
		path + "yoda",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2003))

	fmt.Print("Loading 2004... ")
	path = fmt.Sprintf("2004%s", Separator)
	tournament2004 = []string{
		path + "adam",
		path + "!caos",
		path + "ciclope",
		path + "coyote",
		path + "diodo",
		path + "gostar",
		path + "gotar2",
		path + "gotar",
		path + "irap",
		path + "magneto",
		path + "n3g4_jr",
		path + "new_mini",
		path + "pippo04a",
		path + "pippo04b",
		path + "poldo",
		path + "puma",
		path + "rat-man",
		path + "ravatto",
		path + "rotar",
		path + "selim_b",
		path + "unlimited",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2004))

	fmt.Print("Loading 2007... ")
	path = fmt.Sprintf("2007%s", Separator)
	tournament2007 = []string{
		path + "angel",
		path + "back",
		path + "brontolo",
		path + "electron",
		path + "gongolo",
		path + "microbo1",
		path + "microbo2",
		path + "pippo07a",
		path + "pippo07b",
		path + "pisolo",
		path + "pyro",
		path + "tobey",
		path + "t",
		path + "zigozago",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2007))

	fmt.Print("Loading 2010... ")
	path = fmt.Sprintf("2010%s", Separator)
	tournament2010 = []string{
		path + "copia",
		path + "eurialo",
		path + "macchia",
		path + "niso",
		path + "pippo10a",
		path + "stitch",
		path + "sweat",
		path + "taglia",
		path + "wall-e",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2010))

	fmt.Print("Loading 2011... ")
	path = fmt.Sprintf("2011%s", Separator)
	tournament2011 = []string{
		path + "ataman",
		path + "coeurl",
		path + "minion",
		path + "pain",
		path + "piperita",
		path + "pippo11a",
		path + "pippo11b",
		path + "tannhause",
		path + "unmaldestr",
		path + "wall-e_ii",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2011))

	fmt.Print("Loading 2012... ")
	path = fmt.Sprintf("2012%s", Separator)
	tournament2012 = []string{
		path + "avoider",
		path + "beat",
		path + "china",
		path + "easyjet",
		path + "flash8c",
		path + "flash8e",
		path + "grezbot2",
		path + "lycan",
		path + "pippo12a",
		path + "pippo12b",
		path + "puffomic",
		path + "ryanair",
		path + "silversurf",
		path + "wall-e_iii",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2012))

	fmt.Print("Loading 2013... ")
	path = fmt.Sprintf("2013%s", Separator)
	tournament2013 = []string{
		path + "axolotl",
		path + "destro",
		path + "osvaldo",
		path + "pippo13a",
		path + "pippo13b",
		path + "pray",
		path + "wall-e_iv",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2013))

	fmt.Print("Loading 2015... ")
	path = fmt.Sprintf("2015%s", Separator)
	tournament2015 = []string{
		path + "antman",
		path + "aswhup",
		path + "avoider", // same as 2012
		path + "babadook",
		path + "colour",
		path + "coppi15mc1",
		path + "flash9",
		path + "linabo15",
		path + "mike3",
		path + "pippo15a",
		path + "pippo15b",
		path + "puppet",
		path + "randguard",
		path + "salippo",
		path + "sidewalk",
		path + "tux",
		path + "tyrion",
		path + "wall-e_v",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2015))

	fmt.Print("Loading 2020... ")
	path = fmt.Sprintf("2020%s", Separator)
	tournament2020 = []string{
		path + "antman_20",
		path + "b4b",
		path + "brexit",
		path + "coppi20mc1",
		path + "discotek",
		path + "flash10",
		path + "pippo20a",
		path + "pippo20b",
		path + "wall-e_vi",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2020))

	fmt.Print("Loading 2025... ")
	path = fmt.Sprintf("2025%s", Separator)
	tournament2025 = []string{
		path + "extrasmall",
		path + "flash11",
		path + "kerberos",
		path + "pippo25a",
		path + "sentry2",
		path + "sentry3",
		path + "sgorbio",
		path + "slant6",
		path + "trouble3",
		path + "ultron_25",
		path + "wall-e_vii",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2025))

	fmt.Print("Loading crobs... ")
	path = fmt.Sprintf("crobs%s", Separator)
	crobs = []string{
		path + "adversar",
		path + "agressor",
		path + "assassin",
		path + "b4",
		path + "bishop",
		path + "bouncer",
		path + "cassius",
		path + "catfish3",
		path + "chase",
		path + "chaser",
		path + "cornerkl",
		path + "counter2",
		path + "cruiser",
		path + "cspotrun",
		path + "danimal",
		path + "dave",
		path + "di",
		path + "dirtyh",
		path + "duck",
		path + "etf_kid",
		path + "flyby",
		path + "fred",
		path + "grunt",
		path + "gsmr2",
		path + "hac_atak",
		path + "hak3",
		path + "hitman",
		path + "h-k",
		path + "hunter",
		path + "huntlead",
		path + "intrcptr",
		path + "kamikaze",
		path + "killer",
		path + "leader",
		path + "marvin",
		path + "micro",
		path + "mini",
		path + "ninja",
		path + "nord2",
		path + "nord",
		path + "ogre2",
		path + "ogre",
		path + "pest",
		path + "phantom",
		path + "pingpong",
		path + "pzkmin",
		path + "pzk",
		path + "quack",
		path + "quikshot",
		path + "rabbit10",
		path + "rabbit",
		path + "rambo3",
		path + "rapest",
		path + "reflex",
		path + "rungun",
		path + "scanlock",
		path + "scanner",
		path + "scan",
		path + "sentry",
		path + "silly",
		path + "slead",
		path + "spinner",
		path + "spot",
		path + "squirrel",
		path + "stush-1",
		path + "topgun",
		path + "tracker",
		path + "twedlede",
		path + "twedledm",
		path + "venom",
		path + "watchdog",
		path + "xecutner",
		path + "xhatch",
		path + "yal",
	}
	fmt.Printf("%d robot(s)\n", len(crobs))

	fmt.Print("Loading micro... ")
	path = fmt.Sprintf("micro%s", Separator)
	micro = []string{
		path + "caccola",
		path + "carletto",
		path + "chobin",
		path + "dream",
		path + "ld",
		path + "lucifer",
		path + "marlene",
		path + "md8",
		path + "md9",
		path + "mflash",
		path + "minizai",
		path + "pacoon",
		path + "pikachu",
		path + "pippo00a",
		path + "pippo00",
		path + "pirla",
		path + "p",
		path + "rudy",
		path + "static",
		path + "tanzen",
		path + "uhm",
		path + "zioalfa",
		path + "zzz",
	}
	fmt.Printf("%d robot(s)\n", len(micro))
}

func setup() {
	fmt.Print("Loading others... ")
	others = []string{
		fmt.Sprintf("aminet%santiclock", Separator),
		fmt.Sprintf("aminet%sbeaver", Separator),
		fmt.Sprintf("aminet%sblindschl", Separator),
		fmt.Sprintf("aminet%sblindschl2", Separator),
		fmt.Sprintf("aminet%smirobot", Separator),
		fmt.Sprintf("aminet%sopfer", Separator),
		fmt.Sprintf("aminet%sschwan", Separator),
		fmt.Sprintf("aminet%stron", Separator),
		fmt.Sprintf("cplusplus%sselvaggio", Separator),
		fmt.Sprintf("cplusplus%svikingo", Separator),
	}
	fmt.Printf("%d robot(s)\n", len(others))
	var path string

	fmt.Print("Loading 1990... ")
	path = fmt.Sprintf("1990%s", Separator)
	tournament1990 = []string{
		path + "et_1",
		path + "et_2",
		path + "hunter",
		path + "killer",
		path + "nexus_1",
		path + "rob1",
		path + "scanner",
		path + "york",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1990))

	fmt.Print("Loading 1991... ")
	path = fmt.Sprintf("1991%s", Separator)
	tournament1991 = []string{
		path + "blade3",
		path + "casimiro",
		path + "ccyber",
		path + "clover",
		path + "diagonal",
		path + "et_3",
		path + "f1",
		path + "fdig",
		path + "geltrude",
		path + "genius_j",
		path + "gira",
		path + "gunner",
		path + "jazz",
		path + "nexus_2",
		path + "paolo101",
		path + "paolo77",
		path + "poor",
		path + "qibo",
		path + "robocop",
		path + "runner",
		path + "sara_6",
		path + "seeker",
		path + "warrior2",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1991))

	fmt.Print("Loading 1992... ")
	path = fmt.Sprintf("1992%s", Separator)
	tournament1992 = []string{
		path + "666",
		path + "ap_1",
		path + "assassin",
		path + "baeos",
		path + "banzel",
		path + "bronx-00",
		path + "bry_bry",
		path + "crazy",
		path + "cube",
		path + "cw",
		path + "d47",
		path + "daitan3",
		path + "dancer",
		path + "deluxe",
		path + "dorsai",
		path + "et_4",
		path + "et_5",
		path + "flash",
		path + "genesis",
		path + "hunter",
		path + "ice",
		path + "jack",
		path + "jager",
		path + "johnny",
		path + "lead1",
		path + "marika",
		path + "mimo6new",
		path + "mrcc",
		path + "mut",
		path + "ninus6",
		path + "nl_1a",
		path + "nl_1b",
		path + "ola",
		path + "paolo",
		path + "pavido",
		path + "phobos_1",
		path + "pippo92",
		path + "pippo",
		path + "raid",
		path + "random",
		path + "revenge3",
		path + "robbie",
		path + "robocop2",
		path + "robocop",
		path + "sassy",
		path + "spider",
		path + "sp",
		path + "superv",
		path + "t1000",
		path + "thunder",
		path + "triangol",
		path + "trio",
		path + "uanino",
		path + "warrior3",
		path + "xdraw2",
		path + "zorro",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1992))

	fmt.Print("Loading 1993... ")
	path = fmt.Sprintf("1993%s", Separator)
	tournament1993 = []string{
		path + "am_174",
		path + "ap_2",
		path + "ares",
		path + "argon",
		path + "aspide",
		path + "beast",
		path + "biro",
		path + "blade8",
		path + "boom",
		path + "brain",
		path + "cantor",
		path + "castore",
		path + "casual",
		path + "corner1d",
		path + "corner3",
		path + "courage",
		path + "(c)",
		path + "crob1",
		path + "deluxe_2",
		path + "deluxe_3",
		path + "didimo",
		path + "duke",
		path + "elija",
		path + "fermo",
		path + "flash2",
		path + "food5",
		path + "godel",
		path + "gunnyboy",
		path + "hamp1",
		path + "hamp2",
		path + "hell",
		path + "horse",
		path + "isaac",
		path + "kami",
		path + "lazy",
		path + "mimo13",
		path + "mister2",
		path + "mister3",
		path + "mohawk",
		path + "mutation",
		path + "ninus17",
		path + "nl_2a",
		path + "nl_2b",
		path + "p68",
		path + "p69",
		path + "penta",
		path + "phobos_2",
		path + "pippo93",
		path + "pognant",
		path + "poirot",
		path + "polluce",
		path + "premana",
		path + "puyopuyo",
		path + "raid2",
		path + "rapper",
		path + "r_cyborg",
		path + "r_daneel",
		path + "robocop3",
		path + "spartaco",
		path + "target",
		path + "tm",
		path + "torneo",
		path + "vannina",
		path + "vocus",
		path + "warrior4",
		path + "wassilij",
		path + "wolfgang",
		path + "zulu",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1993))

	fmt.Print("Loading 1994... ")
	path = fmt.Sprintf("1994%s", Separator)
	tournament1994 = []string{
		path + "8bismark",
		path + "anglek2",
		path + "apache",
		path + "bachopin",
		path + "baubau",
		path + "biro",
		path + "blob",
		path + "circlek1",
		path + "corner3b",
		path + "corner4",
		path + "deluxe_4",
		path + "deluxe_5",
		path + "didimo",
		path + "dima10",
		path + "dima9",
		path + "emanuela",
		path + "ematico",
		path + "fastfood",
		path + "flash3",
		path + "funky",
		path + "giali1",
		path + "hal9000",
		path + "heavens",
		path + "horse2",
		path + "iching",
		path + "jet",
		path + "ken",
		path + "lazyii",
		path + "matrox",
		path + "maverick",
		path + "miaomiao",
		path + "nemesi",
		path + "ninus75",
		path + "patcioca",
		path + "pioppo",
		path + "pippo94a",
		path + "pippo94b",
		path + "polipo",
		path + "randwall",
		path + "robot1",
		path + "robot2",
		path + "sdix3",
		path + "sgnaus",
		path + "shadow",
		path + "superfly",
		path + "the_dam",
		path + "t-rex",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1994))

	fmt.Print("Loading 1995... ")
	path = fmt.Sprintf("1995%s", Separator)
	tournament1995 = []string{
		path + "andrea",
		path + "animal",
		path + "apache95",
		path + "archer",
		path + "b115e2",
		path + "b52",
		path + "biro",
		path + "boss",
		path + "camillo",
		path + "carlo",
		path + "circle",
		path + "cri95",
		path + "diablo",
		path + "flash4",
		path + "hal9000",
		path + "heavens",
		path + "horse3",
		path + "kenii",
		path + "losendos",
		path + "mikezhar",
		path + "ninus99",
		path + "paccu",
		path + "passion",
		path + "peribolo",
		path + "pippo95",
		path + "rambo",
		path + "rocco",
		path + "saxy",
		path + "sel",
		path + "skizzo",
		path + "star",
		path + "stinger",
		path + "tabori-1",
		path + "tabori-2",
		path + "tequila",
		path + "tmii",
		path + "tox",
		path + "t-rex",
		path + "tricky",
		path + "twins",
		path + "upv-9596",
		path + "xenon",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1995))

	fmt.Print("Loading 1996... ")
	path = fmt.Sprintf("1996%s", Separator)
	tournament1996 = []string{
		path + "aleph",
		path + "andrea96",
		path + "ap_4",
		path + "carlo96",
		path + "diablo2",
		path + "drago5",
		path + "d_ray",
		path + "fb3",
		path + "gevbass",
		path + "golem",
		path + "gpo2",
		path + "hal9000",
		path + "heavnew",
		path + "hider2",
		path + "infinity",
		path + "jaja",
		path + "memories",
		path + "murdoc",
		path + "natas",
		path + "newb52",
		path + "pacio",
		path + "pippo96a",
		path + "pippo96b",
		path + "!",
		path + "risk",
		path + "robot1",
		path + "robot2",
		path + "rudolf",
		path + "second3",
		path + "s-seven",
		path + "tatank_3",
		path + "tronco",
		path + "uht",
		path + "xabaras",
		path + "yuri",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1996))

	fmt.Print("Loading 1997... ")
	path = fmt.Sprintf("1997%s", Separator)
	tournament1997 = []string{
		path + "1&1",
		path + "abyss",
		path + "ai1",
		path + "andrea97",
		path + "arale",
		path + "belva",
		path + "carlo97",
		path + "ciccio",
		path + "colossus",
		path + "diablo3",
		path + "diabolik",
		path + "drago6",
		path + "erica",
		path + "fable",
		path + "flash5",
		path + "fya",
		path + "gevbass2",
		path + "golem2",
		path + "gundam",
		path + "hal9000",
		path + "jedi",
		path + "kill!",
		path + "me-110c",
		path + "ncmplt",
		path + "paperone",
		path + "pippo97",
		path + "raid3",
		path + "robivinf",
		path + "rudolf_2",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1997))

	fmt.Print("Loading 1998... ")
	path = fmt.Sprintf("1998%s", Separator)
	tournament1998 = []string{
		path + "ai2",
		path + "bartali",
		path + "carla",
		path + "coppi",
		path + "dia",
		path + "dicin",
		path + "eva00",
		path + "eva01",
		path + "freedom",
		path + "fscan",
		path + "goblin",
		path + "goldrake",
		path + "hal9000",
		path + "heavnew",
		path + "maxheav",
		path + "ninja",
		path + "paranoid",
		path + "pippo98",
		path + "plump",
		path + "quarto",
		path + "rattolo",
		path + "rudolf_3",
		path + "son-goku",
		path + "sottolin",
		path + "stay",
		path + "stighy98",
		path + "themicro",
		path + "titania",
		path + "tornado",
		path + "traker1",
		path + "traker2",
		path + "vision",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1998))

	fmt.Print("Loading 1999... ")
	path = fmt.Sprintf("1999%s", Separator)
	tournament1999 = []string{
		path + "11",
		path + "aeris",
		path + "akira",
		path + "alezai17",
		path + "alfa99",
		path + "alien",
		path + "ap_5",
		path + "bastrd!!",
		path + "cancer",
		path + "carlo99",
		path + "#cimice#",
		path + "cortez",
		path + "cyborg",
		path + "dario",
		path + "dav46",
		path + "defender",
		path + "elisir",
		path + "flash6",
		path + "hal9000",
		path + "ilbestio",
		path + "jedi2",
		path + "ka_aroth",
		path + "kakakatz",
		path + "lukather",
		path + "mancino",
		path + "marko",
		path + "mcenrobo",
		path + "m_hingis",
		path + "minatela",
		path + "new",
		path + "nexus_2",
		path + "nl_3a",
		path + "nl_3b",
		path + "obiwan",
		path + "omega99",
		path + "panduro",
		path + "panic",
		path + "pippo99",
		path + "pizarro",
		path + "quarto",
		path + "quingon",
		path + "rudolf_4",
		path + "satana",
		path + "shock",
		path + "songohan",
		path + "stealth",
		path + "storm",
		path + "surrende",
		path + "t1001",
		path + "themicro",
		path + "titania2",
		path + "vibrsper",
		path + "zero",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1999))

	fmt.Print("Loading 2000... ")
	path = fmt.Sprintf("2000%s", Separator)
	tournament2000 = []string{
		path + "7di9",
		path + "bach_2k",
		path + "beholder",
		path + "boom",
		path + "carlo2k",
		path + "coppi_2k",
		path + "daryl",
		path + "dav2000",
		path + "def2",
		path + "defender",
		path + "doppia_g",
		path + "flash7",
		path + "fremen",
		path + "gengis",
		path + "jedi3",
		path + "kongzill",
		path + "mancino",
		path + "marine",
		path + "m_hingis", // same as 1999 - just different head comments
		path + "mrsatan",
		path + "navaho",
		path + "new2",
		path + "newzai17",
		path + "nl_4a",
		path + "nl_4b",
		path + "rudolf_5",
		path + "sharp",
		path + "touch",
		path + "vegeth",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2000))

	fmt.Print("Loading 2001... ")
	path = fmt.Sprintf("2001%s", Separator)
	tournament2001 = []string{
		path + "4ever",
		path + "artu",
		path + "athlon",
		path + "bati",
		path + "bigkarl",
		path + "borg",
		path + "burrfoot",
		path + "charles",
		path + "cisc",
		path + "cobra",
		path + "copter",
		path + "defender",
		path + "disco",
		path + "dnablack",
		path + "dna",
		path + "fizban",
		path + "gers",
		path + "grezbot",
		path + "hammer",
		path + "harris",
		path + "heavnew",
		path + "homer",
		path + "jedi4",
		path + "klr2",
		path + "kyashan",
		path + "max10",
		path + "megazai",
		path + "merlino",
		path + "mflash2",
		path + "microdna",
		path + "midi_zai",
		path + "mnl_1a",
		path + "mnl_1b",
		path + "murray",
		path + "neo0",
		path + "nl_5a",
		path + "nl_5b",
		path + "pentium4",
		path + "pippo1a",
		path + "pippo1b",
		path + "raistlin",
		path + "ridicol",
		path + "risc",
		path + "rudolf_6",
		path + "rudy_xp",
		path + "sdc2",
		path + "sharp2",
		path + "staticii",
		path + "thunder",
		path + "vampire",
		path + "xeon",
		path + "zifnab",
		path + "zombie",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2001))

	fmt.Print("Loading 2002... ")
	path = fmt.Sprintf("2002%s", Separator)
	tournament2002 = []string{
		path + "01",
		path + "adsl",
		path + "anakin",
		path + "asterix",
		path + "attila",
		path + "bruenor",
		path + "colera",
		path + "colosseum",
		path + "copter_2",
		path + "corner5",
		path + "doom2099",
		path + "drizzt",
		path + "dynamite",
		path + "enigma",
		path + "groucho",
		path + "halman",
		path + "harpo",
		path + "idefix",
		path + "jedi5",
		path + "kyash_2",
		path + "marco",
		path + "mazinga",
		path + "medioman",
		path + "mg_one",
		path + "mind",
		path + "moveon",
		path + "neo_sifr",
		path + "obelix",
		path + "ollio",
		path + "padawan",
		path + "peste",
		path + "pippo2a",
		path + "pippo2b",
		path + "regis",
		path + "remus",
		path + "romulus",
		path + "rudolf_7",
		path + "scsi",
		path + "serse",
		path + "ska",
		path + "stanlio",
		path + "staticxp",
		path + "supernov",
		path + "theslayer",
		path + "tifo",
		path + "tigre",
		path + "todos",
		path + "tomahawk",
		path + "vaiolo",
		path + "vauban",
		path + "wulfgar",
		path + "yerba",
		path + "yoyo",
		path + "zorn",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2002))

	fmt.Print("Loading 2003... ")
	path = fmt.Sprintf("2003%s", Separator)
	tournament2003 = []string{
		path + "730",
		path + "adrian",
		path + "aladino",
		path + "alcadia",
		path + "ares",
		path + "barbarian",
		path + "blitz",
		path + "briscolo",
		path + "bruce",
		path + "cadderly",
		path + "cariddi",
		path + "crossover",
		path + "cvirus2",
		path + "cvirus",
		path + "cyborg_2",
		path + "danica",
		path + "dave",
		path + "druzil",
		path + "dynacond",
		path + "elminster",
		path + "falco",
		path + "foursquare",
		path + "frame",
		path + "harlock",
		path + "herpes",
		path + "ici",
		path + "instict",
		path + "irpef",
		path + "janick",
		path + "janu",
		path + "jedi6",
		path + "knt",
		path + "kyash_3c",
		path + "kyash_3m",
		path + "lbr1",
		path + "lbr",
		path + "lebbra",
		path + "maxicond",
		path + "mg_two",
		path + "minicond",
		path + "morituro",
		path + "nautilus",
		path + "nemo",
		path + "neo_sel",
		path + "orione",
		path + "piiico",
		path + "pippo3b",
		path + "pippo3",
		path + "red_wolf",
		path + "rudolf_8",
		path + "scanner",
		path + "scilla",
		path + "sirio",
		path + "sith",
		path + "sky",
		path + "spaceman",
		path + "tartaruga",
		path + "unico",
		path + "valevan",
		path + "virus2",
		path + "virus3",
		path + "virus4",
		path + "virus",
		path + "yoda",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2003))

	fmt.Print("Loading 2004... ")
	path = fmt.Sprintf("2004%s", Separator)
	tournament2004 = []string{
		path + "adam",
		path + "!alien",
		path + "bjt",
		path + "b_selim",
		path + "!caos",
		path + "ciclope",
		path + "confusion",
		path + "coyote",
		path + "diodo",
		path + "!dna",
		path + "fire",
		path + "fisco",
		path + "frankie",
		path + "geriba",
		path + "goofy",
		path + "gostar",
		path + "gotar2",
		path + "gotar",
		path + "irap",
		path + "ire",
		path + "ires",
		path + "jedi7",
		path + "magneto",
		path + "mg_three",
		path + "mosfet",
		path + "m_selim",
		path + "multics",
		path + "mystica",
		path + "n3g4_jr",
		path + "n3g4tivo",
		path + "new_mini",
		path + "pippo04a",
		path + "pippo04b",
		path + "poldo",
		path + "puma",
		path + "rat-man",
		path + "ravatto",
		path + "revo",
		path + "rotar",
		path + "rudolf_9",
		path + "selim_b",
		path + "tempesta",
		path + "unlimited",
		path + "wgdi",
		path + "zener",
		path + "!zeus",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2004))

	fmt.Print("Loading 2007... ")
	path = fmt.Sprintf("2007%s", Separator)
	tournament2007 = []string{
		path + "angel",
		path + "back",
		path + "brontolo",
		path + "colosso",
		path + "electron",
		path + "e",
		path + "gongolo",
		path + "iceman",
		path + "jedi8",
		path + "macro1",
		path + "mammolo",
		path + "microbo1",
		path + "microbo2",
		path + "midi1",
		path + "neutron",
		path + "nustyle",
		path + "pippo07a",
		path + "pippo07b",
		path + "pisolo",
		path + "proton",
		path + "proud",
		path + "pyro",
		path + "rudolf_x",
		path + "rythm",
		path + "tobey",
		path + "t",
		path + "zigozago",
		path + "z",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2007))

	fmt.Print("Loading 2010... ")
	path = fmt.Sprintf("2010%s", Separator)
	tournament2010 = []string{
		path + "buffy",
		path + "cancella",
		path + "change",
		path + "copia",
		path + "enkidu",
		path + "eurialo",
		path + "gantu",
		path + "hal9010",
		path + "incolla",
		path + "jedi9",
		path + "jumba",
		path + "macchia",
		path + "niso",
		path + "party",
		path + "pippo10a",
		path + "reuben",
		path + "stitch",
		path + "suddenly",
		path + "sweat",
		path + "taglia",
		path + "toppa",
		path + "wall-e",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2010))

	fmt.Print("Loading 2011... ")
	path = fmt.Sprintf("2011%s", Separator)
	tournament2011 = []string{
		path + "armin",
		path + "ataman",
		path + "coeurl",
		path + "digitale",
		path + "gerty",
		path + "grendizer",
		path + "gru",
		path + "guntank",
		path + "hal9011",
		path + "jedi10",
		path + "jeeg",
		path + "minion",
		path + "nikita",
		path + "origano",
		path + "ortica",
		path + "pain",
		path + "piperita",
		path + "pippo11a",
		path + "pippo11b",
		path + "smart",
		path + "tannhause",
		path + "tantalo",
		path + "unmaldestr",
		path + "vain",
		path + "vector",
		path + "wall-e_ii",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2011))

	fmt.Print("Loading 2012... ")
	path = fmt.Sprintf("2012%s", Separator)
	tournament2012 = []string{
		path + "avoider",
		path + "beat",
		path + "british",
		path + "camille",
		path + "china",
		path + "cliche",
		path + "crazy96",
		path + "dampyr",
		path + "draka",
		path + "easyjet",
		path + "flash8c",
		path + "flash8e",
		path + "gerty2",
		path + "grezbot2",
		path + "gunnyb29",
		path + "hal9012",
		path + "jedi11",
		path + "life",
		path + "lufthansa",
		path + "lycan",
		path + "mister2b",
		path + "mister3b",
		path + "pippo12a",
		path + "pippo12b",
		path + "power",
		path + "puffomac",
		path + "puffomic",
		path + "puffomid",
		path + "q",
		path + "ryanair",
		path + "silversurf",
		path + "torchio",
		path + "wall-e_iii",
		path + "yeti",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2012))

	fmt.Print("Loading 2013... ")
	path = fmt.Sprintf("2013%s", Separator)
	tournament2013 = []string{
		path + "axolotl",
		path + "destro",
		path + "eternity",
		path + "frisa_13",
		path + "gerty3",
		path + "ghostrider",
		path + "guanaco",
		path + "gunnyb13",
		path + "hal9013",
		path + "jarvis",
		path + "jedi12",
		path + "john_blaze",
		path + "lamela",
		path + "lancia13",
		path + "leopon",
		path + "ncc-1701",
		path + "okapi",
		path + "ortona_13",
		path + "osvaldo",
		path + "pippo13a",
		path + "pippo13b",
		path + "pjanic",
		path + "pray",
		path + "ride",
		path + "ug2k",
		path + "wall-e_iv",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2013))

	fmt.Print("Loading 2015... ")
	path = fmt.Sprintf("2015%s", Separator)
	tournament2015 = []string{
		path + "antman",
		path + "aswhup",
		path + "avoider", // same as 2012
		path + "babadook",
		path + "bttf",
		path + "circles15",
		path + "colour",
		path + "coppi15ma1",
		path + "coppi15ma2",
		path + "coppi15mc1",
		path + "coppi15md1",
		path + "corbu15",
		path + "dlrn",
		path + "flash9",
		path + "frank15",
		path + "g13-14",
		path + "gargantua",
		path + "gerty4",
		path + "hal9015",
		path + "hulk",
		path + "ironman_15",
		path + "jedi13",
		path + "linabo15",
		path + "lluke",
		path + "mcfly",
		path + "mies15",
		path + "mike3",
		path + "misdemeano",
		path + "music",
		path + "one",
		path + "pantagruel",
		path + "pippo15a",
		path + "pippo15b",
		path + "puppet",
		path + "randguard",
		path + "salippo",
		path + "sidewalk",
		path + "the_old",
		path + "thor",
		path + "tux",
		path + "tyrion",
		path + "wall-e_v",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2015))

	fmt.Print("Loading 2020... ")
	path = fmt.Sprintf("2020%s", Separator)
	tournament2020 = []string{
		path + "antman_20",
		path + "b4b",
		path + "brexit",
		path + "carillon",
		path + "coppi20ma1",
		path + "coppi20ma2",
		path + "coppi20mc1",
		path + "coppi20md1",
		path + "discotek",
		path + "dreamland",
		path + "flash10",
		path + "gerty5",
		path + "hal9020",
		path + "hulk_20",
		path + "ironman_20",
		path + "jarvis2",
		path + "jedi14",
		path + "leavy2",
		path + "loneliness",
		path + "pippo20a",
		path + "pippo20b",
		path + "thor_20",
		path + "wall-e_vi",
		path + "wizard2",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2020))

	fmt.Print("Loading 2025... ")
	path = fmt.Sprintf("2025%s", Separator)
	tournament2025 = []string{
		path + "blue", // late entry
		path + "extrasmall",
		path + "flash11",
		path + "gerty6",
		path + "hal9025",
		path + "hulk_25",
		path + "hydra",
		path + "ironman_25",
		path + "jedi15",
		path + "kerberos",
		path + "meeseeks1",
		path + "meeseeks2",
		path + "nova",
		path + "pippo25a",
		path + "rabbitc",
		path + "rotaprinc8",
		path + "sentry2",
		path + "sentry3",
		path + "sgorbio",
		path + "slant6",
		path + "supremo",
		path + "thor_25",
		path + "trouble3",
		path + "ultron_25",
		path + "wall-e_vii",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2025))

	fmt.Print("Loading crobs... ")
	path = fmt.Sprintf("crobs%s", Separator)
	crobs = []string{
		path + "adversar",
		path + "agressor",
		path + "antru",
		path + "assassin",
		path + "b4",
		path + "bishop",
		path + "bouncer",
		path + "boxer",
		path + "cassius",
		path + "catfish3",
		path + "chase",
		path + "chaser",
		path + "cooper1",
		path + "cooper2",
		path + "cornerkl",
		path + "counter",
		path + "counter2",
		path + "cruiser",
		path + "cspotrun",
		path + "danimal",
		path + "dave",
		path + "di",
		path + "dirtyh",
		path + "duck",
		path + "dumbname",
		path + "etf_kid",
		path + "flyby",
		path + "fred",
		path + "friendly",
		path + "grunt",
		path + "gsmr2",
		path + "h-k",
		path + "hac_atak",
		path + "hak3",
		path + "hitnrun",
		path + "hunter",
		path + "huntlead",
		path + "intrcptr",
		path + "jagger",
		path + "jason100",
		path + "kamikaze",
		path + "killer",
		path + "leader",
		path + "leavy",
		path + "lethal",
		path + "maniac",
		path + "marvin",
		path + "mini",
		path + "ninja",
		path + "nord",
		path + "nord2",
		path + "ogre",
		path + "ogre2",
		path + "ogre3",
		path + "perizoom",
		path + "pest",
		path + "phantom",
		path + "pingpong",
		path + "politik",
		path + "pzk",
		path + "pzkmin",
		path + "quack",
		path + "quikshot",
		path + "rabbit10",
		path + "rambo3",
		path + "rapest",
		path + "reflex",
		path + "robbie",
		path + "rook",
		path + "rungun",
		path + "samurai",
		path + "scan",
		path + "scanlock",
		path + "scanner",
		path + "secro",
		path + "sentry",
		path + "shark3",
		path + "shark4",
		path + "silly",
		path + "slead",
		path + "sniper",
		path + "spinner",
		path + "spot",
		path + "squirrel",
		path + "stalker",
		path + "stush-1",
		path + "topgun",
		path + "tracker",
		path + "trial4",
		path + "twedlede",
		path + "twedledm",
		path + "venom",
		path + "watchdog",
		path + "wizard",
		path + "xecutner",
		path + "xhatch",
		path + "yal",
	}
	fmt.Printf("%d robot(s)\n", len(crobs))

	fmt.Print("Loading micro... ")
	path = fmt.Sprintf("micro%s", Separator)
	micro = []string{
		path + "caccola",
		path + "carletto",
		path + "chobin",
		path + "dream",
		path + "ld",
		path + "lucifer",
		path + "marlene",
		path + "md8",
		path + "md9",
		path + "mflash",
		path + "minizai",
		path + "pacoon",
		path + "pikachu",
		path + "pippo00a",
		path + "pippo00",
		path + "pirla",
		path + "p",
		path + "rudy",
		path + "static",
		path + "tanzen",
		path + "uhm",
		path + "zioalfa",
		path + "zzz",
	}
	fmt.Printf("%d robot(s)\n", len(micro))
}
