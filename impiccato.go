/*
Questo programma genera una parola casuale da un elenco di parole italiane e chiede di inserire delle lettere.
Se l'utente indovina una lettera (lettera giusta in posizione giusta), il programma rende visibile la parola con tutte le lettere indovinate. Altrimenti, incrementa il contatore "errore" e al raggiungimento di 10 errori il giocatore perde.
*/

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var clear map[string]func() //create a map for storing clear funcs

func main() {
	var n int
	if !checkFile() {
		createFile()
	}
	fmt.Println("Scegliere la modalità di gioco: 1 per il multiplayer, 2 per il solo")
	fmt.Scan(&n)
	switch n {
	case 1:
		{
			mod1()
		}
	case 2:
		{
			mod2()
		}
	default:
		{
			panic("Inserire una modalità valida")
		}
	}
}

func checkFile() bool {
	// Stat returns file info. It will return an error if there is no file.
	_, err := os.Stat("words.txt")
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func createFile() bool {
	newFile, err := os.Create("words.txt")
	if err != nil {
		return false
	}
	newFile.Close()
	return true
}

func writeFile(word []byte) bool {
	// Open a new file for writing only
	word = append(word, '\n')
	file, err := os.OpenFile(
		"words.txt",
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		0666,
	)
	if err != nil {
		return false
	}
	defer file.Close()
	// Write bytes to file
	_, err = file.Write(word)
	if err != nil {
		return false
	}
	return true
}

func readFile() []string {
	file, err := os.Open("words.txt")
	if err != nil {
		panic("Error while opening the file")
	}
	defer file.Close()
	words := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic("Error while reading")
	}
	return words
}

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") // Windows
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["linux"] = func() {
		cmd := exec.Command("clear") // Linux
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["darwin"] = func() {
		cmd := exec.Command("clear") // Mac
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc. darwin = mac
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func cripta(s string) (a string) {
	// encrypt the secret word
	a = ""
	for range s {
		a += "*"
	}
	return
}

func checkChar(sofar []rune, res []rune, ans rune, err int) (tent []rune, c int) {
	// Check if the char is in the secret word
	ok := false
	tent = sofar
	c = err
	for i, a := range res {
		if strings.ToLower(string(ans)) == strings.ToLower(string(a)) {
			tent[i] = a
			ok = true
		}
	}
	if !ok {
		c++
	}
	return
}

func checkWord(res string, ans string, err int) int {
	// check if the word is equal to the secret word
	if strings.ToLower(res) == strings.ToLower(ans) {
		return -1
	} else {
		return err + 2
	}
}

func game(word string) {
	var s string
	var c int
	cript := cripta(word)
	ris := []rune(word)
	arcript := []rune(cript)
	fmt.Printf("La parola da indovinare è: %s (%d)\n", cript, len(cript))
	for {
		CallClear()
		fmt.Print("La parola da indovinare è: ")
		for _, a := range arcript {
			fmt.Print(string(a))
		}
		fmt.Printf(" (%d)\n", len(arcript))
		drawing(c)
		fmt.Print("Inserire il tentativo: ")
		fmt.Scan(&s)
		a := []rune(s)
		if len(a) == 1 {
			// Carattere
			arcript, c = checkChar(arcript, ris, a[0], c)
			if string(arcript) == word {
				fmt.Printf("Congratulazioni, la parola era %s!\n", word)
				return
			}
		} else {
			// Parola
			c = checkWord(word, s, c)
			if c == -1 {
				fmt.Printf("Congratulazioni, la parola era %s!\n", word)
				return
			}
		}
		if c >= 10 {
			CallClear()
			fmt.Printf("Hai perso! La parola era %s!\n", word)
			drawing(c)
			return
		}
	}
}

func firstUpper(word []byte) []byte {
	// This function set to upper case the first letter of the given word
	word = bytes.ToLower(word)
	word[0] -= 32
	return word
}

func mod1() {
	// Multiplayer mode
	var word string
	CallClear()
	fmt.Print("Inserire la parola da indovinare: ")
	fmt.Scan(&word)
	byword := []byte(word)
	CallClear()
	game(word)
	words := readFile()
	flag := false
	for i := 0; i < len(words); i++ {
		if strings.ToLower(words[i]) == strings.ToLower(word) {
			flag = true
		}
	}
	if !flag {
		writeFile(firstUpper(byword))
	}
	return
}

func mod2() {
	// solo mode
	var words []string
	rand.Seed(time.Now().UnixNano())
	if !checkFile() {
		fmt.Println("Non è stato possibile trovare il file con le parole.")
	} else {
		words = readFile()
	}
	i := rand.Intn(len(words)) // lunghezza array
	word := words[i]
	game(word)
	return
}

func drawing(n int) {
	// This function draws the hangman based on the number of errors made by the gamer
	switch n {
	case 0:
		{
		}
	case 1:
		{
			k := 1
			for i := 0; i < 3; i++ {
				for j := 2; j > i; j-- {
					fmt.Print(" ")
				}
				fmt.Print("*")
				for j := 0; j < k; j++ {
					fmt.Print(" ")
				}
				fmt.Print("*")
				for j := 2; j > i; j-- {
					fmt.Print(" ")
				}
				fmt.Println()
				k += 2
			}
		}
	case 2:
		{
			const n = 20 // altezza dell'asta
			for i := 0; i < n; i++ {
				for j := 0; j < 3; j++ {
					fmt.Print(" ")
				}
				fmt.Print("*")
				fmt.Println()
			}
			drawing(1)
		}
	case 3:
		{
			for j := 0; j < 3; j++ {
				fmt.Print(" ")
			}
			for j := 0; j < 15; j++ {
				fmt.Print("*")
			}
			fmt.Println()
			drawing(2)
		}
	case 4:
		{
			for j := 0; j < 3; j++ {
				fmt.Print(" ")
			}
			for j := 0; j < 15; j++ {
				fmt.Print("*")
			}
			fmt.Println()
			for i := 0; i < 20; i++ {
				for j := 0; j < 3; j++ {
					fmt.Print(" ")
				}
				fmt.Print("*")
				if i < 2 {
					for j := 0; j < 13; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*")
				}
				fmt.Println()
			}
			drawing(1)
		}
	case 5:
		{
			for j := 0; j < 3; j++ {
				fmt.Print(" ")
			}
			for j := 0; j < 15; j++ {
				fmt.Print("*")
			}
			fmt.Println()
			for i := 0; i < 20; i++ {
				for j := 0; j < 3; j++ {
					fmt.Print(" ")
				}
				fmt.Print("*")
				if i < 2 {
					for j := 0; j < 13; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*")
				}
				if i == 2 || i == 4 {
					for j := 0; j < 12; j++ {
						fmt.Print(" ")
					}
					fmt.Print("* *")
				}
				if i == 3 {
					for j := 0; j < 11; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*")
					for j := 0; j < 3; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*")
				}
				fmt.Println()
			}
			drawing(1)
		}
	case 6:
		{
			for j := 0; j < 3; j++ {
				fmt.Print(" ")
			}
			for j := 0; j < 15; j++ {
				fmt.Print("*")
			}
			fmt.Println()
			for i := 0; i < 20; i++ {
				for j := 0; j < 3; j++ {
					fmt.Print(" ")
				}
				fmt.Print("*")
				if (i < 2) || (i > 4 && i < 12) {
					for j := 0; j < 13; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*")
				}
				if i == 2 || i == 4 {
					for j := 0; j < 12; j++ {
						fmt.Print(" ")
					}
					fmt.Print("* *")
				}
				if i == 3 {
					for j := 0; j < 11; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*")
					for j := 0; j < 3; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*")
				}
				fmt.Println()
			}
			drawing(1)
		}
	case 7:
		{
			for j := 0; j < 3; j++ {
				fmt.Print(" ")
			}
			for j := 0; j < 15; j++ {
				fmt.Print("*")
			}
			fmt.Println()
			for i := 0; i < 20; i++ {
				for j := 0; j < 3; j++ {
					fmt.Print(" ")
				}
				fmt.Print("*")
				if (i < 2) || (i == 5) || (i > 8 && i < 12) {
					for j := 0; j < 13; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*")
				}
				if i == 2 || i == 4 {
					for j := 0; j < 12; j++ {
						fmt.Print(" ")
					}
					fmt.Print("* *")
				}
				if i == 3 {
					for j := 0; j < 11; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*")
					for j := 0; j < 3; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*")
				}
				if i == 6 {
					for j := 0; j < 12; j++ {
						fmt.Print(" ")
					}
					for j := 0; j < 2; j++ {
						fmt.Print("*")
					}
				}
				if i == 7 {
					for j := 0; j < 11; j++ {
						fmt.Print(" ")
					}
					fmt.Print("* *")
				}
				if i == 8 {
					for j := 0; j < 10; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*  *")
				}
				fmt.Println()
			}
			drawing(1)
		}
	case 8:
		{
			for j := 0; j < 3; j++ {
				fmt.Print(" ")
			}
			for j := 0; j < 15; j++ {
				fmt.Print("*")
			}
			fmt.Println()
			for i := 0; i < 20; i++ {
				for j := 0; j < 3; j++ {
					fmt.Print(" ")
				}
				fmt.Print("*")
				if (i < 2) || (i == 5) || (i > 8 && i < 12) {
					for j := 0; j < 13; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*")
				}
				if i == 2 || i == 4 {
					for j := 0; j < 12; j++ {
						fmt.Print(" ")
					}
					fmt.Print("* *")
				}
				if i == 3 {
					for j := 0; j < 11; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*")
					for j := 0; j < 3; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*")
				}
				if i == 6 {
					for j := 0; j < 12; j++ {
						fmt.Print(" ")
					}
					for j := 0; j < 3; j++ {
						fmt.Print("*")
					}
				}
				if i == 7 {
					for j := 0; j < 11; j++ {
						fmt.Print(" ")
					}
					fmt.Print("* * *")
				}
				if i == 8 {
					for j := 0; j < 10; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*  *  *")
				}

				fmt.Println()
			}
			drawing(1)
		}
	case 9:
		{
			for j := 0; j < 3; j++ {
				fmt.Print(" ")
			}
			for j := 0; j < 15; j++ {
				fmt.Print("*")
			}
			fmt.Println()
			for i := 0; i < 20; i++ {
				for j := 0; j < 3; j++ {
					fmt.Print(" ")
				}
				fmt.Print("*")
				if (i < 2) || (i == 5) || (i > 8 && i < 12) {
					for j := 0; j < 13; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*")
				}
				if i == 2 || i == 4 {
					for j := 0; j < 12; j++ {
						fmt.Print(" ")
					}
					fmt.Print("* *")
				}
				if i == 3 {
					for j := 0; j < 11; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*")
					for j := 0; j < 3; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*")
				}
				if i == 6 {
					for j := 0; j < 12; j++ {
						fmt.Print(" ")
					}
					for j := 0; j < 3; j++ {
						fmt.Print("*")
					}
				}
				if i == 7 {
					for j := 0; j < 11; j++ {
						fmt.Print(" ")
					}
					fmt.Print("* * *")
				}
				if i == 8 {
					for j := 0; j < 10; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*  *  *")
				}
				if i == 12 {
					for j := 0; j < 12; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*")
				}
				if i == 13 {
					for j := 0; j < 11; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*")
				}
				if i == 14 {
					for j := 0; j < 10; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*")
				}
				fmt.Println()
			}
			drawing(1)
		}
	case 10:
		{
			for j := 0; j < 3; j++ {
				fmt.Print(" ")
			}
			for j := 0; j < 15; j++ {
				fmt.Print("*")
			}
			fmt.Println()
			for i := 0; i < 20; i++ {
				for j := 0; j < 3; j++ {
					fmt.Print(" ")
				}
				fmt.Print("*")
				if (i < 2) || (i == 5) || (i > 8 && i < 12) {
					for j := 0; j < 13; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*")
				}
				if i == 2 || i == 4 {
					for j := 0; j < 12; j++ {
						fmt.Print(" ")
					}
					fmt.Print("* *")
				}
				if i == 3 {
					for j := 0; j < 11; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*")
					for j := 0; j < 3; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*")
				}
				if i == 6 {
					for j := 0; j < 12; j++ {
						fmt.Print(" ")
					}
					for j := 0; j < 3; j++ {
						fmt.Print("*")
					}
				}
				if i == 7 {
					for j := 0; j < 11; j++ {
						fmt.Print(" ")
					}
					fmt.Print("* * *")
				}
				if i == 8 {
					for j := 0; j < 10; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*  *  *")
				}
				if i == 12 {
					for j := 0; j < 12; j++ {
						fmt.Print(" ")
					}
					fmt.Print("* *")
				}
				if i == 13 {
					for j := 0; j < 11; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*   *")
				}
				if i == 14 {
					for j := 0; j < 10; j++ {
						fmt.Print(" ")
					}
					fmt.Print("*     *")
				}

				fmt.Println()
			}
			drawing(1)
		}
	default:
		{
			panic("Error")
		}
	}
}
