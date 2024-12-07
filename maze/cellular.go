package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

// 默认参数常量
const (
	cellularDefaultSize = 55
	defaultProbability  = 0.55
	defaultIterations   = 5

	cellularMaxSize = 1000
	maxIterations   = 20
	minEntranceSize = 3
)

// MazeParams 存储迷宫生成的配置参数
// size: 迷宫大小 (3-100)
// probability: 障碍物生成概率 (0.0-1.0)
// iterations: 细胞自动机迭代次数 (1-20)
type MazeParams struct {
	size        int     // 迷宫尺寸
	probability float64 // 障碍物生成概率
	iterations  int     // 迭代次数
}

// 从请求中解析参数
func parseCellularParams(req *http.Request) MazeParams {
	params := MazeParams{
		size:        cellularDefaultSize,
		probability: defaultProbability,
		iterations:  defaultIterations,
	}

	if sizeStr := req.URL.Query().Get("size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil {
			if s <= 0 {
				log.Printf("Size too small, using default: %d", cellularDefaultSize)
			} else if s > cellularMaxSize {
				log.Printf("Size too large, using max: %d", cellularMaxSize)
				params.size = cellularMaxSize
			} else {
				params.size = s
			}
		}
	}

	if probStr := req.URL.Query().Get("probability"); probStr != "" {
		if p, err := strconv.ParseFloat(probStr, 64); err == nil && p >= 0 && p <= 1 {
			params.probability = p
		}
	}

	if iterStr := req.URL.Query().Get("iterations"); iterStr != "" {
		if iter, err := strconv.Atoi(iterStr); err == nil && iter > 0 && iter <= maxIterations {
			params.iterations = iter
		}
	}

	return params
}

// 初始化迷宫
func initializeMaze(size int, probability float64) [][]int {
	maze := make([][]int, size)
	row := make([]int, size*size) // 一次性分配所有内存
	for i := range maze {
		maze[i] = row[i*size : (i+1)*size]
	}
	for i := range maze {
		for j := range maze[i] {
			if rand.Float64() < probability {
				maze[i][j] = 1
			}
		}
	}

	// 设置入口和出口区域
	setEntranceArea(maze, 0, 0)                                       // 左上角入口
	setEntranceArea(maze, size-minEntranceSize, size-minEntranceSize) // 右下角出口

	return maze
}

// 设置入口区域
func setEntranceArea(maze [][]int, startX, startY int) {
	for i := startX; i < startX+minEntranceSize && i < len(maze); i++ {
		for j := startY; j < startY+minEntranceSize && j < len(maze); j++ {
			maze[i][j] = 0
		}
	}
}

const cellularCSS = `
<style>
    .container {
        display: flex;
        flex-direction: row;
        gap: 20px;
        padding: 20px;
    }
    .maze-box {
        flex: 1;
        border: 1px solid #ccc;
        padding: 10px;
        text-align: center;
    }
</style>`

func cellularHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, cellularCSS)
	fmt.Fprint(w, "<div class='container'>")

	params := parseCellularParams(req)
	maze := initializeMaze(params.size, params.probability)

	if params.size < 160 {
		renderMazeWithTitle(w, maze, fmt.Sprintf("随机迷宫，障碍物率：%d%%", int(params.probability*100)))
	}

	for i := 0; i < params.iterations; i++ {
		cellularMaze(maze)
	}

	if params.size < 260 {
		renderMazeWithTitle(w, maze, fmt.Sprintf("细胞自动机迭代：%d次", params.iterations))
	}

	connectRegionsByBFS(maze)
	renderMazeWithTitle(w, maze, "通过BFS连接所有区域")

	fmt.Fprint(w, "</div>") // 关闭 container div
}

// 渲染带标题的迷宫
func renderMazeWithTitle(w http.ResponseWriter, maze [][]int, title string) {
	fmt.Fprintf(w, "<div class='maze-box'><h3>%s</h3>", title)
	renderCellularMaze(w, maze)
	fmt.Fprint(w, "</div>")
}

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
	size := len(maze)
	// 将左上角3x3区域设置为通道
	for i := 0; i <= 3; i++ {
		for j := 0; j <= 3; j++ {
			maze[i][j] = 0
		}
	}
	// 将右下角3x3区域设置为通道
	for i := size - 4; i <= size-1; i++ {
		for j := size - 4; j <= size-1; j++ {
			maze[i][j] = 0
		}
	}

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

func renderCellularMaze(w http.ResponseWriter, maze [][]int) {
	size := len(maze)
	w.Header().Set("Content-Type", "text/html")

	// 添加表格样式
	fmt.Fprintf(w, `<style>
		table.maze-table { 
			border-collapse: collapse;
			margin: 0 auto;
		}
		table.maze-table td {
			width: 5px;
			height: 5px;
			padding: 0;
		}
		.wall { background-color: #000; }
		.path { background-color: #fff; }
		.start { background-color: #0f0; }
		.end { background-color: #f00; }
	</style>`)

	fmt.Fprintf(w, "<table class='maze-table'>")

	// 渲染顶部边框
	fmt.Fprintf(w, "<tr>")
	for i := 0; i < size+2; i++ {
		fmt.Fprintf(w, "<td class='wall'></td>")
	}
	fmt.Fprintf(w, "</tr>")

	// 渲染迷宫主体
	for i := 0; i < size; i++ {
		fmt.Fprintf(w, "<tr><td class='wall'></td>") // 左边框
		for j := 0; j < size; j++ {
			if maze[i][j] == 1 {
				fmt.Fprintf(w, "<td class='wall'></td>")
			} else if i == 0 && j == 0 {
				fmt.Fprintf(w, "<td class='start'></td>") // 起点
			} else if i == size-1 && j == size-1 {
				fmt.Fprintf(w, "<td class='end'></td>") // 终点
			} else {
				fmt.Fprintf(w, "<td class='path'></td>")
			}
		}
		fmt.Fprintf(w, "<td class='wall'></td></tr>") // 右边框
	}

	// 渲染底部边框
	fmt.Fprintf(w, "<tr>")
	for i := 0; i < size+2; i++ {
		fmt.Fprintf(w, "<td class='wall'></td>")
	}
	fmt.Fprintf(w, "</tr>")

	fmt.Fprintf(w, "</table>")
}

func renderCMaze(w http.ResponseWriter, connectedMaze [][]int) {
	size := len(connectedMaze)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<pre style='line-height:1;font-family:monospace'>")

	for i := 0; i < size+2; i++ {
		fmt.Fprintf(w, "&#9608;")
	}
	fmt.Fprintf(w, "\n")

	for i := 0; i < size; i++ {
		fmt.Fprintf(w, "&#9608;")
		for j := 0; j < size; j++ {
			fmt.Fprintf(w, "%2d", connectedMaze[i][j])
		}
		fmt.Fprintf(w, "&#9608;\n")
	}

	for i := 0; i < size+2; i++ {
		fmt.Fprintf(w, "&#9608;")
	}
	fmt.Fprintf(w, "\n</pre>")
}

// 标记连通区域的函数
func markConnectedRegions(maze [][]int) [][]int {
	size := len(maze)
	// 创建新的二维数组用于标记
	regions := make([][]int, size)
	for i := range regions {
		regions[i] = make([]int, size)
	}

	// 区域编号从1开始
	currentRegion := 1

	// 遍历所有格子
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			// 如果是可通行区域(0)且还未被标记
			if maze[i][j] == 0 && regions[i][j] == 0 {
				// 使用DFS标记整个连通区域
				dfs(maze, regions, i, j, currentRegion)
				currentRegion++
			}
		}
	}

	return regions
}

// DFS辅助函数
func dfs(maze [][]int, regions [][]int, x, y, region int) {
	size := len(maze)

	// 检查边界和是否可访问
	if x < 0 || x >= size || y < 0 || y >= size ||
		maze[x][y] == 1 || regions[x][y] != 0 {
		return
	}

	// 标记当前格子
	regions[x][y] = region

	// 访问四个相邻格子
	dfs(maze, regions, x+1, y, region) // 下
	dfs(maze, regions, x-1, y, region) // 上
	dfs(maze, regions, x, y+1, region) // 右
	dfs(maze, regions, x, y-1, region) // 左
}

// 通过BFS连接所有区域的函数
func connectRegionsByBFS(maze [][]int) {
	size := len(maze)
	regions := markConnectedRegions(maze)

	// 获取所有不同的区域编号
	regionNums := make(map[int]bool)
	for i := range regions {
		for j := range regions[i] {
			if regions[i][j] > 0 {
				regionNums[regions[i][j]] = true
			}
		}
	}

	// 如果只有一个区域或没有区域，直接返回
	if len(regionNums) <= 1 {
		return
	}

	// 将区域编号转换为切片
	regionList := make([]int, 0)
	for num := range regionNums {
		regionList = append(regionList, num)
	}

	// 当还有多个区域时，继续连接
	for len(regionList) > 1 {
		// 对当前第一��区域进行BFS寻路
		region1 := regionList[0]
		path := bfsToNearestRegion(maze, regions, region1)

		if path != nil {
			// 获取连接到的目标区域编号
			targetRegion := regions[path[len(path)-1][0]][path[len(path)-1][1]]

			// 打通路径
			for _, pos := range path {
				maze[pos[0]][pos[1]] = 0
			}

			// 更新regions数组，将targetRegion合并到region1
			for i := 0; i < size; i++ {
				for j := 0; j < size; j++ {
					if regions[i][j] == targetRegion {
						regions[i][j] = region1
					}
				}
			}

			// 从列表中移除已合并的区域
			for i := 0; i < len(regionList); i++ {
				if regionList[i] == targetRegion {
					regionList = append(regionList[:i], regionList[i+1:]...)
					break
				}
			}
		}
	}
}

// 使用BFS寻找到最近的其他区域的路径
func bfsToNearestRegion(maze [][]int, regions [][]int, sourceRegion int) [][2]int {
	size := len(maze)
	visited := make([][]bool, size)
	for i := range visited {
		visited[i] = make([]bool, size)
	}

	// 使用队列存储待访问的点
	queue := [][2]int{}
	parent := make(map[[2]int][2]int)

	// 找到源区域的所有边界点作为起点
	for i := range regions {
		for j := range regions[i] {
			if regions[i][j] == sourceRegion && hasAdjacentWall(regions, i, j) {
				queue = append(queue, [2]int{i, j})
				visited[i][j] = true
			}
		}
	}

	// 定义方向：上、右、下、左
	dirs := [][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

	// BFS搜索
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// 检查四个方向
		for _, dir := range dirs {
			next := [2]int{current[0] + dir[0], current[1] + dir[1]}

			if next[0] >= 0 && next[0] < size &&
				next[1] >= 0 && next[1] < size &&
				!visited[next[0]][next[1]] {

				visited[next[0]][next[1]] = true
				parent[next] = current

				// 如果找到了另一个区域
				if regions[next[0]][next[1]] > 0 &&
					regions[next[0]][next[1]] != sourceRegion {
					// 重建路径
					path := [][2]int{}
					currentPos := next
					for {
						path = append([][2]int{currentPos}, path...)
						if regions[currentPos[0]][currentPos[1]] == sourceRegion {
							break
						}
						currentPos = parent[currentPos]
					}
					return path
				}

				queue = append(queue, next)
			}
		}
	}

	return nil
}

// 检查是否有相邻的墙
func hasAdjacentWall(regions [][]int, x, y int) bool {
	size := len(regions)
	dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for _, dir := range dirs {
		newX, newY := x+dir[0], y+dir[1]
		if newX >= 0 && newX < size && newY >= 0 && newY < size {
			if regions[newX][newY] == 0 { // 0表示墙
				return true
			}
		}
	}
	return false
}
