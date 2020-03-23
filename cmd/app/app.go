package app

import (
	"errors"
	"github.com/shohrukh56/bookFileService/pkg/core/file"
	"github.com/shohrukh56/jwt/pkg/jwt"
	"github.com/shohrukh56/mux/pkg/mux"
	"github.com/shohrukh56/rest/pkg/rest"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	ext = make(map[string]string)

	content = `
This is bookfileservice
`
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
	contentTypeHtml  = "text/html"
	contentTypeText  = "text/plain"
	contentTypePng   = "image/png"
	contentTypeJpg   = "image/jpeg"
	contentTypePdf   = "application/pdf"
)

type Server struct {
	router  *mux.ExactMux
	fileSvc *file.Service
	secret 	jwt.Secret
}

func NewServer(router *mux.ExactMux, fileSvc *file.Service, secret jwt.Secret) *Server {
	return &Server{router: router, fileSvc: fileSvc, secret: secret}
}


func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *Server) Stop() {

}

func (s *Server) Start() {
	ext[".txt"] = contentTypeText
	ext[".pdf"] = contentTypePdf
	ext[".png"] = contentTypePng
	ext[".jpeg"] = contentTypeJpg
	ext[".jpg"] = contentTypeJpg
	ext[".html"] = contentTypeHtml
	s.InitRoutes()
}

func (s *Server) handleIndex() http.HandlerFunc {
	tpl, err := template.ParseFiles("index.gohtml")
	if err != nil {
		panic(err)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		err := tpl.Execute(writer,
			struct {
				Title   string
				Content string
			}{
				Title:   "book file service",
				Content: content,
			})
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
		}
	}
}

func (s *Server) handleSaveFiles() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		err := request.ParseMultipartForm(defaultMaxMemory)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		files := request.MultipartForm.File["file"]
		type FileURL struct {
			Name string `json:"name"`
		}
		fileURLs := make([]FileURL, 0, len(files))

		for _, file := range files {
			contentType, ok := ext[filepath.Ext(file.Filename)]
			if !ok {
				http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			openFile, err := file.Open()
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			newFile, err := s.fileSvc.Save(openFile, contentType)
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			fileURLs = append(fileURLs, FileURL{
				newFile[:len(newFile)-len(filepath.Ext(newFile))],
			})
		}
		err = rest.WriteJSONBody(writer, fileURLs)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) handleGetFile() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		dir, err := ioutil.ReadDir(s.fileSvc.Filepath)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		path, ok := mux.FromContext(request.Context(), "id")
		log.Print(path)
		if ok {
			for _, info := range dir {
				if !info.IsDir() {
					fileName := info.Name()
					fileName = fileName[:len(fileName)-len(filepath.Ext(fileName))]
					if !strings.EqualFold(fileName, path){
						continue
					}

					body, err := ioutil.ReadFile("files/"+info.Name())
					if err != nil {
						http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
						return
					}
					writer.Header().Set("Content-Type", ext[filepath.Ext(info.Name())])
					_, err = writer.Write(body)
					if err != nil {
						log.Println(errors.New("error"))
					}
					return
				}
			}
		}


		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}
