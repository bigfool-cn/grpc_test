package swagger

import (
	"github.com/elazarl/go-bindata-assetfs"
	"log"
	"net/http"
	"path"
	"strings"
)

func ServeSwaggerFile(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, "swagger.json") {
		log.Printf("Not Found: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	p := strings.TrimPrefix(r.URL.Path, "/swagger/")
	p = path.Join("../proto/", p)

	log.Printf("Server swagger-file: %s", p)

	http.ServeFile(w, r, p)
}

func ServeSwaggerUI(mux *http.ServeMux) {
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
		Prefix:   "swagger/swagger-ui",
	})

	prefix := "/swagger-ui/"

	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}
