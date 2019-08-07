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
	fmt.Println("---Note that OPERANDS CANNOT EXCEED ONE DIGIT---")
	fmt.Println("---Type 'exit' at any time to terminate the app---")

	var err error

	for err == nil {

		// Retrieve input
		fmt.Print("Enter a mathematical expression using operators (+,-,/,*,%) and numbers only: ")
		fmt.Scanln(&input)

		// Validate input
		err = validateInput(input)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Determine which operators are present, in PEMDAS order
		operatorOptions := regexp.MustCompile("[+-/*%]")
		matches := operatorOptions.FindAllStringIndex(input, -1)

		// Initialize match slices
		multiplicationMatches := []int{}
		divisionMatches := []int{}
		modulusMatches := []int{}
		additionMatches := []int{}
		subtractionMatches := []int{}

		// Determine the location of each instance of each operator in the string
		for i:=0; i<len(matches); i++ {
			if string(input[matches[i][0]]) == "*" {
				multiplicationMatches = append(multiplicationMatches, matches[i][0])
			} else if string(input[matches[i][0]]) == "/" {
				divisionMatches = append(divisionMatches, matches[i][0])
			} else if string(input[matches[i][0]]) == "%" {
				modulusMatches = append(modulusMatches, matches[i][0])
			} else if string(input[matches[i][0]]) == "+" {
				additionMatches = append(additionMatches, matches[i][0])
			} else if string(input[matches[i][0]]) == "-" {
				subtractionMatches = append(subtractionMatches, matches[i][0])
			} else {
				err = errors.New("operator not supported")
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

			// TODO: determine operation without iterating through each item (possibly with [][]interface{})

			// Determine whether multiplication or division
			for mi:=0; mi<len(multiplicationMatches); mi++ {
				if multiplicationMatches[mi] == multDivModMatches[i] {
					operation = "multiplication"
				}
			}
			// Let's iterate as little as possible by checking for operation first
			if operation == "" {
				for di:=0; di<len(divisionMatches); di++ {
					if divisionMatches[di] == multDivModMatches[i] {
						operation = "division"
					}
				}
				if operation == "" {
					for modi:=0; modi<len(modulusMatches); modi++ {
						if modulusMatches[modi] == multDivModMatches[i] {
							operation = "modulus"
						}
					}
				}
			}

			firstOperand, err := strconv.Atoi(string(input[multDivModMatches[i] - 1]))
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			secondOperand, err := strconv.Atoi(string(input[multDivModMatches[i] + 1]))
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
			input = replaceIndex(input, FloatToString(temp), multDivModMatches[i]-1, multDivModMatches[i]+1)
			fmt.Println("Input update: " + input)

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

	fmt.Println(err.Error())
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
	runes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$^&_={}[]:;'\"<,>.?\\|"
	for _, char := range runes {
		if strings.Index(input, string(char)) != -1 {
			return errors.New("invalid input: " + input)
		}
	}

	return nil

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
