package calculator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"orchestrator/internal/core/enteties"

	"github.com/google/uuid"
)

type Calculator struct {
	addr   string
	client *http.Client
}

type CalculatorJSON struct {
	ID        uuid.UUID `json:"id"`
	Operation string    `json:"operation"`
}

func (c *Calculator) Calculate(operation *enteties.Operation) {
	log.Println("calculate", operation.ID, operation.Op)
	url := fmt.Sprintf("%s/add_task", c.addr)
	reqJSON := &CalculatorJSON{
		ID:        operation.ID,
		Operation: operation.Op,
	}
	reqBody, err := json.Marshal(reqJSON)
	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		log.Println(err)

	}
	defer resp.Body.Close()

}

func NewCalculator(addr string, client *http.Client) *Calculator {
	return &Calculator{
		addr:   addr,
		client: client,
	}
}
