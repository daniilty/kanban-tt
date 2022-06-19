package server

import "net/http"

func getSwaggerHandler(path string) http.HandlerFunc {
	return getFileServerHandler(path, "api")
}

func getFileServerHandler(path string, dir string) http.HandlerFunc {
	fileServer := http.FileServer(http.Dir(dir))

	return http.StripPrefix(path, fileServer).ServeHTTP
}
