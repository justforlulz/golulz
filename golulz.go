package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	SEED        = 72
	WEEK        = 7
	ROWS        = 81
	MAX_COMMITS = 7 // TODO: Should be equal to ~74
	START_DATE  = nil
)

var (
	flSeed       int
	flMaxRows    int
	flMaxCommits int
	flStartDate  string
)

const (
	UsgSeed            = "sets a seed value for commits pattern range 1-100"
	UsgSeedShort       = "short hand version of --seed"
	UsgMaxRows         = "sets the value of maximum rows for contribution graph"
	UsgMaxRowsShort    = "short hand version of --max-rows"
	UsgStartDate       = "sets the starting date of contribution graph. Must be in MM/DD/YYYY format"
	UsgStartDateShort  = "short hand version of --date"
	UsgMaxCommits      = "sets the value of maximum commits for a date"
	UsgMaxCommitsShort = "short hand version of --max-commits"
)

func init() {
	flag.IntVar(&flSeed, "seed", SEED, UsgSeed)
	flag.IntVar(&flSeed, "s", SEED, UsgSeedShort)

	flag.IntVar(&flMaxRows, "max-rows", ROWS, UsgMaxRows)
	flag.IntVar(&flMaxRows, "mr", ROWS, UsgMaxRowsShort)

	flag.StringVar(&flStartDate, "date", START_DATE, UsgStartDate)
	flag.StringVar(&flStartDate, "d", START_DATE, UsgStartDateShort)

	flag.IntVar(&flMaxCommits, "commits", MAX_COMMITS, UsgMaxCommits)
	flag.IntVar(&flMaxCommits, "c", MAX_COMMITS, UsgMaxCommitsShort)

	flag.Parse()

	rand.Seed(time.Now().Unix())

	if flSeed < 1 || flSeed > 100 {
		flSeed = SEED
	}

	_, err := time.Parse("01/02/2006", value)
	if flStartDate == nil || err != nil {
		fmt.Println("please provide a valid date in MM/DD/YYYY format.")
		os.Exit(2)
	}

}

func main() {

	var pattern string

	for i := 0; i < flMaxRows; i++ {
		for j := 0; j < WEEK; j++ {
			if rand.Intn(100) > flSeed {
				pattern += "."
			} else {
				pattern += "0"
			}
		}
		pattern += "\n"
	}

	weeks := strings.Split(pattern, "\n")

	date := parseTime(flStartDate)

	for i := 0; i < len(weeks)-1; i++ {
		week := weeks[i]
		for j := range week {
			cell := (string)(week[j])
			if cell == "0" {
				doCommits(date)
			}
			date = nextDay(date)
		}
	}

}
func doCommits(date time.Time) {
	date = addRandomDuration(date)
	for i := 0; i < rand.Intn(flMaxCommits); i++ {
		formatted := formatTime(date)
		os.Setenv("GIT_AUTHOR_DATE", formatted)
		os.Setenv("GIT_COMMITTER_DATE", formatted)
		cmd := exec.Command("git", "commit", "--allow-empty", "-m "+formatted)
		out, err := cmd.Output()

		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}

		fmt.Println(string(out))
		date = addRandomDuration(date)
	}

}

func parseTime(value string) time.Time {
	val, err := time.Parse("01/02/2006", value)
	if err != nil {
		panic(err)
	}
	return val
}

func formatTime(t time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

}

func nextDay(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day()+1, 0, 0, 0, 0, time.UTC)
}

func addRandomDuration(value time.Time) time.Time {

	h := rand.Intn(24)
	m := rand.Intn(60)
	s := rand.Intn(60)

	d := parseDuration(h, m, s)

	return value.Add(d)
}

func parseDuration(h int, m int, s int) time.Duration {
	d := fmt.Sprintf("%dh%dm%ds", h, m, s)
	duration, err := time.ParseDuration(d)
	if err != nil {
		panic(err)
	}

	return duration
}
