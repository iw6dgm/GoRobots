package count

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

const (
	robotName = `   (2)     ` + string(os.PathSeparator) + `jedi13.r: damage=% 60  `
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

	r1 := result["bouncer"]
	r2 := result["rook"]
	r3 := result["leader"]
	r4 := result["rabbit"]

	if r1.Games != 10 || r2.Games != 10 || r3.Games != 10 || r4.Games != 10 {
		t.Errorf("Invalid games count; want 10 10 10 10, got %d %d %d %d ", r1.Games, r2.Games, r3.Games, r4.Games)
	}

	if r1.Wins != 0 || r2.Wins != 0 || r3.Wins != 10 || r4.Wins != 0 {
		t.Errorf("Invalid wins count; want 3 5 1 0, got %d %d %d %d ", r1.Wins, r2.Wins, r3.Wins, r4.Wins)
	}

	if r1.Ties[0] != 0 || r2.Ties[0] != 0 || r3.Ties[0] != 0 || r4.Ties[0] != 0 {
		t.Errorf("Invalid ties2 count; want 0 1 1 0, got %d %d %d %d ", r1.Ties[0], r2.Ties[0], r3.Ties[0], r4.Ties[0])
	}

	sum := 0
	for i := 1; i < 3; i++ {
		sum += r1.Ties[i] + r2.Ties[i] + r3.Ties[i] + r4.Ties[i]
	}

	if sum > 0 {
		t.Errorf("Invalid ties 3 and 4 count; want 0, got %d", sum)
	}
}
