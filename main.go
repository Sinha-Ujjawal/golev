package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const (
	defaultColor string = "\033[0m"
	red                 = "\033[0;31m"
	green               = "\033[0;32m"
)

type difftype = int

const (
	noop difftype = iota
	replace
	delete
	insert
)

//NOTE: change difftype to below for debugging
//type difftype = string
//
//const (
//	noop difftype = "noop"
//	replace       = "replace"
//	delete        = "delete"
//	insert        = "insert"
//)

type diff struct {
	difftype  difftype
	startIdx1 int
	startIdx2 int
	count     int
}

type diffs struct {
	diffArr []diff
	lines1  []string
	lines2  []string
}

func makeArray2D[T any](nrows int, ncols int) [][]T {
	ret := make([][]T, nrows)
	for row := 0; row < nrows; row++ {
		ret[row] = make([]T, ncols)
	}
	return ret
}

func appendDiff(diffArr []diff, row int, col int, difftype difftype) []diff {
	if len(diffArr) == 0 {
		return append(
			diffArr,
			diff{
				difftype:  difftype,
				startIdx1: row - 1,
				startIdx2: col - 1,
				count:     1,
			},
		)
	}
	lastDiff := diffArr[len(diffArr)-1]
	if lastDiff.difftype == difftype {
		lastDiff.count += 1
		lastDiff.startIdx1 = row - 1
		lastDiff.startIdx2 = col - 1
		diffArr[len(diffArr)-1] = lastDiff
		return diffArr
	} else {
		return append(
			diffArr,
			diff{
				difftype:  difftype,
				startIdx1: row - 1,
				startIdx2: col - 1,
				count:     1,
			},
		)
	}
}

func (diffs diffs) show() {
	lines1 := diffs.lines1
	lines2 := diffs.lines2
	for i := len(diffs.diffArr) - 1; i >= 0; i -= 1 {
		diff := diffs.diffArr[i]
		switch diff.difftype {
		case noop:
			for j := 0; j < diff.count; j++ {
				fmt.Printf("  %s\n", lines1[j+diff.startIdx1])
			}
		case replace:
			for j := 0; j < diff.count; j++ {
				fmt.Printf("%s- %s%s\n", red, lines1[j+diff.startIdx1], defaultColor)
				fmt.Printf("%s+ %s%s\n", green, lines2[j+diff.startIdx2], defaultColor)
			}
		case delete:
			for j := 0; j < diff.count; j++ {
				fmt.Printf("%s- %s%s\n", red, lines1[j+diff.startIdx1], defaultColor)
			}
		case insert:
			for j := 0; j < diff.count; j++ {
				fmt.Printf("%s+ %s%s\n", green, lines2[j+diff.startIdx2], defaultColor)
			}
		}
	}
}

func lavenshteinArr(lines1 []string, lines2 []string) [][]int {
	nrows := len(lines1) + 1
	ncols := len(lines2) + 1
	arr := makeArray2D[int](nrows, ncols)
	arr[0][0] = 0
	for row := 1; row < nrows; row++ {
		arr[row][0] = row
	}
	for col := 1; col < ncols; col++ {
		arr[0][col] = col
	}
	for row := 1; row < nrows; row++ {
		for col := 1; col < ncols; col++ {
			up := arr[row-1][col]
			left := arr[row][col-1]
			diag := arr[row-1][col-1]
			if lines1[row-1] == lines2[col-1] {
				arr[row][col] = diag
			} else {
				best := 1 + int(math.Min(math.Min(float64(up), float64(left)), float64(diag)))
				arr[row][col] = best
			}
		}
	}
	return arr
}

func lavenshteinDistance(lines1 []string, lines2 []string) (int, diffs) {
	arr := lavenshteinArr(lines1, lines2)
	var diffArr []diff
	row, col := len(lines1), len(lines2)
	dist := arr[row][col]
	for row > 0 || col > 0 {
		if row == 0 {
			diffArr = appendDiff(diffArr, row, col, insert)
			col -= 1
		} else if col == 0 {
			diffArr = appendDiff(diffArr, row, col, delete)
			row -= 1
		} else if lines1[row-1] == lines2[col-1] {
			diffArr = appendDiff(diffArr, row, col, noop)
			row -= 1
			col -= 1
		} else {
			left := arr[row][col-1]
			diag := arr[row-1][col-1]
			curr := arr[row][col]
			if curr == diag+1 {
				diffArr = appendDiff(diffArr, row, col, replace)
			} else if curr == left+1 {
				diffArr = appendDiff(diffArr, row, col, insert)
			} else {
				diffArr = appendDiff(diffArr, row, col, delete)
			}
			row -= 1
			col -= 1
		}
	}
	return dist, diffs{diffArr: diffArr, lines1: lines1, lines2: lines2}
}

func readEntireFileAsLines(path string) ([]string, error) {
	var lines []string
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("ERROR: Could not open file: %s for reading: %s\n", path, err)
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return lines, nil
}

func main() {
	var program string
	var file1Path string
	var file2Path string
	program, os.Args = os.Args[0], os.Args[1:]
	if len(os.Args) == 0 {
		fmt.Println("ERROR: <file1path> not provided!.")
		fmt.Printf("Usage: %s <file1path> <file2path>\n", program)
		os.Exit(1)
	}
	file1Path, os.Args = os.Args[0], os.Args[1:]
	if len(os.Args) == 0 {
		fmt.Println("ERROR: <file2path> not provided!.")
		fmt.Printf("Usage: %s <file1path> <file2path>\n", program)
		os.Exit(1)
	}
	file2Path, os.Args = os.Args[0], os.Args[1:]
	lines1, err := readEntireFileAsLines(file1Path)
	if err != nil {
		os.Exit(1)
	}
	lines2, err := readEntireFileAsLines(file2Path)
	if err != nil {
		os.Exit(1)
	}
	_, diffs := lavenshteinDistance(lines1, lines2)
	diffs.show()
}
