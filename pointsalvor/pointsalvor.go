package pointsalvor

import (
	//"context"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
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

	delOpt = "del"
	addOpt = "add"
)

var (
	errUrl         = errors.New("url request is empty")
	errMethod      = errors.New("method is invalid")
	errJsonMarshal = errors.New("eror with marshal request body")
	errRequestForm = errors.New("request is invalid")
	errRequestPerf = errors.New("error with perform request")
	errModelValid  = errors.New("invalid model")
	errSdk         = errors.New("SDK error")
	errOptional    = errors.New("invalid optional updating the bankId")
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

var MethodBank string = fmt.Sprintf("%s-%s-%s-%s", http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut)

func ValidateMethod(method string) bool {
	return strings.Contains(MethodBank, method)
}

var MappingStatusCode = map[string]int{
	"json":       200,
	"no-content": 204,
}

//fucking validate
func ValidateStatusCode(code int) bool {
	for _, cd := range MappingStatusCode {
		if code == cd {
			return true
		}
	}
	return false
}

func (ag *Agent) KnockToApi(ctx context.Context, method string, rout string, reqBody interface{}) (*http.Response, error) {
	//validate rout
	if rout == "" {
		return nil, errUrl
	}

	//validate method
	if method == "" || !ValidateMethod(method) {
		return nil, errMethod
	}

	//						  POST
	//request body is json-encoded
	//data in json-object

	//				GET-PUT-DELETE
	//request body is url-encoded
	//data sends in url-query params

	var jsonBody []byte

	//validate reqBody on nillable
	if reqBody == nil {
		jsonBody = nil
	} else {
		v, err := json.Marshal(reqBody)
		if err != nil {
			return nil, errJsonMarshal
		}
		jsonBody = v
	}

	//creating the request
	req, err := http.NewRequestWithContext(ctx, method, rout, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errRequestForm
	}

	//authorization
	req.Header.Set("Authorization", "Bearer "+ag.token)

	//format-encoding
	req.Header.Set("Content-Type", "application/json; charset=UTF8")

	//perform request
	resp, err := ag.heart.Do(req)
	if err != nil {
		return nil, errRequestPerf
	}

	//validate response body

	if !ValidateStatusCode(resp.StatusCode) {
		return nil, errors.New("SDK error: " + strconv.Itoa(resp.StatusCode))
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

//parse response to model
func DecodeResponseToModels(resp *http.Response, model string) ([]interface{}, error) {
	if valid := ValidateModel(model); !valid {
		return nil, errModelValid
	}

	//mapping interface to model and init bank
	inp := ModelMapping(model)

	//bank for struct
	var (
		input   []interface{}
		storage []interface{}
	)

	defer resp.Body.Close()
	//response to []interface{}
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&input)
	if err != nil {
		return nil, fmt.Errorf("%s - [%s]", "error with decoding input", err.Error())
	}

	//convert []interface{} to map
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

//make rout to api by id: project/id for url-encode
func makeIdRout(id int, url string) string {
	return host + url + fmt.Sprintf(updQuery, strconv.Itoa(id))
}

//Slice contains id of projects
var BankIdProject []int

//Update bankId method
func UpdateBankIdProject(opt string, val int) error {
	if !(opt == delOpt || opt != addOpt) {
		return errOptional
	}

	switch opt {
	case delOpt:

		var (
			stB []string
			res []int
		)

		for _, v := range BankIdProject {
			stB = append(stB, strconv.Itoa(v))
		}

		st := strings.Join(stB, " ")
		rs := strings.Replace(st, strconv.Itoa(val), "", 1)
		bank := strings.Split(rs, " ")

		for _, v := range bank {
			num, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			res = append(res, num)
		}

		BankIdProject = res
	case addOpt:
		BankIdProject = append(BankIdProject, val)
	}

	return nil
}

func (ag *Agent) SiftBankIdProject(ctx context.Context) error {
	res, err := ag.GetAllProjects(ctx)
	if res != nil {
		for _, v := range *res {
			BankIdProject = append(BankIdProject, v.Id)
			return err
		}
	} else {
		return err
	}

	if err != nil {
		return errors.New("error with getting Projects - " + err.Error())
	}

	sort.Ints(BankIdProject)

	return nil
}

//validate existing id in bankId
func ValidateIdProjects(id int) bool {
	sort.Ints(BankIdProject)
	if idx := sort.SearchInts(BankIdProject, id); BankIdProject[idx] != id {
		return false
	} else {
		return true
	}
}
