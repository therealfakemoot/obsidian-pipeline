package obp

import (
	"fmt"
	"io"
	"io/fs"
	// "log"
	"os"
	"path/filepath"
	"strings"
)

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

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

		postDir := filepath.Join(dst, postName)

		err := os.MkdirAll(postDir, 0777)

		if err != nil && !os.IsExist(err) {
			return fmt.Errorf("error creating target directory %q: %w", dst, err)
		}

		_, err = copy(post, filepath.Join(postDir, "index.md"))
		if err != nil {
			return fmt.Errorf("error opening %q for copying: %w", post, err)
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
