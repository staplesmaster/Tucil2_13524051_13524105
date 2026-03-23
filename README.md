# Voxelizer — 3D OBJ Voxelization via Octree (Divide and Conquer)

Tugas Kecil 1 IF2211 Strategi Algoritma — Voxelization Objek 3D menggunakan Octree

## Deskripsi

Program ini mengkonversi model 3D format `.obj` menjadi model voxel menggunakan struktur data **Octree** dengan algoritma **Divide and Conquer**. Transformasi hanya dilakukan pada permukaan objek (*surface voxelization*).

## Struktur File

```
voxelizer/
├── main.go        # Entry point, CLI, output statistik
├── parser.go      # Parser file .obj (vertices + faces)
├── geometry.go    # Vec3, AABB, triangle-AABB intersection (SAT)
├── octree.go      # Octree D&C build algorithm
├── writer.go      # Penulis output .obj (voxel cubes)
└── go.mod
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

### Kompleksitas

- Waktu: `O(T × 8^D)` worst-case, dengan T = jumlah segitiga dan D = kedalaman octree
- Ruang: `O(8^D)` worst-case untuk node, tapi dipangkas secara agresif untuk objek yang memiliki permukaan jauh lebih kecil dari volume total

### Triangle-AABB Intersection (SAT)

Pengujian apakah sebuah segitiga berinterseksi dengan AABB menggunakan **Separating Axis Theorem**:
- 3 axis dari sisi AABB (sumbu koordinat)
- 1 axis dari normal segitiga
- 9 axis dari cross product antara edge segitiga dengan sumbu koordinat

Total 13 sumbu pengujian.

## Cara Menjalankan

### Prasyarat

- Go 1.18+

### Build

```bash
go build -o voxelizer .
```

### Menjalankan

```bash
./voxelizer <input.obj> <max_depth>
```

**Parameter:**
- `input.obj` — path ke file model 3D
- `max_depth` — kedalaman maksimum octree (integer 1–12)

**Contoh:**
```bash
./voxelizer models/pumpkin.obj 6
./voxelizer models/bunny.obj 7
```

### Output

Program menghasilkan:
1. File `<nama>-voxelized.obj` di direktori yang sama dengan input
2. Statistik di CLI:
   - Jumlah voxel, vertex, dan faces yang terbentuk
   - Statistik node octree per kedalaman
   - Statistik node yang dipangkas per kedalaman
   - Kedalaman octree
   - Lama waktu eksekusi
   - Path file output

### Contoh Output CLI

```
=== Voxelizer ===
Input  : models/pumpkin.obj
Depth  : 6

[1/3] Parsing OBJ file...
      Vertices : 5842
      Faces    : 11680

[2/3] Building Octree (max depth = 6)...
[3/3] Writing voxelized OBJ to: models/pumpkin-voxelized.obj

──────────────────────────────────────────────────
HASIL / RESULTS
──────────────────────────────────────────────────
Voxel terbentuk   : 2847
Vertex terbentuk  : 22776
Faces terbentuk   : 34164
Kedalaman octree  : 6

Statistik node octree yang terbentuk:
   1 : 1
   2 : 8
   ...
```

## Format OBJ yang Didukung

Program memperhatikan baris dengan prefix:
- `v x y z` — vertex (koordinat float)
- `f i j k` — face triangular (indeks 1-based)

Baris lain (`vt`, `vn`, `mtl`, `#`, dll.) diabaikan.

## Bonus

- [ ] Concurrency (goroutines per subtree)
- [ ] Interactive 3D viewer

## Anggota Kelompok

| Nama | NIM |
|------|-----|
| ...  | ... |
| ...  | ... |
