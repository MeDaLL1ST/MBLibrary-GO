package mblibrarygo

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

type Mb struct {
	host    string
	api_key string
	scheme  string
	ws      *websocket.Conn
}

func Init(host string, api_key string, scheme string) *Mb {
	mb := &Mb{host: host, api_key: api_key, scheme: scheme}
	return mb
}

func (mb *Mb) Add(key string, value string) error {

	data := map[string]interface{}{
		"key":   key,
		"value": value,
	}

	// Генерируем checksum
	datah := []byte(key + value + mb.api_key)
	hash := sha256.Sum256(datah)
	data["checksum"] = fmt.Sprintf("%x", hash[:])

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

func (mb *Mb) Close() {
	mb.ws.Close()
}

func (mb *Mb) Read() (string, error) {
	_, msg, err := mb.ws.ReadMessage()
	return string(msg), err
}

func (mb *Mb) ReadSync(key string) (string, error) {
	headers := http.Header{}
	headers.Add("Authorization", mb.api_key)
	u1 := url.URL{Scheme: "ws", Host: mb.host, Path: "/subscribe"}
	c, _, err := websocket.DefaultDialer.Dial(u1.String(), headers)
	defer c.Close()
	if err != nil {
		return "", fmt.Errorf("Error by connecting to WebSocket: %v", err)
	}

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
	_, msg, err := c.ReadMessage()
	if err != nil {
		return "", fmt.Errorf("Error reading message from WebSocket: %v", err)
	}
	return string(msg), nil
}
