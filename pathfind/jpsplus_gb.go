package pathfind

import (
	"container/heap"
)

// PreprocessedMaze 存储预处理的迷宫信息
type PreprocessedMaze struct {
	maze       [][]int       // 原始迷宫
	primaryJPs [][][4][2]int // 每个方向的主要跳点
}

// 预处理迷宫，计算跳点和边界
func PreprocessMaze(maze [][]int) *PreprocessedMaze {
	height, width := len(maze), len(maze[0])

	// 创建一个全0的迷宫用于测试
	testMaze := make([][]int, height)
	for i := range testMaze {
		testMaze[i] = make([]int, width)
		// 所有元素默认为0，不需要特别设置
	}

	pm := &PreprocessedMaze{
		maze:       maze, //testMaze, // 使用测试迷宫替代输入的迷宫
		primaryJPs: [][][4][2]int{},
	}

	// 初始化方向数组
	dirs := [][2]int{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1}, // 基本方向 上 下 左 右
		//{-1, -1}, {-1, 1}, {1, -1}, {1, 1}, // 对角线
	}

	// 初始化primaryJPs数组
	pm.primaryJPs = make([][][4][2]int, height)
	for y := 0; y < height; y++ {
		pm.primaryJPs[y] = make([][4][2]int, width)
		// 不需要再为每个x初始化，因为[4][2]int是固定大小的数组，会自动初始化为零值
	}

	// 重新排序循环：先y后x最后是方向i
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if maze[y][x] == 0 {
				for i := range dirs {
					jp := pm.findJumpPoint([2]int{y, x}, dirs[i])
					pm.primaryJPs[y][x][i] = jp
				}
			}
		}
	}

	// // 添加调试输出
	// for y := 0; y < height; y++ {
	// 	for x := 0; x < width; x++ {
	// 		if maze[y][x] == 0 {
	// 			fmt.Printf("\n位置 [%d,%d] 的跳点:\n", y, x)

	// 			// 第一行（上）
	// 			fmt.Printf("[     ] [%2d,%2d] [     ]\n", pm.primaryJPs[y][x][0][0], pm.primaryJPs[y][x][0][1])

	// 			// 第二行（左、中心点、右）
	// 			fmt.Printf("[%2d,%2d] [%2d,%2d] [%2d,%2d]\n", pm.primaryJPs[y][x][2][0], pm.primaryJPs[y][x][2][1], y, x, pm.primaryJPs[y][x][3][0], pm.primaryJPs[y][x][3][1])

	// 			// 第三行（下）
	// 			fmt.Printf("[     ] [%2d,%2d] [     ]\n", pm.primaryJPs[y][x][1][0], pm.primaryJPs[y][x][1][1])
	// 		}
	// 	}
	// }

	return pm
}

// findJumpPoint 在给定方向上寻找跳点
func (pm *PreprocessedMaze) findJumpPoint(pos [2]int, dir [2]int) [2]int {
	next := [2]int{pos[0] + dir[0], pos[1] + dir[1]}

	// 如果不可通行，返回nil
	if !isWalkable(pm.maze, next) {
		return [2]int{-1, -1}
	}
	// 检next在dir方向下一个tile有障碍：撞墙了）
	if !isWalkable(pm.maze, [2]int{next[0] + dir[0], next[1] + dir[1]}) {
		//fmt.Println(strings.Repeat(" ", depth), "found -H", current, next)
		return next
	}

	if dir[0] == 0 { // 水平移动
		// 检查上下是否有强迫邻居（当前tile上下有障碍，next tile上下没有障碍：有上或下的岔路）
		if isWalkable(pm.maze, [2]int{next[0] - 1, next[1]}) && !isWalkable(pm.maze, [2]int{next[0] - 1, pos[1]}) ||
			isWalkable(pm.maze, [2]int{next[0] + 1, next[1]}) && !isWalkable(pm.maze, [2]int{next[0] + 1, pos[1]}) {
			//fmt.Println(strings.Repeat(" ", depth), "found -十", current, next)
			return next
		}

	} else { // 垂直移动
		// 检查左右是否有强迫邻居
		if isWalkable(pm.maze, [2]int{next[0], next[1] - 1}) && !isWalkable(pm.maze, [2]int{pos[0], next[1] - 1}) ||
			isWalkable(pm.maze, [2]int{next[0], next[1] + 1}) && !isWalkable(pm.maze, [2]int{pos[0], next[1] + 1}) {
			//fmt.Println(strings.Repeat(" ", depth), "found |十", current, next)
			return next
		}

	}
	//fmt.Println("continue line move", next, dir)
	return pm.findJumpPoint(next, dir)
}

// FindPathJPSPlus 使用JPS+算法寻找路径
func (pm *PreprocessedMaze) FindPathJPSPlus(start, end [2]int) PathFindResult {
	openList := &JPriorityQueue{}
	heap.Init(openList)
	visited := make(map[[2]int]bool)

	startNode := &JNode{
		pos: start,
		g:   0,
		h:   manhattanDistance(start, end),
	}
	startNode.f = startNode.g + startNode.h

	heap.Push(openList, startNode)

	dirs := [][2]int{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1},
		//{-1, -1}, {-1, 1}, {1, -1}, {1, 1},
	}

	var current *JNode
	var res PathFindResult

	for openList.Len() > 0 {
		current = heap.Pop(openList).(*JNode)
		//fmt.Println("pop", current.pos)

		res.StepRecord.Steps = append(res.StepRecord.Steps, MazeStep{
			Pos:  current.pos,
			Type: "pop", // 标记为已检查
			Dir:  [2]int{0, 0},
		})

		if current.pos == end {
			break
		}

		if visited[current.pos] {
			continue
		}

		visited[current.pos] = true

		for i, dir := range dirs {

			if current.parent != nil {
				// 如果当前方向与来源方向相反，跳过
				if dir[0] == -current.fromDir[0] && dir[1] == -current.fromDir[1] {
					continue
				}
			}

			// 使用预计算的跳点
			jp := pm.primaryJPs[current.pos[0]][current.pos[1]][i]
			if jp == [2]int{-1, -1} || visited[jp] {
				continue
			}
			//fmt.Println("  jp", jp)

			if (dir[0] == 0 && (end[1] > current.pos[1] && end[1] < jp[1] || end[1] < current.pos[1] && end[1] > jp[1])) ||
				(dir[1] == 0 && (end[0] > current.pos[0] && end[0] < jp[0] || end[0] < current.pos[0] && end[0] > jp[0])) {
				jp = end
			}

			res.Check++

			distance := manhattanDistance(current.pos, jp)
			neighbor := &JNode{
				pos:     jp,
				g:       current.g + distance,
				h:       manhattanDistance(jp, end),
				parent:  current,
				fromDir: dir, // 记录来源方向
			}
			neighbor.f = neighbor.g + neighbor.h
			res.Cost++

			//fmt.Println("  push", neighbor.pos)
			heap.Push(openList, neighbor)
			res.StepRecord.Steps = append(res.StepRecord.Steps, MazeStep{
				Pos:  neighbor.pos,
				Type: "push", // 标记为已检查
				Dir:  [2]int{0, 0},
			})
		}
	}

	// 使用jps.go中的rebuildPath函数重建路径
	res.Path = rebuildPath(current)

	return res
}
