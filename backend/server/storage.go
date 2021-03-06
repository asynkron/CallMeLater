package server

type JobStorage interface {
	Pull(count int) ([]Job, error)
	Create(job Job) error
	RescheduleCron(job Job) error
	Complete(job Job) error
	Cancel(job Job) error
	Retry(job Job) error
	Fail(job Job) error
	//GetResults(requestId string) ([]*JobResultEntity, error)

	Read(skip int, limit int, search string) ([]JobEntity, error)
}
