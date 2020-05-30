package main

import (
"fmt"
"os"
"io/ioutil"
"strings"

"github.com/eiannone/keyboard"
)
	
const BRAINFUCK_VERSION = "v1.0"

const (
	KErrorNone 	int				= 0
	KErroSyntaxError 	int		= -1
	KMismatchedBrackets 	int	= -2
	KProblemReadingSourceFile int	= -3

)

const (
	DEFAULT_MEMORY_SIZE = 30000 // byte
	DEFAULT_CODE_SIZE = 10000 // bytes

	TEST_FILE = "a.txt"
)

var (
	memoryPtr int			// Pointer to current location in memory
	instructionPtr int		// Pointer to current instruction
	memory []uint8
	instructions []uint8

	// useful limits and sizes
	memorySize int
	codeSize int
	fileSize int
	
	loopReturns [10]int // statement indexes for loop return
)

func main () {
	displayBanner ()


	memorySize = DEFAULT_MEMORY_SIZE
	codeSize = DEFAULT_CODE_SIZE

	fmt.Printf ("Initialising memory block : %d bytes\n", memorySize)

	if memorySize <= 1000 {
		fmt.Println ("Memory block specificed is too small.")
		os.Exit(-2)
	}

	// init memory block
	memory = make ([]uint8, memorySize, memorySize)
	
	fmt.Println ("Reading File ....")
	if data, err := readFileToMemory (TEST_FILE); err == nil {
		codeSize = len(data)
		
		instructions = make ([]uint8, codeSize, codeSize)
		instructionPtr = 0

		fmt.Println ("Parsing file ...")
		if status := parseFile (data); !status {
			fmt.Println ("Problem parse data file")
		} else {
			fmt.Println ("Executing file ...")
			if status := executeFile (); !status {
				fmt.Println ("Syntax error")
			}

		}
	} else {
		fmt.Println ("Problem reading source file '%s'\n", strings.ToUpper (TEST_FILE))
	}
}

// private methods
func displayBanner () {
	fmt.Printf ("BrainFuck Interpreter %s\n\n", BRAINFUCK_VERSION)
}

func readFileToMemory (filename string) ([]uint8, error) {
	file, err := os.Open (TEST_FILE)
	defer file.Close ()
	
	if err == nil {
		data, err := ioutil.ReadAll (file)
		if err == nil {
	   		return data, nil
    	}
    }
	return nil, err	
}


func parseFile (data []uint8) bool {
	var sourceIndex int
	var targetIndex int
	var bracketCount int

	for sourceIndex < len(data) {
		if validChar (data[sourceIndex]) {
			instructions[targetIndex] = data[sourceIndex]

			// check bracket paring matches up
			if instructions[targetIndex] == '[' {
				bracketCount++
			} else {
				if instructions[targetIndex] == ']' {
					bracketCount--
				}
			}

			if bracketCount < 0 {
				return false // too many closing brackets first
			}

			targetIndex++
		}
		sourceIndex++
	}

	if bracketCount != 0 {
		fileSize = 0
		return false // brackets dont match

	}

	fileSize = targetIndex

	return true
}

func executeFile () bool {
	
	resetPtrs ()
	for instructionPtr < fileSize {
		dump()
		switch instructions[instructionPtr] {
			case '+': incrementValueAtDataPtr ()
 			case '-': decrementValueAtDataPtr ()
			case '<': decrementDataPtr ()
			case '>': incrementDataPtr ()
			case '.': outputByte ()
			case ',': inputByte ()
			case '[': skipForwardIfZero ()
			case ']': skipBackwardIfNotZero ()
			default:
				return false
		}
	}

	return true
}

func resetPtrs () {
	instructionPtr = 0
	memoryPtr = 0
}

// '>' command 
func incrementDataPtr () {
	if memoryPtr < memorySize - 1 {
		memoryPtr++
	}
	instructionPtr++
}

// '<' command
func decrementDataPtr () {
	if memoryPtr > 0 {
		memoryPtr--
	}
	instructionPtr++
}

// '+' command
func incrementValueAtDataPtr () {
	memory[memoryPtr]++
	instructionPtr++

}

// '-'
func decrementValueAtDataPtr () {
	memory[memoryPtr]--
	instructionPtr++
}

// '.'
func outputByte () {
	fmt.Printf("%d ", memory[memoryPtr])
	instructionPtr++
}

// ',' command
func inputByte () {
	char, _, err := keyboard.GetSingleKey ()
	if err == nil {
		memory[memoryPtr] = uint8(char)
	}
	instructionPtr++
}

// '[' command
func skipForwardIfZero () {

	if memory[memoryPtr] != 0 {
		instructionPtr++
	} else {
		for instructions[instructionPtr] != ']'  && (instructionPtr < len(instructions)) {
			instructionPtr++
		}
		instructionPtr++
	}
	
}

// ']' command
func skipBackwardIfNotZero () {
	for instructions[instructionPtr] != '[' && (instructionPtr >= 0) { 
		instructionPtr--
	}	
}

func validChar (ch uint8) bool {
	
	return strings.IndexByte ("+-.,[]<>", ch) != -1
}

func dump (){
fmt.Printf ("%d:", instructionPtr)
	for i:=0; i <5; i++{
		fmt.Printf ("%d ", memory[i])
	}
	fmt.Println ()
}
    
