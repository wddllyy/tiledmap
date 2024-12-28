package pathfind

type MazeStep struct {
	Pos  [2]int // 当前检查的位置
	Type string // 0:普通 1:墙 2:已检查 3:起点 4:终点 5:当前路径
	Dir  [2]int // 方向
}

type MazeStepRecord struct {
	Steps []MazeStep
}
