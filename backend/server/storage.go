package server

type JobStorage interface {
	Pull(count int) ([]Job, error)
	Push(job Job) error
	Complete(job Job) error
	Cancel(job Job) error
	//GetResults(requestId string) ([]*JobResultEntity, error)
}
