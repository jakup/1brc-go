# The One Billion Row Challenge in Go

Just having some fun trying out [1BRC](https://github.com/gunnarmorling/1brc) using Golang.

_2024-03-01: Using file.Read was a lot faster than bufio.Scanner (~1s vs ~11s) reading the file in a single thread.  This version uses a pool of 16 worker goroutines to read and process the file._

_2024-02-28: I created a baseline implementation, it produced the correct output in 1m0.535s on my home workstation.  baseline java implementation took 2m6.798s for comparision._


## Results

| Date | Version | Runtime | Notes |
| --- | --- | --- | --- |
| 2024-03-01 | v1.1 | 0m28.678s | using file.Read with worker threads |
| 2024-02-28 | v1.0 | 1m0.535s | baseline (using int16) |
| 2024-02-29 | v0.9 | 1m18.948s | baseline (using float64) |
| 2024-02-28 | 1brc java baseline | 2m6.798s | - |
