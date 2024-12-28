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
			<link rel="stylesheet" href="/static/style.css">
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

func printBackDiv(w http.ResponseWriter) {
	printSomeBackDiv(1, w)
}

func printSomeBackDiv(count int, w http.ResponseWriter) {
	for i := 0; i < count; i++ {
		fmt.Fprint(w, "</div>") // 关闭 container div
	}
}

func printHtmlHead(w http.ResponseWriter, title string, includepathjs ...bool) {
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>%s</title>
    <link rel="stylesheet" href="/static/style.css">`, title)

	// 检查是否需要包含 path.js
	includeJs := false
	if len(includepathjs) > 0 {
		includeJs = includepathjs[0]
	}

	if includeJs {
		fmt.Fprintf(w, `
    <script src="/static/pathfind.js"></script>`)
	}

	fmt.Fprintf(w, `
</head>
<body>
    `)
}

// 渲染带标题的迷宫
func renderMazeWithTitle(w http.ResponseWriter, maze [][]int, title string) {

	fmt.Fprintf(w, "\n<div class='maze-box'><h3>%s</h3>\n", title)
	fmt.Fprintf(w, "\n<p style='font-size: 12px; margin-top: -15px; color: #666;'>%s</p>\n", "")
	renderMazeWithPath(w, maze, nil, false)
	fmt.Fprintf(w, `</div>`)
}

// 渲染带标题和信息的迷宫
func renderMazePathWithTitle(w http.ResponseWriter, maze [][]int, path [][]bool, title string, info string) {

	fmt.Fprintf(w, "\n<div class='maze-box'><h3>%s</h3>\n", title)
	fmt.Fprintf(w, "\n<p style='font-size: 12px; margin-top: -15px; color: #666;'>%s</p>\n", info)
	renderMazeWithPath(w, maze, path, true)

	fmt.Fprintf(w, `</div>`)
}

func renderMazeWithPath(w http.ResponseWriter, maze [][]int, path [][]bool, showPath bool) {
	size := len(maze)

	fmt.Fprintf(w, `
		<div class="dungeon-container" style="position: relative; width: %dpx; height: %dpx;">`, (size+2)*9, (size+2)*9)

	// dungeon-grid
	fmt.Fprintf(w, `
		<div class="dungeon-grid" style="position: absolute; top: 0; left: 0; display: grid; grid-template-columns: repeat(%d, 8px); grid-template-rows: repeat(%d, 8px);">`, size+2, size+2)

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

	// step-layer
	fmt.Fprintf(w, `
		<div class="step-layer" style="position: absolute; top: 0; left: 0; display: grid; grid-template-columns: repeat(%d, 8px); grid-template-rows: repeat(%d, 8px);">`, size+2, size+2)
	for i := 0; i < size+2; i++ {
		fmt.Fprintf(w, `<div class="step-info"></div>`)
	}
	// 在这里可以添加步数信息或其他辅助信息
	for y := 0; y < size; y++ {
		fmt.Fprintf(w, `<div class="step-info"></div>`)
		for x := 0; x < size; x++ {
			if showPath && path[y][x] {
				fmt.Fprintf(w, `<div class="step-info" style="width: 8px; height: 8px; background-color: rgba(255, 0, 0, 0.5);"></div>`)
			} else {
				fmt.Fprintf(w, `<div class="step-info"></div>`)
			}
		}
		fmt.Fprintf(w, `<div class="step-info"></div>`)
	}
	for i := 0; i < size+2; i++ {
		fmt.Fprintf(w, `<div class="step-info"></div>`)
	}
	fmt.Fprintf(w, `</div>`) // 结束 step-layer
	fmt.Fprintf(w, `</div>`)
}
