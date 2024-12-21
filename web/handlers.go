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
