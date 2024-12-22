package pathfind

import (
	"container/heap"
)

// DNode 表示Dijkstra搜索中的一个节点
type DNode struct {
	pos    [2]int // 位置
	dist   int    // 从起点到当前点的距离
	parent *DNode // 父节点
	index  int    // 在优先队列中的索引
}

// DPriorityQueue 实现堆接口
type DPriorityQueue []*DNode

func (pq DPriorityQueue) Len() int { return len(pq) }

func (pq DPriorityQueue) Less(i, j int) bool {
	return pq[i].dist < pq[j].dist
}

func (pq DPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *DPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*DNode)
	node.index = n
	*pq = append(*pq, node)
}

func (pq *DPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil  // 避免内存泄漏
	node.index = -1 // 标记为已移除
	*pq = old[0 : n-1]
	return node
}

// FindPathDijkstra 使用Dijkstra算法寻找从start到end的路径
func FindPathDijkstra(maze [][]int, start, end [2]int) PathFindResult {
	// 初始化优先队列和访问集合
	openList := &DPriorityQueue{}
	heap.Init(openList)
	visited := make(map[[2]int]bool)

	// 创建起点节点
	startNode := &DNode{
		pos:  start,
		dist: 0,
	}

	heap.Push(openList, startNode)

	// 定义方向：上、右、下、左
	dirs := [][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

	var current *DNode
	var res PathFindResult

	// 主循环
	for openList.Len() > 0 {

		current = heap.Pop(openList).(*DNode)

		// 如果到达终点
		if current.pos == end {
			break
		}

		// 如果已访问过，跳过
		if visited[current.pos] {
			continue
		}

		// 标记为已访问
		visited[current.pos] = true

		// 检查相邻节点
		for _, dir := range dirs {
			res.Check++

			nextPos := [2]int{current.pos[0] + dir[0], current.pos[1] + dir[1]}

			// 检查边界和是否可通行
			if nextPos[0] < 0 || nextPos[0] >= len(maze) ||
				nextPos[1] < 0 || nextPos[1] >= len(maze[0]) ||
				maze[nextPos[0]][nextPos[1]] == 1 ||
				visited[nextPos] {
				continue
			}

			// 创建新节点
			neighbor := &DNode{
				pos:    nextPos,
				dist:   current.dist + 1,
				parent: current,
			}
			res.Cost++

			// 添加到优先队列
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
