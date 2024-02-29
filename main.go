package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

type mma struct {
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

func parseLine(s string) (string, float64) {
	i := 0
	for i = 0; i < len(s); i++ {
		if s[i] == ';' {
			break
		}
	}

	measurement, _ := strconv.ParseFloat(s[i+1:], 64)

	return s[:i], measurement
}

func main() {
	fh, err := os.Open("measurements.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()

	var stations []string
	results := make(map[string]*mma)
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		s := scanner.Text()
		stationName, measurement := parseLine(s)
		v, ok := results[stationName]
		if !ok {
			stations = append(stations, stationName)
			results[stationName] = &mma{
				Min:   measurement,
				Max:   measurement,
				Sum:   measurement,
				Count: 1,
			}
		} else {
			if measurement < v.Min {
				v.Min = measurement
			}
			if measurement > v.Max {
				v.Max = measurement
			}
			v.Sum += measurement
			v.Count++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("scanner: %v", err)
	}

	sort.Strings(stations)
	prefix := "{"
	for _, stationName := range stations {
		m := results[stationName]
		fmt.Printf("%s%s=%.1f/%.1f/%.1f", prefix, stationName,
			m.Min, m.Sum/float64(m.Count), m.Max)
		prefix = ", "
	}
	fmt.Println("}")
}
