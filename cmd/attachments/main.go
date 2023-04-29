package main

import (
	"flag"
	// "fmt"

	"go.uber.org/zap"

	"code.ndumas.com/ndumas/obsidian-pipeline"
)

func main() {
	var (
		source, target, attachmentsDir, blogDir string
		dev                                     bool
	)

	flag.BoolVar(&dev, "dev", false, "developer mode")
	flag.StringVar(&source, "source", "", "source directory containing your vault")
	flag.StringVar(&target, "target", "", "target directory containing your hugo site")
	flag.StringVar(&attachmentsDir, "attachments", "", "directory containing your vault's attachments")
	flag.StringVar(&blogDir, "blog", "", "vault directory containing blog posts to-be-published")

	flag.Parse()
	am := obspipeline.NewPipeline(dev)
	defer am.L.Sync()

	am.Source = source
	am.Target = target
	am.AttachmentsDir = attachmentsDir
	am.BlogDir = blogDir

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

	// fmt.Printf("%#+v\n", am.Attachments)
}
