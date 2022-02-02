package irsdk

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

func readSessionData(r reader, h *header) string {
	// session data (yaml)
	dec := charmap.Windows1252.NewDecoder()
	rbuf := make([]byte, h.sessionInfoLen)
	_, err := r.ReadAt(rbuf, int64(h.sessionInfoOffset))
	if err != nil {
		log.Fatal(err)
	}
	rbuf, err = dec.Bytes(rbuf)
	if err != nil {
		log.Fatal(err)
	}
	yaml := strings.TrimRight(string(rbuf[:h.sessionInfoLen]), "\x00")
	return yaml
}

func getSessionDataPath(lines []string, path string) (string, error) {
	ps := strings.Split(strings.TrimRight(path, ":"), ":")
	lvl := 0
	arrSearch := 0
	arrayCount := 0
	arrayDepth := 0
	for _, l := range lines {
		if len(l) > lvl {
			segment := ps[lvl] + ":"
			if strings.HasPrefix(ps[lvl], "{") {
				endToken := strings.Index(ps[lvl], "}")
				segment = segment[endToken+1:]
				if arrSearch == 0 {
					var err error
					arrayCount, err = strconv.Atoi(ps[lvl][1:endToken])
					if err != nil {
						break
					}
					arrSearch = 1
				}
			}
			if segment == l[lvl+arrSearch:] {
				if arrayCount == 0 {
					lvl++
					if arrSearch == 1 {
						arrayDepth++
						arrSearch = 0
					}
				} else {
					arrayCount--
				}
			} else if strings.HasPrefix(l[lvl+arrSearch+arrayDepth:], segment) {
				if arrayCount == 0 {
					v := strings.Split(l, ": ")
					if len(v) == 2 {
						return v[1], nil
					}
				}
				arrayCount--
			} else if arrayDepth > 0 && countLeadingSpaces(l) < lvl+arrayDepth+1 {
				lvl = 0
				arrSearch = 0
				arrayCount = 0
				arrayDepth = 0
			}
		} else {
			lvl = 0
			arrSearch = 0
			arrayCount = 0
			arrayDepth = 0
		}
	}
	return "", fmt.Errorf("Path not found in Session data")
}

func countLeadingSpaces(line string) int {
	i := 0
	for _, c := range line {
		if c == ' ' || c == '-' {
			i++
		} else {
			break
		}
	}
	return i
}
