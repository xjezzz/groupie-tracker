package artists

import (
	"net/http"
	"text/template"
)

func Error(w http.ResponseWriter, code int) {
	response := struct {
		ErrorCode int
		ErrorText string
	}{
		ErrorCode: code,
		ErrorText: http.StatusText(code),
	}
	w.WriteHeader(code)
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, response)
}
