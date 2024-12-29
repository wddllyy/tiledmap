package pathfind

import (
	"container/heap"
)

// JNode 表示JPS搜索中的一个节点
type JNode struct {
	pos     [2]int // 位置
	g       int    // 从起点到当前点的实际代价
	h       int    // 从当前点到终点的估计代价
	f       int    // f = g + h
	parent  *JNode // 父节点
	index   int    // 在优先队列中的索引
	fromDir [2]int // 从父节点到当前节点的方向
}

// JPriorityQueue 实现堆接口
type JPriorityQueue []*JNode

func (pq JPriorityQueue) Len() int { return len(pq) }

func (pq JPriorityQueue) Less(i, j int) bool {
	return pq[i].f < pq[j].f
}

func (pq JPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *JPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*JNode)
	node.index = n
	*pq = append(*pq, node)
}

func (pq *JPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil  // 避免内存泄漏
	node.index = -1 // 标记为已移除
	*pq = old[0 : n-1]
	return node
}

// 检查位置是否在迷宫范围内且可通行
func isWalkable(maze [][]int, pos [2]int) bool {
	return pos[0] >= 0 && pos[0] < len(maze) &&
		pos[1] >= 0 && pos[1] < len(maze[0]) &&
		maze[pos[0]][pos[1]] == 0
}

// 在给定方向上跳跃，直到找到跳点或碰壁
func jump(maze [][]int, current [2]int, dir [2]int, end [2]int, res *PathFindResult, depth int) [2]int {
	next := [2]int{current[0] + dir[0], current[1] + dir[1]}
	// 如果不可通行，返回nil
	if !isWalkable(maze, next) {
		//fmt.Println("unwalkable return ", next)
		return [2]int{-1, -1}
	}
	res.Check++
	//fmt.Println(strings.Repeat(" ", depth), current, getStrFromDir(dir), next)

	// 如果到达终点，返回当前位置
	if next == end {
		//fmt.Println(strings.Repeat(" ", depth), "found end")
		return next
	}

	// 水平或垂直移动
	if dir[0] == 0 || dir[1] == 0 {
		// 检查强迫邻居
		if dir[0] == 0 { // 水平移动
			// 检查上下是否有强迫邻居
			if isWalkable(maze, [2]int{next[0] - 1, next[1]}) && !isWalkable(maze, [2]int{next[0] - 1, current[1]}) ||
				isWalkable(maze, [2]int{next[0] + 1, next[1]}) && !isWalkable(maze, [2]int{next[0] + 1, current[1]}) {
				//fmt.Println(strings.Repeat(" ", depth), "found -十", current, next)
				return next
			}
			if isWalkable(maze, [2]int{next[0] - 1, next[1]}) && !isWalkable(maze, [2]int{next[0] + dir[0], next[1] + dir[1]}) ||
				isWalkable(maze, [2]int{next[0] + 1, next[1]}) && !isWalkable(maze, [2]int{next[0] + dir[0], next[1] + dir[1]}) {
				//fmt.Println(strings.Repeat(" ", depth), "found -H", current, next)
				return next
			}
		} else { // 垂直移动
			// 检查左右是否有强迫邻居
			if isWalkable(maze, [2]int{next[0], next[1] - 1}) && !isWalkable(maze, [2]int{current[0], next[1] - 1}) ||
				isWalkable(maze, [2]int{next[0], next[1] + 1}) && !isWalkable(maze, [2]int{current[0], next[1] + 1}) {
				//fmt.Println(strings.Repeat(" ", depth), "found |十", current, next)
				return next
			}
			if isWalkable(maze, [2]int{next[0], next[1] - 1}) && !isWalkable(maze, [2]int{next[0] + dir[0], next[1] + dir[1]}) ||
				isWalkable(maze, [2]int{next[0], next[1] + 1}) && !isWalkable(maze, [2]int{next[0] + dir[0], next[1] + dir[1]}) {
				//fmt.Println(strings.Repeat(" ", depth), "found |工", current, next)
				return next
			}
		}
		//fmt.Println("continue line move", next, dir)
		return jump(maze, next, dir, end, res, depth+1)
	}

	// //fmt.Println("continue diagonal line move", next, dir)
	// // 对角线移动
	// // 检查水平和垂直方向是否有跳点
	// if jump(maze, next, [2]int{dir[0], 0}, end, res, depth+1) != [2]int{-1, -1} ||
	// 	jump(maze, next, [2]int{0, dir[1]}, end, res, depth+1) != [2]int{-1, -1} {
	// 	fmt.Println(strings.Repeat(" ", depth), "found:\\", current, next)
	// 	return next
	// }

	// //fmt.Println("continue diagonal check", next, dir)
	// // 继续对角线移动
	return jump(maze, next, dir, end, res, depth+1)
}

func getStrFromDir(dir [2]int) string {

	if dir[0] == 0 && dir[1] == 1 {
		return "→"
	}
	if dir[0] == 0 && dir[1] == -1 {
		return "←"
	}
	if dir[0] == 1 && dir[1] == 0 {
		return "↓"
	}
	if dir[0] == -1 && dir[1] == 0 {
		return "↑"
	}
	if dir[0] == 1 && dir[1] == 1 {
		return "↘"
	}
	if dir[0] == -1 && dir[1] == 1 {
		return "↗"
	}
	if dir[0] == 1 && dir[1] == -1 {
		return "↙"
	}
	if dir[0] == -1 && dir[1] == -1 {
		return "↖"
	}

	return " "
}

// 添加新的辅助函数，用于生成两点之间的路径
func getPointsBetween(start, end [2]int) [][2]int {
	points := make([][2]int, 0)
	dx := end[0] - start[0]
	dy := end[1] - start[1]

	// 获取方向
	stepX := 0
	if dx > 0 {
		stepX = 1
	} else if dx < 0 {
		stepX = -1
	}

	stepY := 0
	if dy > 0 {
		stepY = 1
	} else if dy < 0 {
		stepY = -1
	}

	// 添加起点
	current := start
	points = append(points, current)
	//fmt.Println("start", current)
	// 先横向移动
	for current[0] != end[0] {
		current = [2]int{current[0] + stepX, current[1]}
		points = append(points, current)
	}

	// 再纵向移动
	for current[1] != end[1] {
		current = [2]int{current[0], current[1] + stepY}
		points = append(points, current)
	}

	return points
}

// rebuildPath 从跳点重建完整路径
func rebuildPath(current *JNode) [][2]int {
	if current == nil {
		return make([][2]int, 0)
	}

	// 收集所有跳点
	var nodes []*JNode
	for node := current; node != nil; node = node.parent {
		nodes = append([]*JNode{node}, nodes...)
	}

	path := make([][2]int, 0)
	// 添加起点
	path = append(path, nodes[0].pos)

	// 在每对跳点之间添加中间点
	for i := 0; i < len(nodes)-1; i++ {
		intermediatePoints := getPointsBetween(nodes[i].pos, nodes[i+1].pos)
		path = append(path, intermediatePoints[1:]...)
	}

	return path
}

func FindPathJPS(maze [][]int, start, end [2]int) PathFindResult {
	// 初始化优先队列和访问集合
	openList := &JPriorityQueue{}
	heap.Init(openList)
	visited := make(map[[2]int]bool)
	//fmt.Println("-----------------------------------------------------")
	// 创建起点节点
	startNode := &JNode{
		pos: start,
		g:   0,
		h:   manhattanDistance(start, end),
	}
	startNode.f = startNode.g + startNode.h

	heap.Push(openList, startNode)

	// 定义所有可能的方向（包括对角线）
	dirs := [][2]int{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1}, // 基本方向
		//{-1, -1}, {-1, 1}, {1, -1}, {1, 1}, // 对角线
	}

	var current *JNode
	var res PathFindResult

	// 主循环
	for openList.Len() > 0 {
		current = heap.Pop(openList).(*JNode)
		//fmt.Println("pop from openlist", current.pos)
		res.StepRecord.Steps = append(res.StepRecord.Steps, MazeStep{
			Pos:  current.pos,
			Type: "pop", // 标记为已检查
			Dir:  [2]int{0, 0},
		})
		// 如果到达终点
		if current.pos == end {
			break
		}

		// 标记为已访问
		visited[current.pos] = true

		// 在每个方向上寻找跳点
		for _, dir := range dirs {
			// 跳过与来源方向相反的方向
			if current.parent != nil {
				// 如果当前方向与来源方向相反，跳过
				if dir[0] == -current.fromDir[0] && dir[1] == -current.fromDir[1] {
					continue
				}
			}

			jumpPoint := jump(maze, current.pos, dir, end, &res, 0)
			if jumpPoint == [2]int{-1, -1} || visited[jumpPoint] {
				continue
			}

			distance := manhattanDistance(current.pos, jumpPoint)
			neighbor := &JNode{
				pos:     jumpPoint,
				g:       current.g + distance,
				h:       manhattanDistance(jumpPoint, end),
				parent:  current,
				fromDir: dir, // 记录来源方向
			}
			neighbor.f = neighbor.g + neighbor.h
			res.Cost++

			// 添加到优先队列
			heap.Push(openList, neighbor)
			res.StepRecord.Steps = append(res.StepRecord.Steps, MazeStep{
				Pos:  neighbor.pos,
				Type: "push", // 标记为已检查
				Dir:  [2]int{0, 0},
			})
			//fmt.Println("push to openlist", neighbor.pos, " parent:", current.pos)
		}
	}

	// 使用新函数重建路径
	res.Path = rebuildPath(current)

	return res
}
