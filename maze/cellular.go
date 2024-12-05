package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

func countNeighborsWall(maze [][]int, x, y int) int {
	count := 0
	size := len(maze)
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			checkx := x + i
			checky := y + j
			if i == 0 && j == 0 {
				continue
			}
			if checkx < 0 || checky < 0 || checkx >= size || checky >= size {
				count++
			} else {
				count += maze[checkx][checky]
			}
		}
	}
	return count
}

func cellularMaze(maze [][]int) {
	for i := range maze {
		for j := range maze[i] {
			count := countNeighborsWall(maze, i, j)
			if maze[i][j] == 1 {
				if count < 4 {
					maze[i][j] = 0
				}
			} else {
				if count > 5 {
					maze[i][j] = 1
				}
			}
		}
	}
}

func cellularHandler(w http.ResponseWriter, req *http.Request) {

	size := 55
	maze := make([][]int, size)
	for i := range maze {
		maze[i] = make([]int, size)
		for j := range maze[i] {
			// 将随机生成改为45%的概率生成墙
			if rand.Float64() < 0.550 {
				maze[i][j] = 1
			} else {
				maze[i][j] = 0
			}
		}
	}

	fmt.Fprint(w, htmlTemplate)

	fmt.Fprint(w, "<div class='maze-box'><h3>原始迷宫</h3>")
	renderCellularMaze(w, maze)
	fmt.Fprint(w, "</div>")

	fmt.Fprint(w, "<div class='maze-box'><h3>自动一次</h3>")
	cellularMaze(maze)
	renderCellularMaze(w, maze)
	fmt.Fprint(w, "</div>")

	fmt.Fprint(w, "<div class='maze-box'><h3>自动5次</h3>")
	for i := 0; i < 4; i++ {
		cellularMaze(maze)
	}
	renderCellularMaze(w, maze)
	fmt.Fprint(w, "</div>")

}

func renderCellularMaze(w http.ResponseWriter, maze [][]int) {
	size := len(maze)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<pre style='line-height:1;font-family:monospace'>")

	for i := 0; i < size+2; i++ {
		fmt.Fprintf(w, "&#9608;")
	}
	fmt.Fprintf(w, "\n")

	for i := 0; i < size; i++ {
		fmt.Fprintf(w, "&#9608;")
		for j := 0; j < size; j++ {
			if maze[i][j] == 1 {
				fmt.Fprintf(w, "&#9608;")
			} else {
				if i == 0 && j == 0 {
					fmt.Fprintf(w, "S ")
				} else if i == size-1 && j == size-1 {
					fmt.Fprintf(w, " E")
				} else {
					fmt.Fprintf(w, "  ")
				}
			}
		}
		fmt.Fprintf(w, "&#9608;\n")
	}

	for i := 0; i < size+2; i++ {
		fmt.Fprintf(w, "&#9608;")
	}
	fmt.Fprintf(w, "\n</pre>")
}
