package models

type Vertex struct {
	X, Y, Z float64
}

type Face struct {
	V1, V2, V3 int
}

type Cube struct {
	Center Vertex
	Size	float64
}

type OctreeNode struct {
	CurrentCube   Cube
	Children [8]*OctreeNode
    Face []int
 	IsLeaf   bool
	IsVoxel  bool
}

type OctreeStats struct {
    NodesByDepth    []int 
    SkippedByDepth  []int     
    MaxDepth        int
    LeafNodes       []*OctreeNode 
}


type Model struct {
    Vertices  []Vertex
    Faces     []Face
}
