package obspipeline

import (
	"go.uber.org/zap"
)

func NewPipeline(dev bool) *Pipeline {
	var p Pipeline
	var l *zap.Logger

	l, _ = zap.NewProduction()
	if dev {
		l, _ = zap.NewDevelopment()
	}
	p.L = l
	p.Attachments = make(map[string]string)
	p.Posts = make([]string, 0)

	return &p
}

type Pipeline struct {
	Source, Target          string
	Attachments             map[string]string
	Notes, Posts            []string
	L                       *zap.Logger
	BlogDir, AttachmentsDir string
}

func (p *Pipeline) FindAttachments() error {

	return nil
}

func (p *Pipeline) FindNotes() error {
	return nil
}

func (p *Pipeline) FindPosts() error {
	return nil
}

func (p *Pipeline) CopyPost(post string) error {

	return nil
}

func (p *Pipeline) MoveAttachments(post string) error {

	return nil
}

func (p *Pipeline) SanitizePost(post string) error {

	return nil
}
