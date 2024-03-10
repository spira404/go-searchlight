package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"io/ioutil"
	"path/filepath"
	"github.com/zeebo/xxh3"
)

// information about each file
type info struct {
	size int64
	path string
	name string
	hash string
}

var files = []info{}

func main() {

	home, _ := os.UserHomeDir()
	path := home
	fmt.Print("Enter the path (default: ~/): ")
	fmt.Scanln(&path)

	if strings.Contains(path, "~"){
		path = strings.Replace(path, "~", home, -1)
	}

	size := 0.0
	fmt.Print("Enter the minimal size in MB (default: 0.0): ")
	fmt.Scanln(&size)

	gather_files(path)

	// sort slice by size in bits to bigger
	sort.SliceStable(files, func(i, j int) bool {
		size1 := files[i].size
		size2 := files[j].size
		return size1 < size2
	})

	// exclude first info structs with less size
	for num, val := range files {
		if float64(val.size) >= size*1024*1024 {
			files = files[num:]
			break
		}
	}

	// print out equal files
	process_same(files)
}

// recursively go through every folder and collect info about every file to slice
func gather_files(directory string) {

	var folders = []string{}
	elements, _ := ioutil.ReadDir(directory)

	for _, element := range elements {

		if !element.IsDir() {
			// add file info to the files slice
			files = append(files, info{element.Size(), directory, element.Name(), ""})
		} else {
			// add folder into folders slice
			folders = append(folders, filepath.Join(directory, element.Name()))
		}
	}

	// process remaining folders
	for _, folder := range folders {
		gather_files(folder)
	}
}

// process files with same size
func process_same(slice []info) {
	// for the item
	for i := 0; i < len(slice); i++ {
		// for the item next to item above
		for j := i + 1; j < len(slice); j++ {
			if slice[i].size != slice[j].size {
				// now from i to j-1 there are same-size items
				start := i
				end := j
				// the last not same, item index becomes start
				i = j
				// should be at least two items to be possible to be the equal
				if len(slice[start:end]) > 1 {
					process_equal(slice[start:end])
				}
				// in case of the last items in slice are same, this used, otherwise they're skipping
			} else if j == len(slice)-1 {
				start := i
				end := j
				i = j
				if len(slice[start:end+1]) > 1 {
					process_equal(slice[start : end+1])
				}
			}

		}
	}
}

func process_equal(same []info) {

	for index, _ := range same {
		same[index].hash = mkhash(filepath.Join(same[index].path, same[index].name))
	}

	sort.SliceStable(same, func(i, j int) bool {
		size1 := same[i].hash
		size2 := same[j].hash
		return size1 < size2
	})

	for i := 0; i < len(same); i++ {
		for j := i + 1; j < len(same); j++ {

			if same[i].hash != same[j].hash {

				start := i
				end := j
				i = j

				if len(same[start:end]) > 1 {
					equal := same[start:end]
					size := float64(equal[0].size) / 1024.0 / 1024.0
					fmt.Printf("%.2f MB\n", size)
					fmt.Println(same[0].hash)
					for _, file := range equal {
						fmt.Println(filepath.Join(file.path, file.name))
					}
					fmt.Printf("\n")
				}
				// in case of the last items in slice are same, this used, otherwise they're skipping
			} else if j == len(same)-1 {
				start := i
				end := j
				i = j
				if len(same[start:end+1]) > 1 {
					equal := same[start : end+1]
					size := float64(equal[0].size) / 1024.0 / 1024.0
					fmt.Printf("%.2f MB\n", size)
					fmt.Println(same[0].hash)
					for _, file := range equal {
						fmt.Println(filepath.Join(file.path, file.name))
						
					}
					fmt.Printf("\n")
				}
			}
		}
	}
}

// take md5 hash of two files from input paths and compare them
func mkhash(path string) string {
	plaintext, _ := ioutil.ReadFile(path)
	hexhash := xxh3.Hash([]byte(plaintext))
	hash := fmt.Sprintf("%x", hexhash)
	return hash
}

// remove element of an array by it's index
func remove(slice3 []info, index int) []info {
	return append(slice3[:index], slice3[index+1:]...)
}

// func that can save unused variables
func UNUSED(x ...interface{}) {}
