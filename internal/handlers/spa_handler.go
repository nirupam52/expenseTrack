package handlers

import (
	"io/fs"
	"net/http"
	"os"
)

type SPAHandler struct {
	fileServer http.Handler
	staticFS   http.FileSystem
}

func NewSPAHandler(fsys fs.FS) (*SPAHandler, error) {
	sub, err := fs.Sub(fsys, "build")
	if err != nil {
		return nil, err
	}
	hfs := http.FS(sub)
	return &SPAHandler{
		fileServer: http.FileServer(hfs),
		staticFS:   hfs,
	}, nil
}

func (h *SPAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, err := h.staticFS.Open(r.URL.Path)
	if err == nil {
		f.Close()
		h.fileServer.ServeHTTP(w, r)
		return
	}
	if !os.IsNotExist(err) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	// SPA fallback: serve index.html for client-side routing
	r2 := r.Clone(r.Context())
	r2.URL.Path = "/"
	h.fileServer.ServeHTTP(w, r2)
}
