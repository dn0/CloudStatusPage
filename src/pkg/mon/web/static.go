package web

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

type staticFileSystem struct {
	fs http.FileSystem
}

//nolint:wrapcheck // This is a very thin wrapper.
func (fs staticFileSystem) Open(path string) (http.File, error) {
	file, err := fs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := fs.fs.Open(index); err != nil {
			closeErr := file.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return file, nil
}

func newStaticServer(rootDir, stripPrefix string, middlewares ...middlewareFun) http.Handler {
	r := chi.NewRouter()
	r.Use(middlewares...)
	r.Handle("/*", http.StripPrefix(stripPrefix, http.FileServer(staticFileSystem{http.Dir(rootDir)})))
	return r
}
