package geometry

import (
	"math"
	"src/models"
)

func isIntersecting(cube models.Cube, vertices []models.Vertex, faces []models.Face) bool {
    h := cube.Size / 2

    for _, face := range faces {
        v0 := vertices[face.V1-1]
        v1 := vertices[face.V2-1]
        v2 := vertices[face.V3-1]

        // Geser segitiga agar seolah-olah relatif terhadap pusat kubus di (0, 0, 0)
        v0 = sub(v0, cube.Center)
        v1 = sub(v1, cube.Center)
        v2 = sub(v2, cube.Center)

        // Tes 1, bungkus segitiga dengan balok
        // Cek apakah kubus berada di dalam balok pembungkus
        minX, maxX := findMinMax(v0.X, v1.X, v2.X)
        if minX > h || maxX < -h { continue }

        minY, maxY := findMinMax(v0.Y, v1.Y, v2.Y)
        if minY > h || maxY < -h { continue }

        minZ, maxZ := findMinMax(v0.Z, v1.Z, v2.Z)
        if minZ > h || maxZ < -h { continue }

        // Kalkulasi vektor sisi-sisi segitiga untuk Tes 2 dan Tes 3
        e0 := sub(v1, v0)
        e1 := sub(v2, v1)
        e2 := sub(v0, v2)

        // Tes 2, anggap segitiga sebagai bidang tak hingga
        // Cek apakah kubus berpotongan dengan bidang tersebut atau sepenuhnya berada atas/bawah bidang
        normal := cross(e0, e1)
        planeD := dot(normal, v0)
        rPlane := h*math.Abs(normal.X) + h*math.Abs(normal.Y) + h*math.Abs(normal.Z)
        
        if math.Abs(planeD) > rPlane { continue } // Terpisah oleh cermin!

        // Tes 3
        // Cek apakah kubus A benar-benar berpotongan dengan segitiga (bukan bidang perluasan seperti di Tes 2)
        edges := []models.Vertex{e0, e1, e2}
        isSeparated := false
        
        for _, e := range edges {
            aX := models.Vertex{X: 0, Y: -e.Z, Z: e.Y}
            if axisTest(aX, v0, v1, v2, h) { isSeparated = true; break }

            aY := models.Vertex{X: e.Z, Y: 0, Z: -e.X}
            if axisTest(aY, v0, v1, v2, h) { isSeparated = true; break }

            aZ := models.Vertex{X: -e.Y, Y: e.X, Z: 0}
            if axisTest(aZ, v0, v1, v2, h) { isSeparated = true; break }
        }
        
        if isSeparated { continue } 

        return true
    }

    return false
}


func sub(a, b models.Vertex) models.Vertex {
    return models.Vertex{X: a.X - b.X, Y: a.Y - b.Y, Z: a.Z - b.Z}
}

func cross(a, b models.Vertex) models.Vertex {
    return models.Vertex{
        X: a.Y*b.Z - a.Z*b.Y,
        Y: a.Z*b.X - a.X*b.Z,
        Z: a.X*b.Y - a.Y*b.X,
    }
}

func dot(a, b models.Vertex) float64 {
    return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func findMinMax(x0, x1, x2 float64) (float64, float64) {
    min := math.Min(x0, math.Min(x1, x2))
    max := math.Max(x0, math.Max(x1, x2))
    return min, max
}

func axisTest(axis, v0, v1, v2 models.Vertex, h float64) bool {
    // Proyeksi ortogonal titik-titik segitiga di axis
    p0 := dot(v0, axis)
    p1 := dot(v1, axis)
    p2 := dot(v2, axis)

    // Titik min dan max bayangan segitiga
    min, max := findMinMax(p0, p1, p2)
    
    // Bayangan kubus membentang dari r ke -r
    r := h*math.Abs(axis.X) + h*math.Abs(axis.Y) + h*math.Abs(axis.Z)
    
    // True apabila ada celah antara kedua bayangannya, artinya kubus dan segitiga tidak saling berpotongan
    return min > r || max < -r
}
