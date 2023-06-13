package obp

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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

func (p *Pipeline) Walk() error {
	notesRoot := os.DirFS(p.Source)
	blogRoot := os.DirFS(p.Target)

	err := fs.WalkDir(notesRoot, ".", p.findAttachments)
	if err != nil {
		return fmt.Errorf("error scanning for attachments: %w", err)
	}

	err = fs.WalkDir(notesRoot, ".", p.findNotes)
	if err != nil {
		return fmt.Errorf("error scanning vault for posts: %w", err)
	}

	err = fs.WalkDir(blogRoot, ".", p.findPosts)
	if err != nil {
		return fmt.Errorf("error scanning blog for posts: %w", err)
	}

	return nil
}

func (p *Pipeline) findNotes(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if d.IsDir() {
		return nil
	}

	walkLogger := p.L.Named("FindNotes").With(zap.String("path", path))

	if strings.HasSuffix(path, ".md") && strings.Contains(path, p.BlogDir) {
		walkLogger.Info("found blog post to publish, adding to index")

		p.Notes = append(p.Notes, path)
	}

	return nil
}

func (p *Pipeline) findAttachments(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if d.IsDir() {
		return nil
	}

	walkLogger := p.L.Named("FindAttachments").With(zap.String("path", path))

	if strings.Contains(path, p.AttachmentsDir) {
		walkLogger.Info("found attachment file, adding to index")

		absPath, err := filepath.Abs(filepath.Join(p.Source, path))
		if err != nil {
			return fmt.Errorf("error generating absolute path for attachment %q: %w", path, err)
		}

		walkLogger.Info("adding Attachment",
			zap.String("key", filepath.Base(absPath)),
			zap.String("value", absPath),
		)

		p.Attachments[filepath.Base(absPath)] = absPath
	}

	return nil
}

func (p *Pipeline) findPosts(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if d.IsDir() {
		return nil
	}

	walkLogger := p.L.Named("FindPosts").With(zap.String("path", path))

	if strings.HasSuffix(path, "index.md") {
		walkLogger.Info("found index.md, adding to index")

		p.Posts = append(p.Posts, path)
	}

	return nil
}

func (p *Pipeline) Move() error {
	moveLogger := p.L.Named("Move")
	moveLogger.Info("scanning posts", zap.Strings("posts", p.Posts))

	for _, post := range p.Notes {
		// log.Printf("scanning %q for attachment links", post)
		linkedAttachments, err := extractAttachments(filepath.Join(p.Source, post))
		if err != nil {
			return fmt.Errorf("could not extract attachment links from %q: %w", post, err)
		}

		for _, attachment := range linkedAttachments {
			att, ok := p.Attachments[attachment]
			if !ok {
				return fmt.Errorf("Attachment is linked by post %q but doesn't exist in attachments directory %q", post, p.AttachmentsDir)
			}

			err := moveAttachment(post, att, p.L.Named("moveAttachment"))
			if err != nil {
				return fmt.Errorf("error moving attachments: %w", err)
			}
		}
	}

	return nil
}

func moveAttachment(post, attachment string, l *zap.Logger) error {
	l.Info("moving attachment",
		zap.String("post", post),
		zap.String("attachment", attachment),
	)

	return nil
}

func extractAttachments(post string) ([]string, error) {
	pat := regexp.MustCompile(`\[\[Resources\/attachments\/(.*)?\]\]`)

	attachments := make([]string, 0)

	postBody, err := ioutil.ReadFile(post)
	if err != nil {
		return attachments, fmt.Errorf("error opening post to scan for attachment links: %w", err)
	}

	for _, att := range pat.FindAllSubmatch(postBody, -1) {
		filename := string(att[1])
		attachments = append(attachments, filename)
	}

	return attachments, nil
}

func (p *Pipeline) FindAttachments() error {

	return nil
}

func (p *Pipeline) MoveAttachments(post string) error {

	return nil
}

func (p *Pipeline) FindPosts() error {
	return nil
}

func (p *Pipeline) SanitizePost(post string) error {

	return nil
}

func (p *Pipeline) CopyPost(post string) error {

	return nil
}
