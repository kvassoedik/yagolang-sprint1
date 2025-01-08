package utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// порядок операций
func Precedence(op string) int {
	if op == "*" || op == "/" {
		return 3
	}
	return 2
}

// переделать строку в массив токенов
func ToArray(expr string) []string {
	var arr []string
	var num strings.Builder

	for _, char := range expr {
		switch char {
		case '+', '-', '*', '/', '(', ')':
			if num.Len() > 0 {
				arr = append(arr, num.String())
				num.Reset()
			}
			arr = append(arr, string(char))
		default:
			num.WriteRune(char)
		}
	}
	if num.Len() > 0 {
		arr = append(arr, num.String())
	}
	return arr
}

// StackToOutput appends remaining operators from the stack to the output.
func StackToOutput(output, stack []string) []string {
	for i := len(stack) - 1; i >= 0; i-- {
		output = append(output, stack[i])
	}
	return output
}

// алгоритм для вычисления префиксной формы
func ShuntingYard(arr []string) ([]string, error) {
	var output, stack []string
	prev := ""

	for _, token := range arr {
		if regexp.MustCompile(`^\d+$`).MatchString(token) {
			output = append(output, token)
		} else if strings.Contains("+-*/", token) {
			if prev == "" || strings.Contains("+-*/", prev) {
				return nil, errors.New("Expression is not valid")
			}
			for len(stack) > 0 && stack[len(stack)-1] != "(" && Precedence(stack[len(stack)-1]) >= Precedence(token) {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		} else if token == "(" {
			stack = append(stack, token)
		} else if token == ")" {
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 || stack[len(stack)-1] != "(" {
				return nil, errors.New("Mismatched parentheses")
			}
			stack = stack[:len(stack)-1]
		} else {
			return nil, errors.New("Expression is not valid")
		}
		prev = token
	}

	return StackToOutput(output, stack), nil
}

// вычисление выражения
func Calc(expression string) (float64, error) {
	arr := ToArray(expression)
	postfix, err := ShuntingYard(arr)
	if err != nil {
		return 0, err
	}

	var stack []float64
	for _, token := range postfix {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, errors.New("Expression is not valid")
			}
			n2, n1 := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch token {
			case "+":
				stack = append(stack, n1+n2)
			case "-":
				stack = append(stack, n1-n2)
			case "*":
				stack = append(stack, n1*n2)
			case "/":
				if n2 == 0 {
					return 0, errors.New("division by zero")
				}
				stack = append(stack, n1/n2)
			}
		}
	}
	return stack[0], nil
}
