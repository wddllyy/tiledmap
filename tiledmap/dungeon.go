package tiledmap

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
)

// 参考了https://journal.stuffwithstuff.com/2014/12/21/rooms-and-mazes/

type Room struct {
	X, Y          int // 房间左上角坐标
	Width, Height int // 房间大小
}

type Dungeon struct {
	Width  int
	Height int
	Tiles  [][]int
	Rooms  []Room
}

func NewDungeon(width, height int) *Dungeon {
	d := &Dungeon{
		Width:  width,
		Height: height,
		Tiles:  make([][]int, height),
	}
	// 初始化为墙
	for y := 0; y < height; y++ {
		d.Tiles[y] = make([]int, width)
		for x := 0; x < width; x++ {
			d.Tiles[y][x] = 1
		}
	}
	return d
}

func (d *Dungeon) AddRoom(room Room) bool {
	// 检查房间是否超出边界
	if room.X < 1 || room.Y < 1 ||
		room.X+room.Width >= d.Width-1 ||
		room.Y+room.Height >= d.Height-1 {
		return false
	}

	// 检查是否与其他房间重叠（包括留出1格间距）
	for y := room.Y - 1; y <= room.Y+room.Height; y++ {
		for x := room.X - 1; x <= room.X+room.Width; x++ {
			if d.Tiles[y][x] == 0 {
				return false
			}
		}
	}

	// 添加房间
	for y := room.Y; y < room.Y+room.Height; y++ {
		for x := room.X; x < room.X+room.Width; x++ {
			d.Tiles[y][x] = 0
		}
	}
	d.Rooms = append(d.Rooms, room)
	return true
}

func (d *Dungeon) GenerateMazeBetweenRooms() {
	// 标记所有房间区域（包括边界）
	roomArea := make([][]bool, d.Height)
	for i := range roomArea {
		roomArea[i] = make([]bool, d.Width)
	}

	// 标记房间区域
	for _, room := range d.Rooms {
		// 包括房间周围的一格边界
		for y := room.Y - 1; y <= room.Y+room.Height; y++ {
			if y < 0 || y >= d.Height {
				continue
			}
			for x := room.X - 1; x <= room.X+room.Width; x++ {
				if x < 0 || x >= d.Width {
					continue
				}
				roomArea[y][x] = true
			}
		}
	}

	// 在非房间区域生成迷宫
	// 确保只在2的倍数位置生成墙
	for y := 1; y < d.Height-1; y += 2 {
		for x := 1; x < d.Width-1; x += 2 {
			// 跳过房间区域
			if roomArea[y][x] {
				continue
			}

			// 在当前位置生成迷宫单元
			d.Tiles[y][x] = 0 // 设置为通道
		}
	}

}

func (d *Dungeon) ConnectPassagesByDFS() {
	// 创建访问标记数组
	visited := make([][]bool, d.Height)
	for i := range visited {
		visited[i] = make([]bool, d.Width)
	}

	// 标记房间区域
	roomArea := make([][]bool, d.Height)
	for i := range roomArea {
		roomArea[i] = make([]bool, d.Width)
	}

	// 标记所有房间区域（包括边界）
	for _, room := range d.Rooms {
		for y := room.Y - 1; y <= room.Y+room.Height; y++ {
			if y < 0 || y >= d.Height {
				continue
			}
			for x := room.X - 1; x <= room.X+room.Width; x++ {
				if x < 0 || x >= d.Width {
					continue
				}
				roomArea[y][x] = true
			}
		}
	}

	// 方向数组：上下左右
	dirs := [][2]int{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1},
	}

	// DFS递归函数
	var dfs func(x, y int)
	dfs = func(x, y int) {
		visited[y][x] = true

		// 随机打乱方向
		randDirs := make([][2]int, len(dirs))
		copy(randDirs, dirs)
		rand.Shuffle(len(randDirs), func(i, j int) {
			randDirs[i], randDirs[j] = randDirs[j], randDirs[i]
		})

		// 向四个方向扩展
		for _, dir := range randDirs {
			// 计算新位置（跳两格）
			nx := x + dir[0]*2
			ny := y + dir[1]*2

			// 检查边界
			if nx < 1 || nx >= d.Width-1 || ny < 1 || ny >= d.Height-1 {
				continue
			}

			// 跳过已访问点的和房间区域
			if visited[ny][nx] || roomArea[ny][nx] {
				continue
			}

			// 如果是通道点，连接并继续DFS
			if d.Tiles[ny][nx] == 0 {
				// 打通中间的墙
				mx := x + dir[0]
				my := y + dir[1]
				d.Tiles[my][mx] = 0

				// 继续DFS
				dfs(nx, ny)
			}
		}
	}

	// 遍历地图寻找通道点开始DFS
	for y := 1; y < d.Height-1; y += 2 {
		for x := 1; x < d.Width-1; x += 2 {
			// 跳过已访问的点和房间区域
			if visited[y][x] || roomArea[y][x] {
				continue
			}

			// 如果是通道点，开始DFS
			if d.Tiles[y][x] == 0 {
				dfs(x, y)
			}
		}
	}
}

func (d *Dungeon) findConnectedRegions() []ConnectedRegion {
	visited := make([][]bool, d.Height)
	for i := range visited {
		visited[i] = make([]bool, d.Width)
	}

	var regions []ConnectedRegion
	regionId := 0

	// 遍历地图寻找未访问的通道
	for y := 1; y < d.Height-1; y++ {
		for x := 1; x < d.Width-1; x++ {
			if !visited[y][x] && d.Tiles[y][x] == 0 {
				// 发现新区域，使用BFS填充
				region := ConnectedRegion{id: regionId}
				queue := []struct{ x, y int }{{x, y}}
				visited[y][x] = true

				for len(queue) > 0 {
					curr := queue[0]
					queue = queue[1:]
					region.cells = append(region.cells, curr)

					// 检查四个方向
					dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
					for _, dir := range dirs {
						nx := curr.x + dir[0]
						ny := curr.y + dir[1]

						if nx >= 0 && nx < d.Width && ny >= 0 && ny < d.Height &&
							!visited[ny][nx] && d.Tiles[ny][nx] == 0 {
							visited[ny][nx] = true
							queue = append(queue, struct{ x, y int }{nx, ny})
						}
					}
				}

				regions = append(regions, region)
				regionId++
			}
		}
	}

	return regions
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func GenerateDungeon(width, height, roomCount, minSize, maxSize int) *Dungeon {
	// 确保宽度和高度为奇数
	if width%2 == 0 {
		width++
	}
	if height%2 == 0 {
		height++
	}

	// 确保最小和最大尺寸为奇数
	if minSize%2 == 0 {
		minSize++
	}
	if maxSize%2 == 0 {
		maxSize++
	}

	// 确保尺寸范围合理
	if minSize < 3 {
		minSize = 3
	}
	if maxSize > 15 {
		maxSize = 15
	}
	if minSize > maxSize {
		minSize, maxSize = maxSize, minSize
	}

	dungeon := NewDungeon(width, height)

	// 尝试添加指定数量的房间
	attempts := 0
	for len(dungeon.Rooms) < roomCount && attempts < 10000 {
		// 生成范围内的随机奇数尺寸
		sizeRange := (maxSize - minSize) / 2
		roomWidth := minSize + (rand.Intn(sizeRange+1) * 2)
		roomHeight := minSize + (rand.Intn(sizeRange+1) * 2)

		// 确保房间位置为奇数
		x := int(rand.Intn((width-roomWidth-2)/2))*2 + 1
		y := int(rand.Intn((height-roomHeight-2)/2))*2 + 1

		room := Room{
			X:      x,
			Y:      y,
			Width:  roomWidth,
			Height: roomHeight,
		}

		if dungeon.AddRoom(room) {
			attempts = 0
		} else {
			attempts++
		}
	}

	return dungeon
}

const dungeonCSS = `
<style>
    .dungeon-container {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 20px;
        padding: 20px;
    }
    .dungeon-controls {
        display: flex;
        gap: 10px;
        margin-bottom: 20px;
    }
    .dungeon-grid {
        display: grid;
        gap: 0;
        position: relative;
    }
    .dungeon-cell {
        width: 8px;
        height: 8px;
        display: inline-block;
    }
    .wall { background-color: #666; }
    .floor { background-color: #fff; }
    .room-number {
        position: absolute;
        color: #f00;
        font-size: 12px;
        font-weight: bold;
        z-index: 1;
    }
</style>`

func DungeonHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, dungeonCSS)

	// 解析参数，确保为奇数
	width := 51 // 默认值改为奇数
	if w := r.URL.Query().Get("width"); w != "" {
		if val, err := strconv.Atoi(w); err == nil && val > 12 && val <= 100 {
			if val%2 == 0 {
				val++ // 确保为奇数
			}
			width = val
		}
	}

	height := 51 // 默认值改为奇数
	if h := r.URL.Query().Get("height"); h != "" {
		if val, err := strconv.Atoi(h); err == nil && val > 12 && val <= 100 {
			if val%2 == 0 {
				val++ // 确保为奇数
			}
			height = val
		}
	}

	rooms := 8
	if rm := r.URL.Query().Get("rooms"); rm != "" {
		if val, err := strconv.Atoi(rm); err == nil && val > 0 && val <= 50 {
			rooms = val
		}
	}

	extraPathProb := float32(0.2)
	if ep := r.URL.Query().Get("extraPathProb"); ep != "" {
		if val, err := strconv.ParseFloat(ep, 64); err == nil && val >= 0 && val <= 1 {
			extraPathProb = float32(val)
		}
	}

	// 添加最小和最大尺寸参数
	minSize := 5
	if ms := r.URL.Query().Get("minSize"); ms != "" {
		if val, err := strconv.Atoi(ms); err == nil && val >= 3 && val <= 15 {
			if val%2 == 0 {
				val++ // 确保为奇数
			}
			minSize = val
		}
	}

	maxSize := 11
	if ms := r.URL.Query().Get("maxSize"); ms != "" {
		if val, err := strconv.Atoi(ms); err == nil && val >= 3 && val <= 15 {
			if val%2 == 0 {
				val++ // 确保为奇数
			}
			maxSize = val
		}
	}

	// 控制表单
	fmt.Fprintf(w, `
		<div class="dungeon-container">
			<div class="dungeon-controls">
				<form>
					宽度: <input type="number" name="width" value="%d" min="13" max="99" step="2">
					高度: <input type="number" name="height" value="%d" min="13" max="99" step="2">
					房间数: <input type="number" name="rooms" value="%d" min="2" max="50">
					最小房间尺寸: <input type="number" name="minSize" value="%d" min="3" max="15" step="2">
					最大房间尺寸: <input type="number" name="maxSize" value="%d" min="5" max="15" step="2">
					额外通路概率: <input type="number" name="extraPathProb" value="%0.1f" step="0.1" min="0" max="1">
					<input type="submit" value="生成">
				</form>
			</div>
	`, width, height, rooms, minSize, maxSize, extraPathProb)

	// 创建容器div来水平排列两个视图
	fmt.Fprintf(w, `<div style="display: flex; gap: 20px; justify-content: center;">`)

	// 生成地牢
	dungeon := GenerateDungeon(width, height, rooms, minSize, maxSize)

	// 第一阶段：生成迷宫
	dungeon.GenerateMazeBetweenRooms()
	fmt.Fprintf(w, `
		<div>
			<h3 style="text-align: center">阶段1: 生成迷宫</h3>
			<div class="dungeon-grid" style="grid-template-columns: repeat(%d, 8px);">`, width)
	renderDungeon(w, dungeon)
	fmt.Fprint(w, "</div></div>")

	// 第二阶段：连接通道
	dungeon.ConnectPassagesByDFS()
	fmt.Fprintf(w, `
		<div>
			<h3 style="text-align: center">阶段2: 生成通道</h3>
			<div class="dungeon-grid" style="grid-template-columns: repeat(%d, 8px);">`, width)
	renderDungeon(w, dungeon)
	fmt.Fprint(w, "</div></div>")

	// 第三阶段：连接所有区域
	dungeon.ConnectAllRegions(extraPathProb)
	fmt.Fprintf(w, `
		<div>
			<h3 style="text-align: center">阶段3: 连接区域</h3>
			<div class="dungeon-grid" style="grid-template-columns: repeat(%d, 8px);">`, width)
	renderDungeon(w, dungeon)
	fmt.Fprint(w, "</div></div>")

	// 第四阶段：把死胡同堵上
	dungeon.FillDeadEnds()
	fmt.Fprintf(w, `
		<div>
			<h3 style="text-align: center">阶段4: 堵上死胡同</h3>
			<div class="dungeon-grid" style="grid-template-columns: repeat(%d, 8px);">`, width)
	renderDungeon(w, dungeon)
	fmt.Fprint(w, "</div></div>")
}

func renderDungeon(w http.ResponseWriter, d *Dungeon) {
	// 渲染地牢网格
	for y := 0; y < d.Height; y++ {
		for x := 0; x < d.Width; x++ {
			cellClass := "wall"
			if d.Tiles[y][x] == 0 {
				cellClass = "floor"
			}
			fmt.Fprintf(w, `<div class="dungeon-cell %s"></div>`, cellClass)
		}
	}

	// 渲染房间编号
	for i, room := range d.Rooms {
		centerX := room.X + room.Width/2
		centerY := room.Y + room.Height/2

		pixelX := centerX * 8
		pixelY := centerY * 8

		fmt.Fprintf(w, `<div class="room-number" style="left: %dpx; top: %dpx;">%d</div>`,
			pixelX, pixelY, i+1)
	}
}

// 并查集结构
type UnionFind struct {
	parent []int
	rank   []int
}

// 创建并查集
func NewUnionFind(size int) *UnionFind {
	parent := make([]int, size)
	rank := make([]int, size)
	for i := 0; i < size; i++ {
		parent[i] = i
	}
	return &UnionFind{parent: parent, rank: rank}
}

// 查找根节点
func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x]) // 路径压缩
	}
	return uf.parent[x]
}

// 合并两个集合
func (uf *UnionFind) Union(x, y int) {
	px, py := uf.Find(x), uf.Find(y)
	if px == py {
		return
	}
	// 按秩合并
	if uf.rank[px] < uf.rank[py] {
		uf.parent[px] = py
	} else if uf.rank[px] > uf.rank[py] {
		uf.parent[py] = px
	} else {
		uf.parent[py] = px
		uf.rank[px]++
	}
}

func (d *Dungeon) ConnectAllRegions(extraPathProb float32) {
	// 1. 获取所有区域和连接区信息
	regions, connections := d.FindConnectionInfo()
	if len(regions) <= 1 {
		return
	}
	// fmt.Println("regions:", len(regions))
	// for i, region := range regions {
	// 	fmt.Printf("Region %d:\n", i)
	// 	fmt.Printf("  ID: %d\n", region.id)
	// 	fmt.Printf("  单元格数量: %d\n", len(region.cells))
	// 	fmt.Printf("  单元格: %v\n", region.cells)
	// 	fmt.Println()
	// }
	// fmt.Println("connections:", len(connections))
	// for i, connection := range connections {
	// 	fmt.Printf("Connection %d:\n", i)
	// 	fmt.Printf("  Region1: %d\n", connection.Region1)
	// 	fmt.Printf("  Region2: %d\n", connection.Region2)
	// 	fmt.Printf("  单元格数量: %d\n", len(connection.Cells))
	// 	fmt.Printf("  单元格: %v\n", connection.Cells)
	// 	fmt.Println()
	// }

	// 2. 创建并查集
	uf := NewUnionFind(len(regions))

	// 3. 随机打乱连接区顺序
	randConnections := make([]ConnectionZone, len(connections))
	copy(randConnections, connections)
	rand.Shuffle(len(randConnections), func(i, j int) {
		randConnections[i], randConnections[j] = randConnections[j], randConnections[i]
	})

	// 4. 使用Kruskal算法选择连接区
	for _, conn := range randConnections {
		r1, r2 := conn.Region1, conn.Region2

		// 如果两个区域还未连通
		if uf.Find(r1) != uf.Find(r2) {
			// 打通这个连接区
			for _, cell := range conn.Cells {
				// 随机选择一个格子打通
				if rand.Float32() < extraPathProb { // 20%的概率打通一个格子
					d.Tiles[cell.y][cell.x] = 0
				}
			}
			// 至少确保打通一个格子
			if len(conn.Cells) > 0 {
				randomCell := conn.Cells[rand.Intn(len(conn.Cells))]
				d.Tiles[randomCell.y][randomCell.x] = 0
			}

			// 在并查集中合并这两个区域
			uf.Union(r1, r2)
		}
	}
}

// 添加必要的类型定义
type Cell struct {
	x, y int
}

type ConnectionZone struct {
	Region1, Region2 int
	Cells            []Cell
}

type ConnectedRegion struct {
	id    int
	cells []struct{ x, y int }
}

func (d *Dungeon) FindConnectionInfo() ([]ConnectedRegion, []ConnectionZone) {
	regions := d.findConnectedRegions()
	var connections []ConnectionZone

	// 寻找可能的连接区域
	// TODO：这里可以性能优化
	for i := 0; i < len(regions); i++ {
		for j := i + 1; j < len(regions); j++ {
			// 寻找两个区域之间的连接点
			var connCells []Cell
			for _, cell1 := range regions[i].cells {
				for _, cell2 := range regions[j].cells {
					if abs(cell1.x-cell2.x)+abs(cell1.y-cell2.y) == 2 {
						connCells = append(connCells, Cell{x: (cell1.x + cell2.x) / 2, y: (cell1.y + cell2.y) / 2})
					}
				}
			}
			if len(connCells) > 0 {
				connections = append(connections, ConnectionZone{
					Region1: regions[i].id,
					Region2: regions[j].id,
					Cells:   connCells,
				})
			}
		}
	}
	return regions, connections
}

func (d *Dungeon) FillDeadEnds() {
	changed := true
	for changed {
		changed = false
		for y := 1; y < d.Height-1; y++ {
			for x := 1; x < d.Width-1; x++ {
				if d.Tiles[y][x] == 0 {
					// 计算周围的墙数量
					walls := 0
					if d.Tiles[y-1][x] == 1 {
						walls++
					}
					if d.Tiles[y+1][x] == 1 {
						walls++
					}
					if d.Tiles[y][x-1] == 1 {
						walls++
					}
					if d.Tiles[y][x+1] == 1 {
						walls++
					}

					// 如果是死胡同（三面墙）
					if walls == 3 {
						d.Tiles[y][x] = 1
						changed = true
					}
				}
			}
		}
	}
}
