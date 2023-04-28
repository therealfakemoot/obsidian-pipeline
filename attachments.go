package obspipeline

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

func NewAttachmentMover() *AttachmentMover {
	var am AttachmentMover
	l, _ := zap.NewProduction()
	am.L = l
	am.Attachments = make(map[string]bool)
	am.Posts = make([]string, 0)

	return &am
}

type AttachmentMover struct {
	Source, Target          string
	Attachments             map[string]bool
	Notes, Posts            []string
	L                       *zap.Logger
	BlogDir, AttachmentsDir string
}

func (am *AttachmentMover) Walk() error {
	notesRoot := os.DirFS(am.Source)
	blogRoot := os.DirFS(am.Target)

	err := fs.WalkDir(notesRoot, ".", am.findAttachments)
	if err != nil {
		return fmt.Errorf("error scanning for attachments: %q", err)
	}

	err = fs.WalkDir(notesRoot, ".", am.findNotes)
	if err != nil {
		return fmt.Errorf("error scanning vault for posts: %q", err)
	}

	err = fs.WalkDir(blogRoot, ".", am.findPosts)
	if err != nil {
		return fmt.Errorf("error scanning blog for posts: %q", err)
	}
	return nil
}

func (am *AttachmentMover) findNotes(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if d.IsDir() {
		return nil
	}
	walkLogger := am.L.Named("FindNotes").With(zap.String("path", path))

	if strings.HasSuffix(path, ".md") && strings.Contains(path, am.BlogDir) {
		walkLogger.Info("found blog post to publish, adding to index")
		am.Notes = append(am.Notes, path)
	}
	return nil
}

func (am *AttachmentMover) findAttachments(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if d.IsDir() {
		return nil
	}
	walkLogger := am.L.Named("FindAttachments").With(zap.String("path", path))

	if strings.Contains(path, am.AttachmentsDir) {
		walkLogger.Info("found attachment file, adding to index")
		am.Attachments[path] = true
	}
	return nil
}

func (am *AttachmentMover) findPosts(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if d.IsDir() {
		return nil
	}
	walkLogger := am.L.Named("FindPosts").With(zap.String("path", path))

	if strings.HasSuffix(path, "index.md") {
		walkLogger.Info("found index.md, adding to index")
		am.Posts = append(am.Posts, path)
	}
	return nil
}

func (am *AttachmentMover) Move() error {
	moveLogger := am.L.Named("Move")
	moveLogger.Info("scanning posts", zap.Strings("posts", am.Posts))
	for _, post := range am.Notes {
		// log.Printf("scanning %q for attachment links", post)
		linkedAttachments, err := extractAttachments(filepath.Join(am.Source, post), am.L.Named("extractAttachments"))
		if err != nil {
			return fmt.Errorf("could not extract attachment links from %q: %w", post, err)
		}
		for _, attachment := range linkedAttachments {
			moveAttachment(post, attachment, am.L.Named("moveAttachment"))
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

func extractAttachments(post string, l *zap.Logger) ([]string, error) {

	l.Info("scanning note",
		zap.String("post", post),
	)

	pat := regexp.MustCompile(`\[\[Resources\/attachments\/(.*)?\]\]`)

	attachments := make([]string, 0)
	postBody, err := ioutil.ReadFile(post)
	if err != nil {
		return attachments, fmt.Errorf("error opening post to scan for attachment links: %q", err)
	}

	for _, att := range pat.FindAllSubmatch(postBody, -1) {
		l.Info("found attachment", zap.String("filename", string(att[1])))

	}

	return attachments, nil
}
