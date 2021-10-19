package pointsalvor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
)

const (
	host        = "https://api.todoist.com/rest/v1"
	projectsUrl = "/projects"

	updQuery         = "/%s"
	requestTimeLimit = time.Second * 15
)

var (
	errUrl              = errors.New("url request is empty")
	errMethod           = errors.New("method is invalid")
	errJsonMarshal      = errors.New("eror with marshal request body")
	errRequestForm      = errors.New("request is invalid")
	errRequestPerf      = errors.New("error with perform request")
	errModelValid       = errors.New("invalid model")
	errInvalidNameModel = fmt.Errorf("invalid name of model - empty")
	errSwitchType       = fmt.Errorf("error with switch type")

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

//struct for rename-param in update-methods: renameModel(rename)
type namedModel struct {
	Name string `json:"name"`
}

//Opt-struct used by opt-param in methods
type NamedIdOpt struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

//Validation opt-struct on correct fields with creating ulr-rout fro KnockToApi method 
func validateNamedIdOpt(opt NamedIdOpt, url string) (string, error) {
	var addr int

	if opt.Name == "" {
		return "", errInvalidNameModel
	}

	if id := opt.Id; id != 0 {
		addr = id
	} else {
		return "", errInvalidId
	}

	r, _ := makeIdRout(addr, url)

	return r, nil
}

//Main object provides functional of methods working with todoist-api
type Agent struct {
	Engine *http.Client
	Token string
}

//Create new agent object used by performing sdk-methods to todoist-api
//Required token-api is located in integration-settings of toodist-client in todoist-app
//Contains http.Client struct with 15 second timeout of request-response and mock-testing functional
func NewAgent(tokenApi string) (*Agent, error) {
	if tokenApi == "" {
		return nil, errors.New("tokenApi is empty")
	}

	return &Agent{
		Engine: &http.Client{
			Timeout: requestTimeLimit,
		},
		Token: tokenApi,
	}, nil
}

//Map of string-methods and http.Methods used for validataion 
var MappingMethod = map[string]string{
	"get":  http.MethodGet,
	"post": http.MethodPost,
	"del":  http.MethodDelete,
}

//Validate method string on asserting to used http.Methods
func ValidateMethod(method string) bool {
	for _, meth := range MappingMethod {
		if meth == method {
			return true
		}
	}

	return false
}

//Make rout to api by id: project/id for url-encode
func makeIdRout(id int, url string) (string, bool) {
	if id != 0 {
		return host + url + fmt.Sprintf(updQuery, strconv.Itoa(id)), true
	}

	return "", false
}

//Map encode-type and statusCode for correct response results
var MappingStatusCode = map[string]int{
	"json":       200,
	"no-content": 204,
}

//Validate status code on correct httpStatusCode returned with http.Response
func ValidateStatusCode(code int) bool {
	for _, cd := range MappingStatusCode {
		if code == cd {
			return true
		}
	}
	return false
}

//Main method for send http-requests to toodist-api 
//Requires httpMethod, rout url, reqBody with need-type encoding (json/url-encode)
//Returns response used by Decode method and resp.StatusCode validation 
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
	req.Header.Set("Authorization", "Bearer "+ag.Token)

	//format-encoding
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json; charset=UTF8")
	}

	//perform request
	resp, err := ag.Engine.Do(req)
	if err != nil {
		return nil, errRequestPerf
	}

	return resp, nil
}

//custom codes used by keys in model-codes struct
const (
	pCode uint16 = 0 + iota
	sCode
	tCode
)

//bank of models for validation model in DecodeResponseToModel
var ModelCodes = map[uint16]string{
	pCode: "project",
	sCode: "section",
	tCode: "task",
}

//validate models on existance
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

//Muldi-models
//Decode http.Response to model struct
//Decoding schema: http.Response -> map[string]inteface{} -> model

//One-model decoding
//Decode http.Response to model struct
//Decoding schema: http.Response -> map[string]inteface{} -> model
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

//Muldi-models decoding
//Decode http.Response to model struct
//Decoding schema: http.Response -> map[string]inteface{} -> model
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
