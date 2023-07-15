package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

// information about each file
type info struct {
	size int64
	path string
	name string
}

var files = []info{}

func main() {

	path, _ := os.UserHomeDir()
	size := 0.5

	fmt.Print("Enter the path (default: ~/): ")
	fmt.Scanln(&path)
	fmt.Print("Enter the minimal size in MB (default: 0.5): ")
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
	iterate_same(files)
}

// recursively go through every folder and collect info about every file to slice
func gather_files(directory string) {

	var folders = []string{}
	elements, _ := ioutil.ReadDir(directory)

	for _, element := range elements {

		if !element.IsDir() {
			// add file info to the files slice
			files = append(files, info{element.Size(), directory, element.Name()})
		} else {
			// add folder into folders slice
			folders = append(folders, mkpath(directory, element.Name()))
		}
	}

	// process remaining folders
	for _, folder := range folders {
		gather_files(folder)
	}
}

// process files with same size
func iterate_same(slice []info) {
	// for the item
	for i := 0; i < len(slice); i++ {
		// for the item next to item above
		for j := i + 1; j < len(slice); j++ {
			if slice[i].size == slice[j].size {
				continue
			}
			// now from i to j-1 there are same-size items
			start := i
			end := j
			// the last not same item index becomes start
			i = j
			// should be at least two items to be possible to be the equal
			if len(slice[start:end]) > 1 {
				iterate_equal(slice[start:end])
			}
		}
	}
}

func iterate_equal(same []info) {

	equal := []info{}

	// another "at least 2" check because of recursion
	if len(same) < 2 {
		return
	}

	// for cycle process only first item
	for num, val := range same {

		// adding first item to find it's clone
		equal = append(equal, val)

		for index := num + 1; index < len(same); index++ {

			if check_hash(mkpath(val.path, val.name), mkpath(same[index].path, same[index].name)) {
				equal = append(equal, same[index])
				// remove equal items to process other possible bunch of equal
				same = remove(same, index)
				// go to same index bc of removing above
				index -= 1
			}
		}

		// remove the first item, whether we found it's clone or not
		same = remove(same, 0)

		// if not, leave the cycle
		if len(equal) < 2 {
			break
		}

		// get size in MB
		size := float64(equal[0].size) / 1024.0 / 1024.0
		fmt.Printf("%.2f MB\n", size)

		plaintext, _ := ioutil.ReadFile(mkpath(equal[0].path, equal[0].name))
		hash := md5.Sum([]byte(plaintext))

		fmt.Printf("%x\n", hash)
		for _, file := range equal {
			fmt.Println(mkpath(file.path, file.name))
		}
		fmt.Printf("\n")

		break
	}
	iterate_equal(same)

}

// take md5 hash of two files from input paths and compare them
func check_hash(path1, path2 string) bool {
	plaintext1, _ := ioutil.ReadFile(path1)
	plaintext2, _ := ioutil.ReadFile(path2)
	hash1 := md5.Sum([]byte(plaintext1))
	hash2 := md5.Sum([]byte(plaintext2))
	return hash1 == hash2
}

// make path string from two parts of it
func mkpath(path, name string) string {
	fullpath := path + "/" + name
	return fullpath
}

// remove element of an array by it's index
func remove(slice3 []info, index int) []info {
	return append(slice3[:index], slice3[index+1:]...)
}

// func that can save unused variables
func UNUSED(x ...interface{}) {}
