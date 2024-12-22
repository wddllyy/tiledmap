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
	<div style="display: flex; gap: 20px; justify-content: center;">`,
		size, turnProb, accRatio, erosionRatio)
	// 生成迷宫和寻找路径
	maze := tiledmap.GenerateMaze(size, turnProb)
	path := tiledmap.FindPath(maze)

	// 第一个画布：原始迷宫
	renderMazePathWithTitle(w, maze, path, "原始迷宫", "")

	tiledmap.AccuMaze(maze, path, accRatio)
	renderMazePathWithTitle(w, maze, path, "堆积后", "") // 第二个画布：消除断头路后

	tiledmap.ErosionMaze(maze, erosionRatio)
	renderMazeWithTitle(w, maze, "侵蚀后") // 第三个画布：侵蚀后

	// 结束 HTML
	fmt.Fprint(w, "\n</div></div></body></html>")
}
