package pointsalvor

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

const (
	pModel = "project"
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

//make rout to api by id: project/id for url-encode
func makeIdRout(id int, url string) string {
	return host + url + fmt.Sprintf(updQuery, strconv.Itoa(id))
}

func (ag *Agent) AddProject(ctx context.Context, name string) (*Project, error) {
	if name == "" {
		return nil, errInvalidNameModel
	}

	resp, err := ag.KnockToApi(ctx, http.MethodPost, host+projectsUrl, namedProgect{Name: name})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	model, err := DecodeResponseToModel(resp, pModel)
	if err != nil {
		return nil, err
	}

	res, ok := model.(Project)
	if !ok {
		return nil, fmt.Errorf("error with switch type")
	}

	return &res, nil
}

func (ag *Agent) GetAllProjects(ctx context.Context) ([]Project, error) {
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
			return nil, errConvertTo(pModel)
		}

		input = append(input, res)
	}

	//ag.SiftBankIdProject(ctx)

	return input, nil
}

func (ag *Agent) RenameProject(ctx context.Context, id int, rename string) error {
	if rename == "" {
		return errInvalidNameModel
	}

	resp, err := ag.KnockToApi(ctx, http.MethodPost, makeIdRout(id, projectsUrl), namedProgect{Name: rename})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (ag *Agent) DeleteProject(ctx context.Context, id int) error {
	resp, err := ag.KnockToApi(ctx, http.MethodDelete, makeIdRout(id, projectsUrl), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
