package main

import (
	"fmt"
	"net/http"
)

// handler echoes r.URL.Path
func indexHandler(w http.ResponseWriter, req *http.Request) {
	html := `
	<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; margin: 40px; }
				.algorithms { margin-top: 20px; }
				.algorithms h2 { color: #333; }
				.algorithms ul { list-style-type: none; padding: 0; }
				.algorithms li { margin: 10px 0; }
				.algorithms a {
					display: inline-block;
					padding: 8px 16px;
					background-color: #f0f0f0;
					color: #333;
					text-decoration: none;
					border-radius: 4px;
				}
				.algorithms a:hover { background-color: #ddd; }
			</style>
		</head>
		<body>
			<h1>地图生成算法演示</h1>
			<div class="algorithms">
				<h2>可用算法</h2>
				<ul>
					<li><a href="/cellular">元胞自动机 (Cellular Automata)</a></li>
					<li><a href="/dungeon">地下城生成器 (Dungeon Generator)</a></li>
					<li><a href="/maze">迷宫生成器 (Maze Generator)</a></li>
					<li><a href="/perlin">柏林噪声地图 (Perlin Noise Map)</a></li>
					<li><a href="/wfc">波函数坍缩 (Wave Function Collapse)</a></li>
				</ul>
			</div>
			<div class="algorithms">
				<h2>调试工具</h2>
				<ul>
					<li><a href="/hello">查看 Header 信息</a></li>
					<li><a href="/test/hello">字符串反转测试</a></li>
				</ul>
			</div>
		</body>
	</html>
	`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, html)
}

// handler echoes r.URL.Header
func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}

// handler reverses the subpath string
func testHandler(w http.ResponseWriter, req *http.Request) {
	subPath := req.URL.Path[len("/test/"):]
	runes := []rune(subPath)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	fmt.Fprintf(w, "Reversed path: %s\n", string(runes))
}

const RenderCSS = `
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
	.path { background-color: #339966; }
	.grass { background-color: #228B22; }
    .water { background-color: #4169E1; }
    .sand { background-color: #EED6AF; }
    .forest { background-color: #1B6B1B; }
    .darkwater { background-color: #4169E1; }
	.start { background-color: #0f0; }
	.end { background-color: #f00; }
    .room-number {
        position: absolute;
        color: #f00;
        font-size: 12px;
        font-weight: bold;
        z-index: 1;
    }
</style>`

func printBackDiv(w http.ResponseWriter) {
	printSomeBackDiv(1, w)
}

func printSomeBackDiv(count int, w http.ResponseWriter) {
	for i := 0; i < count; i++ {
		fmt.Fprint(w, "</div>") // 关闭 container div
	}
}

// 渲染带标题的迷宫
func renderMazeWithTitle(w http.ResponseWriter, maze [][]int, title string) {
	fmt.Fprintf(w, "\n<div class='maze-box'><h3>%s</h3>\n", title)
	renderMazeWithPath(w, maze, nil, false)
	fmt.Fprint(w, "\n</div>\n")
}

// 渲染带标题和信息的迷宫
func renderMazePathWithTitle(w http.ResponseWriter, maze [][]int, path [][]bool, title string, info string) {
	fmt.Fprintf(w, "\n<div class='maze-box'><h3>%s</h3>\n", title)
	fmt.Fprintf(w, "\n<p style='font-size: 12px; margin-top: -15px; color: #666;'>%s</p>\n", info)
	renderMazeWithPath(w, maze, path, true)
	fmt.Fprint(w, "\n</div>\n")
}

func renderMazeWithPath(w http.ResponseWriter, maze [][]int, path [][]bool, showPath bool) {
	size := len(maze)

	fmt.Fprintf(w, `
		<div class="dungeon-grid" style="grid-template-columns: repeat(%d, 8px);">`, size+2)

	for i := 0; i < size+2; i++ {
		fmt.Fprintf(w, `<div class="dungeon-cell wall"></div>`)
	}

	for y := 0; y < size; y++ {
		fmt.Fprintf(w, `<div class="dungeon-cell wall"></div>`)
		for x := 0; x < size; x++ {
			cellClass := "wall"
			if maze[y][x] == 0 {
				cellClass = "floor"
			}
			//fmt.Println(y, x)
			if showPath && path[y][x] {
				cellClass = "path"
			}
			fmt.Fprintf(w, `<div class="dungeon-cell %s"></div>`, cellClass)
		}
		fmt.Fprintf(w, `<div class="dungeon-cell wall"></div>`)
	}

	for i := 0; i < size+2; i++ {
		fmt.Fprintf(w, `<div class="dungeon-cell wall"></div>`)
	}
	fmt.Fprintf(w, `</div>`)
}
