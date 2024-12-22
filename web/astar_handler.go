package main

import (
	"fmt"
	"net/http"

	"mazemap/pathfind"
	"mazemap/tiledmap"
)

func astarHandler(w http.ResponseWriter, req *http.Request) {
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
	//renderMazePathWithTitle(w, maze, path, "原始迷宫")

	tiledmap.AccuMaze(maze, path, accRatio)
	tiledmap.ErosionMaze(maze, erosionRatio)

	//renderMazeWithTitle(w, maze, "堆积侵蚀后") // 第三个画布：侵蚀后

	start := [2]int{0, 0}
	end := [2]int{size - 1, size - 1}

	// 使用 A* 寻路
	pathFindRes := pathfind.FindPathAStar(maze, start, end)
	renderPathWithTitle(w, maze, pathFindRes, "A*寻路结果") // 渲染带路径的迷宫

	// 使用 dijkstra 寻路
	pathFindRes = pathfind.FindPathDijkstra(maze, start, end)
	renderPathWithTitle(w, maze, pathFindRes, "Dijkstra寻路结果") // 渲染带路径的迷宫

	// 使用 bestfirst 寻路
	pathFindRes = pathfind.FindPathBestFirst(maze, start, end)
	renderPathWithTitle(w, maze, pathFindRes, "BestFirst寻路结果") // 渲染带路径的迷宫

	fmt.Fprint(w, "\n</div></div></body></html>")
}

func renderPathWithTitle(w http.ResponseWriter, maze [][]int, res pathfind.PathFindResult, title string) {

	title += fmt.Sprintf(" (成本: %d, 长度: %d)", res.Cost, len(res.Path))

	path := res.Path
	// 将路径转换为map以便快速查找
	size := len(maze)
	pathArr := make([][]bool, size)

	//fmt.Println(size)
	for i := range maze {
		sizew := len(maze[0])
		//fmt.Println(sizew)
		pathArr[i] = make([]bool, sizew)
	}

	for _, p := range path {
		pathArr[p[0]][p[1]] = true
	}

	renderMazePathWithTitle(w, maze, pathArr, title)
}
