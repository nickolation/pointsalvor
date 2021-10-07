package pointsalvor

import (
	"context"
	"errors"
	"net/http"
)

const (
	pModel = "project"
	sModel = "section"
	tModel = "task"
)

var (
	errInvalidName = errors.New("invalid name of model")
)

type Project struct {
	Id            int
	Name          string
	Comment_count int
	Order         int
	Color         int
	Shared        bool
	Sync_id       int
	Favorite      bool
	Inbox_project bool
	Url           string
}

type namedProgect struct {
	Name string `json:"name"`
}

func (ag *Agent) AddProject(ctx context.Context, name string) (*Project, error) {
	if name == "" {
		return nil, errInvalidName
	}

	resp, err := ag.KnockToApi(ctx, http.MethodPost, host+projectsUrl, namedProgect{Name: name})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	store, err := DecodeResponseToModels(resp, pModel)
	if err != nil {
		return nil, err
	}

	res, ok := store[0].(Project)
	if !ok {
		return nil, errors.New("error with convert interface to - " + pModel)
	}

	err = UpdateBankIdProject(addOpt, res.Id)
	if err != nil {
		return nil, errors.New("errors with updating bankId - " + err.Error())
	}

	return &res, nil
}

func (ag *Agent) GetAllProjects(ctx context.Context) (*[]Project, error) {
	var input []Project

	resp, err := ag.KnockToApi(ctx, http.MethodGet, host+projectsUrl, nil)
	if err != nil {
		return nil, err
	}

	store, err := DecodeResponseToModels(resp, pModel)
	if err != nil {
		return nil, err
	}

	for _, v := range store {
		res, ok := v.(Project)
		if !ok {
			return nil, errors.New("error with convert interface to - " + pModel)
		}

		input = append(input, res)
	}

	//ag.SiftBankIdProject(ctx)

	return &input, nil
}

func (ag *Agent) RenameProject(ctx context.Context, id int, rename string) error {
	if rename == "" {
		return errInvalidName
	}

	resp, err := ag.KnockToApi(ctx, http.MethodPost, makeIdRout(id, projectsUrl), namedProgect{Name: rename})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (ag *Agent) DeleteProject(ctx context.Context, id int) error {

	return nil
}
