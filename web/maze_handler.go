package main

import (
	"fmt"
	"net/http"
	"strconv"

	"mazemap/tiledmap"
)

const (
	defaultSize         = 19
	defaultTurnProb     = 0.4
	defaultAccRatio     = 0.7
	defaultErosionRatio = 0.5
	maxSize             = 100
)

func parseMazeParams(req *http.Request) (size int, turnProb, accRatio, erosionRatio float64) {
	size = defaultSize
	turnProb = defaultTurnProb
	accRatio = defaultAccRatio
	erosionRatio = defaultErosionRatio

	if sizeStr := req.URL.Query().Get("size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 && s < maxSize {
			size = s
		}
	}

	if turnStr := req.URL.Query().Get("turn"); turnStr != "" {
		if t, err := strconv.ParseFloat(turnStr, 64); err == nil && t >= 0 && t <= 1 {
			turnProb = t
		}
	}

	if accStr := req.URL.Query().Get("acc"); accStr != "" {
		if acc, err := strconv.ParseFloat(accStr, 64); err == nil && acc >= 0 && acc <= 1 {
			accRatio = (acc)
		}
	}

	if erosionStr := req.URL.Query().Get("erosion"); erosionStr != "" {
		if e, err := strconv.ParseFloat(erosionStr, 64); err == nil && e >= 0 && e <= 1 {
			erosionRatio = e
		}
	}

	return
}

func mazeHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, RenderCSS)

	size, turnProb, accRatio, erosionRatio := parseMazeParams(req)

	// 控制表单
	fmt.Fprintf(w, `
		<div class="dungeon-container">
			<div class="dungeon-controls">
				<form>
					尺寸: <input type="number" name="size" value="%d" min="13" max="99" step="2">
					转弯概率: <input type="number" name="turn" value="%0.1f" step="0.1" min="0" max="1">
					堆积系数: <input type="number" name="acc" value="%0.1f" step="0.1" min="0" max="1">
					侵蚀系数: <input type="number" name="erosion" value="%0.1f" step="0.1" min="0" max="1">
					<input type="submit" value="生成">
				</form>
			</div>
	`, size, turnProb, accRatio, erosionRatio)

	fmt.Fprintf(w, `<div style="display: flex; gap: 20px; justify-content: center;">`)
	// 生成迷宫和寻找路径
	maze := tiledmap.GenerateMaze(size, turnProb)
	path := tiledmap.FindPath(maze)

	// 第一个画布：原始迷宫

	fmt.Fprint(w, "<div class='maze-box'><h3>原始迷宫</h3>")
	renderMaze(w, maze, path, true)
	fmt.Fprint(w, "</div>")

	// 第二个画布：消除断头路后
	fmt.Fprint(w, "<div class='maze-box'><h3>堆积后</h3>")
	tiledmap.AccuMaze(maze, path, accRatio)
	renderMaze(w, maze, path, true)
	fmt.Fprint(w, "</div>")

	// 第三个画布：侵蚀后
	fmt.Fprint(w, "<div class='maze-box'><h3>侵蚀后</h3>")
	tiledmap.ErosionMaze(maze, erosionRatio)
	renderMaze(w, maze, path, false)
	fmt.Fprint(w, "</div>")

	// 第四、五、六个画布：空白位置
	fmt.Fprint(w, "<div class='maze-box'></div>")
	fmt.Fprint(w, "<div class='maze-box'></div>")
	fmt.Fprint(w, "<div class='maze-box'></div>")

	// 结束 HTML
	fmt.Fprint(w, "</div></body></html>")
}

func renderMaze(w http.ResponseWriter, maze [][]int, path [][]bool, showPath bool) {
	size := len(maze)

	fmt.Fprintf(w, `
		<div class="dungeon-grid" style="grid-template-columns: repeat(%d, 8px);">`, size+2)

	for i := 0; i < size+2; i++ {
		fmt.Fprintf(w, `<div class="dungeon-cell wall"></div>`)
	}

	for y := 0; y < size; y++ {
		fmt.Fprintf(w, `<div class="dungeon-cell wall"></div>`)
		for x := 0; x < size; x++ {
			cellClass := "wall"
			if maze[y][x] == 0 {
				cellClass = "floor"
			}
			if showPath && path[y][x] {
				cellClass = "path"
			}
			fmt.Fprintf(w, `<div class="dungeon-cell %s"></div>`, cellClass)
		}
		fmt.Fprintf(w, `<div class="dungeon-cell wall"></div>`)
	}

	for i := 0; i < size+2; i++ {
		fmt.Fprintf(w, `<div class="dungeon-cell wall"></div>`)
	}
	fmt.Fprintf(w, `</div>`)
}
