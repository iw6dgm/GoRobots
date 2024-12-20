package count

import (
	"path/filepath"
	"strings"
)

// Robot structure
type Robot struct {
	Name   string
	Games  int
	Wins   int
	Ties   [3]int // ties2, ties3, ties4
	Points int
	Eff    float32
}

// GetName returns robot name without path and extension
func GetName(s string) string {
	return strings.Split(filepath.Base(strings.Trim(s, " ")), ".")[0]
}

func getSurvivor(s string) *Robot {
	return &Robot{Name: GetName(s[8:19]), Wins: 0, Ties: [3]int{0, 0, 0}}
}

func getRobot(s string, robots map[string]*Robot) *Robot {
	n := GetName(s[8:19])
	if r, ok := robots[n]; ok {
		return r
	}
	return &Robot{Name: n, Ties: [3]int{0, 0, 0}}
}

func updateRobot(s string, survivors map[string]*Robot, robots map[string]*Robot) {
	r := getRobot(s, robots)
	n := r.Name
	if s, ok := survivors[n]; ok {
		r.Wins += s.Wins
		for i, v := range s.Ties {
			r.Ties[i] += v
		}
	}
	r.Games++
	robots[n] = r
}

func updateSurvivor(s string, survivors map[string]*Robot) {
	surv := getSurvivor(s)
	survivors[surv.Name] = surv
}

// ParseLogs returns a log parsed into a map of robots
func ParseLogs(lines [][]byte) map[string]*Robot {

	robots := make(map[string]*Robot)
	survivors := make(map[string]*Robot)

	for _, bytes := range lines {

		l := len(bytes)

		if l < 2 {
			continue
		}

		line := string(bytes)

		if strings.HasPrefix(line, "Match") {
			survivors = make(map[string]*Robot)
		} else if strings.Contains(line, "damage=%") {
			if l < 50 {
				updateSurvivor(line, survivors)
			} else {

				if split := strings.Split(line, "\t"); len(split) > 1 {
					for _, s := range split {
						if len(s) > 1 {
							updateSurvivor(s, survivors)
						}
					}
				}

			}
		} else if strings.Contains(line, "Cumulative") {
			s := len(survivors)

			switch s {
			case 0:
				continue
			case 1:
				for _, value := range survivors {
					value.Wins = 1
				}
			default:
				i := s - 2
				for _, value := range survivors {
					value.Ties[i] = 1
				}
			}
		} else if strings.Contains(line, "wins=") {
			if l < 50 {
				updateRobot(line, survivors, robots)
			} else {
				if split := strings.Split(line, "\t"); len(split) > 1 {
					for _, s := range split {
						if len(s) > 1 {
							updateRobot(s, survivors, robots)
						}
					}
				}
			}
		}
	}

	return robots
}
