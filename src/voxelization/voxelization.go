package voxelization

import (
	"math"
	"sync"
	"github.com/staplesmaster/Tucil2_13524051_13524105/src/models"
	"github.com/staplesmaster/Tucil2_13524051_13524105/src/geometry"
)

func StartVoxelization(parsedModel *models.Model, maxDepth int) (*models.OctreeNode, *models.OctreeStats) {
	minX, minY, minZ := math.Inf(1), math.Inf(1), math.Inf(1)
	maxX, maxY, maxZ := math.Inf(-1), math.Inf(-1), math.Inf(-1)

	// Mencari nilai ekstrem untuk setiap sumbu
	for i := range parsedModel.Vertices {
		minX = math.Min(minX, parsedModel.Vertices[i].X)
		minY = math.Min(minY, parsedModel.Vertices[i].Y)
		minZ = math.Min(minZ, parsedModel.Vertices[i].Z)

		maxX = math.Max(maxX, parsedModel.Vertices[i].X)
		maxY = math.Max(maxY, parsedModel.Vertices[i].Y)
		maxZ = math.Max(maxZ, parsedModel.Vertices[i].Z)
	}

	// Menentukan selisih terbesar sebagai ukuran kubus
	dx, dy, dz := maxX-minX, maxY-minY, maxZ-minZ
	side := math.Max(dx, math.Max(dy, dz))

	// Membentuk kubus paling pertama
	rootCube := models.Cube{
		Center: models.Vertex{
			X: (minX + maxX) / 2,
			Y: (minY + maxY) / 2,
			Z: (minZ + maxZ) / 2,
		},
		Size: side,
	}

	// Mengisi daftar face yang valid dari kubus pertama (semua face)
	initialFaces := make([]int, len(parsedModel.Faces))
	for i := range parsedModel.Faces {
		initialFaces[i] = i
	}

	stats := &models.OctreeStats{
		NodesByDepth:   make([]int, maxDepth+1),
		SkippedByDepth: make([]int, maxDepth+1),
		MaxDepth:       maxDepth,
		LeafNodes:      make([]*models.OctreeNode, 0),
	}

	// Membangun pohon Octree secara rekursif (Bottom-Up)
	rootNode := divideConquer(rootCube, parsedModel, initialFaces, 0, maxDepth, stats)

	return rootNode, stats
}

func divideConquer(cube models.Cube, parsedModel *models.Model, parentFaces []int, currentDepth, maxDepth int, stats *models.OctreeStats) *models.OctreeNode {
	// Pencatatatan menggunakan mutex
	stats.Mutex.Lock()
	stats.NodesByDepth[currentDepth]++
	stats.Mutex.Unlock()

	validFaces := geometry.GetIntersectingFaces(cube, parsedModel.Vertices, parsedModel.Faces, parentFaces)

	node := &models.OctreeNode{
		CurrentCube: cube,
		Face:        validFaces,
		IsLeaf:      currentDepth == maxDepth,
		IsVoxel:     false,
	}

	// Tidak ada objek di dalam/bersinggungan dengan kubus ini, maka skip
	if len(validFaces) == 0 {
		// Jika diskip sebelum mencapai maxDepth, catat sebagai skipped
		if currentDepth < maxDepth {
			stats.Mutex.Lock()
			stats.SkippedByDepth[currentDepth]++
			stats.Mutex.Unlock()
		}
		return node
	}

	// Ada objek
	node.IsVoxel = true

	// Sudah mencapai kedalaman maksimum, maka stop
	if node.IsLeaf {
		// Catat node ke dalam LeafNodes
		stats.Mutex.Lock()
		stats.LeafNodes = append(stats.LeafNodes, node)
		stats.Mutex.Unlock()
		
		return node
	}

	// Masih ada objek dan belum maxDepth, maka bagi menjadi 8 sub-kubus
	halfSize := cube.Size / 2    // rusuk sub-kubus
	quarterSize := cube.Size / 4 // pusat sub-kubus

	// Membuat 8 sub-kubus
	offsets := []float64{-quarterSize, quarterSize}

	// Batas kedalaman perekrutan "pekerja"
	const maxConcurrentDepth = 2

	if currentDepth < maxConcurrentDepth {
		// Merekrut "pekerja"
		var wg sync.WaitGroup
		childIndex := 0
		// Membuat 8 kubus
		for _, ox := range offsets {
			for _, oy := range offsets {
				for _, oz := range offsets {
					subCube := models.Cube{
						Center: models.Vertex{
							X: cube.Center.X + ox,
							Y: cube.Center.Y + oy,
							Z: cube.Center.Z + oz,
						},
						Size: halfSize,
					}

					wg.Add(1)

					go func(idx int, c models.Cube) {
						defer wg.Done()
						// Rekursi untuk setiap anak, lalu simpan ke array Children
						// Anak-anaknya hanya memroses validFaces milik parent (node.Face)
						node.Children[idx] = divideConquer(c, parsedModel, node.Face, currentDepth+1, maxDepth, stats)
					}(childIndex, subCube)

					childIndex++
				}
			}
		}

		wg.Wait() // Tunggu ke-8 anak selesai diproses

	} else {
		childIndex := 0

		for _, ox := range offsets {
			for _, oy := range offsets {
				for _, oz := range offsets {
					subCube := models.Cube{
						Center: models.Vertex{
							X: cube.Center.X + ox,
							Y: cube.Center.Y + oy,
							Z: cube.Center.Z + oz,
						},
						Size: halfSize,
					}

					// Rekursi untuk setiap anak, lalu simpan ke array Children
					node.Children[childIndex] = divideConquer(subCube, parsedModel, node.Face, currentDepth+1, maxDepth, stats)
					childIndex++
				}
			}
		}
	}

	return node
}