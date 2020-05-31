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

// error status
const (
	KErrorNone 					 		= 0
	KErrorUnableToOpenFile 		 		= 1
	KErrorMemoryBlockSizeToSmall 		= 2
	KErrorProblemParsingFile     		= 3
	KErrorSyntaxError			 		= 4
	KErrorBracketsNotPaired      		= 5
	KErrorUnexpectedClosingBracketFound = 6
	KErrorInvalidCharacterFound			= 7
)

const (
	DEFAULT_MEMORY_SIZE = 30000 // bytes
	DEFAULT_CODE_SIZE 	= 10000 // bytes
	DEFAULT_SOURCE_FILE = "A.BF" // to be changed
	MINIMUM_MMEORY_BLOCK_SIZE = 1000 // min menory size
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
	nestingLevel   int // count of the bracket nesting level from the root
)


func main () {
	displayBanner ()

	// read command line arguments, 'MEMORY' & 'FILE'
	flag.StringVar (&sourceFileName, "file", DEFAULT_SOURCE_FILE,"Name of default code file")

	
	memorySizePtr = flag.Int ("memory", DEFAULT_MEMORY_SIZE, "Working memory size in bytes")
	flag.Parse ()

	fmt.Printf ("Reading from source file '%s' ...\n", strings.ToUpper (sourceFileName))
	fmt.Printf ("Initialising memory block : %d bytes\n", *memorySizePtr)

	if *memorySizePtr <= MINIMUM_MMEORY_BLOCK_SIZE {
		showStatus (KErrorMemoryBlockSizeToSmall, "", *memorySizePtr)
	}

	// init memory block
	memory = make ([]uint8, *memorySizePtr, *memorySizePtr)
	
	fmt.Println ("Reading File ....")
	if data, err := readFileToMemory (sourceFileName); err == nil {
		codeSize = len(data)
		
		instructions = make ([]uint8, codeSize, codeSize)
		instructionPtr = 0

		fmt.Println ("Parsing file ...")
		if status := parseFile (data); status != KErrorNone {
			showStatus (status, sourceFileName,0)
		} else {
			fmt.Println ("Executing file ...")
			if status := executeFile (); status != KErrorNone {
				showStatus (KErrorSyntaxError, sourceFileName,0)
			}
		}
	} else {
		showStatus (KErrorUnableToOpenFile, sourceFileName,0)
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


func parseFile (data []uint8) int {
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
				showStatus (KErrorUnexpectedClosingBracketFound,"",0)
			}

			targetIndex++
		}
		sourceIndex++
	}

	if bracketCount != 0 {
		showStatus (KErrorBracketsNotPaired,"",0)
	}

	fileSize = targetIndex

	return KErrorNone
}

func executeFile () int {
	
	resetPtrs ()
	for instructionPtr < fileSize {
	
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
				return KErrorInvalidCharacterFound
		}
	}

	return KErrorNone
}

func resetPtrs () {
	instructionPtr = 0
	memoryPtr = 0
	nestingLevel = 0
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
	fmt.Printf("%c", memory[memoryPtr])
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
	if memory[memoryPtr] == 0 {
		// skip to byte after matching closing bracket
		var openBracketCount, closingBracketCount int
		var foundMatchingBracket bool = false

		for !foundMatchingBracket && (instructionPtr < fileSize){
			if instructions[instructionPtr] == '[' {
				openBracketCount++
			} else if instructions[instructionPtr] == ']' {
				closingBracketCount++
			} 
			instructionPtr++
			if (openBracketCount == closingBracketCount) {
				foundMatchingBracket = true;
			}
		}
	} else {
		// move to next byte
		instructionPtr++
	}
}

// ']' command
func skipBackwardIfNotZero () {
	if instructions[instructionPtr] != 0 {
		var foundMatchingBracket bool = false
		var nestingLevel int = 0
		for instructionPtr >= 0 && !foundMatchingBracket {
			if instructions[instructionPtr] == '[' {
				nestingLevel++
			} else if instructions[instructionPtr] == ']' {
				nestingLevel--
			}

			if (nestingLevel == 0) {
				foundMatchingBracket = true
			} else {
				instructionPtr--
			}

		}
	} else {
		instructionPtr++
	}

}

func validChar (ch uint8) bool {
	return strings.IndexByte ("+-.,[]<>", ch) != -1
}

func dump (){
	fmt.Printf ("%d:", instructionPtr)
	for i:=0; i <5; i++{
		fmt.Printf ("%02d ", memory[i])
	}
	fmt.Println ()
}

func showStatus (status int, extraInfoString string, extraInfoValue int) {

	var message string = "\nError: "

	if status != KErrorNone {
		switch (status) {
			case KErrorUnableToOpenFile: 
										message += fmt.Sprintf ("Unable to open file '%s'\n", strings.ToUpper(extraInfoString))
			case KErrorMemoryBlockSizeToSmall:
										message += fmt.Sprintf ("Specified memory size is too small (Minimum is %d bytes)\n", extraInfoValue) 
			case KErrorProblemParsingFile:
										message += fmt.Sprintf ("Problem parsing file '%s'\n", strings.ToUpper(extraInfoString))
			case KErrorSyntaxError:
										message += fmt.Sprintf ("Syntax error in file\n")
			case KErrorBracketsNotPaired:
										message += fmt.Sprintf ("Mismatched opening and closing brackets\n")
			case KErrorUnexpectedClosingBracketFound:
										message += fmt.Sprintf ("Closing backet found before any opening bracket\n")
			case KErrorInvalidCharacterFound:
									    message += fmt.Sprintf ("Invalid character found whilst executing code : '%s'\n", extraInfoString)
			default:
				// shouldnt get here but catch it anyway
				fmt.Println ("Unknown Error Detected (%d)\n", status)
		}

		fmt.Println (message)
		os.Exit(status)
	}	
}

