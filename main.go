package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strconv"
)

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func Precedence(op string) int {
	if op == string('*') || op == string('/') {
		return 3
	} else {
		return 2
	}
}

func ToArray(expr string) []string {
	var arr []string
	var num string = ""
	for _, char := range expr {
		if char == '+' || char == '-' || char == '*' || char == '/' || char == '(' || char == ')' {
			if len(num) > 0 {
				arr = append(arr, num)
				num = ""
			}
			arr = append(arr, string(char))
		} else {
			num += string(char)
		}
	}
	if len(num) > 0 {
		arr = append(arr, num)
		num = ""
	}
	return arr
}

func StackToOutput(output, stack []string) []string {
	for i := len(stack) - 1; i >= 0; i-- {
		output = append(output, stack[i])
	}
	return output
}

func isOperator(op string) bool {
	if op == string('+') || op == string('-') || op == string('*') || op == string('/') {
		return true
	}
	return false
}

func isParenthesis(op string) bool {
	if op == string('(') || op == string(')') {
		return true
	}
	return false
}

func ShuntingYard(arr []string) ([]string, error) {
	var output []string
	var stack []string

	if !regexp.MustCompile(`\d`).MatchString(arr[len(arr)-1]) && string(arr[len(arr)-1]) != ")" {
		return output, errors.New("invalid expression")
	}

	if !regexp.MustCompile(`\d`).MatchString(arr[0]) && string(arr[0]) != "(" {
		return output, errors.New("invalid expression")
	}

	prev := ""

	for _, token := range arr {
		if regexp.MustCompile(`\d`).MatchString(token) {
			output = append(output, token)
		} else if isOperator(token) {
			if isOperator(prev) {
				return output, errors.New("invalid expression")
			}
			for len(stack) > 0 && string(stack[len(stack)-1]) != "(" && Precedence(string(token)) <= Precedence(string(stack[len(stack)-1])) {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		} else if token == "(" {
			if isParenthesis(prev) {
				return output, errors.New("invalid expression")
			}
			stack = append(stack, token)
		} else if token == ")" {
			if isParenthesis(prev) {
				return output, errors.New("invalid expression")
			}
			for stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1]
		}
		prev = token
	}
	if slices.Contains(stack, "(") {
		return output, errors.New("mismatched parenthesis")
	}
	output = StackToOutput(output, stack)
	return output, nil
}

func Calc(expression string) (float64, error) {
	if len(expression) <= 2 {
		return 0, errors.New("invalid expression")
	}

	arr := ToArray(expression)

	postfix, err := ShuntingYard(arr)
	if err != nil {
		return 0, err
	}

	var stack []float64
	for _, elem := range postfix {
		num, err := strconv.ParseFloat(elem, 64)
		if err == nil {
			stack = append(stack, num)
		} else {
			n1 := stack[len(stack)-2]
			n2 := stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			switch elem {
			case string('+'):
				stack = append(stack, n1+n2)
			case string('-'):
				stack = append(stack, n1-n2)
			case string('*'):
				stack = append(stack, n1*n2)
			default:
				if n2 == 0 {
					return 0, errors.New("division by zero")
				}
				stack = append(stack, n1/n2)
			}
		}
	}

	return stack[0], nil
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req Request
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Error: "Internal server error"})
		return
	}

	result, err := Calc(req.Expression)
	if err != nil {
		if err.Error() == "invalid expression" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(Response{Error: "Expression is not valid"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{Error: "Internal server error"})
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Result: fmt.Sprintf("%f", result)})
}

func main() {
	http.HandleFunc("/api/v1/calculate", calculateHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
