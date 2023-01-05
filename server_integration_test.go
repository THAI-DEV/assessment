//go:build integration
// +build integration

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

type Response2 struct {
	*http.Response
	err error
}

var testId = "1"

func Test1CreateData(t *testing.T) {
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

	res := request2(http.MethodPost, uri2("expenses"), body)

	err := res.Decode2(&result)
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

	testId = result.Id
}

func Test2ReadOneData(t *testing.T) {
	var result database.Expense

	res := request2(http.MethodGet, uri2("expenses", testId), nil)

	err := res.Decode2(&result)
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

func Test3UpdateData(t *testing.T) {
	var result database.Expense
	var inputJson handler.ExpenseBody

	idParam := testId

	str := `{
		"id": 0,
		"title": "apple smoothie",
		"amount": 89,
		"note": "no discount",
		"tags": ["beverage"]
	}`

	body := bytes.NewBufferString(str)
	json.Unmarshal([]byte(str), &inputJson)

	res := request2(http.MethodPut, uri2("expenses", idParam), body)

	err := res.Decode2(&result)
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

func Test4ReadAll(t *testing.T) {
	var result []database.Expense

	res := request2(http.MethodGet, uri2("expenses"), nil)

	err := res.Decode2(&result)
	if err != nil {
		t.Fatal("Can not read all expenses", err)
	}

	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.GreaterOrEqual(t, len(result), 1, "Len >= 1")
}

func (rcv *Response2) Decode2(v interface{}) error {
	if rcv.err != nil {
		return rcv.err
	}

	return json.NewDecoder(rcv.Body).Decode(v)
}

func request2(method, url string, body io.Reader) *Response2 {
	req, _ := http.NewRequest(method, url, body)

	req.Header.Add("Authorization", "November 10, 2009")
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)

	return &Response2{res, err}
}

func uri2(paths ...string) string {
	host := "http://container_rest:2565"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)

	return strings.Join(url, "/")
}
