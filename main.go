package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"errors"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// TODO: add parentheses and exponents functionality

func main() {

	var input string
	fmt.Println("--------------------")
	fmt.Println("Welcome to the Go calculator")
	fmt.Println("Type 'exit' at any time to terminate the app")
	fmt.Println("--------------------")

	var err error

	out:
	for err == nil {

		// Retrieve input
		fmt.Print("Enter a mathematical expression using operators (+,-,/,*,%) and numbers only: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			input = scanner.Text()
		} else {
			fmt.Println("Failed to scan input")
			break out
		}

		// Store original input for printing result
		originalInput := input

		// Validate input
		input, err = validateInput(input)
		if err != nil {
			fmt.Println(err.Error())
			if err.Error() == "app terminated" {
				break out
			}
			err = nil
			continue
		}

		// Determine the location of each operator in the input string
		input, multDivModMatches, addSubMatches, multiplicationMatchesMap, divisionMatchesMap, modulusMatchesMap, additionMatchesMap, err := findOperators(input)
		if err != nil {
			fmt.Println(err.Error())
			break out
		}

		// Initialize temp container for results of individual expressions
		var temp float64

		// Perform mult/div/mod operations in operator agnostic order
		for multDivModMatches != nil && len(multDivModMatches) > 0 {

			var operation string

			// Determine which operation is being performed at the current index
			if multiplicationMatchesMap[multDivModMatches[0]] {
				operation = "multiplication"
			} else if divisionMatchesMap[multDivModMatches[0]] {
				operation = "division"
			} else if modulusMatchesMap[multDivModMatches[0]] {
				operation = "modulus"
			}

			// Determine length of operands
			firstOperandLength := findOperandLength(true, 0, input, multDivModMatches)
			secondOperandLength := findOperandLength(false, 0, input, multDivModMatches)
			if firstOperandLength == 0 {
				err = errors.New("first operand not found")
				fmt.Println(err.Error())
				break out
			} else if secondOperandLength == 0 {
				err = errors.New("second operand not found")
				fmt.Println(err.Error())
				break out
			}

			// Retrieve the first and second operand
			firstOperand, err := strconv.ParseFloat(string(input[multDivModMatches[0] - firstOperandLength : multDivModMatches[0]]), 64)
			if err != nil {
				fmt.Println(err.Error())
				break out
			}
			secondOperand, err := strconv.ParseFloat(string(input[multDivModMatches[0] + 1 : multDivModMatches[0] + secondOperandLength + 1]), 64)
			if err != nil {
				fmt.Println(err.Error())
				break out
			}

			// Perform the operation and store result in temp variable
			temp, err = performOperation(operation, float64(firstOperand), float64(secondOperand))
			if err != nil {
				fmt.Println(err.Error())
				break out
			}

			// Edit the original equation, replacing operands and operator with result
			tempStr := floatToString(temp)
			input = replaceIndex(input, tempStr, multDivModMatches[0]-firstOperandLength, multDivModMatches[0]+secondOperandLength)

			// Determine the location of each instance of each operator in the string
			input, multDivModMatches, addSubMatches, multiplicationMatchesMap, divisionMatchesMap, modulusMatchesMap, additionMatchesMap, err = findOperators(input)
			if err != nil {
				fmt.Println(err.Error())
				break out
			}

		}

		// Perform add/sub operations in operator agnostic order
		for addSubMatches != nil && len(addSubMatches) > 0 {

			var operation string

			// Determine which operation is being performed
			if additionMatchesMap[addSubMatches[0]] {
				operation = "addition"
			}

			// Determine length of operands
			firstOperandLength := findOperandLength(true, 0, input, addSubMatches)
			secondOperandLength := findOperandLength(false, 0, input, addSubMatches)
			if firstOperandLength == 0 {
				err = errors.New("first operand not found")
				fmt.Println(err.Error())
				break out
			} else if secondOperandLength == 0 {
				err = errors.New("second operand not found")
				fmt.Println(err.Error())
				break out
			}

			// Retrieve the first and second operand
			firstOperand, err := strconv.ParseFloat(string(input[addSubMatches[0] - firstOperandLength : addSubMatches[0]]), 64)
			if err != nil {
				fmt.Println(err.Error())
				break out
			}
			secondOperand, err := strconv.ParseFloat(string(input[addSubMatches[0] + 1 : addSubMatches[0] + secondOperandLength + 1]), 64)
			if err != nil {
				fmt.Println(err.Error())
				break out
			}

			// Perform the operation and store result in temp variable
			temp, err = performOperation(operation, float64(firstOperand), float64(secondOperand))
			if err != nil {
				fmt.Println(err.Error())
				break out
			}

			// Edit the original equation, replacing operands and operator with result
			tempStr := floatToString(temp)
			input = replaceIndex(input, tempStr, addSubMatches[0]-firstOperandLength, addSubMatches[0]+secondOperandLength)

			// Determine the location of each instance of each operator in the string
			input, multDivModMatches, addSubMatches, multiplicationMatchesMap, divisionMatchesMap, modulusMatchesMap, additionMatchesMap, err = findOperators(input)
			if err != nil {
				fmt.Println(err.Error())
				break out
			}

		}

		// At this point the input should reflect the result
		fmt.Println("--------------------")
		fmt.Println("Result: " + originalInput + " = " + input)
		fmt.Println("--------------------")

	}

	return

}

func performOperation(operation string, firstOperand, secondOperand float64) (temp float64, err error) {

	// Perform the indicated operation and return temp result to be inserted into master result string

	switch operation {

	case "multiplication":
		temp = float64(float64(firstOperand) * float64(secondOperand))
	case "division":
		temp = float64(float64(firstOperand) / float64(secondOperand))
	case "modulus":
		temp = float64(int(firstOperand) % int(secondOperand))
	case "addition":
		temp = float64(float64(firstOperand) + float64(secondOperand))
	case "subtraction":
		temp = float64(float64(firstOperand) - float64(secondOperand))
	default:
		err = errors.New("invalid operation name: " + operation)

	}

	return temp, err

}

func validateInput(input string) (newInput string, err error) {

	// Prepare new input string
	input = removeWhitespace(input)
	newInput = input

	// Check if the user wants to exit app
	if strings.ToLower(input) == "exit" {
		return "", errors.New("app terminated")
	}

	// Validate against alpha runes and special characters
	runes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$^&_={}[]:;'\"<,>?\\|"
	for _, char := range runes {
		if strings.Index(input, string(char)) != -1 {
			return "", errors.New("invalid input: " + input)
		}
	}

	// Subtraction will be handled as addition of negative numbers
	newInput = strings.ReplaceAll(newInput, "-", "+-")

	return newInput, nil

}

func findOperators(input string) (newInput string, multDivModMatches, addSubMatches []int, multiplicationMatchesMap, divisionMatchesMap, modulusMatchesMap, additionMatchesMap map[int]bool, err error) {

	var success bool
	newInput = input

	// Initialize match slices
	var multiplicationMatches []int
	var divisionMatches []int
	var modulusMatches []int
	var additionMatches []int

	// Repeat the mapping process until all subtraction/negative ambiguity is sorted out (this should not require more than 2 iterations!)
	for !success {

		success = true

		// Reset match slices
		multiplicationMatches = []int{}
		divisionMatches = []int{}
		modulusMatches = []int{}
		additionMatches = []int{}

		// Reset match maps
		multiplicationMatchesMap = make(map[int]bool)
		divisionMatchesMap = make(map[int]bool)
		modulusMatchesMap = make(map[int]bool)
		additionMatchesMap = make(map[int]bool)

		// Determine which operators are present, in PEMDAS order
		operatorOptions := regexp.MustCompile("[+-/*%]")
		matches := operatorOptions.FindAllStringIndex(newInput, -1)

		// Find and label each operator instance
		for i:=0; i<len(matches); i++ {
			if string(newInput[matches[i][0]]) == "*" {
				multiplicationMatchesMap[matches[i][0]] = true
				multiplicationMatches = append(multiplicationMatches, matches[i][0])
			} else if string(newInput[matches[i][0]]) == "/" {
				divisionMatchesMap[matches[i][0]] = true
				divisionMatches = append(divisionMatches, matches[i][0])
			} else if string(newInput[matches[i][0]]) == "%" {
				modulusMatchesMap[matches[i][0]] = true
				modulusMatches = append(modulusMatches, matches[i][0])
			} else if string(newInput[matches[i][0]]) == "+" {
				additionMatchesMap[matches[i][0]] = true
				additionMatches = append(additionMatches, matches[i][0])
			} else if string(newInput[matches[i][0]]) == "-" {
				if matches[i][0] > 0 {
					if string(newInput[matches[i][0] - 1]) == "+" {
						// We've got a '+-'
						// Check whether a 'minus' sign is preceded by a multiplication, division, or modulus operator
						if matches[i][0] - 1 > 0 {
							if string(newInput[matches[i][0] - 2]) == "*" || string(newInput[matches[i][0] - 2]) == "/" || string(newInput[matches[i][0] - 2]) == "%" {
								// Replace the '+' that was inserted earlier with a temporary token ('$')
								newInput = replaceIndex(newInput, "$", matches[i][0]-1, matches[i][0]-1)
							}
						} else {
							// Replace the '+' that was inserted earlier with a temporary token ('$')
							newInput = replaceIndex(newInput, "$", matches[i][0]-1, matches[i][0]-1)
						}
					}
				}
			} else if string(newInput[matches[i][0]]) == "." {
				// No need to handle this
			} else {
				err = errors.New("operator not supported")
			}
		}

		// Remove all token ('$') instances
		if strings.Index(newInput, "$") != -1 {
			targetIndex := strings.Index(newInput, "$")
			newInput = replaceIndex(newInput, "", targetIndex, targetIndex)
			success = false
		}

	}

	// Concatenate and sort all related operations
	multDivMatches := append(multiplicationMatches, divisionMatches...)
	multDivModMatches = append(multDivMatches, modulusMatches...)
	sort.Slice(multDivModMatches, func(i, j int) bool {
		return multDivModMatches[i] < multDivModMatches[j]
	})

	// For now, addition and subtraction will be synonymous operations
	addSubMatches = additionMatches
	sort.Slice(addSubMatches, func(i, j int) bool {
		return addSubMatches[i] < addSubMatches[j]
	})

	return

}

func findOperandLength(isFirst bool, i int, input string, matches []int) (operandLength int) {

	var operandComplete bool
	var validOperandRunes = "0123456789.-"

	// Check each successive character until an invalid operand character is found
	for offset:=1; !operandComplete; offset++ {
		var validChar bool
		for _, char := range validOperandRunes {
			if isFirst {
				if matches[i] - offset >= 0 {
					if strings.Index(string(input[matches[i] - offset]), string(char)) != -1 {
						validChar = true
					}
				} else {
					break
				}
			} else {
				if matches[i] + offset < len(input) {
					if strings.Index(string(input[matches[i] + offset]), string(char)) != -1 {
						validChar = true
					}
				} else {
					break
				}
			}
		}

		// If an invalid character is found, then the length has been determined
		if !validChar {
			// Ensure for the second operand that the final 'valid' character isn't a trailing '-'
			/*
			if !isFirst && string(input[matches[i] + offset - 1]) == "-" {
				operandLength--
			}

			 */
			operandComplete = true
			break
		}
		operandLength++
	}

	return

}

func floatToString(num float64) string {

	return strconv.FormatFloat(num, 'f', -1, 64)

}

func replaceIndex(str, sub string, replaceStart, replaceEnd int) string {

	return str[:replaceStart] + sub + str[replaceEnd+1:]

}

func removeWhitespace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}
