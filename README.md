# mux-mongo-project
Test project:

Go webserver

Database: Mongodb

# Build
```docker build -t mux-mongo-project .```

# Run
```docker-compose up -d```

It will start the go webserver at `http://localhost:8080`

# Endpoints
```
func initializeRouter() {
	log.Print("Initializing routes")
	r := mux.NewRouter()
	r.HandleFunc("/", handleCreate).Methods("POST")
	r.HandleFunc("/", handleRead).Methods("GET")
	r.HandleFunc("/download", handleDownload).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}
```
