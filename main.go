package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
)

// files of format [size, path, name]
var files = [][3]string{}

// collect files
// sort them
// filter the same-size files
// filter equal files from same-size files
func main() {
	fmt.Println("Scanning for files")
	//input of directory
	count("/home/jarr/Downloads")
	sort.SliceStable(files, func(i, j int) bool {
		a, _ := strconv.ParseInt(files[i][0], 10, 64)
		b, _ := strconv.ParseInt(files[j][0], 10, 64)
		return a < b
	})
	iterate_same(files)
}

// go recursively from input directory
// collect all files in the "files" slice
func count(directory string) {
	var folders = []string{}
	data, err := ioutil.ReadDir(directory)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range data {
		if !file.IsDir() {
			files = append(files, [3]string{strconv.FormatInt(file.Size(), 10), file.Name(), directory})
		} else {
			folders = append(folders, concatenate(directory, file.Name()))
		}
	}

	for _, folder := range folders {
		count(folder)
	}
}

// check the files with same size
func iterate_same(slice [][3]string) {
	for i := 0; i < len(slice); i++ {
		for j := i + 1; j < len(slice); j++ {
			if slice[i][0] == slice[j][0] {
				continue
			}
			start := i
			end := j
			i = j
			if len(slice[start:end]) > 1 {
				iterate_equal(slice[start:end])
			}
		}
	}
}

// put array of same size files paths
// check if files equal by md5 hash
// print size, hash, files paths
func iterate_equal(slice [][3]string) {
	var same = [][3]string{}
	if len(slice) < 2 {
		return
	}
	for num, val := range slice {
		same = append(same, val)
		for i := num + 1; i < len(slice); i++ {
			if check_hash(concatenate(val[2], val[1]), concatenate(slice[i][2], slice[i][1])) {
				same = append(same, slice[i])
				slice = remove(slice, i)
				i -= 1
			}
		}
		if len(same) < 2 {
			slice = remove(slice, 0)
			return
		}
		slice = remove(slice, 0)
		fmt.Println(same[0][0])
		plaintext, _ := ioutil.ReadFile(concatenate(same[0][2], same[0][1]))
		hash := md5.Sum([]byte(plaintext))
		fmt.Printf("%x\n", hash)
		for _, k := range same {
			fmt.Println(concatenate(k[2], k[1]))
		}
		fmt.Printf("\n")
		break
	}
	iterate_equal(slice)

}

// take md5 hash of two files from input paths and compare them
func check_hash(path1, path2 string) bool {
	plaintext1, _ := ioutil.ReadFile(path1)
	plaintext2, _ := ioutil.ReadFile(path2)
	hash1 := md5.Sum([]byte(plaintext1))
	hash2 := md5.Sum([]byte(plaintext2))
	return hash1 == hash2
}

// concatenate two parts of path
func concatenate(s1, s2 string) string {
	s3 := s1 + "/" + s2
	return s3
}

// remove element of an array by it's index
func remove(slice3 [][3]string, s int) [][3]string {
	return append(slice3[:s], slice3[s+1:]...)
}

// func that can save unused variables
func UNUSED(x ...interface{}) {}
