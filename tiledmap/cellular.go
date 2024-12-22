package tiledmap

import (
	"math/rand"
)

// 修改常量为导出
const (
	CellularDefaultSize = 55
	DefaultProbability  = 0.55
	DefaultIterations   = 5

	CellularMaxSize = 1000
	MaxIterations   = 20
	MinEntranceSize = 3
)

// 导出结构体和字段
type MazeParams struct {
	Size        int     // 迷宫尺寸
	Probability float64 // 障碍物生成概率
	Iterations  int     // 迭代次数
}

// 将函数改为导出
func InitializeMaze(size int, probability float64) [][]int {
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
	setEntranceArea(maze, size-MinEntranceSize, size-MinEntranceSize) // 右下角出口

	return maze
}

func CellularMaze(maze [][]int) {
	size := len(maze)
	setEntranceArea(maze, 0, 0)                                       // 左上角入口
	setEntranceArea(maze, size-MinEntranceSize, size-MinEntranceSize) // 右下角出口

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

func ConnectRegionsByBFS(maze [][]int) {
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
		// 对当前第一区域进行BFS寻路
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

// 设置入口区域
func setEntranceArea(maze [][]int, startX, startY int) {
	for i := startX; i < startX+MinEntranceSize && i < len(maze); i++ {
		for j := startY; j < startY+MinEntranceSize && j < len(maze); j++ {
			maze[i][j] = 0
		}
	}
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
