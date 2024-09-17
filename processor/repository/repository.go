package repository

type Repository interface {
	PutResult(taskId string, result string) error
}
