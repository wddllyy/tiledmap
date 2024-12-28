package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"mazemap/tiledmap"
)

func cellularHandler(w http.ResponseWriter, req *http.Request) {
	printHtmlHead(w, "细胞自动机")

	params := parseCellularParams(req)

	fmt.Fprintf(w, `
<div class="all-container">
	<div class="all-controls">
		<form>
			尺寸: <input type="number" name="size" value="%d" min="13" max="1000" step="2">
			障碍物率: <input type="number" name="prob" value="%0.1f" step="0.1" min="0" max="1">
			迭代次数: <input type="number" name="iter" value="%d" step="0.1" min="0" max="20">
			<input type="submit" value="生成">
		</form>
	</div>
	<div style="display: flex; gap: 20px; justify-content: center;">`,
		params.Size, params.Probability, params.Iterations)

	maze := tiledmap.InitializeMaze(params.Size, params.Probability)

	if params.Size < 160 {
		renderMazeWithTitle(w, maze, fmt.Sprintf("随机迷宫，障碍物率：%d%%", int(params.Probability*100)))
	}

	for i := 0; i < params.Iterations; i++ {
		tiledmap.CellularMaze(maze)
	}

	if params.Size < 260 {
		renderMazeWithTitle(w, maze, fmt.Sprintf("细胞自动机迭代：%d次", params.Iterations))
	}

	tiledmap.ConnectRegionsByBFS(maze)
	renderMazeWithTitle(w, maze, "BFS连接所有区域")

	fmt.Fprint(w, "\n</div></div></body></html>")
}

// 从请求中解析参数
func parseCellularParams(req *http.Request) tiledmap.MazeParams {
	params := tiledmap.MazeParams{
		Size:        tiledmap.CellularDefaultSize,
		Probability: tiledmap.DefaultProbability,
		Iterations:  tiledmap.DefaultIterations,
	}

	if sizeStr := req.URL.Query().Get("size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil {
			if s <= 0 {
				log.Printf("Size too small, using default: %d", tiledmap.CellularDefaultSize)
			} else if s > tiledmap.CellularMaxSize {
				log.Printf("Size too large, using max: %d", tiledmap.CellularMaxSize)
				params.Size = tiledmap.CellularMaxSize
			} else {
				params.Size = s
			}
		}
	}

	if probStr := req.URL.Query().Get("probability"); probStr != "" {
		if p, err := strconv.ParseFloat(probStr, 64); err == nil && p >= 0 && p <= 1 {
			params.Probability = p
		}
	}

	if iterStr := req.URL.Query().Get("iterations"); iterStr != "" {
		if iter, err := strconv.Atoi(iterStr); err == nil && iter > 0 && iter <= tiledmap.MaxIterations {
			params.Iterations = iter
		}
	}

	return params
}
