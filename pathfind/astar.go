package pathfind

import (
	"container/heap"
	"math"
)

// Node 表示搜索中的一个节点
type Node struct {
	pos    [2]int // 位置
	g      int    // 从起点到当前点的实际代价
	h      int    // 从当前点到终点的估计代价
	f      int    // f = g + h
	parent *Node  // 父节点
	index  int    // 在优先队列中的索引
}

type PathFindResult struct {
	Path  [][2]int
	Cost  int
	Check int
}

// PriorityQueue 实现堆接口
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].f < pq[j].f
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*Node)
	node.index = n
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil  // 避免内存泄漏
	node.index = -1 // 标记为已移除
	*pq = old[0 : n-1]
	return node
}

// 计算曼哈顿距离
func manhattanDistance(a, b [2]int) int {
	return int(math.Abs(float64(a[0]-b[0])) + math.Abs(float64(a[1]-b[1])))
}

// FindPathAStar 使用A*算法寻找从start到end的路径
func FindPathAStar(maze [][]int, start, end [2]int) PathFindResult {
	// 初始化开放列表和关闭列表
	openList := &PriorityQueue{}
	heap.Init(openList)
	closedSet := make(map[[2]int]bool)

	// 创建起点节点
	startNode := &Node{
		pos: start,
		g:   0,
		h:   manhattanDistance(start, end),
	}
	startNode.f = startNode.g + startNode.h

	heap.Push(openList, startNode)

	// 定义方向：上、右、下、左
	dirs := [][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

	var current *Node
	var res PathFindResult
	// 主循环
	for openList.Len() > 0 {
		current = heap.Pop(openList).(*Node)

		// 如果到达终点
		if current.pos == end {
			break
		}

		// 将当前节点加入关闭列表
		closedSet[current.pos] = true

		// 检查相邻节点
		for _, dir := range dirs {
			res.Check++

			nextPos := [2]int{current.pos[0] + dir[0], current.pos[1] + dir[1]}

			// 检查边界和是否可通行
			if nextPos[0] < 0 || nextPos[0] >= len(maze) ||
				nextPos[1] < 0 || nextPos[1] >= len(maze[0]) ||
				maze[nextPos[0]][nextPos[1]] == 1 {
				continue
			}

			// 如果在关闭列表中，跳过
			if closedSet[nextPos] {
				continue
			}

			// 计算新的g值
			newG := current.g + 1

			// 创建新节点
			neighbor := &Node{
				pos:    nextPos,
				g:      newG,
				h:      manhattanDistance(nextPos, end),
				parent: current,
			}
			neighbor.f = neighbor.g + neighbor.h
			res.Cost++
			// 添加到开放列表
			heap.Push(openList, neighbor)
		}
	}

	// 重建路径
	res.Path = make([][2]int, 0)
	for node := current; node != nil; node = node.parent {
		res.Path = append([][2]int{node.pos}, res.Path...)
	}
	return res
}
