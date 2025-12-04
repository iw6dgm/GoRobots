package main

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"slices"
	"strings"
)

// Constants
const (
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
	var withConflicts bool

	for {
		attempts++
		fmt.Printf("ATTEMPT : %d\n", attempts)
		shuffle()
		collect()
		withConflicts = alternativePairing()
		if !withConflicts {
			break
		}
	}

	show()
	buildConfigFileYAML()
	//buildSQLInserts() // optional - not needed if using `tournament` scripts
}

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

// Show pairings YAML format
func buildConfigFileYAML() {
	n := 1
	for _, round := range rounds {
		if len(round) > 0 {
			fmt.Printf("------- CFG group%d ------\n", n)
			fmt.Printf("matchF2F: %d\nmatch3VS3: %d\nmatch4VS4: %d\nsourcePath: '.'\n", matchF2F, match3vs3, match4vs4)
			fmt.Printf("label: '%s%d'\n", label, n)

			var sb strings.Builder
			sb.WriteString("listRobots: [\n")
			for i, s := range round {
				if i != 0 {
					sb.WriteString(",\n")
				}
				sb.WriteString("'")
				sb.WriteString(s)
				sb.WriteString("'")
			}
			sb.WriteString("\n]")
			fmt.Println(sb.String())
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

func setupMidi() {
	fmt.Print("Loading others... ")
	others = []string{
		"aminet/anticlock",
		"aminet/beaver",
		"aminet/blindschl",
		"aminet/blindschl2",
		"aminet/mirobot",
		"aminet/opfer",
		"aminet/schwan",
		"aminet/tron",
		"cplusplus/selvaggio",
		"cplusplus/vikingo",
	}
	fmt.Printf("%d robot(s)\n", len(others))

	fmt.Print("Loading 1990... ")
	tournament1990 = []string{
		"1990/et_1",
		"1990/et_2",
		"1990/hunter",
		"1990/killer",
		"1990/nexus_1",
		"1990/rob1",
		"1990/scanner",
		"1990/york",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1990))

	fmt.Print("Loading 1991... ")
	tournament1991 = []string{
		"1991/blade3",
		"1991/casimiro",
		"1991/ccyber",
		"1991/clover",
		"1991/diagonal",
		"1991/et_3",
		"1991/f1",
		"1991/fdig",
		"1991/geltrude",
		"1991/genius_j",
		"1991/gira",
		"1991/gunner",
		"1991/jazz",
		"1991/nexus_2",
		"1991/paolo101",
		"1991/paolo77",
		"1991/poor",
		"1991/qibo",
		"1991/robocop",
		"1991/runner",
		"1991/sara_6",
		"1991/seeker",
		"1991/warrior2",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1991))

	fmt.Print("Loading 1992... ")
	tournament1992 = []string{
		"1992/666",
		"1992/ap_1",
		"1992/assassin",
		"1992/baeos",
		"1992/banzel",
		"1992/bronx-00",
		"1992/bry_bry",
		"1992/crazy",
		"1992/cube",
		"1992/cw",
		"1992/d47",
		"1992/daitan3",
		"1992/dancer",
		"1992/deluxe",
		"1992/dorsai",
		"1992/et_4",
		"1992/et_5",
		"1992/flash",
		"1992/genesis",
		"1992/hunter",
		"1992/ice",
		"1992/jack",
		"1992/jager",
		"1992/johnny",
		"1992/lead1",
		"1992/marika",
		"1992/mimo6new",
		"1992/mrcc",
		"1992/mut",
		"1992/ninus6",
		"1992/nl_1a",
		"1992/nl_1b",
		"1992/ola",
		"1992/paolo",
		"1992/pavido",
		"1992/phobos_1",
		"1992/pippo92",
		"1992/pippo",
		"1992/raid",
		"1992/random",
		"1992/revenge3",
		"1992/robbie",
		"1992/robocop2",
		"1992/robocop",
		"1992/sassy",
		"1992/spider",
		"1992/sp",
		"1992/superv",
		"1992/t1000",
		"1992/thunder",
		"1992/triangol",
		"1992/trio",
		"1992/uanino",
		"1992/warrior3",
		"1992/xdraw2",
		"1992/zorro",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1992))

	fmt.Print("Loading 1993... ")
	tournament1993 = []string{
		"1993/am_174",
		"1993/ap_2",
		"1993/ares",
		"1993/argon",
		"1993/aspide",
		"1993/beast",
		"1993/biro",
		"1993/blade8",
		"1993/boom",
		"1993/brain",
		"1993/cantor",
		"1993/castore",
		"1993/casual",
		"1993/corner1d",
		"1993/corner3",
		"1993/courage",
		"1993/(c)",
		"1993/crob1",
		"1993/deluxe_2",
		"1993/deluxe_3",
		"1993/didimo",
		"1993/duke",
		"1993/elija",
		"1993/fermo",
		"1993/flash2",
		"1993/food5",
		"1993/godel",
		"1993/gunnyboy",
		"1993/hamp1",
		"1993/hamp2",
		"1993/hell",
		"1993/horse",
		"1993/isaac",
		"1993/kami",
		"1993/lazy",
		"1993/mimo13",
		"1993/mister2",
		"1993/mister3",
		"1993/mohawk",
		"1993/mutation",
		"1993/ninus17",
		"1993/nl_2a",
		"1993/nl_2b",
		"1993/p68",
		"1993/p69",
		"1993/penta",
		"1993/phobos_2",
		"1993/pippo93",
		"1993/pognant",
		"1993/poirot",
		"1993/polluce",
		"1993/premana",
		"1993/puyopuyo",
		"1993/raid2",
		"1993/rapper",
		"1993/r_cyborg",
		"1993/r_daneel",
		"1993/robocop3",
		"1993/spartaco",
		"1993/target",
		"1993/tm",
		"1993/torneo",
		"1993/vannina",
		"1993/vocus",
		"1993/warrior4",
		"1993/wassilij",
		"1993/wolfgang",
		"1993/zulu",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1993))

	fmt.Print("Loading 1994... ")
	tournament1994 = []string{
		"1994/8bismark",
		"1994/anglek2",
		"1994/apache",
		"1994/bachopin",
		"1994/baubau",
		"1994/biro",
		"1994/blob",
		"1994/circlek1",
		"1994/corner3b",
		"1994/corner4",
		"1994/deluxe_4",
		"1994/deluxe_5",
		"1994/didimo",
		"1994/dima10",
		"1994/dima9",
		"1994/emanuela",
		"1994/ematico",
		"1994/fastfood",
		"1994/flash3",
		"1994/funky",
		"1994/giali1",
		"1994/hal9000",
		"1994/heavens",
		"1994/horse2",
		"1994/iching",
		"1994/jet",
		"1994/ken",
		"1994/lazyii",
		"1994/matrox",
		"1994/maverick",
		"1994/miaomiao",
		"1994/nemesi",
		"1994/ninus75",
		"1994/patcioca",
		"1994/pioppo",
		"1994/pippo94a",
		"1994/pippo94b",
		"1994/polipo",
		"1994/randwall",
		"1994/robot1",
		"1994/robot2",
		"1994/sdix3",
		"1994/sgnaus",
		"1994/shadow",
		"1994/superfly",
		"1994/the_dam",
		"1994/t-rex",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1994))

	fmt.Print("Loading 1995... ")
	tournament1995 = []string{
		"1995/andrea",
		"1995/animal",
		"1995/apache95",
		"1995/archer",
		"1995/b115e2",
		"1995/b52",
		"1995/biro",
		"1995/boss",
		"1995/camillo",
		"1995/carlo",
		"1995/circle",
		"1995/cri95",
		"1995/diablo",
		"1995/flash4",
		"1995/hal9000",
		"1995/heavens",
		"1995/horse3",
		"1995/kenii",
		"1995/losendos",
		"1995/mikezhar",
		"1995/ninus99",
		"1995/paccu",
		"1995/passion",
		"1995/peribolo",
		"1995/pippo95",
		"1995/rambo",
		"1995/rocco",
		"1995/saxy",
		"1995/sel",
		"1995/skizzo",
		"1995/star",
		"1995/stinger",
		"1995/tabori-1",
		"1995/tabori-2",
		"1995/tequila",
		"1995/tmii",
		"1995/tox",
		"1995/t-rex",
		"1995/tricky",
		"1995/twins",
		"1995/upv-9596",
		"1995/xenon",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1995))

	fmt.Print("Loading 1996... ")
	tournament1996 = []string{
		"1996/aleph",
		"1996/andrea96",
		"1996/ap_4",
		"1996/carlo96",
		"1996/diablo2",
		"1996/drago5",
		"1996/d_ray",
		"1996/fb3",
		"1996/gevbass",
		"1996/golem",
		"1996/gpo2",
		"1996/hal9000",
		"1996/heavnew",
		"1996/hider2",
		"1996/infinity",
		"1996/jaja",
		"1996/memories",
		"1996/murdoc",
		"1996/natas",
		"1996/newb52",
		"1996/pacio",
		"1996/pippo96a",
		"1996/pippo96b",
		"1996/!",
		"1996/risk",
		"1996/robot1",
		"1996/robot2",
		"1996/rudolf",
		"1996/second3",
		"1996/s-seven",
		"1996/tatank_3",
		"1996/tronco",
		"1996/uht",
		"1996/xabaras",
		"1996/yuri",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1996))

	fmt.Print("Loading 1997... ")
	tournament1997 = []string{
		"1997/1&1",
		"1997/abyss",
		"1997/ai1",
		"1997/andrea97",
		"1997/arale",
		"1997/belva",
		"1997/carlo97",
		"1997/ciccio",
		"1997/colossus",
		"1997/diablo3",
		"1997/diabolik",
		"1997/drago6",
		"1997/erica",
		"1997/fable",
		"1997/flash5",
		"1997/fya",
		"1997/gevbass2",
		"1997/golem2",
		"1997/gundam",
		"1997/hal9000",
		"1997/jedi",
		"1997/kill!",
		"1997/me-110c",
		"1997/ncmplt",
		"1997/paperone",
		"1997/pippo97",
		"1997/raid3",
		"1997/robivinf",
		"1997/rudolf_2",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1997))

	fmt.Print("Loading 1998... ")
	tournament1998 = []string{
		"1998/ai2",
		"1998/bartali",
		"1998/carla",
		"1998/coppi",
		"1998/dia",
		"1998/dicin",
		"1998/eva00",
		"1998/eva01",
		"1998/freedom",
		"1998/fscan",
		"1998/goblin",
		"1998/goldrake",
		"1998/hal9000",
		"1998/heavnew",
		"1998/maxheav",
		"1998/ninja",
		"1998/paranoid",
		"1998/pippo98",
		"1998/plump",
		"1998/quarto",
		"1998/rattolo",
		"1998/rudolf_3",
		"1998/son-goku",
		"1998/sottolin",
		"1998/stay",
		"1998/stighy98",
		"1998/themicro",
		"1998/titania",
		"1998/tornado",
		"1998/traker1",
		"1998/traker2",
		"1998/vision",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1998))

	fmt.Print("Loading 1999... ")
	tournament1999 = []string{
		"1999/11",
		"1999/aeris",
		"1999/akira",
		"1999/alezai17",
		"1999/alfa99",
		"1999/alien",
		"1999/ap_5",
		"1999/bastrd!!",
		"1999/cancer",
		"1999/carlo99",
		"1999/#cimice#",
		"1999/cortez",
		"1999/cyborg",
		"1999/dario",
		"1999/dav46",
		"1999/defender",
		"1999/elisir",
		"1999/flash6",
		"1999/hal9000",
		"1999/ilbestio",
		"1999/jedi2",
		"1999/ka_aroth",
		"1999/kakakatz",
		"1999/lukather",
		"1999/mancino",
		"1999/marko",
		"1999/mcenrobo",
		"1999/m_hingis",
		"1999/minatela",
		"1999/new",
		"1999/nexus_2",
		"1999/nl_3a",
		"1999/nl_3b",
		"1999/obiwan",
		"1999/omega99",
		"1999/panduro",
		"1999/panic",
		"1999/pippo99",
		"1999/pizarro",
		"1999/quarto",
		"1999/quingon",
		"1999/rudolf_4",
		"1999/satana",
		"1999/shock",
		"1999/songohan",
		"1999/stealth",
		"1999/storm",
		"1999/surrende",
		"1999/t1001",
		"1999/themicro",
		"1999/titania2",
		"1999/vibrsper",
		"1999/zero",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1999))

	fmt.Print("Loading 2000... ")
	tournament2000 = []string{
		"2000/bach_2k",
		"2000/defender",
		"2000/doppia_g",
		"2000/flash7",
		"2000/jedi3",
		"2000/mancino",
		"2000/marine",
		"2000/m_hingis",
		"2000/navaho",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2000))

	fmt.Print("Loading 2001... ")
	tournament2001 = []string{
		"2001/burrfoot",
		"2001/charles",
		"2001/cisc",
		"2001/cobra",
		"2001/copter",
		"2001/defender",
		"2001/fizban",
		"2001/gers",
		"2001/grezbot",
		"2001/hammer",
		"2001/heavnew",
		"2001/homer",
		"2001/klr2",
		"2001/kyashan",
		"2001/max10",
		"2001/mflash2",
		"2001/microdna",
		"2001/midi_zai",
		"2001/mnl_1a",
		"2001/mnl_1b",
		"2001/murray",
		"2001/neo0",
		"2001/nl_5b",
		"2001/pippo1a",
		"2001/pippo1b",
		"2001/raistlin",
		"2001/ridicol",
		"2001/risc",
		"2001/rudy_xp",
		"2001/sdc2",
		"2001/staticii",
		"2001/thunder",
		"2001/vampire",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2001))

	fmt.Print("Loading 2002... ")
	tournament2002 = []string{
		"2002/01",
		"2002/adsl",
		"2002/anakin",
		"2002/asterix",
		"2002/bruenor",
		"2002/colera",
		"2002/colosseum",
		"2002/copter_2",
		"2002/corner5",
		"2002/doom2099",
		"2002/dynamite",
		"2002/enigma",
		"2002/groucho",
		"2002/halman",
		"2002/harpo",
		"2002/idefix",
		"2002/kyash_2",
		"2002/marco",
		"2002/mazinga",
		"2002/medioman",
		"2002/mg_one",
		"2002/mind",
		"2002/neo_sifr",
		"2002/ollio",
		"2002/padawan",
		"2002/peste",
		"2002/pippo2a",
		"2002/pippo2b",
		"2002/regis",
		"2002/scsi",
		"2002/serse",
		"2002/ska",
		"2002/stanlio",
		"2002/staticxp",
		"2002/supernov",
		"2002/tifo",
		"2002/tigre",
		"2002/todos",
		"2002/tomahawk",
		"2002/vaiolo",
		"2002/vauban",
		"2002/yoyo",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2002))

	fmt.Print("Loading 2003... ")
	tournament2003 = []string{
		"2003/730",
		"2003/adrian",
		"2003/ares",
		"2003/barbarian",
		"2003/blitz",
		"2003/briscolo",
		"2003/bruce",
		"2003/cadderly",
		"2003/cariddi",
		"2003/cvirus2",
		"2003/cvirus",
		"2003/danica",
		"2003/dynacond",
		"2003/falco",
		"2003/foursquare",
		"2003/frame",
		"2003/herpes",
		"2003/ici",
		"2003/instict",
		"2003/irpef",
		"2003/janu",
		"2003/kyash_3c",
		"2003/kyash_3m",
		"2003/lbr1",
		"2003/lbr",
		"2003/lebbra",
		"2003/mg_two",
		"2003/minicond",
		"2003/morituro",
		"2003/nautilus",
		"2003/nemo",
		"2003/neo_sel",
		"2003/piiico",
		"2003/pippo3b",
		"2003/pippo3",
		"2003/red_wolf",
		"2003/scanner",
		"2003/scilla",
		"2003/sirio",
		"2003/sith",
		"2003/sky",
		"2003/spaceman",
		"2003/tartaruga",
		"2003/valevan",
		"2003/virus2",
		"2003/virus",
		"2003/yoda",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2003))

	fmt.Print("Loading 2004... ")
	tournament2004 = []string{
		"2004/adam",
		"2004/b_selim",
		"2004/!caos",
		"2004/ciclope",
		"2004/coyote",
		"2004/diodo",
		"2004/fisco",
		"2004/gostar",
		"2004/gotar2",
		"2004/gotar",
		"2004/irap",
		"2004/ires",
		"2004/magneto",
		"2004/mg_three",
		"2004/mystica",
		"2004/n3g4_jr",
		"2004/n3g4tivo",
		"2004/new_mini",
		"2004/pippo04a",
		"2004/pippo04b",
		"2004/poldo",
		"2004/puma",
		"2004/rat-man",
		"2004/ravatto",
		"2004/rotar",
		"2004/selim_b",
		"2004/unlimited",
		"2004/wgdi",
		"2004/zener",
		"2004/!zeus",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2004))

	fmt.Print("Loading 2007... ")
	tournament2007 = []string{
		"2007/angel",
		"2007/back",
		"2007/brontolo",
		"2007/electron",
		"2007/e",
		"2007/gongolo",
		"2007/iceman",
		"2007/mammolo",
		"2007/microbo1",
		"2007/microbo2",
		"2007/midi1",
		"2007/neutron",
		"2007/pippo07a",
		"2007/pippo07b",
		"2007/pisolo",
		"2007/pyro",
		"2007/rythm",
		"2007/tobey",
		"2007/t",
		"2007/zigozago",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2007))

	fmt.Print("Loading 2010... ")
	tournament2010 = []string{
		"2010/buffy",
		"2010/cancella",
		"2010/copia",
		"2010/enkidu",
		"2010/eurialo",
		"2010/hal9010",
		"2010/macchia",
		"2010/niso",
		"2010/party",
		"2010/pippo10a",
		"2010/reuben",
		"2010/stitch",
		"2010/sweat",
		"2010/taglia",
		"2010/toppa",
		"2010/wall-e",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2010))

	fmt.Print("Loading 2011... ")
	tournament2011 = []string{
		"2011/ataman",
		"2011/coeurl",
		"2011/gerty",
		"2011/hal9011",
		"2011/jeeg",
		"2011/minion",
		"2011/nikita",
		"2011/origano",
		"2011/ortica",
		"2011/pain",
		"2011/piperita",
		"2011/pippo11a",
		"2011/pippo11b",
		"2011/tannhause",
		"2011/unmaldestr",
		"2011/vector",
		"2011/wall-e_ii",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2011))

	fmt.Print("Loading 2012... ")
	tournament2012 = []string{
		"2012/avoider",
		"2012/beat",
		"2012/british",
		"2012/camille",
		"2012/china",
		"2012/dampyr",
		"2012/easyjet",
		"2012/flash8c",
		"2012/flash8e",
		"2012/gerty2",
		"2012/grezbot2",
		"2012/gunnyb29",
		"2012/hal9012",
		"2012/lycan",
		"2012/mister2b",
		"2012/mister3b",
		"2012/pippo12a",
		"2012/pippo12b",
		"2012/power",
		"2012/puffomic",
		"2012/puffomid",
		"2012/q",
		"2012/ryanair",
		"2012/silversurf",
		"2012/torchio",
		"2012/wall-e_iii",
		"2012/yeti",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2012))

	fmt.Print("Loading 2013... ")
	tournament2013 = []string{
		"2013/axolotl",
		"2013/destro",
		"2013/eternity",
		"2013/frisa_13",
		"2013/gerty3",
		"2013/ghostrider",
		"2013/guanaco",
		"2013/gunnyb13",
		"2013/hal9013",
		"2013/jarvis",
		"2013/lamela",
		"2013/leopon",
		"2013/ncc-1701",
		"2013/osvaldo",
		"2013/pippo13a",
		"2013/pippo13b",
		"2013/pray",
		"2013/ug2k",
		"2013/wall-e_iv",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2013))

	fmt.Print("Loading 2015... ")
	tournament2015 = []string{
		"2015/antman",
		"2015/aswhup",
		// "2015/avoider", // Belongs to 2012
		"2015/babadook",
		"2015/circles15",
		"2015/colour",
		"2015/coppi15mc1",
		"2015/coppi15md1",
		"2015/flash9",
		"2015/frank15",
		"2015/gerty4",
		"2015/hal9015",
		"2015/hulk",
		"2015/linabo15",
		"2015/lluke",
		"2015/mcfly",
		"2015/mike3",
		"2015/pippo15a",
		"2015/pippo15b",
		"2015/puppet",
		"2015/randguard",
		"2015/salippo",
		"2015/sidewalk",
		"2015/thor",
		"2015/tux",
		"2015/tyrion",
		"2015/wall-e_v",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2015))

	fmt.Print("Loading 2020... ")
	tournament2020 = []string{
		"2020/antman_20",
		"2020/b4b",
		"2020/brexit",
		"2020/coppi20mc1",
		"2020/coppi20md1",
		"2020/discotek",
		"2020/flash10",
		"2020/gerty5",
		"2020/hal9020",
		"2020/hulk_20",
		"2020/jarvis2",
		"2020/leavy2",
		"2020/loneliness",
		"2020/pippo20a",
		"2020/pippo20b",
		"2020/wall-e_vi",
		"2020/wizard2",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2020))

	fmt.Print("Loading 2025... ")
	tournament2025 = []string{
		"2025/blue", // late entry
		"2025/extrasmall",
		"2025/flash11",
		"2025/gerty6",
		"2025/hal9025",
		"2025/hulk_25",
		"2025/hydra",
		"2025/kerberos",
		"2025/meeseeks1",
		"2025/meeseeks2",
		"2025/nova",
		"2025/pippo25a",
		"2025/rabbitc",
		"2025/rotaprinc8",
		"2025/sentry2",
		"2025/sentry3",
		"2025/sgorbio",
		"2025/slant6",
		"2025/supremo",
		"2025/trouble3",
		"2025/ultron_25",
		"2025/wall-e_vii",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2025))

	fmt.Print("Loading crobs... ")
	crobs = []string{
		"crobs/adversar",
		"crobs/agressor",
		"crobs/antru",
		"crobs/assassin",
		"crobs/b4",
		"crobs/bishop",
		"crobs/bouncer",
		"crobs/boxer",
		"crobs/cassius",
		"crobs/catfish3",
		"crobs/chase",
		"crobs/chaser",
		"crobs/cooper1",
		"crobs/cooper2",
		"crobs/cornerkl",
		"crobs/counter",
		"crobs/counter2",
		"crobs/cruiser",
		"crobs/cspotrun",
		"crobs/danimal",
		"crobs/dave",
		"crobs/di",
		"crobs/dirtyh",
		"crobs/duck",
		"crobs/dumbname",
		"crobs/etf_kid",
		"crobs/flyby",
		"crobs/fred",
		"crobs/friendly",
		"crobs/grunt",
		"crobs/gsmr2",
		"crobs/h-k",
		"crobs/hac_atak",
		"crobs/hak3",
		"crobs/hitnrun",
		"crobs/hunter",
		"crobs/huntlead",
		"crobs/intrcptr",
		"crobs/jagger",
		"crobs/jason100",
		"crobs/kamikaze",
		"crobs/killer",
		"crobs/leader",
		"crobs/leavy",
		"crobs/lethal",
		"crobs/maniac",
		"crobs/marvin",
		"crobs/mini",
		"crobs/ninja",
		"crobs/nord",
		"crobs/nord2",
		"crobs/ogre",
		"crobs/ogre2",
		"crobs/ogre3",
		"crobs/perizoom",
		"crobs/pest",
		"crobs/phantom",
		"crobs/pingpong",
		"crobs/politik",
		"crobs/pzk",
		"crobs/pzkmin",
		"crobs/quack",
		"crobs/quikshot",
		"crobs/rabbit10",
		"crobs/rambo3",
		"crobs/rapest",
		"crobs/reflex",
		"crobs/robbie",
		"crobs/rook",
		"crobs/rungun",
		"crobs/samurai",
		"crobs/scan",
		"crobs/scanlock",
		"crobs/scanner",
		"crobs/secro",
		"crobs/sentry",
		"crobs/shark3",
		"crobs/shark4",
		"crobs/silly",
		"crobs/slead",
		"crobs/sniper",
		"crobs/spinner",
		"crobs/spot",
		"crobs/squirrel",
		"crobs/stalker",
		"crobs/stush-1",
		"crobs/topgun",
		"crobs/tracker",
		"crobs/trial4",
		"crobs/twedlede",
		"crobs/twedledm",
		"crobs/venom",
		"crobs/watchdog",
		"crobs/wizard",
		"crobs/xecutner",
		"crobs/xhatch",
		"crobs/yal",
	}
	fmt.Printf("%d robot(s)\n", len(crobs))

	fmt.Print("Loading micro... ")
	micro = []string{
		"micro/caccola",
		"micro/carletto",
		"micro/chobin",
		"micro/dream",
		"micro/ld",
		"micro/lucifer",
		"micro/marlene",
		"micro/md8",
		"micro/md9",
		"micro/mflash",
		"micro/minizai",
		"micro/pacoon",
		"micro/pikachu",
		"micro/pippo00a",
		"micro/pippo00",
		"micro/pirla",
		"micro/p",
		"micro/rudy",
		"micro/static",
		"micro/tanzen",
		"micro/uhm",
		"micro/zioalfa",
		"micro/zzz",
	}
	fmt.Printf("%d robot(s)\n", len(micro))
}

func setupMicro() {
	fmt.Print("Loading others... ")
	others = []string{
		"aminet/anticlock",
		"aminet/mirobot",
		"aminet/schwan",
		"aminet/tron",
		"cplusplus/selvaggio",
		"cplusplus/vikingo",
	}
	fmt.Printf("%d robot(s)\n", len(others))

	fmt.Print("Loading 1990... ")
	tournament1990 = []string{
		"1990/et_1",
		"1990/et_2",
		"1990/hunter",
		"1990/nexus_1",
		"1990/scanner",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1990))

	fmt.Print("Loading 1991... ")
	tournament1991 = []string{
		"1991/blade3",
		"1991/ccyber",
		"1991/diagonal",
		"1991/et_3",
		"1991/fdig",
		"1991/genius_j",
		"1991/gira",
		"1991/gunner",
		"1991/jazz",
		"1991/nexus_2",
		"1991/paolo101",
		"1991/paolo77",
		"1991/poor",
		"1991/robocop",
		"1991/runner",
		"1991/seeker",
		"1991/warrior2",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1991))

	fmt.Print("Loading 1992... ")
	tournament1992 = []string{
		"1992/ap_1",
		"1992/assassin",
		"1992/baeos",
		"1992/banzel",
		"1992/bronx-00",
		"1992/bry_bry",
		"1992/crazy",
		"1992/d47",
		"1992/daitan3",
		"1992/dancer",
		"1992/deluxe",
		"1992/et_4",
		"1992/et_5",
		"1992/flash",
		"1992/genesis",
		"1992/hunter",
		"1992/ice",
		"1992/johnny",
		"1992/mimo6new",
		"1992/mut",
		"1992/ninus6",
		"1992/nl_1a",
		"1992/nl_1b",
		"1992/ola",
		"1992/paolo",
		"1992/pavido",
		"1992/phobos_1",
		"1992/pippo",
		"1992/raid",
		"1992/random",
		"1992/revenge3",
		"1992/robbie",
		"1992/robocop2",
		"1992/robocop",
		"1992/superv",
		"1992/t1000",
		"1992/thunder",
		"1992/trio",
		"1992/uanino",
		"1992/warrior3",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1992))

	fmt.Print("Loading 1993... ")
	tournament1993 = []string{
		"1993/ares",
		"1993/argon",
		"1993/aspide",
		"1993/beast",
		"1993/biro",
		"1993/boom",
		"1993/casual",
		"1993/corner1d",
		"1993/corner3",
		"1993/courage",
		"1993/(c)",
		"1993/crob1",
		"1993/deluxe_2",
		"1993/didimo",
		"1993/elija",
		"1993/fermo",
		"1993/flash2",
		"1993/gunnyboy",
		"1993/hell",
		"1993/horse",
		"1993/isaac",
		"1993/kami",
		"1993/lazy",
		"1993/mimo13",
		"1993/mohawk",
		"1993/ninus17",
		"1993/nl_2a",
		"1993/nl_2b",
		"1993/phobos_2",
		"1993/pippo93",
		"1993/pognant",
		"1993/premana",
		"1993/raid2",
		"1993/rapper",
		"1993/r_cyborg",
		"1993/r_daneel",
		"1993/robocop3",
		"1993/spartaco",
		"1993/target",
		"1993/tournament",
		"1993/vannina",
		"1993/wassilij",
		"1993/wolfgang",
		"1993/zulu",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1993))

	fmt.Print("Loading 1994... ")
	tournament1994 = []string{
		"1994/anglek2",
		"1994/baubau",
		"1994/biro",
		"1994/circlek1",
		"1994/corner3b",
		"1994/didimo",
		"1994/dima10",
		"1994/dima9",
		"1994/emanuela",
		"1994/ematico",
		"1994/heavens",
		"1994/iching",
		"1994/jet",
		"1994/nemesi",
		"1994/ninus75",
		"1994/pioppo",
		"1994/pippo94b",
		"1994/robot1",
		"1994/robot2",
		"1994/superfly",
		"1994/the_dam",
		"1994/t-rex",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1994))

	fmt.Print("Loading 1995... ")
	tournament1995 = []string{
		"1995/andrea",
		"1995/b115e2",
		"1995/carlo",
		"1995/circle",
		"1995/diablo",
		"1995/flash4",
		"1995/heavens",
		"1995/mikezhar",
		"1995/ninus99",
		"1995/rocco",
		"1995/sel",
		"1995/skizzo",
		"1995/tmii",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1995))

	fmt.Print("Loading 1996... ")
	tournament1996 = []string{
		"1996/andrea96",
		"1996/carlo96",
		"1996/drago5",
		"1996/d_ray",
		"1996/gpo2",
		"1996/murdoc",
		"1996/natas",
		"1996/risk",
		"1996/tronco",
		"1996/yuri",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1996))

	fmt.Print("Loading 1997... ")
	tournament1997 = []string{
		"1997/ciccio",
		"1997/drago6",
		"1997/erica",
		"1997/fya",
		"1997/pippo97",
		"1997/raid3",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1997))

	fmt.Print("Loading 1998... ")
	tournament1998 = []string{
		"1998/carla",
		"1998/fscan",
		"1998/maxheav",
		"1998/pippo98",
		"1998/plump",
		"1998/themicro",
		"1998/traker1",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1998))

	fmt.Print("Loading 1999... ")
	tournament1999 = []string{
		"1999/ap_5",
		"1999/flash6",
		"1999/mcenrobo",
		"1999/nexus_2",
		"1999/surrende",
		"1999/themicro",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1999))

	fmt.Print("Loading 2000... ")
	// Empty for micro setup
	fmt.Printf("%d robot(s)\n", len(tournament2000))

	fmt.Print("Loading 2001... ")
	tournament2001 = []string{
		"2001/burrfoot",
		"2001/charles",
		"2001/cisc",
		"2001/cobra",
		"2001/copter",
		"2001/gers",
		"2001/grezbot",
		"2001/hammer",
		"2001/homer",
		"2001/klr2",
		"2001/kyashan",
		"2001/max10",
		"2001/mflash2",
		"2001/microdna",
		"2001/midi_zai",
		"2001/mnl_1a",
		"2001/mnl_1b",
		"2001/murray",
		"2001/neo0",
		"2001/pippo1a",
		"2001/pippo1b",
		"2001/raistlin",
		"2001/ridicol",
		"2001/risc",
		"2001/rudy_xp",
		"2001/sdc2",
		"2001/staticii",
		"2001/thunder",
		"2001/vampire",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2001))

	fmt.Print("Loading 2002... ")
	tournament2002 = []string{
		"2002/01",
		"2002/adsl",
		"2002/anakin",
		"2002/copter_2",
		"2002/corner5",
		"2002/doom2099",
		"2002/groucho",
		"2002/idefix",
		"2002/kyash_2",
		"2002/marco",
		"2002/mazinga",
		"2002/mind",
		"2002/neo_sifr",
		"2002/pippo2a",
		"2002/pippo2b",
		"2002/regis",
		"2002/scsi",
		"2002/ska",
		"2002/stanlio",
		"2002/staticxp",
		"2002/supernov",
		"2002/tigre",
		"2002/vaiolo",
		"2002/vauban",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2002))

	fmt.Print("Loading 2003... ")
	tournament2003 = []string{
		"2003/730",
		"2003/barbarian",
		"2003/blitz",
		"2003/briscolo",
		"2003/bruce",
		"2003/cvirus",
		"2003/danica",
		"2003/falco",
		"2003/foursquare",
		"2003/frame",
		"2003/herpes",
		"2003/ici",
		"2003/instict",
		"2003/janu",
		"2003/kyash_3m",
		"2003/lbr1",
		"2003/lbr",
		"2003/lebbra",
		"2003/minicond",
		"2003/morituro",
		"2003/nemo",
		"2003/neo_sel",
		"2003/piiico",
		"2003/pippo3b",
		"2003/pippo3",
		"2003/red_wolf",
		"2003/scilla",
		"2003/sirio",
		"2003/tartaruga",
		"2003/valevan",
		"2003/virus",
		"2003/yoda",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2003))

	fmt.Print("Loading 2004... ")
	tournament2004 = []string{
		"2004/adam",
		"2004/!caos",
		"2004/ciclope",
		"2004/coyote",
		"2004/diodo",
		"2004/gostar",
		"2004/gotar2",
		"2004/gotar",
		"2004/irap",
		"2004/magneto",
		"2004/n3g4_jr",
		"2004/new_mini",
		"2004/pippo04a",
		"2004/pippo04b",
		"2004/poldo",
		"2004/puma",
		"2004/rat-man",
		"2004/ravatto",
		"2004/rotar",
		"2004/selim_b",
		"2004/unlimited",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2004))

	fmt.Print("Loading 2007... ")
	tournament2007 = []string{
		"2007/angel",
		"2007/back",
		"2007/brontolo",
		"2007/electron",
		"2007/gongolo",
		"2007/microbo1",
		"2007/microbo2",
		"2007/pippo07a",
		"2007/pippo07b",
		"2007/pisolo",
		"2007/pyro",
		"2007/tobey",
		"2007/t",
		"2007/zigozago",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2007))

	fmt.Print("Loading 2010... ")
	tournament2010 = []string{
		"2010/copia",
		"2010/eurialo",
		"2010/macchia",
		"2010/niso",
		"2010/pippo10a",
		"2010/stitch",
		"2010/sweat",
		"2010/taglia",
		"2010/wall-e",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2010))

	fmt.Print("Loading 2011... ")
	tournament2011 = []string{
		"2011/ataman",
		"2011/coeurl",
		"2011/minion",
		"2011/pain",
		"2011/piperita",
		"2011/pippo11a",
		"2011/pippo11b",
		"2011/tannhause",
		"2011/unmaldestr",
		"2011/wall-e_ii",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2011))

	fmt.Print("Loading 2012... ")
	tournament2012 = []string{
		"2012/avoider",
		"2012/beat",
		"2012/china",
		"2012/easyjet",
		"2012/flash8c",
		"2012/flash8e",
		"2012/grezbot2",
		"2012/lycan",
		"2012/pippo12a",
		"2012/pippo12b",
		"2012/puffomic",
		"2012/ryanair",
		"2012/silversurf",
		"2012/wall-e_iii",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2012))

	fmt.Print("Loading 2013... ")
	tournament2013 = []string{
		"2013/axolotl",
		"2013/destro",
		"2013/osvaldo",
		"2013/pippo13a",
		"2013/pippo13b",
		"2013/pray",
		"2013/wall-e_iv",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2013))

	fmt.Print("Loading 2015... ")
	tournament2015 = []string{
		"2015/antman",
		"2015/aswhup",
		// "2015/avoider", // Belongs to 2012
		"2015/babadook",
		"2015/colour",
		"2015/coppi15mc1",
		"2015/flash9",
		"2015/linabo15",
		"2015/mike3",
		"2015/pippo15a",
		"2015/pippo15b",
		"2015/puppet",
		"2015/randguard",
		"2015/salippo",
		"2015/sidewalk",
		"2015/tux",
		"2015/tyrion",
		"2015/wall-e_v",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2015))

	fmt.Print("Loading 2020... ")
	tournament2020 = []string{
		"2020/antman_20",
		"2020/b4b",
		"2020/brexit",
		"2020/coppi20mc1",
		"2020/discotek",
		"2020/flash10",
		"2020/pippo20a",
		"2020/pippo20b",
		"2020/wall-e_vi",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2020))

	fmt.Print("Loading 2025... ")
	tournament2025 = []string{
		"2025/extrasmall",
		"2025/flash11",
		"2025/kerberos",
		"2025/pippo25a",
		"2025/sentry2",
		"2025/sentry3",
		"2025/sgorbio",
		"2025/slant6",
		"2025/trouble3",
		"2025/ultron_25",
		"2025/wall-e_vii",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2025))

	fmt.Print("Loading crobs... ")
	crobs = []string{
		"crobs/adversar",
		"crobs/agressor",
		"crobs/assassin",
		"crobs/b4",
		"crobs/bishop",
		"crobs/bouncer",
		"crobs/cassius",
		"crobs/catfish3",
		"crobs/chase",
		"crobs/chaser",
		"crobs/cornerkl",
		"crobs/counter2",
		"crobs/cruiser",
		"crobs/cspotrun",
		"crobs/danimal",
		"crobs/dave",
		"crobs/di",
		"crobs/dirtyh",
		"crobs/duck",
		"crobs/etf_kid",
		"crobs/flyby",
		"crobs/fred",
		"crobs/grunt",
		"crobs/gsmr2",
		"crobs/hac_atak",
		"crobs/hak3",
		"crobs/hitman",
		"crobs/h-k",
		"crobs/hunter",
		"crobs/huntlead",
		"crobs/intrcptr",
		"crobs/kamikaze",
		"crobs/killer",
		"crobs/leader",
		"crobs/marvin",
		"crobs/micro",
		"crobs/mini",
		"crobs/ninja",
		"crobs/nord2",
		"crobs/nord",
		"crobs/ogre2",
		"crobs/ogre",
		"crobs/pest",
		"crobs/phantom",
		"crobs/pingpong",
		"crobs/pzkmin",
		"crobs/pzk",
		"crobs/quack",
		"crobs/quikshot",
		"crobs/rabbit10",
		"crobs/rabbit",
		"crobs/rambo3",
		"crobs/rapest",
		"crobs/reflex",
		"crobs/rungun",
		"crobs/scanlock",
		"crobs/scanner",
		"crobs/scan",
		"crobs/sentry",
		"crobs/silly",
		"crobs/slead",
		"crobs/spinner",
		"crobs/spot",
		"crobs/squirrel",
		"crobs/stush-1",
		"crobs/topgun",
		"crobs/tracker",
		"crobs/twedlede",
		"crobs/twedledm",
		"crobs/venom",
		"crobs/watchdog",
		"crobs/xecutner",
		"crobs/xhatch",
		"crobs/yal",
	}
	fmt.Printf("%d robot(s)\n", len(crobs))

	fmt.Print("Loading micro... ")
	micro = []string{
		"micro/caccola",
		"micro/carletto",
		"micro/chobin",
		"micro/dream",
		"micro/ld",
		"micro/lucifer",
		"micro/marlene",
		"micro/md8",
		"micro/md9",
		"micro/mflash",
		"micro/minizai",
		"micro/pacoon",
		"micro/pikachu",
		"micro/pippo00a",
		"micro/pippo00",
		"micro/pirla",
		"micro/p",
		"micro/rudy",
		"micro/static",
		"micro/tanzen",
		"micro/uhm",
		"micro/zioalfa",
		"micro/zzz",
	}
	fmt.Printf("%d robot(s)\n", len(micro))
}

func setup() {
	fmt.Print("Loading others... ")
	others = []string{
		"aminet/anticlock",
		"aminet/beaver",
		"aminet/blindschl",
		"aminet/blindschl2",
		"aminet/mirobot",
		"aminet/opfer",
		"aminet/schwan",
		"aminet/tron",
		"cplusplus/selvaggio",
		"cplusplus/vikingo",
	}
	fmt.Printf("%d robot(s)\n", len(others))

	fmt.Print("Loading 1990... ")
	tournament1990 = []string{
		"1990/et_1",
		"1990/et_2",
		"1990/hunter",
		"1990/killer",
		"1990/nexus_1",
		"1990/rob1",
		"1990/scanner",
		"1990/york",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1990))

	fmt.Print("Loading 1991... ")
	tournament1991 = []string{
		"1991/blade3",
		"1991/casimiro",
		"1991/ccyber",
		"1991/clover",
		"1991/diagonal",
		"1991/et_3",
		"1991/f1",
		"1991/fdig",
		"1991/geltrude",
		"1991/genius_j",
		"1991/gira",
		"1991/gunner",
		"1991/jazz",
		"1991/nexus_2",
		"1991/paolo101",
		"1991/paolo77",
		"1991/poor",
		"1991/qibo",
		"1991/robocop",
		"1991/runner",
		"1991/sara_6",
		"1991/seeker",
		"1991/warrior2",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1991))

	fmt.Print("Loading 1992... ")
	tournament1992 = []string{
		"1992/666",
		"1992/ap_1",
		"1992/assassin",
		"1992/baeos",
		"1992/banzel",
		"1992/bronx-00",
		"1992/bry_bry",
		"1992/crazy",
		"1992/cube",
		"1992/cw",
		"1992/d47",
		"1992/daitan3",
		"1992/dancer",
		"1992/deluxe",
		"1992/dorsai",
		"1992/et_4",
		"1992/et_5",
		"1992/flash",
		"1992/genesis",
		"1992/hunter",
		"1992/ice",
		"1992/jack",
		"1992/jager",
		"1992/johnny",
		"1992/lead1",
		"1992/marika",
		"1992/mimo6new",
		"1992/mrcc",
		"1992/mut",
		"1992/ninus6",
		"1992/nl_1a",
		"1992/nl_1b",
		"1992/ola",
		"1992/paolo",
		"1992/pavido",
		"1992/phobos_1",
		"1992/pippo92",
		"1992/pippo",
		"1992/raid",
		"1992/random",
		"1992/revenge3",
		"1992/robbie",
		"1992/robocop2",
		"1992/robocop",
		"1992/sassy",
		"1992/spider",
		"1992/sp",
		"1992/superv",
		"1992/t1000",
		"1992/thunder",
		"1992/triangol",
		"1992/trio",
		"1992/uanino",
		"1992/warrior3",
		"1992/xdraw2",
		"1992/zorro",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1992))

	fmt.Print("Loading 1993... ")
	tournament1993 = []string{
		"1993/am_174",
		"1993/ap_2",
		"1993/ares",
		"1993/argon",
		"1993/aspide",
		"1993/beast",
		"1993/biro",
		"1993/blade8",
		"1993/boom",
		"1993/brain",
		"1993/cantor",
		"1993/castore",
		"1993/casual",
		"1993/corner1d",
		"1993/corner3",
		"1993/courage",
		"1993/(c)",
		"1993/crob1",
		"1993/deluxe_2",
		"1993/deluxe_3",
		"1993/didimo",
		"1993/duke",
		"1993/elija",
		"1993/fermo",
		"1993/flash2",
		"1993/food5",
		"1993/godel",
		"1993/gunnyboy",
		"1993/hamp1",
		"1993/hamp2",
		"1993/hell",
		"1993/horse",
		"1993/isaac",
		"1993/kami",
		"1993/lazy",
		"1993/mimo13",
		"1993/mister2",
		"1993/mister3",
		"1993/mohawk",
		"1993/mutation",
		"1993/ninus17",
		"1993/nl_2a",
		"1993/nl_2b",
		"1993/p68",
		"1993/p69",
		"1993/penta",
		"1993/phobos_2",
		"1993/pippo93",
		"1993/pognant",
		"1993/poirot",
		"1993/polluce",
		"1993/premana",
		"1993/puyopuyo",
		"1993/raid2",
		"1993/rapper",
		"1993/r_cyborg",
		"1993/r_daneel",
		"1993/robocop3",
		"1993/spartaco",
		"1993/target",
		"1993/tm",
		"1993/tournament",
		"1993/vannina",
		"1993/vocus",
		"1993/warrior4",
		"1993/wassilij",
		"1993/wolfgang",
		"1993/zulu",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1993))

	fmt.Print("Loading 1994... ")
	tournament1994 = []string{
		"1994/8bismark",
		"1994/anglek2",
		"1994/apache",
		"1994/bachopin",
		"1994/baubau",
		"1994/biro",
		"1994/blob",
		"1994/circlek1",
		"1994/corner3b",
		"1994/corner4",
		"1994/deluxe_4",
		"1994/deluxe_5",
		"1994/didimo",
		"1994/dima10",
		"1994/dima9",
		"1994/emanuela",
		"1994/ematico",
		"1994/fastfood",
		"1994/flash3",
		"1994/funky",
		"1994/giali1",
		"1994/hal9000",
		"1994/heavens",
		"1994/horse2",
		"1994/iching",
		"1994/jet",
		"1994/ken",
		"1994/lazyii",
		"1994/matrox",
		"1994/maverick",
		"1994/miaomiao",
		"1994/nemesi",
		"1994/ninus75",
		"1994/patcioca",
		"1994/pioppo",
		"1994/pippo94a",
		"1994/pippo94b",
		"1994/polipo",
		"1994/randwall",
		"1994/robot1",
		"1994/robot2",
		"1994/sdix3",
		"1994/sgnaus",
		"1994/shadow",
		"1994/superfly",
		"1994/the_dam",
		"1994/t-rex",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1994))

	fmt.Print("Loading 1995... ")
	tournament1995 = []string{
		"1995/andrea",
		"1995/animal",
		"1995/apache95",
		"1995/archer",
		"1995/b115e2",
		"1995/b52",
		"1995/biro",
		"1995/boss",
		"1995/camillo",
		"1995/carlo",
		"1995/circle",
		"1995/cri95",
		"1995/diablo",
		"1995/flash4",
		"1995/hal9000",
		"1995/heavens",
		"1995/horse3",
		"1995/kenii",
		"1995/losendos",
		"1995/mikezhar",
		"1995/ninus99",
		"1995/paccu",
		"1995/passion",
		"1995/peribolo",
		"1995/pippo95",
		"1995/rambo",
		"1995/rocco",
		"1995/saxy",
		"1995/sel",
		"1995/skizzo",
		"1995/star",
		"1995/stinger",
		"1995/tabori-1",
		"1995/tabori-2",
		"1995/tequila",
		"1995/tmii",
		"1995/tox",
		"1995/t-rex",
		"1995/tricky",
		"1995/twins",
		"1995/upv-9596",
		"1995/xenon",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1995))

	fmt.Print("Loading 1996... ")
	tournament1996 = []string{
		"1996/aleph",
		"1996/andrea96",
		"1996/ap_4",
		"1996/carlo96",
		"1996/diablo2",
		"1996/drago5",
		"1996/d_ray",
		"1996/fb3",
		"1996/gevbass",
		"1996/golem",
		"1996/gpo2",
		"1996/hal9000",
		"1996/heavnew",
		"1996/hider2",
		"1996/infinity",
		"1996/jaja",
		"1996/memories",
		"1996/murdoc",
		"1996/natas",
		"1996/newb52",
		"1996/pacio",
		"1996/pippo96a",
		"1996/pippo96b",
		"1996/!",
		"1996/risk",
		"1996/robot1",
		"1996/robot2",
		"1996/rudolf",
		"1996/second3",
		"1996/s-seven",
		"1996/tatank_3",
		"1996/tronco",
		"1996/uht",
		"1996/xabaras",
		"1996/yuri",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1996))

	fmt.Print("Loading 1997... ")
	tournament1997 = []string{
		"1997/1&1",
		"1997/abyss",
		"1997/ai1",
		"1997/andrea97",
		"1997/arale",
		"1997/belva",
		"1997/carlo97",
		"1997/ciccio",
		"1997/colossus",
		"1997/diablo3",
		"1997/diabolik",
		"1997/drago6",
		"1997/erica",
		"1997/fable",
		"1997/flash5",
		"1997/fya",
		"1997/gevbass2",
		"1997/golem2",
		"1997/gundam",
		"1997/hal9000",
		"1997/jedi",
		"1997/kill!",
		"1997/me-110c",
		"1997/ncmplt",
		"1997/paperone",
		"1997/pippo97",
		"1997/raid3",
		"1997/robivinf",
		"1997/rudolf_2",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1997))

	fmt.Print("Loading 1998... ")
	tournament1998 = []string{
		"1998/ai2",
		"1998/bartali",
		"1998/carla",
		"1998/coppi",
		"1998/dia",
		"1998/dicin",
		"1998/eva00",
		"1998/eva01",
		"1998/freedom",
		"1998/fscan",
		"1998/goblin",
		"1998/goldrake",
		"1998/hal9000",
		"1998/heavnew",
		"1998/maxheav",
		"1998/ninja",
		"1998/paranoid",
		"1998/pippo98",
		"1998/plump",
		"1998/quarto",
		"1998/rattolo",
		"1998/rudolf_3",
		"1998/son-goku",
		"1998/sottolin",
		"1998/stay",
		"1998/stighy98",
		"1998/themicro",
		"1998/titania",
		"1998/tornado",
		"1998/traker1",
		"1998/traker2",
		"1998/vision",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1998))

	fmt.Print("Loading 1999... ")
	tournament1999 = []string{
		"1999/11",
		"1999/aeris",
		"1999/akira",
		"1999/alezai17",
		"1999/alfa99",
		"1999/alien",
		"1999/ap_5",
		"1999/bastrd!!",
		"1999/cancer",
		"1999/carlo99",
		"1999/#cimice#",
		"1999/cortez",
		"1999/cyborg",
		"1999/dario",
		"1999/dav46",
		"1999/defender",
		"1999/elisir",
		"1999/flash6",
		"1999/hal9000",
		"1999/ilbestio",
		"1999/jedi2",
		"1999/ka_aroth",
		"1999/kakakatz",
		"1999/lukather",
		"1999/mancino",
		"1999/marko",
		"1999/mcenrobo",
		"1999/m_hingis",
		"1999/minatela",
		"1999/new",
		"1999/nexus_2",
		"1999/nl_3a",
		"1999/nl_3b",
		"1999/obiwan",
		"1999/omega99",
		"1999/panduro",
		"1999/panic",
		"1999/pippo99",
		"1999/pizarro",
		"1999/quarto",
		"1999/quingon",
		"1999/rudolf_4",
		"1999/satana",
		"1999/shock",
		"1999/songohan",
		"1999/stealth",
		"1999/storm",
		"1999/surrende",
		"1999/t1001",
		"1999/themicro",
		"1999/titania2",
		"1999/vibrsper",
		"1999/zero",
	}
	fmt.Printf("%d robot(s)\n", len(tournament1999))

	fmt.Print("Loading 2000... ")
	tournament2000 = []string{
		"2000/7di9",
		"2000/bach_2k",
		"2000/beholder",
		"2000/boom",
		"2000/carlo2k",
		"2000/coppi_2k",
		"2000/daryl",
		"2000/dav2000",
		"2000/def2",
		"2000/defender",
		"2000/doppia_g",
		"2000/flash7",
		"2000/fremen",
		"2000/gengis",
		"2000/jedi3",
		"2000/kongzill",
		"2000/mancino",
		"2000/marine",
		"2000/m_hingis",
		"2000/mrsatan",
		"2000/navaho",
		"2000/new2",
		"2000/newzai17",
		"2000/nl_4a",
		"2000/nl_4b",
		"2000/rudolf_5",
		"2000/sharp",
		"2000/touch",
		"2000/vegeth",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2000))

	fmt.Print("Loading 2001... ")
	tournament2001 = []string{
		"2001/4ever",
		"2001/artu",
		"2001/athlon",
		"2001/bati",
		"2001/bigkarl",
		"2001/borg",
		"2001/burrfoot",
		"2001/charles",
		"2001/cisc",
		"2001/cobra",
		"2001/copter",
		"2001/defender",
		"2001/disco",
		"2001/dnablack",
		"2001/dna",
		"2001/fizban",
		"2001/gers",
		"2001/grezbot",
		"2001/hammer",
		"2001/harris",
		"2001/heavnew",
		"2001/homer",
		"2001/jedi4",
		"2001/klr2",
		"2001/kyashan",
		"2001/max10",
		"2001/megazai",
		"2001/merlino",
		"2001/mflash2",
		"2001/microdna",
		"2001/midi_zai",
		"2001/mnl_1a",
		"2001/mnl_1b",
		"2001/murray",
		"2001/neo0",
		"2001/nl_5a",
		"2001/nl_5b",
		"2001/pentium4",
		"2001/pippo1a",
		"2001/pippo1b",
		"2001/raistlin",
		"2001/ridicol",
		"2001/risc",
		"2001/rudolf_6",
		"2001/rudy_xp",
		"2001/sdc2",
		"2001/sharp2",
		"2001/staticii",
		"2001/thunder",
		"2001/vampire",
		"2001/xeon",
		"2001/zifnab",
		"2001/zombie",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2001))

	fmt.Print("Loading 2002... ")
	tournament2002 = []string{
		"2002/01",
		"2002/adsl",
		"2002/anakin",
		"2002/asterix",
		"2002/attila",
		"2002/bruenor",
		"2002/colera",
		"2002/colosseum",
		"2002/copter_2",
		"2002/corner5",
		"2002/doom2099",
		"2002/drizzt",
		"2002/dynamite",
		"2002/enigma",
		"2002/groucho",
		"2002/halman",
		"2002/harpo",
		"2002/idefix",
		"2002/jedi5",
		"2002/kyash_2",
		"2002/marco",
		"2002/mazinga",
		"2002/medioman",
		"2002/mg_one",
		"2002/mind",
		"2002/moveon",
		"2002/neo_sifr",
		"2002/obelix",
		"2002/ollio",
		"2002/padawan",
		"2002/peste",
		"2002/pippo2a",
		"2002/pippo2b",
		"2002/regis",
		"2002/remus",
		"2002/romulus",
		"2002/rudolf_7",
		"2002/scsi",
		"2002/serse",
		"2002/ska",
		"2002/stanlio",
		"2002/staticxp",
		"2002/supernov",
		"2002/theslayer",
		"2002/tifo",
		"2002/tigre",
		"2002/todos",
		"2002/tomahawk",
		"2002/vaiolo",
		"2002/vauban",
		"2002/wulfgar",
		"2002/yerba",
		"2002/yoyo",
		"2002/zorn",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2002))

	fmt.Print("Loading 2003... ")
	tournament2003 = []string{
		"2003/730",
		"2003/adrian",
		"2003/aladino",
		"2003/alcadia",
		"2003/ares",
		"2003/barbarian",
		"2003/blitz",
		"2003/briscolo",
		"2003/bruce",
		"2003/cadderly",
		"2003/cariddi",
		"2003/crossover",
		"2003/cvirus2",
		"2003/cvirus",
		"2003/cyborg_2",
		"2003/danica",
		"2003/dave",
		"2003/druzil",
		"2003/dynacond",
		"2003/elminster",
		"2003/falco",
		"2003/foursquare",
		"2003/frame",
		"2003/harlock",
		"2003/herpes",
		"2003/ici",
		"2003/instict",
		"2003/irpef",
		"2003/janick",
		"2003/janu",
		"2003/jedi6",
		"2003/knt",
		"2003/kyash_3c",
		"2003/kyash_3m",
		"2003/lbr1",
		"2003/lbr",
		"2003/lebbra",
		"2003/maxicond",
		"2003/mg_two",
		"2003/minicond",
		"2003/morituro",
		"2003/nautilus",
		"2003/nemo",
		"2003/neo_sel",
		"2003/orione",
		"2003/piiico",
		"2003/pippo3b",
		"2003/pippo3",
		"2003/red_wolf",
		"2003/rudolf_8",
		"2003/scanner",
		"2003/scilla",
		"2003/sirio",
		"2003/sith",
		"2003/sky",
		"2003/spaceman",
		"2003/tartaruga",
		"2003/unico",
		"2003/valevan",
		"2003/virus2",
		"2003/virus3",
		"2003/virus4",
		"2003/virus",
		"2003/yoda",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2003))

	fmt.Print("Loading 2004... ")
	tournament2004 = []string{
		"2004/adam",
		"2004/!alien",
		"2004/bjt",
		"2004/b_selim",
		"2004/!caos",
		"2004/ciclope",
		"2004/confusion",
		"2004/coyote",
		"2004/diodo",
		"2004/!dna",
		"2004/fire",
		"2004/fisco",
		"2004/frankie",
		"2004/geriba",
		"2004/goofy",
		"2004/gostar",
		"2004/gotar2",
		"2004/gotar",
		"2004/irap",
		"2004/ire",
		"2004/ires",
		"2004/jedi7",
		"2004/magneto",
		"2004/mg_three",
		"2004/mosfet",
		"2004/m_selim",
		"2004/multics",
		"2004/mystica",
		"2004/n3g4_jr",
		"2004/n3g4tivo",
		"2004/new_mini",
		"2004/pippo04a",
		"2004/pippo04b",
		"2004/poldo",
		"2004/puma",
		"2004/rat-man",
		"2004/ravatto",
		"2004/revo",
		"2004/rotar",
		"2004/rudolf_9",
		"2004/selim_b",
		"2004/tempesta",
		"2004/unlimited",
		"2004/wgdi",
		"2004/zener",
		"2004/!zeus",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2004))

	fmt.Print("Loading 2007... ")
	tournament2007 = []string{
		"2007/angel",
		"2007/back",
		"2007/brontolo",
		"2007/colosso",
		"2007/electron",
		"2007/e",
		"2007/gongolo",
		"2007/iceman",
		"2007/jedi8",
		"2007/macro1",
		"2007/mammolo",
		"2007/microbo1",
		"2007/microbo2",
		"2007/midi1",
		"2007/neutron",
		"2007/nustyle",
		"2007/pippo07a",
		"2007/pippo07b",
		"2007/pisolo",
		"2007/proton",
		"2007/proud",
		"2007/pyro",
		"2007/rudolf_x",
		"2007/rythm",
		"2007/tobey",
		"2007/t",
		"2007/zigozago",
		"2007/z",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2007))

	fmt.Print("Loading 2010... ")
	tournament2010 = []string{
		"2010/buffy",
		"2010/cancella",
		"2010/change",
		"2010/copia",
		"2010/enkidu",
		"2010/eurialo",
		"2010/gantu",
		"2010/hal9010",
		"2010/incolla",
		"2010/jedi9",
		"2010/jumba",
		"2010/macchia",
		"2010/niso",
		"2010/party",
		"2010/pippo10a",
		"2010/reuben",
		"2010/stitch",
		"2010/suddenly",
		"2010/sweat",
		"2010/taglia",
		"2010/toppa",
		"2010/wall-e",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2010))

	fmt.Print("Loading 2011... ")
	tournament2011 = []string{
		"2011/armin",
		"2011/ataman",
		"2011/coeurl",
		"2011/digitale",
		"2011/gerty",
		"2011/grendizer",
		"2011/gru",
		"2011/guntank",
		"2011/hal9011",
		"2011/jedi10",
		"2011/jeeg",
		"2011/minion",
		"2011/nikita",
		"2011/origano",
		"2011/ortica",
		"2011/pain",
		"2011/piperita",
		"2011/pippo11a",
		"2011/pippo11b",
		"2011/smart",
		"2011/tannhause",
		"2011/tantalo",
		"2011/unmaldestr",
		"2011/vain",
		"2011/vector",
		"2011/wall-e_ii",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2011))

	fmt.Print("Loading 2012... ")
	tournament2012 = []string{
		"2012/avoider",
		"2012/beat",
		"2012/british",
		"2012/camille",
		"2012/china",
		"2012/cliche",
		"2012/crazy96",
		"2012/dampyr",
		"2012/draka",
		"2012/easyjet",
		"2012/flash8c",
		"2012/flash8e",
		"2012/gerty2",
		"2012/grezbot2",
		"2012/gunnyb29",
		"2012/hal9012",
		"2012/jedi11",
		"2012/life",
		"2012/lufthansa",
		"2012/lycan",
		"2012/mister2b",
		"2012/mister3b",
		"2012/pippo12a",
		"2012/pippo12b",
		"2012/power",
		"2012/puffomac",
		"2012/puffomic",
		"2012/puffomid",
		"2012/q",
		"2012/ryanair",
		"2012/silversurf",
		"2012/torchio",
		"2012/wall-e_iii",
		"2012/yeti",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2012))

	fmt.Print("Loading 2013... ")
	tournament2013 = []string{
		"2013/axolotl",
		"2013/destro",
		"2013/eternity",
		"2013/frisa_13",
		"2013/gerty3",
		"2013/ghostrider",
		"2013/guanaco",
		"2013/gunnyb13",
		"2013/hal9013",
		"2013/jarvis",
		"2013/jedi12",
		"2013/john_blaze",
		"2013/lamela",
		"2013/lancia13",
		"2013/leopon",
		"2013/ncc-1701",
		"2013/okapi",
		"2013/ortona_13",
		"2013/osvaldo",
		"2013/pippo13a",
		"2013/pippo13b",
		"2013/pjanic",
		"2013/pray",
		"2013/ride",
		"2013/ug2k",
		"2013/wall-e_iv",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2013))

	fmt.Print("Loading 2015... ")
	tournament2015 = []string{
		"2015/antman",
		"2015/aswhup",
		// "2015/avoider", // Belongs to 2012
		"2015/babadook",
		"2015/bttf",
		"2015/circles15",
		"2015/colour",
		"2015/coppi15ma1",
		"2015/coppi15ma2",
		"2015/coppi15mc1",
		"2015/coppi15md1",
		"2015/corbu15",
		"2015/dlrn",
		"2015/flash9",
		"2015/frank15",
		"2015/g13-14",
		"2015/gargantua",
		"2015/gerty4",
		"2015/hal9015",
		"2015/hulk",
		"2015/ironman_15",
		"2015/jedi13",
		"2015/linabo15",
		"2015/lluke",
		"2015/mcfly",
		"2015/mies15",
		"2015/mike3",
		"2015/misdemeano",
		"2015/music",
		"2015/one",
		"2015/pantagruel",
		"2015/pippo15a",
		"2015/pippo15b",
		"2015/puppet",
		"2015/randguard",
		"2015/salippo",
		"2015/sidewalk",
		"2015/the_old",
		"2015/thor",
		"2015/tux",
		"2015/tyrion",
		"2015/wall-e_v",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2015))

	fmt.Print("Loading 2020... ")
	tournament2020 = []string{
		"2020/antman_20",
		"2020/b4b",
		"2020/brexit",
		"2020/carillon",
		"2020/coppi20ma1",
		"2020/coppi20ma2",
		"2020/coppi20mc1",
		"2020/coppi20md1",
		"2020/discotek",
		"2020/dreamland",
		"2020/flash10",
		"2020/gerty5",
		"2020/hal9020",
		"2020/hulk_20",
		"2020/ironman_20",
		"2020/jarvis2",
		"2020/jedi14",
		"2020/leavy2",
		"2020/loneliness",
		"2020/pippo20a",
		"2020/pippo20b",
		"2020/thor_20",
		"2020/wall-e_vi",
		"2020/wizard2",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2020))

	fmt.Print("Loading 2025... ")
	tournament2025 = []string{
		"2025/blue", // late entry
		"2025/extrasmall",
		"2025/flash11",
		"2025/gerty6",
		"2025/hal9025",
		"2025/hulk_25",
		"2025/hydra",
		"2025/ironman_25",
		"2025/jedi15",
		"2025/kerberos",
		"2025/meeseeks1",
		"2025/meeseeks2",
		"2025/nova",
		"2025/pippo25a",
		"2025/rabbitc",
		"2025/rotaprinc8",
		"2025/sentry2",
		"2025/sentry3",
		"2025/sgorbio",
		"2025/slant6",
		"2025/supremo",
		"2025/thor_25",
		"2025/trouble3",
		"2025/ultron_25",
		"2025/wall-e_vii",
	}
	fmt.Printf("%d robot(s)\n", len(tournament2025))

	fmt.Print("Loading crobs... ")
	crobs = []string{
		"crobs/adversar",
		"crobs/agressor",
		"crobs/antru",
		"crobs/assassin",
		"crobs/b4",
		"crobs/bishop",
		"crobs/bouncer",
		"crobs/boxer",
		"crobs/cassius",
		"crobs/catfish3",
		"crobs/chase",
		"crobs/chaser",
		"crobs/cooper1",
		"crobs/cooper2",
		"crobs/cornerkl",
		"crobs/counter",
		"crobs/counter2",
		"crobs/cruiser",
		"crobs/cspotrun",
		"crobs/danimal",
		"crobs/dave",
		"crobs/di",
		"crobs/dirtyh",
		"crobs/duck",
		"crobs/dumbname",
		"crobs/etf_kid",
		"crobs/flyby",
		"crobs/fred",
		"crobs/friendly",
		"crobs/grunt",
		"crobs/gsmr2",
		"crobs/h-k",
		"crobs/hac_atak",
		"crobs/hak3",
		"crobs/hitnrun",
		"crobs/hunter",
		"crobs/huntlead",
		"crobs/intrcptr",
		"crobs/jagger",
		"crobs/jason100",
		"crobs/kamikaze",
		"crobs/killer",
		"crobs/leader",
		"crobs/leavy",
		"crobs/lethal",
		"crobs/maniac",
		"crobs/marvin",
		"crobs/mini",
		"crobs/ninja",
		"crobs/nord",
		"crobs/nord2",
		"crobs/ogre",
		"crobs/ogre2",
		"crobs/ogre3",
		"crobs/perizoom",
		"crobs/pest",
		"crobs/phantom",
		"crobs/pingpong",
		"crobs/politik",
		"crobs/pzk",
		"crobs/pzkmin",
		"crobs/quack",
		"crobs/quikshot",
		"crobs/rabbit10",
		"crobs/rambo3",
		"crobs/rapest",
		"crobs/reflex",
		"crobs/robbie",
		"crobs/rook",
		"crobs/rungun",
		"crobs/samurai",
		"crobs/scan",
		"crobs/scanlock",
		"crobs/scanner",
		"crobs/secro",
		"crobs/sentry",
		"crobs/shark3",
		"crobs/shark4",
		"crobs/silly",
		"crobs/slead",
		"crobs/sniper",
		"crobs/spinner",
		"crobs/spot",
		"crobs/squirrel",
		"crobs/stalker",
		"crobs/stush-1",
		"crobs/topgun",
		"crobs/tracker",
		"crobs/trial4",
		"crobs/twedlede",
		"crobs/twedledm",
		"crobs/venom",
		"crobs/watchdog",
		"crobs/wizard",
		"crobs/xecutner",
		"crobs/xhatch",
		"crobs/yal",
	}
	fmt.Printf("%d robot(s)\n", len(crobs))

	fmt.Print("Loading micro... ")
	micro = []string{
		"micro/caccola",
		"micro/carletto",
		"micro/chobin",
		"micro/dream",
		"micro/ld",
		"micro/lucifer",
		"micro/marlene",
		"micro/md8",
		"micro/md9",
		"micro/mflash",
		"micro/minizai",
		"micro/pacoon",
		"micro/pikachu",
		"micro/pippo00a",
		"micro/pippo00",
		"micro/pirla",
		"micro/p",
		"micro/rudy",
		"micro/static",
		"micro/tanzen",
		"micro/uhm",
		"micro/zioalfa",
		"micro/zzz",
	}
	fmt.Printf("%d robot(s)\n", len(micro))
}
