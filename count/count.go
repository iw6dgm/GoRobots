package count

import (
	"os"
	"strings"
)

const (
	// Separator OS dependent path separator
	Separator = string(os.PathSeparator)
)

// Robot structure
type Robot struct {
	Name   string
	Games  int
	Wins   int
	Ties   []int // ties2, ties3, ties4
	Points int
	Eff    float32
}

func getName(s string) string {
	v := strings.Trim(s, " ")

	if strings.HasPrefix(v, Separator) {
		return v[1 : len(v)-1]
	}

	return v
}

func getSurvivor(s string) Robot {
	return Robot{Name: getName(s[8:19]), Wins: 0, Ties: []int{0, 0, 0}}
}

func getRobot(s string, robots map[string]Robot) Robot {
	n := getName(s[8:19])
	if r, ok := robots[n]; ok {
		return r
	}
	return Robot{Name: n, Ties: []int{0, 0, 0}}
}

func updateRobot(s string, survivors map[string]Robot, robots map[string]Robot) {
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

func updateSurvivor(s string, survivors map[string]Robot) {
	surv := getSurvivor(s)
	survivors[surv.Name] = surv
}

func ParseLogs(lines []string) map[string]Robot {

	robots := make(map[string]Robot)
	survivors := make(map[string]Robot)

	for _, line := range lines {

		l := len(line)

		if l < 2 {
			continue
		}

		if strings.HasPrefix(line, "Match") {
			survivors = make(map[string]Robot)
		} else if strings.Index(line, "damage=%") != -1 {
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
		} else if strings.Index(line, "Cumulative") != -1 {
			s := len(survivors)

			switch s {
			case 0:
				continue
			case 1:
				for key, value := range survivors {
					value.Wins = 1
					survivors[key] = value
				}
			default:
				i := s - 2
				for key, value := range survivors {
					value.Ties[i] = 1
					survivors[key] = value
				}
			}
		} else if strings.Index(line, "wins=") != -1 {
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
