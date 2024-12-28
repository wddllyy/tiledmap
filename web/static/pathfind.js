let currentPlayback = null;
let playbackSpeed = 200;

let isPlaying = false;

function togglePlayback() {
    const playbackBtn = document.getElementById('playback-btn');
    
    if (isPlaying) {
        // 暂停播放
        clearInterval(currentPlayback);
        currentPlayback = null;  // 添加这行
        playbackBtn.textContent = '播放';
        console.log('1 isPlaying:', isPlaying);
    } else {
        // 开始播放
        startPlayback();
        playbackBtn.textContent = '暂停';
        console.log('2 isPlaying:', isPlaying);
    }
    
    isPlaying = !isPlaying;
}

function updateSpeed(value) {
    playbackSpeed = parseInt(value);
    document.getElementById('speed-value').textContent = value;
    
    if (currentPlayback) {
        clearInterval(currentPlayback);
        startPlayback();
    }
}

function startPlayback() {
    if (currentPlayback) {
        clearInterval(currentPlayback);
    }
    
    const allStepsData = window.allStepsData;
    
    let step = 0
    //console.log('allStepsData :', allStepsData);

    // 创建单个 renderStep 函数
    function renderStep() {
        Object.keys(allStepsData).forEach(title => {
            const stepsData = allStepsData[title];
            //console.log('title: ', title,' stepsData :', stepsData," stepsData.Steps:", stepsData.Steps);
            if (!stepsData || !stepsData.Steps) {
                //console.error('无效的 stepsData 或 stepsData.Steps:', title);
                return;
            }
            if (step >= stepsData.Steps.length) {
                return;
            }
        
            // 查找或创建步骤显示元素
            const mazeContainer = document.querySelector(`.dungeon-container[data-title="${title}"]`);
            let stepDisplay = document.getElementById(`step-display-${title}`);
            if (!stepDisplay) {
                stepDisplay = document.createElement('div');
                stepDisplay.id = `step-display-${title}`;
                stepDisplay.className = 'step-display';
                mazeContainer.appendChild(stepDisplay);
            }
            
            //console.log('当前步骤数据:', stepsData[step]);
            stepDisplay.textContent = `当前步骤: ${step + 1} / ${stepsData.length}`;
            
            // 绘制当前步骤的单元格
            const currentStepData = stepsData.Steps[step];
            //console.log('stepsData :', stepsData, " step :", step, " currentStepData :", currentStepData);
            updateMazeCell(title, currentStepData);
        });
        step++;
    }
    
    // 使用单个计时器
    currentPlayback = setInterval(renderStep, playbackSpeed);
}

function updateMazeCell(title, stepData) {
    // 1. 首先找到对应的迷宫容器
    const mazeContainer = document.querySelector(`.dungeon-container[data-title="${title}"]`);
    if (!mazeContainer) {
        console.error('找不到迷宫容器:', title);
        return;
    }

    // 2. 找到迷宫网格
    const dungeonGrid = mazeContainer.querySelector('.dungeon-grid');
    if (!dungeonGrid) {
        console.error('找不到迷宫网格:', title);
        return;
    }
    //console.log('dungeonGrid :', dungeonGrid);

    // 3. 从grid-template-columns样式中获取迷宫的宽度
    const gridStyle = window.getComputedStyle(dungeonGrid);
    const columnsCount = gridStyle.gridTemplateColumns.split(' ').length;
    //console.log('gridStyle :', gridStyle, ' columnsCount :', columnsCount);

    // 3. 计算实际的单元格索引（考虑边框）
    // 由于周围有一圈墙壁，所以需要调整位置
    const row = stepData.Pos[0];
    const col = stepData.Pos[1];
    const index = (row + 1) * columnsCount + (col + 1); // +1 是因为周围有墙

    console.log('row :', row, ' col :', col,' index :', index);
    // 4. 获取对应的单元格
    const cell = dungeonGrid.children[index];
    if (!cell) {
        console.error('找不到单元格:', row, col);
        return;
    }

    // 5. 更新 MazeStep 层
    const mazeStepLayer = cell.querySelector('.maze-step');
    if (mazeStepLayer) {
        mazeStepLayer.className = 'maze-step ' + getClassForType(stepData.Type);
    }
}

function getClassForType(type) {
    switch (type) {
        case 0: return 'floor';   // 空地
        case 1: return 'wall';    // 墙
        case 2: return 'checked'; // 已检查的格子
        case 3: return 'start';   // 起点
        case 4: return 'end';     // 终点
        case 5: return 'path';    // 最终路径
        default: return 'floor';
    }
}