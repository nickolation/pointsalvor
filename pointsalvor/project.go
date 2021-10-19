package pointsalvor

import (
	"context"
	"fmt"
	"log"
	"net/http"
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

func (ag *Agent) AddProject(ctx context.Context, name string) (*Project, error) {
	if name == "" {
		return nil, errInvalidNameModel
	}

	resp, err := ag.KnockToApi(ctx, http.MethodPost, host+projectsUrl, namedModel{Name: name})
	if err != nil {
		return nil, err
	}

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

func (ag *Agent) RenameProject(ctx context.Context, opt NamedIdOpt) error {
	rout, err := validateNamedIdOpt(opt, projectsUrl)
	if err != nil {
		return err
	}
	log.Printf("Rout: [%s]", rout)

	resp, err := ag.KnockToApi(ctx, http.MethodPost, rout, namedModel{Name: opt.Name})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (ag *Agent) DeleteProject(ctx context.Context, id int) error {
	rout, ok := makeIdRout(id, projectsUrl)
	if !ok {
		return errInvalidId
	}
	log.Printf("Rout: [%s]", rout)

	resp, err := ag.KnockToApi(ctx, http.MethodDelete, rout, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
