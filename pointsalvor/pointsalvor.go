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
	projectsUrl = "/projects"

	updQuery  = "/%s"
	taskClose = updQuery + "/close"

	requestTimeLimit = time.Second * 15

	delOpt = "del"
	addOpt = "add"
)

var (
	errUrl              = errors.New("url request is empty")
	errMethod           = errors.New("method is invalid")
	errJsonMarshal      = errors.New("eror with marshal request body")
	errRequestForm      = errors.New("request is invalid")
	errRequestPerf      = errors.New("error with perform request")
	errModelValid       = errors.New("invalid model")
	errOptional         = errors.New("invalid optional updating the bankId")
	errInvalidNameModel = fmt.Errorf("invalid name of model - empty")

	errDecodeIf = func(e error) error {
		return fmt.Errorf("error with decoding interface - [%v]", e)
	}

	errDecodeMap = func(e error) error {
		return fmt.Errorf("error with decoding map - [%v]", e)
	}

	errConvertTo = func(s string) error {
		return fmt.Errorf("invalid convertation to interface - [%s]", s)
	}

	errKnockTo = func(e error) error {
		return fmt.Errorf("invalid knock to api - [%v]", e)
	}
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

var MappingMethod = map[string]string{
	"get":  http.MethodGet,
	"post": http.MethodPost,
	"del":  http.MethodDelete,
}

func ValidateMethod(method string) bool {
	for _, meth := range MappingMethod {
		if meth == method {
			return true
		}
	}

	return false
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
	if !ValidateMethod(method) {
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
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json; charset=UTF8")
	}

	//perform request
	resp, err := ag.heart.Do(req)
	if err != nil {
		return nil, errRequestPerf
	}

	return resp, nil
}

const (
	pCode uint16 = 0 + iota
	sCode
	tCode
)

//bank of models
var ModelCodes = map[uint16]string{
	pCode: "project",
	sCode: "section",
	tCode: "task",
}

//validate models on exist
func ValidateModel(model string) bool {
	for _, mod := range ModelCodes {
		if mod == model {
			return true
		}
	}

	return false
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

func DecodeResponseToModel(resp *http.Response, model string) (interface{}, error) {
	if valid := ValidateModel(model); !valid {
		return nil, errModelValid
	}

	//mapping interface to model and init bank
	inp := ModelMapping(model)

	//bank for struct
	var input interface{}

	defer resp.Body.Close()
	//response to []interface{}
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&input)
	if err != nil {
		return nil, errDecodeIf(err)
	}

	//convert []interface{} to map
	//decode map to struct
	//add storage of struct
	res, ok := input.(map[string]interface{})
	if ok {
		err = mapstructure.Decode(res, &inp)
		if err != nil {
			return nil, errDecodeMap(err)
		}
	}

	return inp, nil
}

//parse response to models
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
		return nil, errDecodeIf(err)
	}

	//convert []interface{} to map
	//decode map to struct
	//add storage of struct
	for _, part := range input {
		res, ok := part.(map[string]interface{})
		if ok {
			err = mapstructure.Decode(res, &inp)
			if err != nil {
				return nil, errDecodeMap(err)
			}
			storage = append(storage, inp)
		}
	}

	return storage, nil
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

/*
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
*/
