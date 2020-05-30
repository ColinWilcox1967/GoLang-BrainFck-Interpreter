package main

import (
"fmt"
"os"
"flag"
"io/ioutil"
"strings"

"github.com/eiannone/keyboard"
)
	
const BRAINFUCK_VERSION = "v1.0"

const (

	KErrorNone = 0
)

const (
	DEFAULT_MEMORY_SIZE = 30000 // bytes
	DEFAULT_CODE_SIZE 	= 10000 // bytes
	DEFAULT_SOURCE_FILE = "A.BF" // to be changed
)

var (
	sourceFileName string = DEFAULT_SOURCE_FILE
	memorySizePtr	     *int    
)

var (
	memoryPtr int			// Pointer to current location in memory
	instructionPtr int		// Pointer to current instruction
	memory []uint8
	instructions []uint8

	// useful limits and sizes
	codeSize int 		= DEFAULT_CODE_SIZE
	sourceFileSize int
	fileSize 	   int
)


func main () {
	displayBanner ()

	// read command line arguments, 'MEMORY' & 'FILE'
	flag.StringVar (&sourceFileName, "file", DEFAULT_SOURCE_FILE,"Name of default code file")

	
	memorySizePtr = flag.Int ("memory", DEFAULT_MEMORY_SIZE, "Working memory size in bytes")
	flag.Parse ()

	fmt.Printf ("Reading from source file '%s' ...\n", strings.ToUpper (sourceFileName))
	fmt.Printf ("Initialising memory block : %d bytes\n", *memorySizePtr)

	if *memorySizePtr <= 1000 {
		fmt.Println ("Memory block specified is too small.")
		os.Exit(-2)
	}

	// init memory block
	memory = make ([]uint8, *memorySizePtr, *memorySizePtr)
	
	fmt.Println ("Reading File ....")
	if data, err := readFileToMemory (sourceFileName); err == nil {
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
		fmt.Println ("Problem reading source file '%s'\n", strings.ToUpper (sourceFileName))
	}
}

// private methods
func displayBanner () {
	fmt.Printf ("BrainFuck Interpreter %s\n\n", BRAINFUCK_VERSION)
}

func readFileToMemory (filename string) ([]uint8, error) {
	file, err := os.Open (sourceFileName)
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
	if memoryPtr < *memorySizePtr - 1 {
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

func showStatus (status int) {
	
	if (status != KErrorNone) {
		os.Exit (-status)
	}
}
    
