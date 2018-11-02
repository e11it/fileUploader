package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/caarlos0/env"
)
type config struct {

	Port         int        `env:"PORT" envDefault:"8080"`
	Folder	  string		`env:"FOLDER" envExpand:"true"`
}

func (c *config) getListenAddress() string {
	return fmt.Sprintf(":%d", c.Port)
}

func (c * config) joinPath(subPath string) string {
	return path.Join(c.Folder,subPath)
}

func main() {
	var err error
	cfg := config{}
	err = env.Parse(&cfg)
	if len(cfg.Folder) == 0 {
		cfg.Folder, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}

	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	fmt.Printf("%+v\n", cfg)

	http.HandleFunc("/", uploadHandler(&cfg))
	log.Printf("Listen on %s\n", cfg.getListenAddress())
	err = http.ListenAndServe(cfg.getListenAddress(),nil)
	if err != nil {
		log.Println(err)
	}


}

func writeResponse(w http.ResponseWriter, statusCode int, msg string ) {
	w.WriteHeader(statusCode)
	w.Write([]byte(msg))
}

func uploadHandler(cfg *config) func(w http.ResponseWriter,r *http.Request) {

	return func(w http.ResponseWriter,r *http.Request) {
		var (
			subpath string
			fullpath string
			file string
		)

		subpath = path.Clean(r.URL.Path)
		_, file = path.Split(subpath)
		if len(file) == 0 {
			writeResponse(w, 406, "No filename provided")
			return
		}

		fullpath = cfg.joinPath(subpath)

		if _, err := os.Stat(path.Dir(fullpath)); os.IsNotExist(err) {
			writeResponse(w, 400, "Directory doesn't exists")
			return
		}

		f, err := os.OpenFile(fullpath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Printf("[Err] Error open file for wtite: %s\n", err)
			writeResponse(w, 400, "Cant open file")
			return
		}
		defer f.Close()

		log.Println("Upload file: " + fullpath)
		_, err = io.Copy(f, r.Body)
		if err != nil {
			writeResponse(w, 400, "Cant write file content")
			return
		}


		w.WriteHeader(200)
	}
}
