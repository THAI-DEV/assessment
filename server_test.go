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
	"github.com/THAI-DEV/assessment/handler"
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

	res := request(http.MethodPost, uri("expenses"), body)

	err := res.Decode(&result)
	if err != nil {
		t.Fatal("Can not create expense", err)
	}

	assert.EqualValues(t, http.StatusCreated, res.StatusCode)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.EqualValues(t, inputJson.Title, result.Title)
	assert.EqualValues(t, inputJson.Amount, result.Amount)
	assert.EqualValues(t, inputJson.Note, result.Note)
	assert.EqualValues(t, inputJson.Tags, result.Tags)

	id, _ := strconv.Atoi(result.Id)
	assert.Greater(t, id, 0)
}

func TestReadOneDataCaseFound(t *testing.T) {
	var result database.Expense

	res := request(http.MethodGet, uri("expenses", "1"), nil)

	err := res.Decode(&result)
	if err != nil {
		t.Fatal("Can not read expenses", err)
	}

	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.NotEqualValues(t, "", result.Title)
	assert.NotEqualValues(t, "", result.Amount)
	assert.NotEqualValues(t, "", result.Note)
	assert.NotEqualValues(t, "", result.Tags)

	assert.Greater(t, len(result.Tags), 0, "Len > 0")
}

func TestReadOneDataCaseNotFound(t *testing.T) {
	res := request(http.MethodGet, uri("expenses", "0"), nil)
	assert.EqualValues(t, http.StatusBadRequest, res.StatusCode)
}

func (rcv *Response) Decode(v interface{}) error {
	if rcv.err != nil {
		return rcv.err
	}

	return json.NewDecoder(rcv.Body).Decode(v)
}

func TestUpdateDataCaseNoId(t *testing.T) {
	var result database.Expense
	var inputJson handler.ExpenseBody

	idParam := "2"

	str := `{
		"title": "apple smoothie",
		"amount": 89,
		"note": "no discount",
		"tags": ["beverage"]
	}`

	body := bytes.NewBufferString(str)
	json.Unmarshal([]byte(str), &inputJson)

	res := request(http.MethodPut, uri("expenses", idParam), body)

	err := res.Decode(&result)
	if err != nil {
		t.Fatal("Can not create expense", err)
	}

	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.EqualValues(t, inputJson.Title, result.Title)
	assert.EqualValues(t, inputJson.Amount, result.Amount)
	assert.EqualValues(t, inputJson.Note, result.Note)
	assert.EqualValues(t, inputJson.Tags, result.Tags)

	id, _ := strconv.Atoi(result.Id)
	assert.Greater(t, id, 0)
}

func TestUpdateDataCaseId(t *testing.T) {
	var result database.Expense
	var inputJson handler.ExpenseBody

	idParam := "2"

	str := `{
		"id": 0,
		"title": "apple smoothie",
		"amount": 89,
		"note": "no discount",
		"tags": ["beverage"]
	}`

	body := bytes.NewBufferString(str)
	json.Unmarshal([]byte(str), &inputJson)

	res := request(http.MethodPut, uri("expenses", idParam), body)

	err := res.Decode(&result)
	if err != nil {
		t.Fatal("Can not create expense", err)
	}

	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.EqualValues(t, inputJson.Title, result.Title)
	assert.EqualValues(t, inputJson.Amount, result.Amount)
	assert.EqualValues(t, inputJson.Note, result.Note)
	assert.EqualValues(t, inputJson.Tags, result.Tags)

	id, _ := strconv.Atoi(result.Id)
	assert.Greater(t, id, 0)
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
