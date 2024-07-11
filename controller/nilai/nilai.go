package nilai

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "strconv"

    "mahasiswa/database"
    "mahasiswa/model/nilai"

    "github.com/gorilla/mux"
)

func GetNilai(w http.ResponseWriter, r *http.Request) {
    rows, err := database.DB.Query("SELECT * FROM nilai")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var nilailist []nilai.Nilai
    for rows.Next() {
        var c nilai.Nilai
        if err := rows.Scan(&c.IDNilai, &c.IDMahasiswa, &c.IDMatakuliah, &c.Nilai); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        nilailist = append(nilailist, c)
    }

    if err := rows.Err(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(nilailist)
}

func GetNilaiByID(w http.ResponseWriter, r *http.Request) {
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

    var nilai nilai.Nilai
    query := "SELECT id_nilai, id_mahasiswa, id_matakuliah, nilai FROM nilai WHERE id_nilai = ?"
    err = database.DB.QueryRow(query, id).Scan(&nilai.IDNilai, &nilai.IDMahasiswa, &nilai.IDMatakuliah, &nilai.Nilai)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Nilai not found", http.StatusNotFound)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(nilai)
}

func PostNilai(w http.ResponseWriter, r *http.Request) {
    var pc nilai.Nilai
    if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Prepare the SQL statement for inserting a new course
    query := `
    INSERT INTO nilai (id_mahasiswa, id_matakuliah, nilai) 
    VALUES (?, ?, ?)`

    // Execute the SQL statement
    res, err := database.DB.Exec(query, pc.IDMahasiswa, pc.IDMatakuliah, pc.Nilai)
    if err != nil {
        http.Error(w, "Failed to insert nilai: "+err.Error(), http.StatusInternalServerError)
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
        "message": "Nilai added successfully",
        "id":      id,
    })
}

func PutNilai(w http.ResponseWriter, r *http.Request) {
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
    var pc nilai.Nilai
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
    UPDATE nilai
    SET id_mahasiswa=?, id_matakuliah=?, nilai=?
    WHERE id_nilai=?`

    // Execute the SQL statement
    result, err := database.DB.Exec(query, pc.IDMahasiswa, pc.IDMatakuliah, pc.Nilai, id)
    if err != nil {
        http.Error(w, "Failed to update nilai: "+err.Error(), http.StatusInternalServerError)
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
        "message": "Nilai updated successfully",
    })
}

func DeleteNilai(w http.ResponseWriter, r *http.Request) {
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

    // Prepare the SQL statement for deleting a nilai
    query := `
    DELETE FROM nilai
    WHERE id_nilai = ?`

    // Execute the SQL statement
    result, err := database.DB.Exec(query, id)
    if err != nil {
        http.Error(w, "Failed to delete nilai: "+err.Error(), http.StatusInternalServerError)
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
        "message": "Nilai deleted successfully",
    })
}