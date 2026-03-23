package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/staplesmaster/Tucil2_13524051_13524105/src/voxelization"
	"github.com/staplesmaster/Tucil2_13524051_13524105/src/file"

)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Format pemanggilan salah!")
		fmt.Println("Cara penggunaan: go run main.go <path_file.obj> <kedalaman_maksimum>")
		os.Exit(1)
	}

	filePath := os.Args[1]

	maxDepth, err := strconv.Atoi(os.Args[2])
	if err != nil || maxDepth < 0 {
		fmt.Println("Error: Kedalaman maksimum harus berupa angka bulat positif")
		os.Exit(1)
	}

	fmt.Println("Mulai memproses file:", filePath)
	fmt.Println("Target kedalaman Octree:", maxDepth)

	startTime := time.Now()

	// 1. Parse file menjadi struct Model
	parsedModel, err := file.ParseOBJ(filePath)
	if err != nil {
		fmt.Printf("Gagal memproses file .obj: %v\n", err)
		os.Exit(1)
	}

	// 2. Mulai Voxelization (Kita abaikan rootNode dengan '_', ambil stats-nya saja!)
	_, stats := voxelization.StartVoxelization(parsedModel, maxDepth)

	// 3. Kalkulasi total instan langsung dari stats.LeafNodes
	totalVoxels := len(stats.LeafNodes)
	totalVertices := totalVoxels * 8
	totalFaces := totalVoxels * 12

	// 4. Manajemen Output Path
	baseName := filepath.Base(filePath)
	ext := filepath.Ext(baseName)
	queryFile := strings.TrimSuffix(baseName, ext)

	err = os.MkdirAll("test", os.ModePerm)
	if err != nil {
		fmt.Printf("Gagal membuat folder test: %v\n", err)
		os.Exit(1)
	}

	outputFileName := fmt.Sprintf("%sSol_%d.obj", queryFile, maxDepth)
	outputPath := filepath.Join("test", outputFileName)
	
	// 5. Tulis file langsung menggunakan daftar daun (LeafNodes)
	err = file.WriteObj(outputPath, stats.LeafNodes)
	if err != nil {
		fmt.Printf("Gagal menyimpan file .obj: %v\n", err)
		os.Exit(1)
	}

	elapsedTime := time.Since(startTime)

	// 6. Cetak Statistik dengan elegan
	fmt.Println("\n=== Statistik Voxelization ===")
	fmt.Printf("Banyaknya voxel yang terbentuk: %d\n", totalVoxels)
	fmt.Printf("Banyaknya vertex yang terbentuk: %d\n", totalVertices)
	fmt.Printf("Banyaknya faces yang terbentuk: %d\n", totalFaces)

	fmt.Println("Statistik node octree yang terbentuk")
	for i := 1; i <= maxDepth; i++ {
		fmt.Printf("%d : %d\n", i, stats.NodesByDepth[i])
	}

	fmt.Println("Statistik node yang tidak perlu ditelusuri")
	for i := 1; i <= maxDepth; i++ {
		fmt.Printf("%d : %d\n", i, stats.SkippedByDepth[i])
	}

	fmt.Printf("Kedalaman octree: %d\n", maxDepth)
	fmt.Printf("Lama waktu program berjalan: %s\n", elapsedTime)
	fmt.Printf("Path dimana file .obj disimpan: %s\n", outputPath)
}