package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type mma struct {
	Min   int16
	Max   int16
	Sum   int
	Count int
}

func parseLine(s string) (stationName string, measurement int16) {
	var i int = 0

	for i = 0; i < len(s); i++ {
		if s[i] == ';' {
			break
		}
	}
	stationName = s[:i]

	neg := false
	for i++; i < len(s); i++ {
		switch s[i] {
		case '-':
			neg = true
		case '.':
		default:
			measurement = (measurement * 10) + int16(s[i]-'0')
		}
	}
	if neg {
		measurement = -measurement
	}
	return
}

func main() {
	var stations []string
	results := make(map[string]*mma)

	fh, _ := os.Open("measurements.txt")
	defer fh.Close()

	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		stationName, measurement := parseLine(scanner.Text())
		v, ok := results[stationName]
		if !ok {
			stations = append(stations, stationName)
			results[stationName] = &mma{
				Min:   measurement,
				Max:   measurement,
				Sum:   int(measurement),
				Count: 1,
			}
			continue
		} else {
			if measurement < v.Min {
				v.Min = measurement
			}
			if measurement > v.Max {
				v.Max = measurement
			}
			v.Sum += int(measurement)
			v.Count++
		}
	}

	sort.Strings(stations)
	prefix := "{"
	for _, stationName := range stations {
		m := results[stationName]
		fmt.Printf("%s%s=%.1f/%.1f/%.1f", prefix, stationName,
			float64(m.Min)*0.1,
			math.Round(float64(m.Sum)/float64(m.Count))*0.1,
			float64(m.Max)*0.1)
		prefix = ", "
	}
	fmt.Println("}")
}
