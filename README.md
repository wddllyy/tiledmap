# tiledmap
一些生成tiledmap的小尝试：

***



## maze：
基于dfs生成maze，然后在此基础上生成tiledmap

### 参数解释:

"http://127.0.0.1:9999/maze?size=59&turn=0.8&acc=0.95&erosion=0.4"


1. size控制生成的地图大小，暂时只支持奇数
2. turn用来控制生成maze时，路径转向的偏好，越大则越倾向于转向，越小则越倾向于沿用之前的方向
3. acc用来控制堆积系数（在不影响连通性的前提下），越大则堆积的障碍物越多越集中，越小则堆积的障碍物越少越分散
4. erosion用来控制侵蚀系数，越大则侵蚀的越厉害（越空旷），越小则侵蚀的越少


### 原理：
1. 先生成一个间隔挖空的初始地图，然后利用dfs联通所有的挖空区域。这样形成了起点到终点之间有一个唯一路径的maze。（这个过程中turn用来控制生成maze时，路径转向的偏好，该数值越大，每次联通一个tile，下一个tile更倾向于转向（而不是沿用之前的方向），所以该参数越小，生成的迷宫的路径就越倾向于直来直去，如果该数值越大，则倾向于生成的路径越扭来扭去。）
2. 在此maze的基础上，通过堆积，填充所有断头路（三面都是死路的路径点），来生成一些大面积的障碍物，这里通过acc参数控制，acc为1时，就是把所有的岔路都填充满；acc为0时，不填充任何岔路
3. 在填充完之后，用erosion来模拟侵蚀的过程，把所有的路径侵蚀的宽一些。更像自然地貌。侵蚀的时候，越孤立的障碍物越容易被侵蚀。erosion为1时，则所有的障碍物都被侵蚀掉了

总体上是一个：1.生成maze和唯一路径；2.填充岔路,堆积障碍；3.侵蚀障碍，扩宽通路的过程。整体成本较低，可控参数较简单

![image](https://github.com/wddllyy/tiledmap/blob/main/doc/IMG/Screenshot_maze.png)

### TODO: 
 1. maze生成时，控制生成的唯一路径的长度
 2. 侵蚀时，控制不要对唯一路径长度影响太大
 3. 利用细胞自动机来做侵蚀
   
   
***


## cellular：
基于细胞自动机的方案来生成tiledmap：

### 参数解释:



"http://127.0.0.1:9999/cellular?size=85&probability=0.6&iterations=5"

1. size 是尺寸
2. probability 是随机生成的初始地图里障碍块的占比
3. iterations 是细胞自动机的迭代次数，次数越多地图越规整

###原理：
1. 初始化一个size*size的地图，随机生成障碍块，障碍块占比为probability
2. 迭代计算每个cell的邻居cell的障碍块数量，如果邻居cell的障碍块数量小于4，则该cell变为空白块；若大于5，则该cell变障碍物
3. 上述过程迭代iterations次
4. 然后计算地图中所有的联通区域，这些区域可能彼此不联通。因此需要一个方法将他们自然地联通起来
    1. 计算所有的联通区域
    2. 找到联通区域的边界，然后BFS向外探索，直到和其他联通区域相连，然后回溯重建最短路径
    3. 重复上述过程直到把所有区域连接起来

![image](https://github.com/wddllyy/tiledmap/blob/main/doc/IMG/Screenshot_cellular.png)

### TODO:
1. 利用分块不同随机率来做到一些对比更鲜明的堆积和稀疏

***

## perlin：
基于perlin噪声的方案来生成tiledmap：

### 参数解释:
"http://127.0.0.1:9999/perlin?size=200&scale=4.0&threshold=0&fbm=true"

1. size 是尺寸
2. scale 是perlin噪声的缩放系数
3. threshold 是perlin噪声的阈值，大于该阈值的值会被设置为障碍物
4. fbm 是是否使用fbm，fbm是分形布朗运动，可以生成更自然的噪声

### 原理：
1. 生成perlin噪声图
2. 根据阈值生成障碍物（如果使用fbm，则对噪声图进行fbm处理，然后根据阈值生成障碍物）
3. 计算所有联通区域，并连接他们

![image](https://github.com/wddllyy/tiledmap/blob/main/doc/IMG/Screenshot_perlin.png)
![image](https://github.com/wddllyy/tiledmap/blob/main/doc/IMG/Screenshot_perlinFBM.png)



### TODO:
1 增加更多可控性，比如保证起点终点有平缓自然的路径

***

## dungeon：
基于随机房间的方案来生成tiledmap：

### 参数解释:
"http://localhost:9999/dungeon?width=51&height=51&rooms=18&minSize=5&maxSize=9&extraPathProb=0.0"

1. size 是尺寸
2. rooms 是房间数量
3. minSize 是最小房间尺寸
4. maxSize 是最大房间尺寸
5. extraPathProb 是额外路径的概率

### 原理：
1. 生成指定尺寸的初始地图，地图中随机散布房间，房间尺寸在minSize和maxSize之间（因为是随机散布，所以不能保证生成足够的房间数量）
2. 在房间之间的空地上，随机生成路径，路径的宽度为1
3. 计算所有联通区域（包括房间和路径），并随机连接他们
4. 把所有死胡同堵上，降低maze部分的难度

更详细的过程参考: https://journal.stuffwithstuff.com/2014/12/21/rooms-and-mazes/

![image](https://github.com/wddllyy/tiledmap/blob/main/doc/IMG/Screenshot_dungeon.png)

### TODO:
1. 性能优化

