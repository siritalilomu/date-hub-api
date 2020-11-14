package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

// Server . . .
type Server struct {
	Router     *mux.Router
	Middleware negroni.Handler
}

// Route . . .
type Route struct {
	Path    string
	Handler http.Handler
	Methods []string
	Recover bool
}

type serverError struct {
	Message string
	Status  int
}

// NewRoute . . .
func NewRoute(path string, handler func(http.ResponseWriter, *http.Request), methods ...string) Route {
	return Route{Path: path, Handler: http.HandlerFunc(handler), Methods: methods}
}

// AddRoutes . . .
func (s *Server) AddRoutes(routes []Route) {
	for _, r := range routes {
		h := s.Router.Handle(r.Path, r.Handler)
		if len(r.Methods) > 0 {
			h.Methods(r.Methods...)
		}
	}
}

// NewServer . . .
func NewServer() *Server {
	s := Server{Router: mux.NewRouter()}
	s.Router.Use(recoveryMiddleware)
	return &s
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				switch t := err.(type) {
				case string:
					http.Error(w, t, http.StatusInternalServerError)
				case error:
					http.Error(w, t.Error(), http.StatusInternalServerError)
				case serverError:
					http.Error(w, t.Message, t.Status)
				default:
					http.Error(w, "unknown error", http.StatusInternalServerError)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// GetStringVar . . .
func GetStringVar(req *http.Request, varName string) (param string) {
	if param = mux.Vars(req)[varName]; param == "" {
		PanicWithStatus(fmt.Errorf("missing url variable: %s", varName), http.StatusBadRequest)
	}
	return param
}

// GetIntVar . . .
func GetIntVar(req *http.Request, varName string) (param int) {
	var err error
	if param, err = strconv.Atoi(mux.Vars(req)[varName]); err != nil {
		PanicWithStatus(fmt.Errorf("missing or invalid url variable: %s\nerror: %s", varName, err.Error()), http.StatusBadRequest)
	}
	return param
}

// GetStringParam . . .
func GetStringParam(req *http.Request, paramName string, optional ...bool) (param string) {
	if len(optional) > 1 {
		panic("too many parameters being passed to server.GetStringParam func")
	}
	optional = append(optional, false)
	isOptional := optional[0]

	var paramStrings []string
	var ok bool
	if paramStrings, ok = req.URL.Query()[paramName]; !isOptional && (!ok || len(paramStrings) != 1) {
		PanicWithStatus(fmt.Errorf("missing or invalid url parameter: %s", paramName), http.StatusBadRequest)
	}

	if len(paramStrings) == 1 {
		if param = paramStrings[0]; !isOptional && param == "" {
			PanicWithStatus(fmt.Errorf("parameter '%s' cannot be empty", paramName), http.StatusBadRequest)
		}
	}
	return param
}

// GetStringParams . . .
func GetStringParams(req *http.Request, paramName string, optional ...bool) (params []string) {
	if len(optional) > 1 {
		panic("too many parameters being passed to server.GetStringParams func")
	}
	return strings.Split(GetStringParam(req, paramName, optional...), ",")
}

// GetBoolParam . . .
func GetBoolParam(req *http.Request, paramName string, optional ...bool) (param bool) {
	if len(optional) > 1 {
		panic("too many parameters being passed to server.GetBoolParam func")
	}
	optional = append(optional, false)
	isOptional := optional[0]

	var paramStrings []string
	var ok bool
	if paramStrings, ok = req.URL.Query()[paramName]; !isOptional && (!ok || len(paramStrings) != 1) {
		PanicWithStatus(fmt.Errorf("missing or invalid url parameter: %s", paramName), http.StatusBadRequest)
	}

	if len(paramStrings) == 1 {
		var err error
		if param, err = strconv.ParseBool(paramStrings[0]); !isOptional && err != nil {
			PanicWithStatus(fmt.Errorf("missing or invalid url parameter: %s\nerror: %s", paramName, err.Error()), http.StatusBadRequest)
		}
	}

	return param
}

// GetIntParam . . .
func GetIntParam(req *http.Request, paramName string) (param int) {
	var paramStrings []string
	var ok bool
	var err error
	if paramStrings, ok = req.URL.Query()[paramName]; !ok || len(paramStrings) != 1 {
		PanicWithStatus(fmt.Errorf("missing or invalid url parameter: %s", paramName), http.StatusBadRequest)
	}
	if param, err = strconv.Atoi(paramStrings[0]); err != nil {
		PanicWithStatus(fmt.Errorf("missing or invalid url parameter: %s\nerror: %s", paramName, err.Error()), http.StatusBadRequest)
	}
	return param
}

// PanicWithStatus . . .
func PanicWithStatus(err error, status int) {
	panic(serverError{Message: err.Error(), Status: status})
}

// DoExternalAPIRequest does api call to external websites
func DoExternalAPIRequest(method, baseURL, url string, requestBody []byte) <-chan io.ReadCloser {

	b := make(chan io.ReadCloser)

	go func() {
		client := &http.Client{}
		defer close(b)
		req, err := http.NewRequest(method, fmt.Sprintf(`%s%s`, baseURL, url), bytes.NewBuffer(requestBody))
		if err != nil {
			panic(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		b <- resp.Body
	}()
	return b
}
