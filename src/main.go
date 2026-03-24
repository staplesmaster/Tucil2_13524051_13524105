package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/staplesmaster/Tucil2_13524051_13524105/src/file"
	"github.com/staplesmaster/Tucil2_13524051_13524105/src/voxelization"
)

func main() {
	red := "\033[31m"
    reset := "\033[0m" 
	
	for {filePath, maxDepth :=file.ValidateInput()
		printFilePath := filepath.Join("test", filepath.Base(filePath))
		fmt.Println("Mulai memproses file:", printFilePath)
		fmt.Println("Target kedalaman Octree:", maxDepth)

		startTime := time.Now()

		// Parse file menjadi struct Model
		parsedModel, err := file.ParseOBJ(filePath)
		if err != nil {
			fmt.Printf(red + "Error : Gagal memproses file .obj: %v\n" + reset, err)
			os.Exit(1)
		}

		// Voxelization
		_, stats := voxelization.StartVoxelization(parsedModel, maxDepth)

		// Total face dan vertex hasil
		totalVoxels := len(stats.LeafNodes)
		totalVertices := totalVoxels * 8
		totalFaces := totalVoxels * 12

		// Output

		outputPath, err := file.WriteObj(filePath, maxDepth, stats.LeafNodes)

		if err != nil {
			fmt.Printf(red + "Error : Gagal menyimpan file .obj: %v\n" + reset, err)
			os.Exit(1)
		}


		elapsedTime := time.Since(startTime)

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
		printOutputPath := filepath.Join("test", filepath.Base(outputPath))
		fmt.Printf("Path dimana file .obj disimpan: %s\n\n", printOutputPath)

		fmt.Print("Apakah masih mau melakukan voxelization? (y/n) : ")
		var lanjut string
		fmt.Scan(&lanjut)
		
		for strings.Compare(lanjut, "y") != 0 && strings.Compare(lanjut, "n") != 0{
			fmt.Print(red + "Error : Input tidak valid, masukan yang valid (y/n) : " + reset)
			fmt.Scan(&lanjut)
			lanjut = strings.ToLower(lanjut)
		}
		if strings.Compare(lanjut, "y") == 0 {
			continue
		}
		
		os.Exit(0)
		
	}
}

