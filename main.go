package main

import (
	"fmt"
	"log"
	"mahasiswa/controller/auth"
	"mahasiswa/controller/matakuliah"
	"mahasiswa/controller/nilai"
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
	router.HandleFunc("/matakuliah/{id}", matakuliah.GetMatakuliahByID).Methods("GET")
	router.HandleFunc("/matakuliah", auth.JWTAuth(matakuliah.PostMatakuliah)).Methods("POST")
	router.HandleFunc("/matakuliah/{id}", auth.JWTAuth(matakuliah.PutMatakuliah)).Methods("PUT")
	router.HandleFunc("/matakuliah/{id}", auth.JWTAuth(matakuliah.DeleteMatakuliah)).Methods("DELETE")

	//Router handler Nilai
	router.HandleFunc("/nilai", nilai.GetNilai).Methods("GET")
	router.HandleFunc("/nilai/{id}", nilai.GetNilaiByID).Methods("GET")
	router.HandleFunc("/nilai", auth.JWTAuth(nilai.PostNilai)).Methods("POST")
	router.HandleFunc("/nilai/{id}", auth.JWTAuth(nilai.PutNilai)).Methods("PUT")
	router.HandleFunc("/nilai/{id}", auth.JWTAuth(nilai.DeleteNilai)).Methods("DELETE")

	//mengaktifkan CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		Debug:          true,
	})

	handler := c.Handler(router)

	fmt.Println("Server is running on http://localhost:8026")
	log.Fatal(http.ListenAndServe(":8026", handler))
}
