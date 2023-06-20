package obp

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CopyPosts(src, dst string) error {
	posts := make([]string, 0)

	srcRoot := os.DirFS(src)
	err := fs.WalkDir(srcRoot, ".", func(path string, d fs.DirEntry, err error) error {
		// here's where I walk through the source directory and collect all the markdown notes
		if err != nil {
			return fmt.Errorf("could not walk %q: %w", path, err)
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".md") {
			posts = append(posts, filepath.Join(src, path))
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("walkfunc failed: %w", err)
	}

	for _, post := range posts {
		base := filepath.Base(post)

		splitPostName := strings.Split(base, ".")

		postName := strings.Join(splitPostName[:len(splitPostName)-1], ".")

		log.Printf("provided dst: %q\n", dst)
		postDir := filepath.Join(dst, postName)
		log.Printf("postDir: %q\n", postDir)
		return nil

		err := os.MkdirAll(postDir, 0777)

		if err != nil && !os.IsExist(err) {
			return fmt.Errorf("error creating target directory %q: %w", dst, err)
		}
	}

	return nil
}

func Sanitize(src string) error {
	return nil
}

func GatherMedia(src string) error {
	return nil
}
