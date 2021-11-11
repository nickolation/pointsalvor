package pointsalvor

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

//	constant value in other struct
//	overhead
const (
	pModel = "project"
)

//Project entity struct used by model in methods
type Project struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Comment_count int    `json:"comment_count"`
	Order         int    `json:"order"`
	Color         int    `json:"color"`
	Shared        bool   `json:"shared"`
	Sync_id       int    `json:"sync_id"`
	Favorite      bool   `json:"favorite"`
	Inbox_project bool   `json:"inbox_project"`
	Url           string `json:"url"`
}

//Make project in todoist-client
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
		//	unvariable custom error
		//	strange
		return nil, fmt.Errorf("error with switch type")
	}

	return &res, nil
}

//Get all project in client
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
			return nil, errConvertTo(pModel)
		}

		input = append(input, res)
	}

	return &input, nil
}

//Rename project with opt-struct by id-project and name
func (ag *Agent) RenameProject(ctx context.Context, opt NamedIdOpt) error {
	rout, err := validateNamedIdOpt(opt, projectsUrl)
	if err != nil {
		return err
	}

	//	optional log
	log.Printf("Rout: [%s]", rout)

	resp, err := ag.KnockToApi(ctx, http.MethodPost, rout, namedModel{Name: opt.Name})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

//Delete project by id
func (ag *Agent) DeleteProject(ctx context.Context, id int) error {
	rout, ok := makeIdRout(id, projectsUrl)
	if !ok {
		return errInvalidId
	}

	//	optional log
	log.Printf("Rout: [%s]", rout)

	resp, err := ag.KnockToApi(ctx, http.MethodDelete, rout, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
