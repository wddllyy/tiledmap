# tiledmap
一些生成tiledmap的小尝试
##1. maze
   这是最基本的，基于生成tiledmaze，进一步来生成地图。
   下面是一个参数模拟"http://127.0.0.1:9999/maze?size=59&turn=0.8&acc=0.95&erosion=0.4"
   
   
    1. size控制生成的地图大小，暂时只支持奇数
    2. turn用来控制生成maze是，路径转向的偏好，该数值越大，每次联通一个tile，下一个tile更倾向于转向（而不是沿用之前的方向），所以该参数越小，生成的迷宫的路径就越倾向于直来直去，如果该数值越大，则倾向于生成的路径越扭来扭去。
    3. 当生成完maze之后，就是一个有唯一通路的一个maze
    4. 在此maze的基础上，通过堆积，填充所有断头路（三面都是死路的路径点），来生成一些大面积的障碍物，这里通过acc参数控制，acc为1时，就是把所有的岔路都填充满；acc为0时，不填充任何岔路。
    5. 在填充完之后，用erosion来模拟侵蚀的过程，把所有的路径侵蚀的宽一些。更像自然地貌。侵蚀的时候，越孤立的障碍物越容易被侵蚀。erosion为1时，则所有的障碍物都被侵蚀掉了。

    总体上是一个：1.生成maze和唯一路径；2.填充岔路,堆积障碍；3.侵蚀障碍，扩宽通路的过程。整体成本较低，可控参数较简单。
    ![image](https://github.com/wddllyy/tiledmap/blob/main/doc/IMG/Screenshot_maze.png)

    TODO: 
        1. maze生成时，控制生成的唯一路径的长度。
        2. 侵蚀时，控制不要对唯一路径长度影响太大。
   
   
    
##2. cellular
    是一个基于细胞自动机的方案来生成tiledmap
    下面是一个参数模拟："http://127.0.0.1:9999/cellular?size=85&probability=0.6&iterations=5"

    1. size 是尺寸
    2. probability 是随机生成的初始地图里障碍块的占比
    3. iterations 是细胞自动机的迭代次数，次数越多地图越规整

    原理：
    1. 初始化一个size*size的地图，随机生成障碍块，障碍块占比为probability
    2. 迭代计算每个cell的邻居cell的障碍块数量，如果邻居cell的障碍块数量小于4，则该cell变为空白块；若大于5，则该cell变障碍物
    3. 上述过程迭代iterations次
    4. 然后计算地图中所有的联通区域，这些区域可能彼此不联通。因此需要一个方法将他们自然地联通起来
       1. 计算所有的联通区域
       2. 找到联通区域的边界，然后BFS向外探索，直到和其他联通区域相连，然后回溯重建最短路径
       3. 重复上述过程直到把所有区域连接起来
    
    ![image](https://github.com/wddllyy/tiledmap/blob/main/doc/IMG/Screenshot_cellular.png)

