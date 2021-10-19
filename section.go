package pointsalvor

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const (
	sModel       = "section"
	sectionsUrl  = "/sections"
	projectQuery = "?project_id=%s"
)


//Section entity struct used by model in methods
type Section struct {
	Id         int    `json:"id"`
	Project_id int    `json:"project_id"`
	Order      byte   `json:"order"`
	Name       string `json:"name"`
}

//Query-rout used by get-method
//Create url-encode rout for knock to todoist-api
func makeQueryRout(id int) (string, bool) {
	if id != 0 {
		return host + sectionsUrl + fmt.Sprintf(projectQuery, strconv.Itoa(id)), true
	}
	return "", false
}

//Get all sections by id project
func (ag *Agent) GetAllSections(ctx context.Context, id int) (*[]Section, error) {
	var input []Section

	rout, ok := makeQueryRout(id)
	if !ok {
		return nil, errInvalidId
	}
	log.Printf("Rout: [%s]", rout)

	resp, err := ag.KnockToApi(ctx, http.MethodGet, rout, nil)
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

//Struct-opt for AddSection-method 
//Validates on value of name/id
type NewSectionOpt struct {
	Project_id int    `json:"project_id"`
	Name       string `json:"name"`
}

//Checking for valid opt struct
func validateSectionOpt(opt NewSectionOpt) error {
	if opt.Name == "" {
		return errInvalidNameModel
	}

	if opt.Project_id == 0 {
		return errInvalidId
	}

	return nil
}

//Make new section in toodist-client by opt-struct with project-id and name
func (ag *Agent) AddSection(ctx context.Context, opt NewSectionOpt) (*Section, error) {
	if err := validateSectionOpt(opt); err != nil {
		return nil, err
	}

	resp, err := ag.KnockToApi(ctx, http.MethodPost, host+sectionsUrl, opt)

	if err != nil {
		return nil, errKnockTo(err)
	}

	model, err := DecodeResponseToModel(resp, sModel)
	if err != nil {
		return nil, err
	}

	res, ok := model.(Section)
	if !ok {
		return nil, errSwitchType
	}

	return &res, nil
}

//Rename goal section by opt-struct with id-section and new name 
func (ag *Agent) RenameSection(ctx context.Context, opt NamedIdOpt) error {
	rout, err := validateNamedIdOpt(opt, sectionsUrl)
	if err != nil {
		return err
	}
	log.Printf("Rout: [%s]", rout)

	resp, err := ag.KnockToApi(ctx, http.MethodPost, rout, namedModel{
		Name: opt.Name,
	})

	if err != nil {
		return errKnockTo(err)
	}
	defer resp.Body.Close()

	return nil
}

//Delete existing section by id 
func (ag *Agent) DeleteSection(ctx context.Context, id int) error {
	rout, ok := makeIdRout(id, sectionsUrl)
	if !ok {
		return errInvalidId
	}
	log.Printf("Rout: [%s]", rout)

	resp, err := ag.KnockToApi(ctx, http.MethodDelete, rout, nil)
	if err != nil {
		return errKnockTo(err)
	}
	defer resp.Body.Close()

	return nil
}
