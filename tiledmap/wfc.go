// wfc.go
package tiledmap

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
)

// Tile类型常量
const (
	TILE_GRASS = iota
	TILE_WATER
	TILE_SAND
	TILE_FOREST
	TILE_DARKWATER // 新增
)

// WFC结构体
type WFC struct {
	width, height int
	cells         [][]WFCCell
	tileRules     map[int][]int
	maxEntropy    int
}

type WFCCell struct {
	collapsed bool
	options   []int
	entropy   int
}

func NewWFC(width, height int) *WFC {
	wfc := &WFC{
		width:      width,
		height:     height,
		cells:      make([][]WFCCell, height),
		tileRules:  make(map[int][]int),
		maxEntropy: 5,
	}

	wfc.initTileRules()

	for y := 0; y < height; y++ {
		wfc.cells[y] = make([]WFCCell, width)
		for x := 0; x < width; x++ {
			wfc.cells[y][x] = WFCCell{
				collapsed: false,
				options:   []int{TILE_GRASS, TILE_WATER, TILE_SAND, TILE_FOREST, TILE_DARKWATER},
				entropy:   5,
			}
		}
	}

	return wfc
}

func (w *WFC) initTileRules() {
	w.tileRules[TILE_GRASS] = []int{TILE_GRASS, TILE_FOREST, TILE_SAND}
	w.tileRules[TILE_WATER] = []int{TILE_WATER, TILE_SAND, TILE_DARKWATER}
	w.tileRules[TILE_SAND] = []int{TILE_SAND, TILE_GRASS, TILE_WATER}
	w.tileRules[TILE_FOREST] = []int{TILE_FOREST, TILE_GRASS}
	w.tileRules[TILE_DARKWATER] = []int{TILE_DARKWATER, TILE_WATER}
}

func (w *WFC) Generate() [][]int {
	for !w.isFullyCollapsed() {
		x, y := w.findLowestEntropy()
		w.collapseCell(x, y)
		w.propagate(x, y)
	}

	result := make([][]int, w.height)
	for y := 0; y < w.height; y++ {
		result[y] = make([]int, w.width)
		for x := 0; x < w.width; x++ {
			result[y][x] = w.cells[y][x].options[0]
		}
	}

	return result
}

func (w *WFC) isFullyCollapsed() bool {
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			if !w.cells[y][x].collapsed {
				return false
			}
		}
	}
	return true
}

func (w *WFC) findLowestEntropy() (int, int) {
	minEntropy := w.maxEntropy + 1
	var candidates [][2]int

	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			cell := &w.cells[y][x]
			if !cell.collapsed && cell.entropy > 0 {
				if cell.entropy < minEntropy {
					minEntropy = cell.entropy
					candidates = candidates[:0]
					candidates = append(candidates, [2]int{x, y})
				} else if cell.entropy == minEntropy {
					candidates = append(candidates, [2]int{x, y})
				}
			}
		}
	}

	if len(candidates) > 0 {
		chosen := candidates[rand.Intn(len(candidates))]
		return chosen[0], chosen[1]
	}
	return 0, 0
}

func (w *WFC) collapseCell(x, y int) {
	cell := &w.cells[y][x]
	if len(cell.options) > 0 {
		// 简单随机选择，不考虑权重
		chosenIndex := rand.Intn(len(cell.options))
		chosenValue := cell.options[chosenIndex]
		cell.options = []int{chosenValue}
		cell.collapsed = true

		// 修正：熵应该是剩余选项的数量
		cell.entropy = len(cell.options)
	}
}

func (w *WFC) propagate(startX, startY int) {
	queue := [][2]int{{startX, startY}}
	dx := []int{0, 1, 0, -1}
	dy := []int{-1, 0, 1, 0}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		x, y := current[0], current[1]

		currentOptions := w.cells[y][x].options

		for i := 0; i < 4; i++ {
			newX := x + dx[i]
			newY := y + dy[i]

			if newX >= 0 && newX < w.width && newY >= 0 && newY < w.height {
				neighbor := &w.cells[newY][newX]
				if !neighbor.collapsed {
					oldLen := len(neighbor.options)
					newOptions := make([]int, 0)
					for _, option := range neighbor.options {
						valid := false
						for _, currentOption := range currentOptions {
							if w.canBeNeighbors(currentOption, option) {
								valid = true
								break
							}
						}
						if valid {
							newOptions = append(newOptions, option)
						}
					}

					neighbor.options = newOptions
					neighbor.entropy = len(newOptions)

					if oldLen != len(newOptions) {
						queue = append(queue, [2]int{newX, newY})
					}
				}
			}
		}
	}
}

func (w *WFC) canBeNeighbors(tile1, tile2 int) bool {
	validNeighbors := w.tileRules[tile1]
	for _, valid := range validNeighbors {
		if valid == tile2 {
			return true
		}
	}
	return false
}

// HTML模板常量
const wfcHtmlTemplate = `
    <style>
        .container {
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 20px;
            padding: 20px;
        }
        .controls {
            display: flex;
            gap: 10px;
            margin-bottom: 20px;
        }
        table {
            border-collapse: collapse;
        }
        td {
            width: 20px;
            height: 20px;
            padding: 0;
        }
        .grass { background-color: #228B22; }
        .water { background-color: #4169E1; }
        .sand { background-color: #EED6AF; }
        .forest { background-color: #1B6B1B; }
        .darkwater { background-color: #4169E1; }
    </style>
<div class="container">
    <div class="controls">
        <form>
            宽度: <input type="number" name="width" value="%d" min="1" max="100">
            高度: <input type="number" name="height" value="%d" min="1" max="100">
            <input type="submit" value="生成">
        </form>
    </div>
`

const htmlFooter = `
    </table>
    </div>
`

func WFCHandler(w http.ResponseWriter, r *http.Request) {
	// 从URL参数获取宽度和高度
	width, err := strconv.Atoi(r.URL.Query().Get("width"))
	if err != nil || width <= 0 {
		width = 32 // 默认值
	}
	height, err := strconv.Atoi(r.URL.Query().Get("height"))
	if err != nil || height <= 0 {
		height = 32 // 默认值
	}

	wfc := NewWFC(width, height)
	tileMap := wfc.Generate()

	// HTML头部
	fmt.Fprintf(w, wfcHtmlTemplate, width, height)

	// 生成表格
	fmt.Fprintf(w, "<table>")
	for y := 0; y < height; y++ {
		fmt.Fprintf(w, "<tr>")
		for x := 0; x < width; x++ {
			var class string
			switch tileMap[y][x] {
			case TILE_GRASS:
				class = "grass"
			case TILE_WATER:
				class = "water"
			case TILE_SAND:
				class = "sand"
			case TILE_FOREST:
				class = "forest"
			case TILE_DARKWATER:
				class = "darkwater"
			}
			fmt.Fprintf(w, `<td class="%s"></td>`, class)
		}
		fmt.Fprintf(w, "</tr>")
	}

	fmt.Fprint(w, htmlFooter)
}
