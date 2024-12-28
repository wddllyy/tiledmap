package main

import (
	"fmt"
	"net/http"
	"strconv"

	"mazemap/tiledmap"
)

func wfcHandler(w http.ResponseWriter, r *http.Request) {
	printHtmlHead(w, "波函数坍塌")

	// 从URL参数获取宽度和高度
	width := 32 // 默认值
	if w := r.URL.Query().Get("width"); w != "" {
		if val, err := strconv.Atoi(w); err == nil && val > 0 && val <= 100 {
			width = val
		}
	}

	height := 32 // 默认值
	if h := r.URL.Query().Get("height"); h != "" {
		if val, err := strconv.Atoi(h); err == nil && val > 0 && val <= 100 {
			height = val
		}
	}

	// 控制表单
	fmt.Fprintf(w, `
<div class="all-container">
	<div class="all-controls">
		<form>
			宽度: <input type="number" name="width" value="%d" min="1" max="100">
			高度: <input type="number" name="height" value="%d" min="1" max="100">
			<input type="submit" value="生成">
		</form>
	</div>
	<div style="display: flex; gap: 20px; justify-content: center;">`,
		width, height)

	wfc := tiledmap.NewWFC(width, height)
	tileMap := wfc.Generate()

	renderWFCWithTitle(w, tileMap, "波函数坍缩生成")

	fmt.Fprint(w, "\n</div></div></body></html>")
}

func renderWFCWithTitle(w http.ResponseWriter, tileMap [][]int, title string) {
	fmt.Fprintf(w, `
		<div>
			<h3 style="text-align: center">%s</h3>
			<div class="wfcout-grid">`, title)
	renderWFC(w, tileMap)
	fmt.Fprint(w, "</div></div>")
}

func renderWFC(w http.ResponseWriter, tileMap [][]int) {
	height := len(tileMap)
	width := len(tileMap[0])

	fmt.Fprintf(w, `
	<div class="wfc-grid" style="grid-template-columns: repeat(%d, 8px);">`, width)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var class string
			switch tileMap[y][x] {
			case tiledmap.TILE_GRASS:
				class = "grass"
			case tiledmap.TILE_WATER:
				class = "water"
			case tiledmap.TILE_SAND:
				class = "sand"
			case tiledmap.TILE_FOREST:
				class = "forest"
			case tiledmap.TILE_DARKWATER:
				class = "darkwater"
			}
			fmt.Fprintf(w, `<div class="wfc-cell %s"></div>`, class)
		}
	}
}
