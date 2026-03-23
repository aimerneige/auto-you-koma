package agent

import "context"

type PipelinePayload struct {
	ProjectID     string
	Input         interface{}
	Output        interface{}
	CurrentStatus string
}

type Agent interface {
	Name() string
	Execute(ctx context.Context, input *PipelinePayload) (*PipelinePayload, error)
}
