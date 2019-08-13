package main

import (
	"fmt"
	"regexp"
	"errors"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// TODO: add parentheses, exponents functionality, allow for more than 1-digit numbers

func main() {

	var input string
	fmt.Println("---Welcome to the Go calculator---")
	fmt.Println("---Type 'exit' at any time to terminate the app---")

	var err error

	out:
	for err == nil {

		// Retrieve input
		fmt.Print("Enter a mathematical expression using operators (+,-,/,*,%) and numbers only: ")
		fmt.Scanln(&input)

		// Validate input
		err = validateInput(input)
		if err != nil {
			fmt.Println(err.Error())
			if err.Error() == "app terminated" {
				break out
			}
			err = nil
			continue
		}

		// Determine which operators are present, in PEMDAS order
		operatorOptions := regexp.MustCompile("[+-/*%]")
		matches := operatorOptions.FindAllStringIndex(input, -1)

		// Initialize match slices
		var multiplicationMatches []int
		var divisionMatches []int
		var modulusMatches []int
		var additionMatches []int
		var subtractionMatches []int

		// Initialize match maps
		multiplicationMatchesMap := map[int]bool{}
		divisionMatchesMap := map[int]bool{}
		modulusMatchesMap := map[int]bool{}
		additionMatchesMap := map[int]bool{}
		subtractionMatchesMap := map[int]bool{}

		// Determine the location of each instance of each operator in the string
		for i:=0; i<len(matches); i++ {
			if string(input[matches[i][0]]) == "*" {
				multiplicationMatchesMap[matches[i][0]] = true
				multiplicationMatches = append(multiplicationMatches, matches[i][0])
			} else if string(input[matches[i][0]]) == "/" {
				divisionMatchesMap[matches[i][0]] = true
				divisionMatches = append(divisionMatches, matches[i][0])
			} else if string(input[matches[i][0]]) == "%" {
				modulusMatchesMap[matches[i][0]] = true
				modulusMatches = append(modulusMatches, matches[i][0])
			} else if string(input[matches[i][0]]) == "+" {
				additionMatchesMap[matches[i][0]] = true
				additionMatches = append(additionMatches, matches[i][0])
			} else if string(input[matches[i][0]]) == "-" {
				subtractionMatchesMap[matches[i][0]] = true
				subtractionMatches = append(subtractionMatches, matches[i][0])
			} else if string(input[matches[i][0]]) == "." {
				// No need to handle this
			} else {
				err = errors.New("operator not supported")
				fmt.Println(err.Error())
				break out
			}
		}

		// Initialize temp container for results of individual expressions
		var temp float64



		// Perform multiplication, division, and modulus first

		// Order occurrences of '*', '/', and '%' in order of index
		multDivMatches := append(multiplicationMatches, divisionMatches...)
		multDivModMatches := append(multDivMatches, modulusMatches...)
		sort.Slice(multDivModMatches, func(i, j int) bool {
			return multDivModMatches[i] < multDivModMatches[j]
		})
		// Perform mult/div/mod operations in operator agnostic order
		for i:=0; i<len(multDivModMatches); i++ {

			var operation string

			// Determine which operation is being performed at the current index
			if multiplicationMatchesMap[multDivModMatches[i]] {
				operation = "multiplication"
			} else if divisionMatchesMap[multDivModMatches[i]] {
				operation = "division"
			} else if modulusMatchesMap[multDivModMatches[i]] {
				operation = "modulus"
			}

			// Determine length of operands
			firstOperandLength := findOperandLength(true, i, input, multDivModMatches)
			secondOperandLength := findOperandLength(false, i, input, multDivModMatches)
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
			firstOperand, err := strconv.ParseFloat(string(input[multDivModMatches[i] - firstOperandLength : multDivModMatches[i]]), 64)
			if err != nil {
				fmt.Println(err.Error())
				break out
			}
			secondOperand, err := strconv.ParseFloat(string(input[multDivModMatches[i] + 1 : multDivModMatches[i] + secondOperandLength + 1]), 64)
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
			tempStr := FloatToString(temp)
			oldInputLen := len(input)
			input = replaceIndex(input, tempStr, multDivModMatches[i]-firstOperandLength, multDivModMatches[i]+secondOperandLength)
			newInputLen := len(input)
			fmt.Println("Input update: " + input)

			// Update indexes of other operators based on input modification
			if i != len(multDivModMatches) - 1 {
				for a:=i+1; a<len(multDivModMatches); a++ {
					// For each item after i in the slice, modify location index
					multDivModMatches[a] += newInputLen - oldInputLen
				}
			}

			// Update operation maps
			for k, _ := range multiplicationMatchesMap {
				multiplicationMatchesMap[k + newInputLen - oldInputLen] = true
				delete(multiplicationMatchesMap, k)
			}
			for k, _ := range divisionMatchesMap {
				divisionMatchesMap[k + newInputLen - oldInputLen] = true
				delete(divisionMatchesMap, k)
			}
			for k, _ := range modulusMatchesMap {
				modulusMatchesMap[k + newInputLen - oldInputLen] = true
				delete(modulusMatchesMap, k)
			}

		}



		// Perform addition and subtraction last

		addSubMatches := append(additionMatches, subtractionMatches...)
		sort.Slice(multDivModMatches, func(i, j int) bool {
			return multDivModMatches[i] < multDivModMatches[j]
		})
		// Perform add/sub operations in operator agnostic order
		for i:=0; i<len(addSubMatches); i++ {

			var operation string

			// Determine whether addition or subtraction
			for ai:=0; ai<len(additionMatches); ai++ {
				if additionMatches[ai] == addSubMatches[i] {
					operation = "addition"
				}
			}
			// Let's iterate as little as possible by checking for operation first
			if operation == "" {
				for si:=0; si<len(subtractionMatches); si++ {
					if subtractionMatches[si] == addSubMatches[i] {
						operation = "subtraction"
					}
				}
			}

			firstOperand, err := strconv.Atoi(string(input[addSubMatches[i] - 1]))
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			secondOperand, err := strconv.Atoi(string(input[addSubMatches[i] + 1]))
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			temp, err = performOperation(operation, float64(firstOperand), float64(secondOperand))
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			// Edit the original equation, replacing operands and operator with result
			input = replaceIndex(input, FloatToString(temp), addSubMatches[i]-1, addSubMatches[i]+1)
			fmt.Println("Input update: " + input)

		}

		// At this point the input should reflect the result
		fmt.Println("Result: " + input)

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

func validateInput(input string) (err error) {

	// Check if the user wants to exit app
	if strings.ToLower(input) == "exit" {
		return errors.New("app terminated")
	}

	// Validate against alpha runes and special characters
	runes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$^&_={}[]:;'\"<,>?\\|"
	for _, char := range runes {
		if strings.Index(input, string(char)) != -1 {
			return errors.New("invalid input: " + input)
		}
	}

	return nil

}

func findOperandLength(isFirst bool, i int, input string, matches []int) (operandLength int) {

	var operandComplete bool
	var validOperandRunes = "0123456789."

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
		if !validChar {
			operandComplete = true
			break
		}
		operandLength++
	}

	return

}

func FloatToString(num float64) string {

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
