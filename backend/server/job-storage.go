package server

type JobStorage interface {
	Pull(count int) ([]*RequestData, error)
	Push(data *RequestData) error
	Complete(requestId string) error
}
