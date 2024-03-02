# The One Billion Row Challenge in Go

Just having some fun trying out [1BRC](https://github.com/gunnarmorling/1brc) using Golang.

_2024-03-01: Using file.Read was a lot faster than bufio.Scanner (~1s vs ~11s) reading the file in a single thread.  This version uses a pool of 16 worker goroutines to read and process the file._
_2024-02-28: I created a baseline implementation, it produced the correct output in 1m0.535s on my home workstation.  baseline java implementation took 2m6.798s for comparision._


## Results

| Version | Runtime |
| --- | --- |
| 1brc java baseline | 2m6.798s |
| v0.9 | 1m18.948s |
| v1.0 | 1m0.535s |
| v1.1 | 0m28.678s |
