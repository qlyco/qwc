package core

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/qlyco/qwc/internal/res"
)

func walkFiles(directory string, extension string, visit func(name string, source []byte) error) error {
	return walkPaths(directory, func(relativePath string, filePath string) error {
		if path.Ext(relativePath) != extension {
			return nil
		}

		source, err := os.ReadFile(filePath)

		if err != nil {
			return err
		}

		return visit(trimExt(relativePath), source)
	})
}

func walkPaths(directory string, visit func(relativePath string, filePath string) error) error {
	err := fs.WalkDir(os.DirFS(directory), res.WalkRoot, func(relativePath string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}

		filePath := filepath.Join(directory, filepath.FromSlash(relativePath))
		err = visit(relativePath, filePath)

		if err != nil {
			return fmt.Errorf(res.FileErrorFormat, filePath, err)
		}

		return nil
	})

	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}

	return err
}

func writeAtomic(filePath string, source io.Reader, modTime time.Time) error {
	directory := filepath.Dir(filePath)
	err := os.MkdirAll(directory, os.ModePerm)

	if err != nil {
		return err
	}

	temp, err := os.CreateTemp(directory, res.TempFilePattern)

	if err != nil {
		return err
	}

	defer os.Remove(temp.Name())

	_, err = io.Copy(temp, source)
	err = errors.Join(err, temp.Close(), os.Chtimes(temp.Name(), modTime, modTime))

	if err != nil {
		return err
	}

	return os.Rename(temp.Name(), filePath)
}

func copyFile(targetFile string, sourceFile string, modTime time.Time) error {
	source, err := os.Open(sourceFile)

	if err != nil {
		return err
	}

	defer source.Close()

	return writeAtomic(targetFile, source, modTime)
}

func isFile(filePath string) bool {
	info, err := os.Stat(filePath)

	return err == nil && info.Mode().IsRegular()
}
