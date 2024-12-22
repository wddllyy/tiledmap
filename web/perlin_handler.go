package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"mazemap/tiledmap"
)

func perlinHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, RenderCSS)

	// 解析参数
	size := 50
	if s := req.URL.Query().Get("size"); s != "" {
		if val, err := strconv.Atoi(s); err == nil && val > 0 && val <= 512 {
			size = val
		}
	}

	scale := 5.0
	if s := req.URL.Query().Get("scale"); s != "" {
		if val, err := strconv.ParseFloat(s, 64); err == nil && val > 0 {
			scale = val
		}
	}

	threshold := 0.0
	if t := req.URL.Query().Get("threshold"); t != "" {
		if val, err := strconv.ParseFloat(t, 64); err == nil {
			threshold = val
		}
	}

	// 添加 FBM 参数
	useFBM := req.URL.Query().Get("fbm") == "true"

	// 修改控制表单，添加 FBM 选项
	fmt.Fprintf(w, `
<div class="dungeon-container">
	<div class="dungeon-controls">
		<form>
			大小: <input type="number" name="size" value="%d" min="20" max="512">
			缩放: <input type="number" name="scale" value="%.1f" step="0.1" min="0.1" max="19">
			阈值: <input type="number" name="threshold" value="%.2f" step="0.05" min="-1" max="1">
			<label><input type="checkbox" name="fbm" value="true" %s> 使用FBM</label>
			<input type="submit" value="生成">
		</form>
	</div>
	<div style="display: flex; gap: 20px; justify-content: center;">`,
		size, scale, threshold,
		func() string {
			if useFBM {
				return "checked"
			}
			return ""
		}())

	// 生成迷宫

	maze := tiledmap.GeneratePerlinMaze(size, scale, threshold, useFBM)
	renderMazeWithTitle(w, maze, "柏林噪声地图")

	tiledmap.ConnectRegionsByBFS(maze)
	renderMazeWithTitle(w, maze, "BFS连接所有区域")

	fmt.Fprint(w, "\n</div></div></body></html>")
}

func perlinGrayHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, RenderCSS)

	// 解析参数
	size := 256
	if s := req.URL.Query().Get("size"); s != "" {
		if val, err := strconv.Atoi(s); err == nil && val > 0 && val <= 1024 {
			size = val
		}
	}

	scale := 5.0
	if s := req.URL.Query().Get("scale"); s != "" {
		if val, err := strconv.ParseFloat(s, 64); err == nil && val > 0 {
			scale = val
		}
	}

	// 添加 FBM 参数
	useFBM := req.URL.Query().Get("fbm") == "true"
	octaves := 4
	lacunarity := 2.0
	persistence := 0.5

	// 修改表单，添加 FBM 选项
	fmt.Fprintf(w, `
<div class="gray-container" style="flex-direction: row; justify-content: center;">
	<div class="gray-controls" style="position: absolute; top: 20px;">
		<form>
			尺寸: <input type="number" name="size" value="%d" min="64" max="1024">
			缩放: <input type="number" name="scale" value="%.1f" step="0.1" min="0.1" max="20">
			<label><input type="checkbox" name="fbm" value="true" %s> 使用FBM</label>
			<input type="submit" value="生成">
		</form>
	</div>`,
		size, scale,
		func() string {
			if useFBM {
				return "checked"
			}
			return ""
		}())

	// 生成三个Canvas
	scales := []float64{scale}
	for i, currentScale := range scales {
		fmt.Fprintf(w, `
			<div style="margin: 60px 10px 0 10px; text-align: center;">
				<div style="margin-bottom: 10px;">Scale: %.1f</div>
				<canvas id="perlinCanvas%d" width="%d" height="%d"></canvas>
			</div>`, currentScale, i, size, size)
	}

	// 生成JavaScript代码
	fmt.Fprint(w, "<script>")

	// 为每个Canvas生成和渲染噪声数据
	perlin := tiledmap.NewPerlinNoise(rand.Int63())
	for i, currentScale := range scales {
		fmt.Fprintf(w, `
			{
				const canvas = document.getElementById('perlinCanvas%d');
				const ctx = canvas.getContext('2d');
				ctx.imageSmoothingEnabled = false; // 禁用平滑处理
				const pixelSize = 1; // 设置每个像素的实际大小
				canvas.width = %d * pixelSize;
				canvas.height = %d * pixelSize;
				const imageData = ctx.createImageData(canvas.width, canvas.height);
				const data = imageData.data;
				const noiseData = [`, i, size, size)

		for y := 0; y < size; y++ {
			for x := 0; x < size; x++ {
				nx := float64(x) / float64(size) * currentScale
				ny := float64(y) / float64(size) * currentScale
				var value float64
				if useFBM {
					value = (perlin.FBM(nx, ny, octaves, lacunarity, persistence) + 1) / 2
				} else {
					value = (perlin.Noise2D(nx, ny) + 1) / 2
				}
				fmt.Fprintf(w, "%.4f,", value)
			}
		}

		fmt.Fprintf(w, `];
				for (let y = 0; y < %d; y++) {
					for (let x = 0; x < %d; x++) {
						const value = Math.floor(noiseData[y * %d + x] * 255);
						// 填充像素
						for (let py = 0; py < pixelSize; py++) {
							for (let px = 0; px < pixelSize; px++) {
								const idx = ((y * pixelSize + py) * (%d * pixelSize) + (x * pixelSize + px)) * 4;
								data[idx] = value;     // R
								data[idx + 1] = value; // G
								data[idx + 2] = value; // B
								data[idx + 3] = 255;   // A
							}
						}
					}
				}
				ctx.putImageData(imageData, 0, 0);
			}`, size, size, size, size)
	}

	fmt.Fprint(w, "</script></div>")
}
