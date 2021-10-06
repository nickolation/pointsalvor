package pointsalvor

import "context"

type Section struct {
	Id         int
	Project_id int64
	Order      byte
	Name       string
}

func (ag *Agent) GetAllSections(ctx context.Context) ([]Section, error) {
	return nil, nil
}

func (ag *Agent) AddSection(ctx context.Context, name string) (Section, error) {
	return Section{}, nil
}

func (ag *Agent) RenameSection(ctx context.Context, name string) (Section, error) {
	return Section{}, nil
}

func (ag *Agent) DeleteSection(ctx context.Context, name string) (Section, error) {
	return Section{}, nil
}
