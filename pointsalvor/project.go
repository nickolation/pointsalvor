package pointsalvor

import (
	"context"
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
	return nil, nil
}

func (ag *Agent) RenameProject(ctx context.Context, rename string) (Project, error) {
	return Project{}, nil
}

func (ag *Agent) DeleteProject(ctx context.Context, id int) (Project, error) {
	return Project{}, nil
}
