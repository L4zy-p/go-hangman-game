package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

// var inputReader = bufio.NewReader(os.Stdin)
var dictionary = []string{
	"Zombies",
	"Mango",
	"Programming",
	"Noodle",
	"Youtuber",
	"Computer",
	"Mouse",
}

func main() {
	rand.Seed(time.Now().UnixNano()) // set seed เพื่อจะเอามาใช้ random

	targetWord := getRandomWord()                        // สุ่มคำจาก dictionary
	guessedLetters := initializeGuessedWords(targetWord) // intial ตัวอักษรจาก คำที่ได้
	hangmanState := 0                                    // ตัว count ในการวาด hangman
	hintCount := 0                                       // ตัว count ในการใช้ hint

	// วนลูปจนกว่าเกมจะ over
	for !isGameOver(targetWord, guessedLetters, hangmanState) {

		// print การแสดง game word A _ _ _ _
		printGameState(targetWord, guessedLetters, hangmanState)
		input := readInput() // อ่าน input ที่พิมพ์เข้ามา

		// ถ้า input เข้ามาเป็น hint และ hint ยังไม่ได้ใช้
		if strings.TrimSpace(input) == "hint" && checkHintLimit(hintCount) {
			// แสดง hint มา 1 ตัวอักษร
			fmt.Printf("Show hint [%s]\n", getHint(targetWord, guessedLetters))
			fmt.Println()
			hintCount++ // เพิ่ม hintCount ไป 1 เพื่อเป็นตัวเช็คว่า hint limit ครบหรือยัง
			continue
		} else if strings.TrimSpace(input) == "hint" && !checkHintLimit(hintCount) {
			// ถ้า input เข้ามาเป็น hint และ hint ใช้ไปแล้ว
			fmt.Println("You aleady show hint!!!")
			fmt.Println()
			continue
		} else if len(input) != 1 {
			// ถ้า input ที่ส่งเข้ามาไม่ถูกต้อง
			fmt.Println("Invalid input. Please use letters only...")
			fmt.Println()
			continue
		}

		letter := rune(input[0])
		// ถ้าตัวอักษรที่ส่งมาอยู่ในคำนี้
		if isCorrectGuess(targetWord, letter) {
			// ถ้าตัวกษรนี้ใส่ไปแล้ว
			if aleadyUsedLetter(guessedLetters, letter) {
				fmt.Println("You've already used that letter")
				fmt.Println()
				fmt.Println()
			} else {
				// แสดงตัวอักษร หรือใส่ตัวอักษร แทน _
				guessedLetters[letter] = true
			}
		} else {
			hangmanState++
		}
	}

	printGameState(targetWord, guessedLetters, hangmanState) // print game state เพื่อแสดง game state สุดท้าย
	fmt.Print("Game Over...")
	// ถ้าเรา ใส่ตัวอักษรทุกตัวถูกหมด
	if isWordGuessed(targetWord, guessedLetters) {
		fmt.Print("You Win!")
	} else if isHangmanComplete(hangmanState) {
		// ถ้า วาด hangman ครบทุก state แล้ว
		fmt.Print("You lose!")
	} else {
		panic("invalid state. Game is over and there is no winner")
	}
}

func getRandomWord() string {
	// random คำ rand.Intn ทั้งหมด เท่ากับ length ของ dictionary
	targetWord := dictionary[rand.Intn(len(dictionary))]
	return targetWord
}

func printGameState(targetWord string, guessedLetters map[rune]bool, hangmanState int) {
	// แสดงผลของ Game
	println(getWordGuessingProgress(targetWord, guessedLetters))

	fmt.Println()
	fmt.Println()

	fmt.Println(getHangmanDrawing(hangmanState))
}

func initializeGuessedWords(targetWord string) map[rune]bool {
	gussedLetters := map[rune]bool{}                                           // สร้าง map {[rune]: bool}
	gussedLetters[unicode.ToLower(rune(targetWord[0]))] = true                 // ให้ตัวอักษรตัวแรกของคำแสดง สมมุติว่าถ้ามี a ด้านหน้า a ทุกในคำจะแสดง
	gussedLetters[unicode.ToLower(rune(targetWord[len(targetWord)-1]))] = true // ให้ตัวอักษรที่อยู่ตัวสุดท้ายแสดง สมมุติว่าถ้ามี a ด้านหลัง a ทุกในคำจะแสดง

	return gussedLetters
}

func getHangmanDrawing(hangmanState int) string {
	// อ่าน file เพื่อวาดตัว hangman จาก states
	data, err := ioutil.ReadFile(fmt.Sprintf("states/hangman%d", hangmanState))
	if err != nil {
		panic(err)
	}

	return string(data)
}

func getWordGuessingProgress(targetWord string, guessedLetters map[rune]bool) string {
	// แสดงคำ A _ _ _ e
	result := ""
	for _, ch := range targetWord {
		if ch == ' ' {
			// ถ้ามีการเว้นวรรคในคำ ให้ใส่ช่องว่าง
			result += " "
		} else if guessedLetters[unicode.ToLower(ch)] == true {
			// ถ้าเป็น true ให้แสดงตัวอักษรขึ้นมา
			result += fmt.Sprintf("%c", ch)
		} else {
			// ถ้าไม่ ให้แสดง _
			result += "_"
		}
		// เว้นช่องว่างระหว่างคำ
		result += " "
	}
	return result
}

func readInput() string {
	// อ่าน input ที่ส่งเข้ามา
	// fmt.Print("> ")
	// input, err := inputReader.ReadString('\n')
	// if err != nil {
	// 	panic(err)
	// }

	var textInput string
	fmt.Print("> ")
	fmt.Scan(&textInput)

	return strings.TrimSpace(textInput)
}

func isCorrectGuess(targetWord string, letter rune) bool {
	// เช็คว่า ตัวอักษรที่ส่งมามีอยู่ใน คำนั้นหรือป่าว
	return strings.ContainsRune(targetWord, letter)
}

func isWordGuessed(targetWord string, guessedLetters map[rune]bool) bool {
	// วนว่า ตัวอักษรที่ใส่มาถูกทุกตัวหรือยัง
	for _, ch := range targetWord {
		if !guessedLetters[unicode.ToLower(ch)] {
			return false
		}
	}
	return true
}

func isHangmanComplete(hangmanState int) bool {
	// เช็คว่า hangman state ครบหรือยัง
	return hangmanState >= 9
}

func isGameOver(targetWord string, guessedLetters map[rune]bool, hangmanState int) bool {
	// เช็คว่า กรอกตัวอักษรครบทุกตัวหรือยัง และ hangman state ครบหรือยัง
	return isWordGuessed(targetWord, guessedLetters) || isHangmanComplete(hangmanState)
}

func aleadyUsedLetter(guessedLetters map[rune]bool, letter rune) bool {
	// เช็คว่าตัวอักษรนี้ กรอกมาหรือยัง
	return guessedLetters[unicode.ToLower(letter)]
}

func getHint(targetWord string, guessedLetters map[rune]bool) string {
	// วนหา ตัวอักษรที่ยังไม่กรอก หรือ ไม่ถูกต้องเพื่อที่จะนำไปแสดง hint
	hint := []string{}
	for _, ch := range targetWord {
		if !guessedLetters[unicode.ToLower(ch)] {
			hint = append(hint, string(ch))
		}
	}
	return randomHint(hint)
}

func randomHint(hints []string) string {
	// random hint
	randHint := hints[rand.Intn(len(hints))]
	return randHint
}

func checkHintLimit(hintCount int) bool {
	// check hint limit ไม่เกิน 1
	return hintCount < 1
}
