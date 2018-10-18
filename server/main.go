package main

import (
	"net/http"

	"github.com/chinput/InputMethodService/server/config"
	"github.com/chinput/InputMethodService/server/model"
)

func main() {
	config.Init("../config.toml")

	model.InitByHand(model.InitConf{
		Dbtype:    "mongo",
		Dbhost:    config.DBUrl(),
		Dbname:    config.DBName(),
		Findlimit: 10,
	})

	http.Handle("/tmp/", http.StripPrefix("/tmp/", http.FileServer(http.Dir(config.TmpPath()))))
	http.HandleFunc("/api", CommonHandler)
	http.Handle("/debug/", http.StripPrefix("/debug/", http.FileServer(http.Dir("test"))))

	http.ListenAndServe(":8080", nil)
}
