package swagger

import (
	"github.com/shurcooL/vfsgen"
	"log"
	"net/http"
)

func main() {

	var fs http.FileSystem = http.Dir("./swagger-ui")

	err := vfsgen.Generate(fs, vfsgen.Options{
		PackageName:  "swagger",
		Filename:     "./datafile.go",
		VariableName: "Swagger",
	})

	if err != nil {
		log.Fatalln(err)
	}
}
