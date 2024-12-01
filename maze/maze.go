package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
)

type Point struct {
	x, y int
}

const (
	defaultSize         = 19
	defaultTurnProb     = 0.4
	defaultAccRatio     = 0.7
	defaultErosionRatio = 0.5
	maxSize             = 100
)

const htmlTemplate = `
<html>
<head>
    <style>
        .container {
            display: grid;
            grid-template-columns: repeat(3, 1fr);
            grid-template-rows: repeat(2, auto);
            gap: 10px;
            padding: 10px;
        }
        .maze-box {
            border: 1px solid #ccc;
            padding: 10px;
        }
    </style>
</head>
<body>
<div class="container">
`

func parseParams(req *http.Request) (size int, turnProb, accRatio, erosionRatio float64) {
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
	size, turnProb, accRatio, erosionRatio := parseParams(req)

	fmt.Fprint(w, htmlTemplate)

	// 生成迷宫和寻找路径
	maze := generateMaze(size, turnProb)
	path := findPath(maze)

	// 第一个画布：原始迷宫
	fmt.Fprint(w, "<div class='maze-box'><h3>原始迷宫</h3>")
	renderMaze(w, maze, path, true)
	fmt.Fprint(w, "</div>")

	// 第二个画布：消除断头路后
	fmt.Fprint(w, "<div class='maze-box'><h3>堆积后</h3>")
	accuMaze(maze, path, accRatio)
	renderMaze(w, maze, path, true)
	fmt.Fprint(w, "</div>")

	// 第三个画布：侵蚀后
	fmt.Fprint(w, "<div class='maze-box'><h3>侵蚀后</h3>")
	erosionMaze(maze, erosionRatio)
	renderMaze(w, maze, path, false)
	fmt.Fprint(w, "</div>")

	// 第四、五、六个画布：空白位置
	fmt.Fprint(w, "<div class='maze-box'></div>")
	fmt.Fprint(w, "<div class='maze-box'></div>")
	fmt.Fprint(w, "<div class='maze-box'></div>")

	// 结束 HTML
	fmt.Fprint(w, "</div></body></html>")
}

func generateMaze(size int, turnProb float64) [][]int {
	maze := make([][]int, size)
	for i := range maze {
		maze[i] = make([]int, size)
		for j := range maze[i] {
			maze[i][j] = 1 // 1表示墙
		}
	}

	var dfs func(p Point, lastp int)
	dfs = func(p Point, lastp int) {
		maze[p.x][p.y] = 0 // 0表示路径
		dirs := []Point{{-2, 0}, {2, 0}, {0, -2}, {0, 2}}

		pos := []int{0, 1, 2, 3}

		if rand.Float64() < turnProb {
			if lastp <= 1 {
				pos = []int{pos[2], pos[3], pos[0], pos[1]}
			}
		} else {
			if lastp > 1 {
				pos = []int{pos[2], pos[3], pos[0], pos[1]}
			}
		}

		pos1 := pos[:2]
		pos2 := pos[2:]

		rand.Shuffle(len(pos1), func(i, j int) {
			pos1[i], pos1[j] = pos1[j], pos1[i]
		})

		rand.Shuffle(len(pos2), func(i, j int) {
			pos2[i], pos2[j] = pos2[j], pos2[i]
		})

		for _, pp := range pos {
			next := Point{p.x + dirs[pp].x, p.y + dirs[pp].y}
			if next.x >= 0 && next.x < size && next.y >= 0 && next.y < size && maze[next.x][next.y] == 1 {
				maze[p.x+dirs[pp].x/2][p.y+dirs[pp].y/2] = 0
				dfs(next, pp)
			}
		}
	}

	dfs(Point{0, 0}, 0)
	maze[0][0] = 0
	maze[size-1][size-1] = 0

	return maze
}

func accuMaze(maze [][]int, path [][]bool, accPrecent float64) int {
	size := len(maze)
	dirs := []Point{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	totalCount := 0

	// 计算初始断头路数量
	initialDeadEnds := countDeadEnds(maze, size, path)
	if initialDeadEnds == 0 {
		return 0
	}

	// 创建队列存储断头路点
	queue := make([]Point, 0)

	// 首次遍历找出所有断头路
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if isDeadEnd(Point{i, j}, maze) {
				queue = append(queue, Point{i, j})
			}
		}
	}

	// 处理队列中的断头路
	for len(queue) > 0 && float64(totalCount)/float64(initialDeadEnds) < accPrecent {
		cur := queue[0]
		queue = queue[1:]

		// 如果是起点或终点，跳过
		if (cur.x == 0 && cur.y == 0) || (cur.x == size-1 && cur.y == size-1) {
			continue
		}

		// 如果当前点仍然是断头路（可能在处理其他点时被改变）
		if isDeadEnd(cur, maze) {
			// 填充当前断头路
			maze[cur.x][cur.y] = 1
			totalCount++

			// 检查周围的点是否变成新的断头路
			for _, d := range dirs {
				next := Point{cur.x + d.x, cur.y + d.y}
				if next.x >= 0 && next.x < size && next.y >= 0 && next.y < size &&
					maze[next.x][next.y] == 0 &&
					isDeadEnd(next, maze) {
					queue = append(queue, next)
				}
			}
		}
	}

	return totalCount
}

func erosionMaze(maze [][]int, erosionPercent float64) int {
	size := len(maze)
	dirs := []Point{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	// 初始化候选点队列和权重映射
	candidates := make([]Point, 0)
	weightMap := make(map[Point]float64)

	// 计算初始可侵蚀点和总墙数
	totalWalls := 0
	for i := 1; i < size-1; i++ {
		for j := 1; j < size-1; j++ {
			if maze[i][j] == 1 {
				totalWalls++
				p := Point{i, j}
				if weight := calculateWeight(p, maze, dirs); weight > 0 {
					candidates = append(candidates, p)
					weightMap[p] = weight
				}
			}
		}
	}

	if totalWalls == 0 {
		return 0
	}

	targetCount := int(float64(totalWalls) * erosionPercent)
	if targetCount == 0 {
		return 0
	}

	eroded := 0
	for eroded < targetCount && len(candidates) > 0 {
		// 根据权重选择一个点
		selectedIdx := selectPointByWeight(candidates, weightMap)
		if selectedIdx < 0 {
			break
		}

		// 侵蚀选中的点
		selected := candidates[selectedIdx]
		maze[selected.x][selected.y] = 0
		eroded++

		// 从候选列表中移除已侵蚀的点
		candidates[selectedIdx] = candidates[len(candidates)-1]
		candidates = candidates[:len(candidates)-1]
		delete(weightMap, selected)

		// 更新受影响点的权重
		for _, d := range dirs {
			nx, ny := selected.x+d.x, selected.y+d.y
			if nx >= 1 && nx < size-1 && ny >= 1 && ny < size-1 && maze[nx][ny] == 1 {
				p := Point{nx, ny}
				if weight := calculateWeight(p, maze, dirs); weight > 0 {
					weightMap[p] = weight
					if !containsPoint(candidates, p) {
						candidates = append(candidates, p)
					}
				}
			}
		}
	}

	return eroded
}

// 计算点的权重
func calculateWeight(p Point, maze [][]int, dirs []Point) float64 {
	emptyCount := 0
	for _, d := range dirs {
		nx, ny := p.x+d.x, p.y+d.y
		if maze[nx][ny] == 0 {
			emptyCount++
		}
	}

	switch emptyCount {
	case 1:
		return 1.0
	case 2:
		return 4.0
	case 3:
		return 39.0
	case 4:
		return 160.0
	default:
		return 0.0
	}
}

// 根据权重选择点
func selectPointByWeight(candidates []Point, weightMap map[Point]float64) int {
	if len(candidates) == 0 {
		return -1
	}

	totalWeight := 0.0
	for _, p := range candidates {
		totalWeight += weightMap[p]
	}

	randWeight := rand.Float64() * totalWeight
	cumWeight := 0.0

	for i, p := range candidates {
		cumWeight += weightMap[p]
		if randWeight <= cumWeight {
			return i
		}
	}

	return len(candidates) - 1
}

// 检查点是否在候选列表中
func containsPoint(candidates []Point, p Point) bool {
	for _, c := range candidates {
		if c == p {
			return true
		}
	}
	return false
}

// 辅助函数���判断一个点是否是断头路
func isDeadEnd(p Point, maze [][]int) bool {
	size := len(maze)
	if maze[p.x][p.y] != 0 {
		return false
	}

	// 如果是起点或终点，不算断头路
	if (p.x == 0 && p.y == 0) || (p.x == size-1 && p.y == size-1) {
		return false
	}

	dirs := []Point{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	pathCount := 0
	for _, d := range dirs {
		nx, ny := p.x+d.x, p.y+d.y
		if nx >= 0 && nx < size && ny >= 0 && ny < size && maze[nx][ny] == 0 {
			pathCount++
		}
	}
	return pathCount == 1
}

// 计算非最短路径上的空白格子数量
func countDeadEnds(maze [][]int, size int, path [][]bool) int {
	count := 0

	// 遍历整个迷宫
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			// 如果是通路(值为0)且不在最短路径上
			if maze[i][j] == 0 && !path[i][j] {
				count++
			}
		}
	}

	return count
}

func findPath(maze [][]int) [][]bool {
	size := len(maze)
	path := make([][]bool, size)
	for i := range path {
		path[i] = make([]bool, size)
	}

	// 使用队列进行BFS搜索
	queue := []Point{{0, 0}}
	// 记录每个点的前驱节点,用于回溯路径
	parent := make(map[Point]Point)

	dirs := []Point{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	found := false

	for len(queue) > 0 && !found {
		cur := queue[0]
		queue = queue[1:]

		if cur.x == size-1 && cur.y == size-1 {
			found = true
			break
		}

		for _, d := range dirs {
			next := Point{cur.x + d.x, cur.y + d.y}
			if next.x >= 0 && next.x < size && next.y >= 0 && next.y < size &&
				maze[next.x][next.y] == 0 && !path[next.x][next.y] {
				queue = append(queue, next)
				parent[next] = cur
				path[next.x][next.y] = true
			}
		}
	}
	// 重置path为全false
	for i := range path {
		for j := range path[i] {
			path[i][j] = false
		}
	}
	if found {
		// 从终点回溯到起点,标记路径
		cur := Point{size - 1, size - 1}
		for cur.x != 0 || cur.y != 0 {
			path[cur.x][cur.y] = true
			cur = parent[cur]
		}
		path[0][0] = true
	}

	return path
}

func renderMaze(w http.ResponseWriter, maze [][]int, path [][]bool, showPath bool) {
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
				} else if showPath && path[i][j] {
					fmt.Fprintf(w, "##")
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
