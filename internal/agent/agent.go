package agent

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Smer4k/slow-web-calculator/internal/datatypes"
	"github.com/gorilla/mux"
)

type Agent struct {
	Router         *mux.Router
	AddrMainServer string
	AddrAgent      string
	Status         datatypes.Status
	CurrentTask    datatypes.Task
}

func NewAgent(orchestratorPort, port string) *Agent {
	return &Agent{
		Router:         mux.NewRouter(),
		AddrMainServer: "http://localhost" + orchestratorPort + "/",
		AddrAgent:      "http://localhost" + port + "/",
		Status:         datatypes.Idle,
	}
}

func (a *Agent) InitAgent() {
	a.Router.HandleFunc("/", a.redirectToMainServer).Methods(http.MethodGet)
	http.Handle("/", a.Router)
	a.AddAgentToMainServer(true)
	a.PingMainServer()
}

// добавляет адрес агента в список оркестра
func (a *Agent) AddAgentToMainServer(check bool) {
	vals := url.Values{}
	vals.Add("server", a.AddrAgent)
	_, err := http.PostForm(a.AddrMainServer+"addServer", vals)
	if err != nil {
		panic(err)
	}
	if !check {
		a.Status = datatypes.Idle
	}
}

// запрашивает новую задачу
func (a *Agent) GetNewTask() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sgetExpression?%s=%s", a.AddrMainServer, "agent", a.AddrAgent), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("agent", a.AddrAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	var task datatypes.Task

	if contentType := resp.Header.Get("Content-Type"); contentType == "application/json" {
		if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
			fmt.Println("Ошибка при чтении тела ответа:", err)
			return
		}
		a.CurrentTask = task
		a.Status = datatypes.Work
		fmt.Printf("Получил %s%s%s приступаю к работе\n", task.Expression.Left, task.Expression.Operator, task.Expression.Right)
		go a.SolveExpression()
	}
}

// отправляет решение
func (a *Agent) PostAnswer() {
	vals := url.Values{}
	jsonData, err := json.Marshal(a.CurrentTask)
	if err != nil {
		fmt.Println(err)
		return
	}
	vals.Add("answer", string(jsonData))
	_, err = http.PostForm(a.AddrMainServer+"postAnswer", vals)
	if err != nil {
		fmt.Println("Главный сервер не отвечает, пробую повторно отправить решение")
		ticker := time.NewTicker(2 * time.Second)
		for range ticker.C {
			_, err = http.PostForm(a.AddrMainServer+"postAnswer", vals)
			if err != nil {
				fmt.Println("Главный сервер не отвечает, пробую повторно отправить решение")
			} else {
				ticker.Stop()
				break
			}
		}
	}

	a.CurrentTask = datatypes.Task{}
	a.Status = datatypes.Idle
	fmt.Println("Отправил решение, ожидаю новой задачи")
}

// выполняет решение выражения
func (a *Agent) SolveExpression() {
	time.Sleep(time.Duration(a.CurrentTask.TimeExec) * time.Second)
	var leftNum, rightNum, total float64
	leftNum, _ = strconv.ParseFloat(a.CurrentTask.Expression.Left, 64)
	rightNum, _ = strconv.ParseFloat(a.CurrentTask.Expression.Right, 64)
	switch a.CurrentTask.Expression.Operator {
	case "+":
		total = leftNum + rightNum
	case "-":
		total = leftNum - rightNum
	case "*":
		total = leftNum * rightNum
	case "/":
		total = leftNum / rightNum
	}
	a.CurrentTask.Answer = strconv.FormatFloat(total, 'g', -1, 64)
	fmt.Println("Задача решена, ответ: ", total)
	go a.PostAnswer()
}

// проверяет работоспособность оркестра и делает запрос на получение задачи если агент стоит без дела
func (a *Agent) PingMainServer() {

	ticker := time.NewTicker(time.Duration((float32(5) + rand.Float32())) * time.Second)
	failConnect := false
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			_, err := http.Get(a.AddrMainServer)
			if err != nil {
				if !failConnect {
					a.Status = datatypes.Reconnect
					fmt.Printf("Главный сервер не отвечает, ошибка:\n%s\n", err)
					failConnect = true
				}
				continue
			} else {
				if failConnect {
					a.AddAgentToMainServer(false)
					failConnect = false
				}
				if a.Status == datatypes.Idle {
					a.AddAgentToMainServer(true)
					go a.GetNewTask()
				}
			}
		}
	}()
}
