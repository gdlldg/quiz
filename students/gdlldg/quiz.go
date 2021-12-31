package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var problemsPath = flag.String("problems_path", "./problems.csv", "Set the file to read problems data, default to problems.csv in current direcotory")

type problem struct {
	prompt     string
	answer     string
	userAnswer string
}

var problems []problem

func checkErr(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func main() {
	// flag will only be recoganized after parsing
	flag.Parse()
	loadProblems()
	askProblems()
	checkProblems()
}

func loadProblems() {
	f, err := os.Open(*problemsPath)
	checkErr(err)
	csvReader := csv.NewReader(f)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		checkErr(err)
		problems = append(problems, problem{
			prompt: record[0],
			answer: record[1],
		})
	}
}

func askProblems() {
	r := bufio.NewReader(os.Stdin)
	for i, p := range problems {
		fmt.Printf("Q: %s, what's your answer?\n", p.prompt)
		asn, err := r.ReadString('\n')
		checkErr(err)
		// NOTE:
		// 1. remember slice gives you the copy, not the pointer to the original value
		// 2. values read from bufio contains newline
		problems[i].userAnswer = strings.TrimSpace(asn)
	}
}

func checkProblems() {
	correctCount := 0
	for _, p := range problems {
		if p.answer == p.userAnswer {
			correctCount++
		}
	}
	fmt.Printf("You got %d of %d problems correct.\n", correctCount, len(problems))
}
