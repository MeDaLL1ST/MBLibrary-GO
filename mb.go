package mblibrarygo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type Mb struct {
	host    string
	api_key string
	scheme  string
	ws      *websocket.Conn
}

type Gw struct {
	host    string
	api_key string
	scheme  string
}

type node struct {
	Id     int      `json:"Id"`
	Topics []string `json:"Topics"`
	IP     string   `json:"IP"`
	Scheme string   `json:"Scheme"`
	APIKey string   `json:"APIKey"`
}

func InitMb(host string, api_key string, scheme string) *Mb {
	mb := &Mb{host: host, api_key: api_key, scheme: scheme}
	return mb
}

func InitGw(host string, api_key string, scheme string) *Gw {
	gw := &Gw{host: host, api_key: api_key, scheme: scheme}
	return gw
}

func (gw *Gw) Add(key string, value string, topic ...string) error {

	data := make(map[string]interface{})

	if len(topic) > 0 {
		data["key"] = key
		data["value"] = value
		data["topic"] = topic[0]

	} else {
		data["key"] = key
		data["value"] = value
	}

	// Преобразуем data в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}
	// Отправляем POST-запрос
	req, err := http.NewRequest("GET", gw.scheme+"://"+gw.host+"/add", nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Добавляем заголовок Content-Type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", gw.api_key)
	// Добавляем данные в тело запроса
	req.Body = ioutil.NopCloser(bytes.NewReader(jsonData))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err //fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	return nil
}

func (mb *Mb) Add(key string, value string) error {

	data := map[string]interface{}{
		"key":   key,
		"value": value,
	}

	// Преобразуем data в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}
	// Отправляем POST-запрос
	req, err := http.NewRequest("GET", mb.scheme+"://"+mb.host+"/add", nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Добавляем заголовок Content-Type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", mb.api_key)
	// Добавляем данные в тело запроса
	req.Body = ioutil.NopCloser(bytes.NewReader(jsonData))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err //fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	return nil
}

func (gw *Gw) Info() ([]node, error) {

	// Отправляем POST-запрос
	req, err := http.NewRequest("GET", gw.scheme+"://"+gw.host+"/info", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Добавляем заголовок Content-Type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", gw.api_key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}
	// Парсим JSON
	var jsonData map[string][]node
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return jsonData["nodes"], nil
}

func (mb *Mb) List() ([]interface{}, error) {

	// Отправляем POST-запрос
	req, err := http.NewRequest("GET", mb.scheme+"://"+mb.host+"/list", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Добавляем заголовок Content-Type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", mb.api_key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}
	// Парсим JSON
	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	// Получаем массив ключей
	keys := jsonData["keys"].([]interface{})

	return keys, nil
}

func (mb *Mb) Info(key string) (int, error) {

	data := map[string]interface{}{
		"key": key,
	}

	// Преобразуем data в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return -1, fmt.Errorf("error marshaling JSON: %v", err)
	}
	// Отправляем POST-запрос
	req, err := http.NewRequest("GET", mb.scheme+"://"+mb.host+"/info", nil)
	if err != nil {
		return -1, fmt.Errorf("error creating request: %v", err)
	}

	// Добавляем заголовок Content-Type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", mb.api_key)

	// Добавляем данные в тело запроса
	req.Body = ioutil.NopCloser(bytes.NewReader(jsonData))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return -1, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, fmt.Errorf("error reading response: %v", err)
	}
	// Парсим JSON
	var jsonData2 map[string]int
	err = json.Unmarshal(body, &jsonData2)
	if err != nil {
		return -1, fmt.Errorf("error unmarshaling JSON: %v", err)
	}
	if _, ok := jsonData2["subs"]; !ok {
		fmt.Println("Field \"subs\" not found")
		return -1, fmt.Errorf("JSON dont have field \"subs\"")
	}
	// Получаем значение subs
	subscriptionCount, ok := jsonData2["subs"]
	if !ok {
		return -1, fmt.Errorf("invalid JSON structure: expected \"subs\" field")
	}
	return subscriptionCount, nil
}

func (mb *Mb) Subscribe(key string) error {
	headers := http.Header{}
	headers.Add("Authorization", mb.api_key)
	u1 := url.URL{Scheme: "ws", Host: mb.host, Path: "/subscribe"}
	c, _, err := websocket.DefaultDialer.Dial(u1.String(), headers)
	if err != nil {
		return fmt.Errorf("Error by connecting to WebSocket: %v", err)
	}

	// Создаем JSON-пayload для отправки
	payload := map[string]string{
		"key": key,
	}

	// Преобразуем payload в байтовый массив
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Error marshaling JSON: %v", err)
	}

	// Отправляем сообщение
	err = c.WriteMessage(websocket.TextMessage, jsonBytes)
	if err != nil {
		return fmt.Errorf("Error sending message to WebSocket: %v", err)
	}
	if mb.ws == nil {
		mb.ws = c
	} else {
		mb.Close()
		mb.ws = c
	}
	return nil
}

func (mb *Mb) ReSubscribe(key string) error {
	payload := map[string]string{
		"key": key,
	}

	// Преобразуем payload в байтовый массив
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Error marshaling JSON: %v", err)
	}

	// Отправляем сообщение
	err = mb.ws.WriteMessage(websocket.TextMessage, jsonBytes)
	if err != nil {
		return fmt.Errorf("Error sending message to WebSocket: %v", err)
	}
	return nil
}

func (mb *Mb) Close() {
	mb.ws.Close()
}

func (mb *Mb) Read() (string, error) {
	_, msg, err := mb.ws.ReadMessage()
	return string(msg), err
}

func (mb *Mb) ReadSync(key string, functions ...func()) (string, error) {
	headers := http.Header{}
	headers.Add("Authorization", mb.api_key)
	u1 := url.URL{Scheme: "ws", Host: mb.host, Path: "/subscribe"}
	c, _, err := websocket.DefaultDialer.Dial(u1.String(), headers)
	if err != nil {
		return "", fmt.Errorf("Error by connecting to WebSocket: %v", err)
	}
	defer c.Close()
	// Создаем JSON-пayload для отправки
	payload := map[string]string{
		"key": key,
	}

	// Преобразуем payload в байтовый массив
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("Error marshaling JSON: %v", err)
	}

	// Отправляем сообщение
	err = c.WriteMessage(websocket.TextMessage, jsonBytes)
	if err != nil {
		return "", fmt.Errorf("Error sending message to WebSocket: %v", err)
	}
	for _, fn := range functions {
		go fn()
	}
	_, msg, err := c.ReadMessage()
	if err != nil {
		return "", fmt.Errorf("Error reading message from WebSocket: %v", err)
	}
	return string(msg), nil
}

func (gw *Gw) ReadSync(key string, topic string, functions ...func()) (string, error) {

	nodes, err := gw.Info()
	if err != nil {
		return "", fmt.Errorf("Error getting info: %v", err)
	}

	var target_node node
	exists := false
	for _, nd := range nodes {
		for _, top := range nd.Topics {
			if top == topic {
				exists = true
				target_node = nd
			}
		}
	}

	if !exists {
		return "", fmt.Errorf("No some topic")
	}

	if target_node.APIKey == "" {
		target_node.APIKey = gw.api_key
	}

	headers := http.Header{}
	headers.Add("Authorization", target_node.APIKey)
	u1 := url.URL{Scheme: "ws", Host: target_node.IP, Path: "/subscribe"}
	c, _, err := websocket.DefaultDialer.Dial(u1.String(), headers)
	if err != nil {
		return "", fmt.Errorf("Error by connecting to WebSocket: %v", err)
	}
	defer c.Close()
	// Создаем JSON-пayload для отправки
	payload := map[string]string{
		"key": key,
	}

	// Преобразуем payload в байтовый массив
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("Error marshaling JSON: %v", err)
	}

	// Отправляем сообщение
	err = c.WriteMessage(websocket.TextMessage, jsonBytes)
	if err != nil {
		return "", fmt.Errorf("Error sending message to WebSocket: %v", err)
	}
	time.Sleep(time.Millisecond * 50)
	for _, fn := range functions {
		go fn()
	}

	_, msg, err := c.ReadMessage()
	if err != nil {
		return "", fmt.Errorf("Error reading message from WebSocket: %v", err)
	}
	return string(msg), nil
}