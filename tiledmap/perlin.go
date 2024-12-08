package main

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
)

type PerlinNoise struct {
	permutation []int
	gradients   [][2]float64
}

func NewPerlinNoise(seed int64) *PerlinNoise {
	p := &PerlinNoise{
		permutation: make([]int, 256),
		gradients:   make([][2]float64, 256),
	}

	// 初始化置换表
	for i := 0; i < 256; i++ {
		p.permutation[i] = i
	}
	// Fisher-Yates 洗牌算法
	for i := 255; i > 0; i-- {
		j := rand.Intn(i + 1)
		p.permutation[i], p.permutation[j] = p.permutation[j], p.permutation[i]
	}

	// 生成随机梯度向量
	for i := 0; i < 256; i++ {
		angle := rand.Float64() * 2 * math.Pi
		p.gradients[i] = [2]float64{
			math.Cos(angle),
			math.Sin(angle),
		}
	}

	return p
}

func (p *PerlinNoise) Noise2D(x, y float64) float64 {
	// 获取整数坐标
	x0 := int(math.Floor(x))
	y0 := int(math.Floor(y))
	x1 := x0 + 1
	y1 := y0 + 1

	// 计算小数部分
	sx := x - float64(x0)
	sy := y - float64(y0)

	//fmt.Println(x0, y0, x1, y1, sx, sy, x, y)

	// 计算四个角的贡献
	n00 := p.dotGridGradient(x0, y0, x, y)
	n10 := p.dotGridGradient(x1, y0, x, y)
	n01 := p.dotGridGradient(x0, y1, x, y)
	n11 := p.dotGridGradient(x1, y1, x, y)
	//fmt.Println("-", n00, n10, n01, n11)
	// 插值
	sx = p.fade(sx)
	sy = p.fade(sy)

	nx0 := p.lerp(n00, n10, sx)
	nx1 := p.lerp(n01, n11, sx)

	ret := p.lerp(nx0, nx1, sy)
	//fmt.Println("---", nx0, nx1, sy, ret)
	return ret
}

func (p *PerlinNoise) dotGridGradient(ix, iy int, x, y float64) float64 {
	// 获取梯度向量
	idx := p.permutation[ix&255] + p.permutation[iy&255]
	gradient := p.gradients[idx&255]

	// 计算距离向量
	dx := x - float64(ix)
	dy := y - float64(iy)

	// 计算点积
	return dx*gradient[0] + dy*gradient[1]
}

func (p *PerlinNoise) fade(t float64) float64 {
	// 使用平滑函数 6t^5 - 15t^4 + 10t^3
	return t * t * t * (t*(t*6-15) + 10)
}

func (p *PerlinNoise) lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

// 生成柏林噪声迷宫
func GeneratePerlinMaze(size int, scale float64, threshold float64) [][]int {
	perlin := NewPerlinNoise(rand.Int63())
	maze := make([][]int, size)
	for i := range maze {
		maze[i] = make([]int, size)
	}

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			// 生成柏林噪声值
			nx := float64(x) / float64(size) * scale
			ny := float64(y) / float64(size) * scale
			value := perlin.Noise2D(nx, ny)
			//fmt.Println("--", nx, ny, value)
			// 根据阈值确定是墙还是路
			if value > threshold {
				maze[y][x] = 1 // 墙
			} else {
				maze[y][x] = 0 // 路
			}
		}
	}

	// 确保入口和出口是通路
	setEntranceArea(maze, 0, 0)           // 左上角入口
	setEntranceArea(maze, size-3, size-3) // 右下角出口

	return maze
}

const perlinCSS = `
<style>
    .container {
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
        gap: 20px;
        padding: 20px;
        justify-content: center;
    }
    .controls {
        display: flex;
        gap: 10px;
        margin-bottom: 20px;
        width: 100%;
        justify-content: center;
    }
    .maze-box {
        border: 1px solid #ccc;
        padding: 10px;
        flex-shrink: 0;
    }
</style>`

func perlinHandler(w http.ResponseWriter, req *http.Request) {

	fmt.Fprint(w, perlinCSS)

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
	octaves := 4
	lacunarity := 2.0
	persistence := 0.5

	// 修改控制表单，添加 FBM 选项
	fmt.Fprint(w, `<div class="controls">
		<form>
			大小: <input type="number" name="size" value="`+strconv.Itoa(size)+`" min="20" max="512">
			缩放: <input type="number" name="scale" value="`+fmt.Sprintf("%.1f", scale)+`" step="0.1" min="0.1" max="19">
			阈值: <input type="number" name="threshold" value="`+fmt.Sprintf("%.2f", threshold)+`" step="0.05" min="-1" max="1">
			<label><input type="checkbox" name="fbm" value="true" `+func() string {
		if useFBM {
			return "checked"
		}
		return ""
	}()+`> 使用FBM</label>
			<input type="submit" value="生成">
		</form>
	</div>`)

	// 修改生成迷宫的部分
	maze := make([][]int, size)
	for i := range maze {
		maze[i] = make([]int, size)
	}

	perlin := NewPerlinNoise(rand.Int63())
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			nx := float64(x) / float64(size) * scale
			ny := float64(y) / float64(size) * scale
			var value float64
			if useFBM {
				value = perlin.FBM(nx, ny, octaves, lacunarity, persistence)
			} else {
				value = perlin.Noise2D(nx, ny)
			}
			if value > threshold {
				maze[y][x] = 1 // 墙
			} else {
				maze[y][x] = 0 // 路
			}
		}
	}

	// 生成并渲染迷宫
	fmt.Fprint(w, "<div class='container'>")

	renderMazeWithTitle(w, maze, "柏林噪声地图")

	connectRegionsByBFS(maze)
	renderMazeWithTitle(w, maze, "BFS连接所有区域")

	fmt.Fprint(w, "</div>")
}

const perlinGrayCSS = `
<style>
    .gray-container {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 20px;
        padding: 20px;
    }
    .gray-controls {
        display: flex;
        gap: 10px;
        margin-bottom: 20px;
    }
    canvas {
        border: 1px solid #ccc;
    }
</style>`

// 添加 FBM 实现函数
func (p *PerlinNoise) FBM(x, y float64, octaves int, lacunarity, persistence float64) float64 {
	total := 0.0
	frequency := 1.0
	amplitude := 1.0
	maxValue := 0.0

	for i := 0; i < octaves; i++ {
		total += p.Noise2D(x*frequency, y*frequency) * amplitude
		maxValue += amplitude
		frequency *= lacunarity  // 频率增加
		amplitude *= persistence // 振幅减小
	}

	return total / maxValue // 归一化结果
}

func perlinGrayHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, perlinGrayCSS)

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
	fmt.Fprint(w, `<div class="gray-container" style="flex-direction: row; justify-content: center;">
		<div class="gray-controls" style="position: absolute; top: 20px;">
			<form>
				尺寸: <input type="number" name="size" value="`+strconv.Itoa(size)+`" min="64" max="1024">
				缩放: <input type="number" name="scale" value="`+fmt.Sprintf("%.1f", scale)+`" step="0.1" min="0.1" max="20">
				<label><input type="checkbox" name="fbm" value="true" `+func() string {
		if useFBM {
			return "checked"
		}
		return ""
	}()+`> 使用FBM</label>
				<input type="submit" value="生成">
			</form>
		</div>`)

	// 生成三个Canvas
	//scales := []float64{scale / 4, scale / 2, scale, scale * 2, scale * 4}
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
	perlin := NewPerlinNoise(rand.Int63())
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
						// 填充 8x8 的像素块
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

// http://kitfox.com/projects/perlinNoiseMaker/
// https://hmxs.games/posts/107/index.html
// https://zhuanlan.zhihu.com/p/260771031
// https://www.bilibili.com/video/BV19f42197ME
// https://juejin.cn/post/7085186517588525092
// https://omo.moe/archives/394/
// https://replay923.github.io/2018/06/04/PerlinNoise/
// https://www.cnblogs.com/KillerAery/p/10765897.html
// https://juejin.cn/post/7367997561510723620

// todo:
// https://indienova.com/indie-game-development/tinykeepdev-procedural-dungeon-generation-algorithm/
// https://www.gcores.com/articles/168310
