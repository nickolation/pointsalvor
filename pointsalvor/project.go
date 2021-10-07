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

func (ag *Agent) AddProject(ctx context.Context, name string) (Project, error) {
	return Project{}, nil
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
			return nil, errors.New("error with convert interface to - " + pModel)
		}

		input = append(input, res)
	}

	return input, nil
}

func (ag *Agent) RenameProject(ctx context.Context, rename string) (Project, error) {
	return Project{}, nil
}

func (ag *Agent) DeleteProject(ctx context.Context, id int) (Project, error) {
	return Project{}, nil
}
