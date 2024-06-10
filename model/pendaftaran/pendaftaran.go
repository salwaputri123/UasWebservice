package pendaftaran

type Pendaftaran struct{
	IdPendaftaran int `json:"id_pendaftaran"`
	IdMahasiswa int `json:"id_mahasiswa"`
	IdMatakuliah int `json:"id_matakuliah"`
	Nilai string `json:"nilai"`
}