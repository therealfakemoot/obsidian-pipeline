package main

import (
	"context"
	"flag"

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

	m := fsm.NewStateMachine(&fsm.NoteFound)
	note := "bleep"
	ctx := context.WithValue(context.Background(), "note", note)
	_, err := m.Transition(ctx, "CopyPost")
	if err != nil {
		l.Fatal("could not transition from NoteFound to CopyPost", zap.Error(err))

	}
}
