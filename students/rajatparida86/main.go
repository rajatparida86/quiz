package main

import (
	"encoding/csv"
	"fmt"
	flag "github.com/ogier/pflag"
	"os"
)

var (
	file = flag.String("file", "problems.csv", "CSV file to parse")
	//limit = flag.Int("limit", 10, "Time limit for quiz")
)

type problem struct {
	question string
	answer   string
}

func newProblem(question string, answer string) *problem {
	return &problem{
		question,
		answer,
	}
}

func main() {
	flag.Parse()
	for _, p := range getQuiz() {
		fmt.Printf("%v=%v \n", p.question, p.answer)
	}
}

func getQuiz() []*problem {
	var quiz []*problem
	csvFile, err := os.Open(*file)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	lines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		panic(err)
	}
	for _, line := range lines {
		quiz = append(quiz, newProblem(line[0], line[1]))
	}
	return quiz
}
