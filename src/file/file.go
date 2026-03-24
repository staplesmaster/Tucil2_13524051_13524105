package file

import (
	"bufio"
	"fmt"
	"os"
	"github.com/staplesmaster/Tucil2_13524051_13524105/src/models"
	"strconv"
	"strings"
	"path/filepath"
)

func WriteObj(filePath string, maxDepth int, leaves []*models.OctreeNode) (outputPath string, err error) {

	baseName := filepath.Base(filePath)
	ext := filepath.Ext(baseName)
	queryFile := strings.TrimSuffix(baseName, ext)

	err = os.MkdirAll(resolveTestDir("result"), os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("gagal membuat file output: %v", err)
	}

	outputFileName := fmt.Sprintf("%s-Voxelization_Result_Depth-%d.obj", queryFile, maxDepth)
	outputPath = filepath.Join(resolveTestDir("result"), outputFileName)	
	

	file, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("gagal membuat file output: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Mulai dari indeks 1
	vertexOffset := 1

	// Satu kubus terdiri dari 8 vertex dan 12 face
	for _, node := range leaves {
		cube := node.CurrentCube
		h := cube.Size / 2
		cx, cy, cz := cube.Center.X, cube.Center.Y, cube.Center.Z

		v1 := fmt.Sprintf("v %f %f %f\n", cx-h, cy-h, cz+h)
		v2 := fmt.Sprintf("v %f %f %f\n", cx+h, cy-h, cz+h)
		v3 := fmt.Sprintf("v %f %f %f\n", cx+h, cy+h, cz+h)
		v4 := fmt.Sprintf("v %f %f %f\n", cx-h, cy+h, cz+h)
		v5 := fmt.Sprintf("v %f %f %f\n", cx-h, cy-h, cz-h)
		v6 := fmt.Sprintf("v %f %f %f\n", cx+h, cy-h, cz-h)
		v7 := fmt.Sprintf("v %f %f %f\n", cx+h, cy+h, cz-h)
		v8 := fmt.Sprintf("v %f %f %f\n", cx-h, cy+h, cz-h)

		writer.WriteString(v1 + v2 + v3 + v4 + v5 + v6 + v7 + v8)

		writer.WriteString(fmt.Sprintf("f %d %d %d\n", vertexOffset, vertexOffset+1, vertexOffset+2))
		writer.WriteString(fmt.Sprintf("f %d %d %d\n", vertexOffset, vertexOffset+2, vertexOffset+3))
		
		writer.WriteString(fmt.Sprintf("f %d %d %d\n", vertexOffset+5, vertexOffset+4, vertexOffset+7))
		writer.WriteString(fmt.Sprintf("f %d %d %d\n", vertexOffset+5, vertexOffset+7, vertexOffset+6))
		
		writer.WriteString(fmt.Sprintf("f %d %d %d\n", vertexOffset+3, vertexOffset+2, vertexOffset+6))
		writer.WriteString(fmt.Sprintf("f %d %d %d\n", vertexOffset+3, vertexOffset+6, vertexOffset+7))
		
		writer.WriteString(fmt.Sprintf("f %d %d %d\n", vertexOffset+4, vertexOffset+5, vertexOffset+1))
		writer.WriteString(fmt.Sprintf("f %d %d %d\n", vertexOffset+4, vertexOffset+1, vertexOffset+0))
		
		writer.WriteString(fmt.Sprintf("f %d %d %d\n", vertexOffset+1, vertexOffset+5, vertexOffset+6))
		writer.WriteString(fmt.Sprintf("f %d %d %d\n", vertexOffset+1, vertexOffset+6, vertexOffset+2))
		
		writer.WriteString(fmt.Sprintf("f %d %d %d\n", vertexOffset+4, vertexOffset+0, vertexOffset+3))
		writer.WriteString(fmt.Sprintf("f %d %d %d\n", vertexOffset+4, vertexOffset+3, vertexOffset+7))

		// Langsung increment 8
		vertexOffset += 8
	}

	return outputPath,nil


}

func ValidateInput () (inputPath string, maxDepth int){
	for {
		fmt.Print("Masukan nama file yang ada di /test (misal : cow.obj) : ")
		n, _ := fmt.Scan(&inputPath)
		if (n != 1){
			fmt.Println("Error : Masukan nama file tidak valid")
			continue
		}
		if !strings.HasSuffix(strings.ToLower(inputPath), ".obj"){
			fmt.Println(("Error : File tidak berekstensi .obj"))
			continue
		}
		break

	}
	for {
		fmt.Print("Masukan max depth voxelization (rentang valid [1,12] inklusif) : ")
		n, _ := fmt.Scan(&maxDepth)
		if (n != 1){
			fmt.Println("Error : Masukan 1 angka dari [1,12]")
			continue
		}
		if maxDepth < 1 || maxDepth > 12{
			fmt.Println(("Error : Masukan hanya boleh dari rentang [1,12]"))
			continue
		}
		break

	}
	temp := inputPath
	inputPath = filepath.Join(resolveTestDir("test"), temp)
	 
	return inputPath, maxDepth
	
}

func ParseOBJ(inputPath string) (*models.Model, error){
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("Gagal membuka file %v", err)
	}
	defer file.Close()

	var model models.Model
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		
		parts := strings.Fields(line) 

		if len(parts) == 0 {
			continue
		}

		// Membaca vertex
		if parts[0] == "v" {
			if len(parts) != 4 {
				return nil, fmt.Errorf("Format vertex tidak valid di baris %d: %s", lineNum, line)
			}
			
			x, errX := strconv.ParseFloat(parts[1], 64)
			y, errY := strconv.ParseFloat(parts[2], 64)
			z, errZ := strconv.ParseFloat(parts[3], 64)
			
			if errX != nil || errY != nil || errZ != nil {
				return nil, fmt.Errorf("Tipe data koordinat tidak valid di baris %d: %s", lineNum, line)
			}
			
			model.Vertices = append(model.Vertices, models.Vertex{X: x, Y: y, Z: z})

		// Membaca face
		} else if parts[0] == "f" {
			if len(parts) != 4 {
				return nil, fmt.Errorf("Format face tidak valid di baris %d: %s", lineNum, line)
			}
			
			v1, err1 := strconv.Atoi(parts[1])
			v2, err2 := strconv.Atoi(parts[2])
			v3, err3 := strconv.Atoi(parts[3])
			
			if err1 != nil || err2 != nil || err3 != nil {
				return nil, fmt.Errorf("Tipe data indeks face tidak valid di baris %d: %s", lineNum, line)
			}
			
			model.Faces = append(model.Faces, models.Face{V1: v1, V2: v2, V3: v3})

		} else {
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error saat membaca isi file %v", err)
	}

	if len(model.Vertices) == 0 && len(model.Faces) == 0 {
		return nil, fmt.Errorf("File kosong sehingga tidak valid")
	}

	if len(model.Vertices) < 3 || len(model.Faces) == 0 {
		return nil, fmt.Errorf("File tidak valid! Setidaknya harus ada tiga vertex dan satu face")
	}

	for _, f := range model.Faces {
        if f.V1 < 1 || f.V1 > len(model.Vertices) ||
           f.V2 < 1 || f.V2 > len(model.Vertices) ||
           f.V3 < 1 || f.V3 > len(model.Vertices) {
            return nil, fmt.Errorf("File tidak valid: face memiliki indeks vertex yang tidak ada")
        }
    }

	return &model, nil

}


func resolveTestDir(types string) string {
    candidates := []string{}

    if (types == "test"){
		if exePath, err := os.Executable(); err == nil {
			exeDir := filepath.Dir(exePath)
			candidates = append(candidates,
				filepath.Join(exeDir, "test"),      // test di folder binary
				filepath.Join(exeDir, "..", "test"), // test di parent binary
			)
   		}

		if cwd, err := os.Getwd(); err == nil {
			candidates = append(candidates,
				filepath.Join(cwd, "test"),
				filepath.Join(cwd, "..", "test"), // saat go run dari src
			)
		}

		for _, dir := range candidates {
			if info, err := os.Stat(dir); err == nil && info.IsDir() {
				return dir
			}
    	}
	}

	if types == "result"{
		if exePath, err := os.Executable(); err == nil {
			exeDir := filepath.Dir(exePath)
			candidates = append(candidates,
				filepath.Join(exeDir, "result"),      // result di folder binary
				filepath.Join(exeDir, "..", "result"), // result di parent binary
			)
   		}

		if cwd, err := os.Getwd(); err == nil {
			candidates = append(candidates,
				filepath.Join(cwd, "result"),
				filepath.Join(cwd, "..", "result"), // saat go run dari src
			)
		}

		for _, dir := range candidates {
			if info, err := os.Stat(dir); err == nil && info.IsDir() {
				return dir
			}
    	}
	}

	
    return candidates[0]
}