package main

import (
	"context"
	"errors"
	"flag"
	"io/fs"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"

	"code.ndumas.com/ndumas/obsidian-pipeline/fsm"
)

func main() {
	var (
		source, target, attachmentsDir, blogDir string
		dev                                     bool
	)
	l, _ := zap.NewProduction()

	flag.BoolVar(&dev, "dev", false, "developer mode")
	flag.StringVar(&source, "source", "", "source directory containing your vault")
	flag.StringVar(&target, "target", "", "target directory containing your hugo site")
	flag.StringVar(&attachmentsDir, "attachments", "", "directory containing your vault's attachments")
	flag.StringVar(&blogDir, "blog", "", "vault directory containing blog posts to-be-published")

	flag.Parse()

	fileNames := make(chan string)

	walkFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}
		fileNames <- path

		return nil
	}

	var wg sync.WaitGroup
	go func() {
		defer wg.Done()
		wg.Add(1)
		for fn := range fileNames {
			if !strings.HasSuffix(fn, ".md") {
				continue
			}
			m := fsm.NewStateMachine(&fsm.NoteFound)
			ctx := context.WithValue(context.Background(), "note", fn)
			_, err := m.Transition(ctx, "CopyPost")
			if err != nil {
				l.Fatal("could not transition from NoteFound to CopyPost", zap.Error(err))
			}
		}
	}()

	root := os.DirFS(source)
	err := fs.WalkDir(root, ".", walkFunc)
	if err != nil {
		l.Fatal("error walking for files", zap.Error(err))
	}
	close(fileNames)

	wg.Wait()
}
