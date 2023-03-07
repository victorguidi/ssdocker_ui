package main

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"
	"strings"
)

type Api struct {
	Addr      string
	SshServer *SshServer
}

func New(addr string, s *SshServer) *Api {
	return &Api{
		Addr:      addr,
		SshServer: s,
	}
}

func (a *Api) Start() error {
	a.SshServer.Connect()
	mime.AddExtensionType(".js", "application/javascript")
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/api/dockers", a.GetDockers)
	http.HandleFunc("/api/servers", a.GetServer)
	http.HandleFunc("/api/send", a.SendCommands)
	log.Println("Server started on port", a.Addr)
	log.Fatal(http.ListenAndServe(a.Addr, nil))
	return nil
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).WriteHeader(http.StatusOK)
}

func WriteJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}

func (a *Api) GetServer(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	type response struct {
		Servers []string `json:"servers"`
	}
	resp := response{
		Servers: a.SshServer.Hosts,
	}
	WriteJSON(w, resp)
}

func (a *Api) GetDockers(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	type Response struct {
		Containers []Container `json:"containers"`
	}

	resp := Response{}

	for _, conn := range a.SshServer.SshConnection {
		conn.Channel <- "docker ps -a"
		response := <-conn.Channel

		containers, err := ProcessCommandOne(response)
		if err != nil {
			panic(err)
		}
		for i := range containers {
			containers[i].Server = conn.Host
		}
		resp.Containers = append(resp.Containers, containers...)
	}
	w.Header().Set("Content-Type", "application/json")
	WriteJSON(w, resp)
}

func (a *Api) SendCommands(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	type Request struct {
		Server  string `json:"server"`
		Command string `json:"command"`
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	type ResponseContainers struct {
		Containers []Container `json:"containers"`
	}

	type ResponseLogs struct {
		Logs string `json:"logs"`
	}

	resp := ResponseContainers{}
	respLogs := ResponseLogs{}

	for _, conn := range a.SshServer.SshConnection {
		if conn.Host == req.Server {
			conn.Channel <- req.Command
			response := <-conn.Channel
			if req.Command == "docker ps -a" {
				containers, err := ProcessCommandOne(response)
				if err != nil {
					panic(err)
				}
				for i := range containers {
					containers[i].Server = conn.Host
				}
				resp.Containers = append(resp.Containers, containers...)
				w.Header().Set("Content-Type", "application/json")
				WriteJSON(w, resp)
				return
			} else if strings.Contains(req.Command, "logs") {
				respLogs.Logs = response
				w.Header().Set("Content-Type", "application/json")
				WriteJSON(w, respLogs)
				return
			}
		}
	}
}
