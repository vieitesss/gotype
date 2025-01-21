package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	//go:embed static/es_500.txt
	fileData      []byte
	fileWords     []string
	numberOfWords = 20
)

func updateStatus(status GameStatus) tea.Cmd {
	return func() tea.Msg {
		return GameStatusUpdatedMsg{status: status}
	}
}

// TODO: make this function select an amount of words randomly from a list of
// words.
func getRandomTextToWords() tea.Msg {
	rand.Seed(time.Now().UnixNano())

	if fileWords == nil {
		res, err := loadWordsFromBytes(fileData)
		if err != nil {
			fmt.Println("Could not load words from bytes.")
		}
		fileWords = res
	}

	words := make([]string, numberOfWords)
	for i := 0; i < numberOfWords; i++ {
		random := rand.Intn(len(fileWords))
		words[i] = string(fileWords[random])
	}

	return TextToWriteMsg(words)
}

func loadWordsFromBytes(data []byte) ([]string, error) {
	var words []string
	scanner := bufio.NewScanner(bytes.NewReader(data))

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}
