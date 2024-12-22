package tiledmap

import (
	"math"
	"math/rand"
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
func GeneratePerlinMaze(size int, scale float64, threshold float64, useFBM bool) [][]int {

	octaves := 4
	lacunarity := 2.0
	persistence := 0.5

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
			var value float64
			if useFBM {
				value = perlin.FBM(nx, ny, octaves, lacunarity, persistence)
			} else {
				value = perlin.Noise2D(nx, ny)
			}
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
// https://juejin.cn/post/7119679952575954975?from=search-suggest
