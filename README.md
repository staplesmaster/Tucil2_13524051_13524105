# Tugas Kecil 1 IF2211 Strategi Algoritma — Voxelization Objek 3D menggunakan Octree

## Deskripsi

Program ini mengkonversi model 3D format `.obj` menjadi model voxel menggunakan struktur data **Octree** dengan algoritma **Divide and Conquer**. Transformasi hanya dilakukan pada permukaan objek (*surface voxelization*).

## Struktur File

```
Tucil2_13524051_13524105/
├── go.mod
├── README.md
├── bin/                  # binary hasil build
├── docs/
│   └── Laporan Tucil 2 - 13524103 - 13524105      
├── src/
│   ├── main.go           # entry point program
│   ├── file/
│   │   └── file.go       # parsing/IO file .obj
│   ├── geometry/
│   │   └── geometry.go   # operasi geometri 3D
│   ├── models/
│   │   └── models.go     # struktur data model
│   └── voxelization/
│       └── voxelization.go  # algoritma octree voxelization
│ 
├── test/                 # data uji
```

## Algoritma Divide and Conquer

### Ide Utama

Octree membagi ruang 3D menjadi 8 oktan secara rekursif. Setiap node merepresentasikan sebuah kubus di ruang 3D. Subdivisi berhenti ketika:
1. Tidak ada segitiga yang berinterseksi dengan node (node dipangkas / pruned), atau  
2. Kedalaman maksimum tercapai → node menjadi **leaf voxel**

### Langkah-Langkah

1. **Inisialisasi**: Hitung *bounding box* seluruh model, lalu bentuk menjadi kubus kubik sebagai root.

2. **Divide**: Bagi kubus menjadi 8 anak (oktan) dengan membelah di titik tengah (midpoint) setiap sumbu X, Y, Z.

3. **Filter**: Untuk setiap oktan anak, uji apakah ada segitiga dari model yang berinterseksi dengan kotak anak tersebut (menggunakan Separating Axis Theorem / SAT).
   - Jika tidak ada → **pangkas** oktan ini (tidak perlu ditelusuri lebih dalam).

4. **Conquer**: Jika ada segitiga yang berinterseksi dan kedalaman belum maksimum → rekursi ke anak-anak.

5. **Base Case**: Jika kedalaman = maxDepth → tandai node sebagai **leaf voxel** (voxel permukaan).


### Triangle-AABB Intersection (SAT)

Pengujian apakah sebuah segitiga berinterseksi dengan AABB menggunakan **Separating Axis Theorem**:
- 3 axis dari sisi AABB (sumbu koordinat)
- 1 axis dari normal segitiga
- 9 axis dari cross product antara edge segitiga dengan sumbu koordinat

Total 13 sumbu pengujian.

## Cara Menjalankan

### Prasyarat

- Go 1.18+

### Clone Repository

```bash
git clone https://github.com/staplesmaster/Tucil2_13524051_13524105.git
cd Tucil2_13524051_13524105
```

### Menjalankan di Linux

#### Build & run
Menjalankan di root folder 
```bash
go build -o voxelizer ./src
./voxelizer 
```

#### Menjalankan di bin

```bash
./voxelizer 
```

### Menjalankan di Windows

#### Build (.exe) & run

```powershell
go build -o voxelizer.exe ./src
.\voxelizer.exe 
```

#### Menjalankan di bin

```powershell
.\voxelizer.exe 
```
Catatan: di Windows executable harus `.exe`

### Output

Program menghasilkan:
1. File `<nama>-Voxelization_Result_Depth-<max depth>.obj` di folder result
2. Statistik di CLI:
   - Jumlah voxel, vertex, dan faces yang terbentuk
   - Statistik node octree per kedalaman
   - Statistik node yang dipangkas per kedalaman
   - Kedalaman octree
   - Lama waktu eksekusi
   - Path file output

## Format OBJ yang Didukung

Program memperhatikan baris dengan prefix:
- `v x y z` — vertex (koordinat float)
- `f i j k` — face triangular (indeks 1-based)

## Anggota Kelompok

| Nama | NIM |
|------|-----|
| Mikhael Andrian Yonatan | 13524051 |
| Nicholas Luis Chandra | 13524105 |
