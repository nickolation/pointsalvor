package pointsalvor

import "context"

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

type TaskData struct {
	Project_id, Section_id int
	Content, Due_string    string
}

func (ag *Agent) AddTaskToProject(ctx context.Context, input TaskData) (Task, error) {
	return Task{}, nil
}

func (ag *Agent) GetAllTasksBySection(ctx context.Context, sectionId int) ([]Task, error) {
	return nil, nil
}

func (ag *Agent) GetAllTasksByProject(ctx context.Context, projectId int) ([]Task, error) {
	return nil, nil
}

func (ag *Agent) CloseTaskById(ctx context.Context, id int) (Task, error) {
	return Task{}, nil
}
