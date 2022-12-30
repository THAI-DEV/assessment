package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/THAI-DEV/assessment/database"
	"github.com/stretchr/testify/assert"
)

type Response struct {
	*http.Response
	err error
}

func TestCreateData(t *testing.T) {
	var result database.Expense
	var inputJson database.Expense

	str := `{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`

	body := bytes.NewBufferString(str)
	json.Unmarshal([]byte(str), &inputJson)

	err := request(http.MethodPost, uri("expenses"), body).Decode(&result)
	if err != nil {
		t.Fatal("Can not create expense", err)
	}

	assert.EqualValues(t, 201, http.StatusCreated)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.EqualValues(t, inputJson.Title, result.Title)
	assert.EqualValues(t, inputJson.Amount, result.Amount)
	assert.EqualValues(t, inputJson.Note, result.Note)
	assert.EqualValues(t, inputJson.Tags, result.Tags)

	id, _ := strconv.Atoi(result.Id)
	assert.Greater(t, id, 0)
}

func (rcv *Response) Decode(v interface{}) error {
	if rcv.err != nil {
		return rcv.err
	}

	return json.NewDecoder(rcv.Body).Decode(v)
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)

	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)

	return &Response{res, err}
}

func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)

	return strings.Join(url, "/")
}
