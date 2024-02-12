package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/hamao0820/sudoku/detect"
	"gocv.io/x/gocv"
)

type SudokuID int

const (
	_ SudokuID = iota
	SudokuIDOne
	SudokuIDTwo
	SudokuIDThree
	SudokuIDFour
	SudokuIDFive
	SudokuIDSix
	SudokuIDSeven
	SudokuIDEight
	SudokuIDNine
	SudokuIDTen
)

var IDs = []SudokuID{
	SudokuIDOne,
	SudokuIDTwo,
	SudokuIDThree,
	SudokuIDFour,
	SudokuIDFive,
	SudokuIDSix,
	SudokuIDSeven,
	SudokuIDEight,
	SudokuIDNine,
	SudokuIDTen,
}

type Answer struct {
	Data Data `json:"data"`
}

type Data []struct {
	Id   SudokuID
	Name string
	Cell Cell
}

type Cell [][]int

var (
	data = loadJSON()
)

func main() {
	n := flag.Int("n", 0, "number of images")
	flag.Parse()
	if *n == 0 {
		panic("set: -n [number of images]")
	}
	for _, id := range IDs {
		fmt.Println("id: ", id)
		fmt.Println("name: ", getName(id))
		collectSquareFromCamera(id, *n)

		fmt.Printf("id: %d finished\n", id)
		fmt.Println("wait 3 seconds")
		time.Sleep(3 * time.Second)
	}

	collectCell()
}

func collectSquareFromCamera(id SudokuID, n int) {
	webcam, _ := gocv.OpenVideoCapture(0)
	webcam.Set(gocv.VideoCaptureFPS, 30)
	defer webcam.Close()
	window := gocv.NewWindow("collect")
	img := gocv.NewMat()
	defer img.Close()

	i := 0
	frame := 0
	for {
		time.Sleep(100 * time.Millisecond)
		webcam.Read(&img)

		origin := gocv.NewMat()
		img.CopyTo(&origin)
		if drawSquare(&img) {
			frame++
			if frame >= 2 {
				square, err := detect.GetSquare(origin)
				if err != nil {
					continue
				}
				i++
				fmt.Printf("i: %02d\n", i)
				dir := fmt.Sprintf("data/squares/%d", id)
				if _, err := os.Stat(dir); os.IsNotExist(err) {
					os.Mkdir(dir, 0755)
				}
				filename := filepath.Join(dir, fmt.Sprintf("%d.png", rand.Int31()))
				gocv.IMWrite(filename, square)
				frame = 0

				if i >= n {
					return
				}
			}
		} else {
			frame = 0
		}

		window.IMShow(img)
		window.WaitKey(1)
	}
}

func drawSquare(img *gocv.Mat) bool {
	detect.FitSize(img, 500, 500)

	gray := detect.ToGray(*img)
	defer gray.Close()

	edge := detect.FindEdge(gray)
	defer edge.Close()

	contours, _ := detect.FindContours(edge)
	defer contours.Close()

	min_area := float64(img.Rows()*img.Cols()) * 0.2
	largeContours := detect.FilterContours(contours, min_area)

	convexes := detect.GetConvexes(largeContours)

	polies := detect.GetPolygons(largeContours, convexes)
	defer polies.Close()

	// 正方形に近いものを選ぶ
	selectedIndex, _ := detect.SelectNearestSquareIndex(polies)
	if selectedIndex == -1 {
		return false
	}

	gocv.DrawContours(img, polies, selectedIndex, color.RGBA{255, 0, 0, 255}, 3)
	return true
}

func collectCell() {
	paths, err := filepath.Glob("data/squares/*")
	if err != nil {
		panic(err)
	}
	counts := make(map[int]int)
	for _, path := range paths {
		if !isDir(path) {
			continue
		}
		dir := filepath.Base(path)
		id_, err := strconv.Atoi(dir)
		if err != nil {
			panic(err)
		}
		ansCells, err := getCell(SudokuID(id_))
		if err != nil {
			fmt.Println(err)
			continue
		}
		images, err := filepath.Glob(filepath.Join(path, "*.png"))
		if err != nil {
			panic(err)
		}
		for _, image := range images {
			img := gocv.IMRead(image, gocv.IMReadColor)
			if img.Empty() {
				fmt.Println("Error reading image from: ", image)
				continue
			}
			defer img.Close()

			cells := splitCells(img)
			for y := 0; y < 9; y++ {
				for x := 0; x < 9; x++ {
					cell := cells[y][x]
					if cell.Empty() {
						continue
					}
					ans := ansCells[y][x]
					dir := fmt.Sprintf("data/cells/%d", ans)
					if _, err := os.Stat(dir); os.IsNotExist(err) {
						os.Mkdir(dir, 0755)
					}
					filename := filepath.Join(dir, fmt.Sprintf("%05d.png", counts[ans]+1))
					gocv.IMWrite(filename, cell)
					counts[ans]++
				}
			}
		}
	}
	for k, v := range counts {
		fmt.Printf("%d: %d\n", k, v)
	}
}

func splitCells(square gocv.Mat) [][]gocv.Mat {
	cells := make([][]gocv.Mat, 9)
	dx := float64(square.Cols()) / 9
	dy := float64(square.Rows()) / 9

	padding := 1.0
	for y := 0; y < 9; y++ {
		cells[y] = make([]gocv.Mat, 9)
		for x := 0; x < 9; x++ {
			cells[y][x] = square.Region(image.Rect(int(float64(x)*dx+padding), int(float64(y)*dy+padding), int(float64(x+1)*dx-padding), int(float64(y+1)*dy-padding)))
		}
	}
	return cells
}

func loadJSON() Data {
	f, err := os.Open("data.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var ans Answer
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&ans); err != nil {
		panic(err)
	}

	return ans.Data
}

func getName(id SudokuID) string {
	for _, d := range data {
		if d.Id == id {
			return d.Name
		}
	}
	return ""
}

func getCell(id SudokuID) (Cell, error) {
	for _, d := range data {
		if d.Id == id {
			return d.Cell, nil
		}
	}
	return [][]int{}, fmt.Errorf("not find")
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}
