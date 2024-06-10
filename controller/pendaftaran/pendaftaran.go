package pendaftaran

import (
	"encoding/json"
	"net/http"
	"strconv"
	"log"

	"mahasiswa/database"
	"mahasiswa/model/pendaftaran"
	"github.com/gorilla/mux"
)

func GetPendaftaran(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT * FROM pendaftaran")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pendaftaranlist []pendaftaran.Pendaftaran
	for rows.Next() {
	var c pendaftaran.Pendaftaran
	if err := rows.Scan(&c.IdPendaftaran,&c.IdMahasiswa,&c.IdMatakuliah, &c.Nilai); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pendaftaranlist = append(pendaftaranlist, c)
}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pendaftaranlist)
}
		
func PostPendaftaran(w http.ResponseWriter, r *http.Request) {
	var pc pendaftaran.Pendaftaran
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	
	// Prepare the SQL statement for inserting a new course
	query := `
	INSERT INTO pendaftaran (id_mahasiswa, id_matakuliah, nilai) 
	VALUES (?, ?, ?)`
	
	// Execute the SQL statement
	res, err := database.DB.Exec(query, pc.IdMahasiswa, pc.IdMatakuliah, pc.Nilai)
	if err != nil {
		http.Error(w, "Failed to insert pendaftaran: "+err.Error(), http.StatusInternalServerError)
		return
	}
		
	// Get the last inserted ID
	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve last insert ID: "+err.Error(), http.StatusInternalServerError)
		return
	}
		
	// Return the newly created ID in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Course added successfully",
		"id":      id,
	})
}

func PutPendaftaran(w http.ResponseWriter, r *http.Request) {
    // Ambil ID dari URL
    vars := mux.Vars(r)
    idStr, ok := vars["id"]
    if !ok {
        http.Error(w, "ID not provided", http.StatusBadRequest)
        return
    }
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    // Log the received ID for debugging
    log.Printf("Received ID for update: %d", id)

    // Decode JSON body
    var pc pendaftaran.Pendaftaran
    decoder := json.NewDecoder(r.Body)
    err = decoder.Decode(&pc)
    if err != nil {
        http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Log decoded struct for debugging
    log.Printf("Decoded struct: %+v", pc)

    // Prepare the SQL statement for updating the category admin
    query := `
        UPDATE pendaftaran
        SET id_mahasiswa=?, id_matakuliah=?, nilai=?
        WHERE id_pendaftaran=?`

    // Execute the SQL statement
    result, err := database.DB.Exec(query, pc.IdMahasiswa, pc.IdMatakuliah, pc.Nilai, id)
    if err != nil {
        http.Error(w, "Failed to update pendaftaran: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Get the number of rows affected
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Check if any rows were updated
    if rowsAffected == 0 {
        http.Error(w, "No rows were updated", http.StatusNotFound)
        return
    }

    // Return success message
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Pendaftaran updated successfully",
    })
}

func DeletePendaftaran(w http.ResponseWriter, r *http.Request) {
    // Extract ID from URL
    vars := mux.Vars(r)
    idStr, ok := vars["id"]
    if !ok {
        http.Error(w, "ID not provided", http.StatusBadRequest)
        return
    }
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    // Prepare the SQL statement for deleting a pendaftaran
    query := `
        DELETE FROM pendaftaran
        WHERE id_pendaftaran = ?`

    // Execute the SQL statement
    result, err := database.DB.Exec(query, id)
    if err != nil {
        http.Error(w, "Failed to delete pendaftaran: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Check if any rows were affected
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
        return
    }

    if rowsAffected == 0 {
        http.Error(w, "No rows were deleted", http.StatusNotFound)
        return
    }

    // Return the response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Pendaftaran deleted successfully",
    })
}