package count

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

const (
	robotName = `   (2)     /jedi13.r: damage=% 60  `
)

func TestGetName(t *testing.T) {

	if name := getName(robotName[8:19]); name != "jedi13" {
		t.Errorf("Error while parsing robot name; want jedi13, got [%s]", name)
	}
}

func TestParseLog(t *testing.T) {
	// Read entire file content, giving us little control but
	// making it very simple. No need to close the file.
	content, err := ioutil.ReadFile("test.log")
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string and print to screen
	text := string(content)

	result := ParseLogs(strings.Split(text, "\n"))

	if l := len(result); l != 4 {
		t.Errorf("Error while parsing logs; want 4, got %d, result = %v+", l, result)
	}

	r1 := result["!"]
	r2 := result["son-goku"]
	r3 := result["1_1"]
	r4 := result["null"]

	if r1.games != 10 || r2.games != 10 || r3.games != 10 || r4.games != 10 {
		t.Errorf("Invalid games count; want 10 10 10 10, got %d %d %d %d ", r1.games, r2.games, r3.games, r4.games)
	}

	if r1.wins != 3 || r2.wins != 5 || r3.wins != 1 || r4.wins != 0 {
		t.Errorf("Invalid wins count; want 3 5 1 0, got %d %d %d %d ", r1.wins, r2.wins, r3.wins, r4.wins)
	}

	if r1.ties[0] != 0 || r2.ties[0] != 1 || r3.ties[0] != 1 || r4.ties[0] != 0 {
		t.Errorf("Invalid ties2 count; want 0 1 1 0, got %d %d %d %d ", r1.ties[0], r2.ties[0], r3.ties[0], r4.ties[0])
	}

	sum := 0
	for i := 1; i < 3; i++ {
		sum += r1.ties[i] + r2.ties[i] + r3.ties[i] + r4.ties[i]
	}

	if sum > 0 {
		t.Errorf("Invalid ties 3 and 4 count; want 0, got %d", sum)
	}
}
