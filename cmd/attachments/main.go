package main

import (
	"flag"

	"go.uber.org/zap"

	"code.ndumas.com/ndumas/obsidian-pipeline"
)

func main() {
	am := obspipeline.NewAttachmentMover()
	defer am.L.Sync()

	flag.StringVar(&am.Source, "source", "", "source directory containing your vault")
	flag.StringVar(&am.Target, "target", "", "target directory containing your hugo site")
	flag.StringVar(&am.AttachmentsDir, "attachments", "", "directory containing your vault's attachments")
	flag.StringVar(&am.BlogDir, "blog", "", "vault directory containing blog posts to-be-published")

	flag.Parse()

	switch {
	case am.Source == "":
		am.L.Fatal("please provide -source")
		fallthrough
	case am.Target == "":
		am.L.Fatal("please provide -target")
		fallthrough
	case am.AttachmentsDir == "":
		am.L.Fatal("please provide -attachments")
		fallthrough
	case am.BlogDir == "":
		am.L.Fatal("please provide -blog")
	}

	err := am.Walk()
	if err != nil {
		am.L.Fatal("error walking blog or notes dir to gather file names", zap.Error(err))
	}

	err = am.Move()
	if err != nil {
		am.L.Fatal("error moving notes", zap.Error(err))
	}
}
