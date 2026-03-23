package src

import (
	"fmt"
	"os"
	"strings"
)

func main () {
	inputPath, maxDepth := inputValidation()

}

func inputValidation () (inputPath string, maxDepth int ){

	for {
		fmt.Print("Masukan nama file .obj: ")
		fmt.Scanln(&inputPath)
		if !strings.HasSuffix(strings.ToLower((inputPath)), ".obj"){
			fmt.Println("Error : File harus berekstensi .obj!")
			continue
		}
		if _,err := os.Stat(inputPath) ; os.IsNotExist(err) {
			fmt.Println("Error : File tidak ada atau tidak valid!")
			continue
		}

		break
	}
	for {
		fmt.Print("Masukan depth voxelization: ")
		n, err := fmt.Scanln(&maxDepth)

		if err != nil || n != 1 {
			fmt.Println("Error: input harus angka bulat.")
			var buang string
			fmt.Scanln(&buang) 
			continue
		}
		
		if maxDepth > 12 || maxDepth < 1{
			fmt.Println("Error : Depth harus berada pada rentang [1,12] inklusif!")
			continue
		}

		break
	}
	return inputPath, maxDepth
}