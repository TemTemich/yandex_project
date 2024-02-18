package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/Pramod-Devireddy/go-exprtk"
	"github.com/google/uuid"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	orcHost := os.Getenv("SERVICE_URL")

	client := &http.Client{}
	orch := NewOrchestratorAdapter(
		orcHost,
		client,
	)

	workers := NewWorkers(10)
	go workers.Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/add_task", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var buf bytes.Buffer
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		type Request struct {
			ID        uuid.UUID `json:"id"`
			Operation string    `json:"operation"`
		}

		var req Request
		if err := json.Unmarshal(buf.Bytes(), &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println(orch)
		t := NewTask(
			orch,
			req.ID,
			req.Operation,
		)
		workers.AddTask(t)
		w.WriteHeader(http.StatusOK)
	})

	return http.ListenAndServe(":3030", mux)
}

// Конфигурация системы
type Config struct {
	Addr           string
	CountGoroutine int
}

func NewConfig() *Config {
	addr := os.Getenv("ORCHESTRATOR_ADDR")
	if addr == "" {
		addr = "http://localhost:8000"
	}
	countGoroutine := os.Getenv("DAEMON_COUNT_GOROUTINE")
	if countGoroutine == "" {
		countGoroutine = "10"
	}
	count, _ := strconv.Atoi(countGoroutine)

	return &Config{
		Addr:           addr,
		CountGoroutine: count,
	}
}

type orchestratorAdapter struct {
	client *http.Client
	addr   string
}

type OrhestratorAPIJSON struct {
	ID        uuid.UUID `json:"id"`
	Operation string    `json:"operation"`
}

func (o *orchestratorAdapter) UpdateOperation(id uuid.UUID, operation string) {
	log.Printf("id: %s, operation: %s\n", id, operation)
	url := fmt.Sprintf("%s/operation/update", o.addr)
	reqJSON := &OrhestratorAPIJSON{
		ID:        id,
		Operation: operation,
	}
	resBody, err := json.Marshal(reqJSON)
	if err != nil {
		log.Println("Не удалось сформировать запрос", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(resBody))
	if err != nil {
		log.Println("Не удалось создать запрос на обновление операции", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := o.client.Do(req)
	if err != nil {
		log.Println("Не удалось отправить запрос на обновление операции", err)
	}
	defer resp.Body.Close()

	//o.client.Get(fmt.Sprintf("%s/update_operation?id=%s&operation=%s", o.addr, id, operation))
}

func NewOrchestratorAdapter(addr string, client *http.Client) *orchestratorAdapter {
	return &orchestratorAdapter{
		client: client,
		addr:   addr,
	}
}

type orchestartor interface {
	UpdateOperation(id uuid.UUID, result string)
}

type Task struct {
	id        uuid.UUID
	operation string

	orchestartor orchestartor
}

func (t *Task) Do() {
	// Получение ответ
	exprtkObj := exprtk.NewExprtk()
	defer exprtkObj.Delete()
	exprtkObj.SetExpression(t.operation)
	exprtkObj.CompileExpression()
	ans := exprtkObj.GetEvaluatedValue()

	//time.Sleep(5 * time.Second)
	// Переводим ответ в строку, а затем отправляем его на сервер
	ansString := fmt.Sprintf("%v", ans) // float64 -> string
	t.orchestartor.UpdateOperation(t.id, ansString)
}

func NewTask(orchestartor orchestartor, id uuid.UUID, operation string) *Task {
	return &Task{
		id:           id,
		operation:    operation,
		orchestartor: orchestartor,
	}
}

type Workers struct {
	queueCapacity int
	secondQueue   []*Task

	done  chan struct{}
	queue chan *Task

	wg *sync.WaitGroup
	mu *sync.Mutex
}

func (w *Workers) AddTask(task *Task) {
	select {
	case w.queue <- task:
		return
	default:
		w.mu.Lock()
		log.Printf("count task in queue: %d count task in second queue: %d\n", len(w.queue), len(w.secondQueue))
		w.secondQueue = append(w.secondQueue, task)
		w.mu.Unlock()
	}
}

func (w *Workers) doneTask() {
	w.done <- struct{}{}
}

func (w *Workers) Run() {
	w.wg.Add(w.queueCapacity)
	for i := 0; i < w.queueCapacity; i++ {
		go func() {
			defer w.wg.Done()
			for {
				select {
				// case <-w.done:
				// 	if len(w.secondQueue) != 0 {
				// 		w.mu.Lock()
				// 		elem := w.secondQueue[0]
				// 		w.secondQueue = w.secondQueue[1:]
				// 		w.mu.Unlock()
				// 		w.queue <- elem
				// 	}
				case task := <-w.queue:
					//log.Printf("goroutine id: %v start task id: %v", i, task.id)
					//log.Printf("secoend queue size: %d\n", len(w.secondQueue))
					task.Do()
					//go w.doneTask()
				}
			}
		}()
	}
	w.wg.Wait()
}

func NewWorkers(queueCapacity int) *Workers {

	return &Workers{
		queueCapacity: queueCapacity,
		done:          make(chan struct{}),
		queue:         make(chan *Task, queueCapacity),
		wg:            &sync.WaitGroup{},
		mu:            &sync.Mutex{},
	}
}

//
