package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"mazemap/pathfind"
	"mazemap/tiledmap"
)

func astarHandler(w http.ResponseWriter, req *http.Request) {
	printHtmlHead(w, "迷宫寻路演示", true)

	size, turnProb, accRatio, erosionRatio := parseMazeParams(req)

	// 控制表单
	fmt.Fprintf(w, `
<div class="all-container">
	<div class="all-controls">
		<form>
			尺寸: <input type="number" name="size" value="%d" min="5" max="99" step="2">
			转弯概率: <input type="number" name="turn" value="%0.1f" step="0.1" min="0" max="1">
			堆积系数: <input type="number" name="acc" value="%0.1f" step="0.1" min="0" max="1">
			侵蚀系数: <input type="number" name="erosion" value="%0.1f" step="0.1" min="0" max="1">
			<input type="submit" value="生成">
		</form>
	</div>
	<div class="playback-controls">
		<button onclick="togglePlayback()" id="playback-btn">播放</button>
		<input type="range" min="50" max="1000" value="200" 
			   onchange="updateSpeed(this.value)" id="speed-control">
		<span>更新间隔: <span id="speed-value">200</span>ms</span>
		<button onclick="stepPlayback()" id="playback-btn">Step</button>
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

	// 使用 jps 寻路
	pathFindRes = pathfind.FindPathJPS(maze, start, end)
	renderPathWithTitle(w, maze, pathFindRes, "JPS寻路结果") // 渲染带路径的迷宫

	// 使用 JPS+ 寻路
	preprocessedMaze := pathfind.PreprocessMaze(maze)          // 预处理迷宫
	pathFindRes = preprocessedMaze.FindPathJPSPlus(start, end) // 调用 JPS+ 寻路
	renderPathWithTitle(w, maze, pathFindRes, "JPS+寻路结果")

	fmt.Fprint(w, "\n</div></div></body></html>")
}

func renderPathWithTitle(w http.ResponseWriter, maze [][]int, res pathfind.PathFindResult, title string) {

	info := fmt.Sprintf(" (成本:%d,检查:%d,长度:%d)", res.Cost, res.Check, len(res.Path))

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

	// 只保留步骤数据
	stepsJSON, _ := json.Marshal(res.StepRecord)
	fmt.Fprintf(w, `
	<script>
	if (!window.allStepsData) window.allStepsData = {};
	window.allStepsData["%s"] = %s;
	//console.log("加载步骤数据:", "%s", window.allStepsData["%s"]);
	</script>
	`, title, string(stepsJSON), title, title)

	renderMazePathWithTitle(w, maze, pathArr, title, info)
}
