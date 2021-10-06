package pointsalvor

import (
	//"context"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

const (
	host = "https://api.todoist.com/rest/v1"

	tasksUrl    = "/tasks"
	sectionsUrl = "/sections"
	projectsUrl = "/projects"

	sectionQuery = "/projects_id=%s"
	taskQuery    = sectionQuery + "&section_id=%s"
	updQuery     = "/%s"
	taskClose    = updQuery + "/close"

	requestTimeLimit = time.Second * 15
)

var (
	errUrl         = errors.New("url request is empty")
	errMethod      = errors.New("method is invalid")
	errJsonMarshal = errors.New("eror with marshal request body")
	errRequestForm = errors.New("request is invalid")
	errRequestPerf = errors.New("error with perform request")
	errModelValid  = errors.New("invalid model")
)

type Agent struct {
	heart *http.Client
	token string
}

func NewAgent(tokenApi string) (*Agent, error) {
	if tokenApi == "" {
		return nil, errors.New("tokenApi is empty")
	}

	return &Agent{
		heart: &http.Client{
			Timeout: requestTimeLimit,
		},
		token: tokenApi,
	}, nil
}

var MethodBank string = fmt.Sprintf(http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut)

func (ag *Agent) KnockToApi(ctx context.Context, method string, rout string, reqBody interface{}) (*http.Response, error) {
	//validate rout
	if rout == "" {
		return nil, errUrl
	}

	//validate method
	if method == "" || !strings.Contains(MethodBank, method) {
		return nil, errMethod
	}

	//						  POST
	//request body is json-encoded
	//data in json-object

	//				GET-PUT-DELETE
	//request body is url-encoded
	//data sends in url-query params
	reqBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errJsonMarshal
	}

	//creating the request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, host+projectsUrl, bytes.NewBuffer(nil))
	if err != nil {
		return nil, errRequestForm
	}

	//authorization
	req.Header.Set("Authorization", "Bearer "+ag.token)

	//perform request
	resp, err := ag.heart.Do(req)
	if err != nil {
		return nil, errRequestPerf
	}

	defer resp.Body.Close()

	//validate response body
	if code := resp.StatusCode; code != http.StatusOK {
		return nil, errors.New("SDK eror: " + strconv.Itoa(code))
	}

	return resp, nil
}

//bank of models
var ModelBank string = fmt.Sprintf("%s-%s-%s", "project", "section", "task")

//validate models on exist
func ValidateModel(model string) bool {
	return strings.Contains(ModelBank, model)
}

//model-string map
func ModelMapping(model string) interface{} {
	switch model {
	case "project":
		return Project{}
	case "section":
		return Section{}
	case "task":
		return Task{}
	}

	return nil
}

func RePointer(bank []*interface{}) []interface{} {
	var store []interface{}
	for _, v := range bank {
		store = append(store, *v)
	}

	return store
}

func DecodeResponseToModels(resp *http.Response, model string) ([]interface{}, error) {
	if valid := ValidateModel(model); !valid {
		return nil, errModelValid
	}

	//mapping interface to model and init bank
	inp := ModelMapping(model)
	/*
		switch _ := inp.(type) {
		case Project:
			var storage []Project
		case Section:
			var storage []Section
		case Task:
			var storage []Task
		} */

	//bank for struct
	var (
		input   []interface{}
		storage []interface{}
	)

	//response to []interface{}
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&input)
	if err != nil {
		return nil, fmt.Errorf("%s - [%s]", "error with decoding input", err.Error())
	}

	//convert []interface}{} to map
	//decode map to struct
	//add storage of struct
	for _, part := range input {
		res, ok := part.(map[string]interface{})
		if ok {
			err = mapstructure.Decode(res, &inp)
			if err != nil {
				return nil, fmt.Errorf("%s - [%s]", "error with decoding map to struct", err.Error())
			}
			storage = append(storage, inp)
		}
	}

	return storage, nil
}
