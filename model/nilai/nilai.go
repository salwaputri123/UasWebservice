package nilai

type Nilai struct {
    IDNilai      int `json:"id_nilai"`
    IDMahasiswa  int `json:"id_mahasiswa"`
    IDMatakuliah int `json:"id_matakuliah"`
    Nilai        string `json:"nilai"`
}