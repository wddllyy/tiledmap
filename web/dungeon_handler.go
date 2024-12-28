package main

import (
	"fmt"
	"net/http"
	"strconv"

	"mazemap/tiledmap"
)

func dungeonHandler(w http.ResponseWriter, r *http.Request) {

	printHtmlHead(w, "迷宫生成算法")

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

	extraPathProb := float32(0.2)
	if ep := r.URL.Query().Get("extraPathProb"); ep != "" {
		if val, err := strconv.ParseFloat(ep, 64); err == nil && val >= 0 && val <= 1 {
			extraPathProb = float32(val)
		}
	}

	// 控制表单
	fmt.Fprintf(w, `
<div class="all-container">
	<div class="all-controls">
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
	<div style="display: flex; gap: 20px; justify-content: center;">`,
		width, height, rooms, minSize, maxSize, extraPathProb)

	// 生成地牢
	dungeon := tiledmap.GenerateDungeon(width, height, rooms, minSize, maxSize)

	// 第一阶段：生成迷宫
	dungeon.GenerateMazeBetweenRooms()
	renderDungeonWithTitle(w, dungeon, "阶段1: 生成迷宫")

	// 第二阶段：连接通道
	dungeon.ConnectPassagesByDFS()
	renderDungeonWithTitle(w, dungeon, "阶段2: 生成通道")

	// 第三阶段：连接所有区域
	dungeon.ConnectAllRegions(extraPathProb)
	renderDungeonWithTitle(w, dungeon, "阶段3: 连接区域")

	// 第四阶段：把死胡同堵上
	dungeon.FillDeadEnds()
	renderDungeonWithTitle(w, dungeon, "阶段4: 堵上死胡同")

	fmt.Fprint(w, "\n</div></div></body></html>")
}

func renderDungeonWithTitle(w http.ResponseWriter, d *tiledmap.Dungeon, title string) {
	fmt.Fprintf(w, `
		<div>
			<h3 style="text-align: center">%s</h3>
			<div class="wfc-grid" style="grid-template-columns: repeat(%d, 8px);">`, title, d.Width)
	renderDungeon(w, d)
	fmt.Fprint(w, "</div></div>")
}

func renderDungeon(w http.ResponseWriter, d *tiledmap.Dungeon) {
	// 渲染地牢网格
	for y := 0; y < d.Height; y++ {
		for x := 0; x < d.Width; x++ {
			cellClass := "wall"
			if d.Tiles[y][x] == 0 {
				cellClass = "floor"
			}
			fmt.Fprintf(w, `<div class="wfc-cell %s"></div>`, cellClass)
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
