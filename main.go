package main

import (
	"fmt"
	"log"
	"mahasiswa/controller/auth"
	"mahasiswa/controller/matakuliah"
	"mahasiswa/controller/pendaftaran"
	"mahasiswa/database"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	database.InitDB()
	fmt.Println("Hello World")

	router := mux.NewRouter()

	router.HandleFunc("/regis", auth.Registration).Methods("POST")
	router.HandleFunc("/login", auth.Login).Methods("POST")
	
	//Router handler Matakuliah
	router.HandleFunc("/matakuliah", matakuliah.GetMatakuliah).Methods("GET")
	router.HandleFunc("/matakuliah", auth.JWTAuth(matakuliah.PostMatakuliah)).Methods("POST")
	router.HandleFunc("/matakuliah/{id}", auth.JWTAuth(matakuliah.PutMatakuliah)).Methods("PUT")
	router.HandleFunc("/matakuliah/{id}", auth.JWTAuth(matakuliah.DeleteMatakuliah)).Methods("DELETE")

	//Router handler Pendaftaran
	router.HandleFunc("/pendaftaran", pendaftaran.GetPendaftaran).Methods("GET")
	router.HandleFunc("/pendaftaran", auth.JWTAuth(pendaftaran.PostPendaftaran)).Methods("POST")
	router.HandleFunc("/pendaftaran/{id}", auth.JWTAuth(pendaftaran.PutPendaftaran)).Methods("PUT")
	router.HandleFunc("/pendaftaran/{id}", auth.JWTAuth(pendaftaran.DeletePendaftaran)).Methods("DELETE")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:5500"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		Debug: true,
	})

	handler := c.Handler(router)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}