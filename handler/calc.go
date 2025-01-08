package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"final/utils"
)

type ExpressionRequest struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

// работа с сервером
func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ExpressionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusUnprocessableEntity)
		return
	}

	result, err := utils.Calc(req.Expression)
	if err != nil {
		status := http.StatusUnprocessableEntity
		if err.Error() == "division by zero" {
			status = http.StatusInternalServerError
		}

		resp := Response{Error: err.Error()}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := Response{Result: strconv.FormatFloat(result, 'f', -1, 64)}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
