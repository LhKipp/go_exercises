package main

import (
	"fmt"
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/dog":        "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc": "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml :=

`- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution`
	yamlHandler, err := YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redirectPath, has_val := pathsToUrls[r.URL.Path]
		if has_val {
			http.Redirect(w, r, redirectPath, http.StatusMovedPermanently)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

type PathToUrl struct {
	Path string `yaml:"path" json:"path"`
	Url  string `yaml:"url" json:"url"`
}

func parseYAML(yml []byte) (pathToUrls []PathToUrl, err error) {
	err = yaml.Unmarshal([]byte(yml), &pathToUrls)
	fmt.Println(pathToUrls)
	return pathToUrls, err
}

func pathsToUrlsAsMap(pathToUrls []PathToUrl) map[string]string {
	result := make(map[string]string)
	for _, pathToUrl := range pathToUrls {
		result[pathToUrl.Path] = pathToUrl.Url
	}
	return result
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathToUrls, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	mappings := pathsToUrlsAsMap(pathToUrls)

	return func(w http.ResponseWriter, r *http.Request) {
		redirectPath, has_val := mappings[r.URL.Path]
		if has_val {
			http.Redirect(w, r, redirectPath, http.StatusMovedPermanently)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}, nil
}
