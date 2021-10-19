package pointsalvor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Task struct {
	Id            int
	Project_id    int
	Section_id    int
	Parent_id     int
	Content       string
	Description   string
	Comment_count byte
	Assignee      int
	Assigner      int
	Order         byte
	Priority      byte
	Url           string
}

const (
	tModel       = "task"
	sectionQuery = "?section_id=%s"

	tasksUrl  = "/tasks"
	taskClose = updQuery + "/close"

	sectDir = "section"
	projDir = "project"
)

var (
	errInvalidDirect = errors.New("invalid direct")
	errInvalidOpt    = errors.New("invalid option - request body")
	errInvalidId     = errors.New("invalid id - empty value")
)

//Struct used by params-request body
//Being enttite of AddTask method
type NewTaskOpt struct {
	Section_id int    `json:"section_id"`
	Content    string `json:"content"`
	Due_string string `json:"due_string"`
}

//Validate option-request body to opt{} and bool
//Make request-body struct to send by params opt
//Bool params is response to succes validation
func ValidateTaskOpt(opt NewTaskOpt) (interface{}, bool) {
	var (
		text, due string
	)

	if t := opt.Content; t != "" {
		text = t
	} else {
		return nil, false
	}

	var don bool

	if d := opt.Due_string; d != "" {
		due = d
		don = true
	}

	if id := opt.Section_id; id != 0 {
		if don {
			return struct {
				Content    string `json:"content"`
				Section_id int    `json:"section_id"`
				Due_string string `json:"due_string"`
			}{
				Content:    text,
				Section_id: id,
				Due_string: due,
			}, true
		} else {
			return struct {
				Content    string `json:"content"`
				Section_id int    `json:"section_id"`
			}{
				Content:    opt.Content,
				Section_id: id,
			}, true
		}
	}

	return nil, false
}

//Create task to todoist-client by opt-struct
func (ag *Agent) AddTask(ctx context.Context, opt NewTaskOpt) (*Task, error) {
	reqBody, ok := ValidateTaskOpt(opt)
	if !ok {
		return nil, errInvalidOpt
	}

	resp, err := ag.KnockToApi(ctx, http.MethodPost, host+tasksUrl, reqBody)
	if err != nil {
		return nil, errKnockTo(err)
	}

	model, err := DecodeResponseToModel(resp, tModel)
	if err != nil {
		return nil, err
	}

	res, ok := model.(Task)
	if !ok {
		return nil, errSwitchType
	}

	return &res, nil

}

//Option-struct for get-method
//Used by validate-method to create the rout
type GetTaskOpt struct {
	Id     int    `json:"id"`
	Direct string `json:"direct"`
}

//Make rout to request using direct in todoist-API enpdoints
//Bool-params is response to succes validataion
func makeUrlDirectly(opt GetTaskOpt) (string, error) {
	var (
		dir  string
		addr int
	)

	if d := opt.Direct; d != "" {
		dir = d
	} else {
		return "", errInvalidDirect
	}

	if id := opt.Id; id != 0 {
		addr = id
	} else {
		return "", errInvalidId
	}

	if dir == projDir {
		return host + tasksUrl + fmt.Sprintf(projectQuery, strconv.Itoa(addr)), nil
	}
	if dir == sectDir {
		return host + tasksUrl + fmt.Sprintf(sectionQuery, strconv.Itoa(addr)), nil
	}

	return "", nil
}

//Get all tasks in todoist-client by params
func (ag *Agent) GetAllTasks(ctx context.Context, opt GetTaskOpt) (*[]Task, error) {
	var input []Task

	rout, err := makeUrlDirectly(opt)
	if err != nil {
		return nil, errInvalidDirect
	}
	log.Printf("Rout: [%s]", rout)

	resp, err := ag.KnockToApi(ctx, http.MethodGet, rout, nil)
	if err != nil {
		return nil, errKnockTo(err)
	}

	store, err := DecodeResponseToModels(resp, tModel)
	if err != nil {
		return nil, err
	}

	for _, v := range store {
		res, ok := v.(Task)
		if !ok {
			return nil, errConvertTo(tModel)
		}

		input = append(input, res)
	}

	return &input, nil
}

//Make url for delete-method rout
//Validate on id value
func makeUrlDelete(id int) (string, bool) {
	if id != 0 {
		return host + tasksUrl + fmt.Sprintf(taskClose, strconv.Itoa(id)), true
	}

	return "", false
}

//Make task done
func (ag *Agent) CloseTask(ctx context.Context, id int) error {
	rout, ok := makeUrlDelete(id)
	if !ok {
		return errInvalidId
	}

	resp, err := ag.KnockToApi(ctx, http.MethodPost, rout, nil)
	if err != nil {
		return errKnockTo(err)
	}
	defer resp.Body.Close()

	return nil
}
