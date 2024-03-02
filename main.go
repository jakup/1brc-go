package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"sync"
)

const BUF_SIZE = 1048576

type mma struct {
	Min   int16
	Max   int16
	Sum   int
	Count int
}

type row struct {
	StationName string
	Measurement int16
}

func parseLine(s []byte) (r row) {
	var i int = 0

	for i = 0; i < len(s); i++ {
		if s[i] == ';' {
			break
		}
	}
	r.StationName = string(s[:i])

	neg := false
	for i++; i < len(s); i++ {
		switch s[i] {
		case '-':
			neg = true
		case '.':
		default:
			r.Measurement = (r.Measurement * 10) + int16(s[i]-'0')
		}
	}
	if neg {
		r.Measurement = -r.Measurement
	}
	return
}

func main() {
	var wg sync.WaitGroup
	var readlock sync.Mutex
	var maplock sync.RWMutex
	results := make(map[string]*mma)

	eof := false
	pbuf := make([]byte, BUF_SIZE)
	pi := 0

	fh, _ := os.Open("measurements.txt")
	defer fh.Close()

	for workerId := 0; workerId < 16; workerId++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			buf := make([]byte, BUF_SIZE)
			for !eof {
				readlock.Lock()
				n, _ := fh.Read(buf)
				if n == 0 {
					eof = true
					readlock.Unlock()
					break
				}

				// Copy partial values from beginning and end of buf somewhere for later processing..
				var firstnl int
				for firstnl = 0; firstnl < n; firstnl++ {
					if buf[firstnl] == '\n' {
						break
					}
				}
				copy(pbuf[pi:], buf[:firstnl+1])
				pi += firstnl + 1

				var lastnl int
				for lastnl = n - 1; lastnl >= 0; lastnl-- {
					if buf[lastnl] == '\n' {
						break
					}
				}
				copy(pbuf[pi:], buf[lastnl+1:n])
				pi += n - lastnl - 1

				readlock.Unlock()

				for i, j := firstnl+1, firstnl+1; j <= lastnl; j++ {
					if buf[j] == '\n' {
						r := parseLine(buf[i:j])
						i = j + 1
						maplock.RLock()
						v, ok := results[r.StationName]
						maplock.RUnlock()
						if !ok {
							maplock.Lock()
							results[r.StationName] = &mma{
								Min:   r.Measurement,
								Max:   r.Measurement,
								Sum:   int(r.Measurement),
								Count: 1,
							}
							maplock.Unlock()
							continue
						} else {
							if r.Measurement < v.Min {
								v.Min = r.Measurement
							}
							if r.Measurement > v.Max {
								v.Max = r.Measurement
							}
							v.Sum += int(r.Measurement)
							v.Count++
						}
					}
				}
			}
		}()
	}

	wg.Wait()
	for i, j := 0, 0; j < pi; j++ {
		if pbuf[j] == '\n' {
			r := parseLine(pbuf[i:j])
			i = j + 1
			v, ok := results[r.StationName]
			if !ok {
				results[r.StationName] = &mma{
					Min:   r.Measurement,
					Max:   r.Measurement,
					Sum:   int(r.Measurement),
					Count: 1,
				}
				continue
			} else {
				if r.Measurement < v.Min {
					v.Min = r.Measurement
				}
				if r.Measurement > v.Max {
					v.Max = r.Measurement
				}
				v.Sum += int(r.Measurement)
				v.Count++
			}
		}
	}

	stations := make([]string, len(results))
	i := 0
	for k := range results {
		stations[i] = k
		i++
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
