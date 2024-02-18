package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"orchestrator/internal/core/ports"

	"github.com/google/uuid"
)

type Handlers struct {
	service ports.Service
}

type ExpressionJSON struct {
	Expression string `json:"expression"`
}

type ExpressionResponseJSON struct {
	ID uuid.UUID `json:"id"`
}

// Добавить вычисление арифметического выражения
func (h *Handlers) AddExpression(res http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	var expression ExpressionJSON
	if err := json.Unmarshal(buf.Bytes(), &expression); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	expressionID, err := h.service.AddExpression(expression.Expression)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	expressionResponse := ExpressionResponseJSON{
		ID: expressionID,
	}
	ans, err := json.Marshal(expressionResponse)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(ans)
}

type ExpressionsJSON struct {
	ID          uuid.UUID `json:"id"`
	Expressions string    `json:"expression"`
	Status      string    `json:"status"`
	Result      string    `json:"result"`
}

// - `/get_all` - получение списка выражений со статусами
func (h *Handlers) GetExpressions(res http.ResponseWriter, req *http.Request) {
	expressions, err := h.service.GetExpressions(req.Context())
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return

	}
	expressionsJSON := make([]*ExpressionsJSON, len(expressions))
	for i, expression := range expressions {
		expressionsJSON[i] = &ExpressionsJSON{
			ID:          expression.ID,
			Expressions: expression.Expression,
			Status:      expression.Status,
			Result:      expression.Result,
		}
	}

	ans, err := json.Marshal(expressionsJSON)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(ans)
}

type ExpressionByIDJSON struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
	Result string    `json:"result"`
}

// - получение значения выражения по его идентификатору
func (h *Handlers) GetExpressionByID(res http.ResponseWriter, req *http.Request) {
	id := req.Context().Value("id")
	ID, _ := uuid.Parse(id.(string))
	log.Println("id", ID)
	expression, err := h.service.GetExpression(ID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return

	}
	log.Println(expression)
	expressionResponse := ExpressionByIDJSON{
		ID:     expression.ID,
		Result: expression.Result,
		Status: expression.Status,
	}
	ans, err := json.Marshal(expressionResponse)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(ans)
}

type OperationTimeJSON struct {
	ID        uuid.UUID `json:"id"`
	Operation string    `json:"operation"`
	Duration  string    `json:"duration"`
	Result    string    `json:"result"`
}

// - `/get_operations` - получение списка доступных операций со временем их выполнения
func (h *Handlers) GetOperations(res http.ResponseWriter, req *http.Request) {
	operations, err := h.service.GetOperations(req.Context())
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	operationResponse := make([]OperationTimeJSON, len(operations))
	for i, operation := range operations {
		operationResponse[i] = OperationTimeJSON{
			ID:        operation.ID,
			Operation: operation.Op,
			Duration:  operation.Leadtime,
			Result:    operation.Result,
		}
	}
	ans, err := json.Marshal(operationResponse)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(ans)
}

// - `/get_task` - получение задачи для выполнения
func (h *Handlers) GetOperation(res http.ResponseWriter, req *http.Request) {}

type OperationJSON struct {
	ID        uuid.UUID `json:"id"`
	Operation string    `json:"operation"`
}

// - `/` - прием результата обработки данных
func (h *Handlers) UpdateOperation(res http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(req.Body); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	var operation OperationJSON
	if err := json.Unmarshal(buf.Bytes(), &operation); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	h.service.UpdateOperation(operation.ID, operation.Operation)
	log.Println("Результат операции:", operation.ID, ":", operation.Operation)
	res.WriteHeader(http.StatusOK)
}

func NewHandlers(service ports.Service) *Handlers {
	return &Handlers{
		service: service,
	}
}
