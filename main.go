package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	// "path/filepath"
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
	Source, Target string
	Attachments    map[string]bool
	Posts          []string
	L              *zap.Logger
}

func (am *AttachmentMover) Walk() error {
	return nil
}

func (am *AttachmentMover) walk(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if d.IsDir() {
		return nil
	}
	walkLogger := am.L.Named("Walk()")

	walkLogger.Info("scanning for relevance", zap.String("path", path))
	if strings.HasSuffix(path, "index.md") {
		walkLogger.Info("found index.md, adding to index", zap.String("path", path))
		am.Posts = append(am.Posts, path)
	}

	if !strings.HasSuffix(path, ".md") && strings.Contains(path, "attachments") {
		walkLogger.Info("found  attachment file, adding to index", zap.String("path", path))
		am.Attachments[path] = true
	}
	return nil
}

func (am *AttachmentMover) Move() error {
	moveLogger := am.L.Named("Move()")
	moveLogger.Info("scanning posts", zap.Strings("posts", am.Posts))
	for _, post := range am.Posts {
		// log.Printf("scanning %q for attachment links", post)
		linkedAttachments, err := extractAttachments(post)
		if err != nil {
			return fmt.Errorf("could not extract attachment links from %q: %w", post, err)
		}
		for _, attachment := range linkedAttachments {
			moveAttachment(post, attachment)
		}
	}

	return nil
}

func moveAttachment(post, attachment string) error {

	return nil
}

func extractAttachments(fn string) ([]string, error) {
	attachments := make([]string, 0)

	return attachments, nil
}

func main() {
	am := NewAttachmentMover()

	flag.StringVar(&am.Source, "source", "", "source directory containing your vault")
	flag.StringVar(&am.Target, "target", "", "target directory containing your hugo site")

	flag.Parse()

	if am.Source == "" || am.Target == "" {
		am.L.Fatal("flags not provided")
	}

	err := am.Walk()
	if err != nil {
		// log.Fatalf("error walking blog dir to gather file names: %s\n", err)
	}

	err = am.Move()
	if err != nil {
		// log.Fatalf("error walking blog dir to gather file names: %s\n", err)
	}
}
