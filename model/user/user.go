package user

type Mahasiswa struct {
	IdMahasiswa int `json:"id_mahasiswa"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nama string `json:"nama"`
}