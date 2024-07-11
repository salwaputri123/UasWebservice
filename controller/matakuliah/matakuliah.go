package matakuliah

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strconv"
    "log"

    "mahasiswa/database"
    "github.com/gorilla/mux"
    "mahasiswa/model/matakuliah" 
)

func GetMatakuliah(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT * FROM matakuliah")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var matakuliahList[]matakuliah.MataKuliah
    for rows.Next() {
        var c matakuliah.MataKuliah
        if err := rows.Scan(&c.IdMatakuliah, &c.NamaMatakuliah); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        matakuliahList = append(matakuliahList, c)
    }
	
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matakuliahList)
}

func GetMatakuliahByID(w http.ResponseWriter, r *http.Request) {
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

    var matakuliah matakuliah.MataKuliah
    query := "SELECT id_matakuliah, nama_matakuliah FROM matakuliah WHERE id_matakuliah = ?"
    err = database.DB.QueryRow(query, id).Scan(&matakuliah.IdMatakuliah, &matakuliah.NamaMatakuliah)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Matakuliah not found", http.StatusNotFound)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(matakuliah)
}


func PostMatakuliah(w http.ResponseWriter, r *http.Request) {
	var pc matakuliah.MataKuliah
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for inserting a new course
	query := `
		INSERT INTO matakuliah (nama_matakuliah) 
		VALUES (?)`

	// Execute the SQL statement
	res, err := database.DB.Exec(query, pc.NamaMatakuliah)
	if err != nil {
		http.Error(w, "Failed to insert matakuliah: "+err.Error(), http.StatusInternalServerError)
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
		"message": "matakuliah added successfully",
		"id":      id,
	})
}
func PutMatakuliah(w http.ResponseWriter, r *http.Request) {
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

    // Decode JSON body
    var pc matakuliah.MataKuliah
    if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Prepare the SQL statement for updating the category admin
    query := `
        UPDATE matakuliah
        SET nama_matakuliah=?
        WHERE id_matakuliah=?`

    // Execute the SQL statement
    result, err := database.DB.Exec(query, pc.NamaMatakuliah, id)
	if err != nil {
        http.Error(w, "Failed to update matakuliah: "+err.Error(), http.StatusInternalServerError)
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
        "message": "Course updated successfully",
    })
}

func DeleteMatakuliah(w http.ResponseWriter, r *http.Request) {
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

    // Log the received ID for debugging
    log.Printf("Received ID for deletion: %d", id)

    // Prepare the SQL statement for deleting a matakuliah
    query := `
        DELETE FROM matakuliah
        WHERE id_matakuliah = ?`

    // Execute the SQL statement
    result, err := database.DB.Exec(query, id)
    if err != nil {
        http.Error(w, "Failed to delete matakuliah: "+err.Error(), http.StatusInternalServerError)
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
        "message": "Matakuliah deleted successfully",
    })
}


