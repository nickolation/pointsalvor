package pointsalvor

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

var (
	errSwitch = fmt.Errorf("error with switch type")
)

const (
	sModel       = "section"
	sectionsUrl  = "/sections"
	sectionQuery = "?project_id=%s"
)

type Section struct {
	Id         int    `json:"id"`
	Project_id int    `json:"project_id"`
	Order      byte   `json:"order"`
	Name       string `json:"name"`
}

func makeQueryRout(id int) string {
	return host + sectionsUrl + fmt.Sprintf(sectionQuery, strconv.Itoa(id))
}

func (ag *Agent) GetAllSections(ctx context.Context, id int) (*[]Section, error) {
	var input []Section

	resp, err := ag.KnockToApi(ctx, http.MethodGet, makeQueryRout(id), nil)
	if err != nil {
		return nil, err
	}

	store, err := DecodeResponseToModels(resp, sModel)
	if err != nil {
		return nil, err
	}

	for _, v := range store {
		res, ok := v.(Section)
		if !ok {
			return nil, errConvertTo(sModel)
		}

		input = append(input, res)
	}

	return &input, nil
}

type inputSection struct {
	Project_id int    `json:"project_id"`
	Name       string `json:"name"`
}

func (ag *Agent) AddSection(ctx context.Context, name string, projId int) (*Section, error) {
	if name == "" {
		return nil, errInvalidNameModel
	}

	resp, err := ag.KnockToApi(ctx, http.MethodPost, host+sectionsUrl, inputSection{
		Project_id: projId,
		Name:       name,
	})

	if err != nil {
		return nil, errKnockTo(err)
	}

	model, err := DecodeResponseToModel(resp, sModel)
	if err != nil {
		return nil, err
	}

	res, ok := model.(Section)
	if !ok {
		return nil, errSwitch
	}

	return &res, nil
}

func (ag *Agent) RenameSection(ctx context.Context, name string) (Section, error) {
	return Section{}, nil
}

func (ag *Agent) DeleteSection(ctx context.Context, name string) (Section, error) {
	return Section{}, nil
}
