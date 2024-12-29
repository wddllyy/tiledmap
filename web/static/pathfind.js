let currentPlayback = null;
let playbackSpeed = 200;

let isPlaying = false;
window.step = 0;  // 添加初始化
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

function stepPlayback() {
    renderAllSteps();
}

function updateSpeed(value) {
    playbackSpeed = parseInt(value);
    document.getElementById('speed-value').textContent = value;
    
    if (currentPlayback) {
        clearInterval(currentPlayback);
        startPlayback();
    }
}

function renderAllSteps() {
    Object.keys(window.allStepsData).forEach(title => {
        const stepsData = window.allStepsData[title];
        if (!stepsData || !stepsData.Steps || window.step >= stepsData.Steps.length) {
            return;
        }
    
        // 查找迷宫容器
        const mazeContainer = document.querySelector(`.maze-box[data-title="${title}"]`);
        if (!mazeContainer) {
            console.error('找不到迷宫容器:', title);
            return;
        }

        // 更新步骤显示
        let stepDisplay = mazeContainer.querySelector('.step-display');
        if (!stepDisplay) {
            stepDisplay = document.createElement('div');
            stepDisplay.className = 'step-display';
            mazeContainer.appendChild(stepDisplay);
        }
        stepDisplay.textContent = `步骤: ${window.step + 1} / ${stepsData.Steps.length}`;
        
        // 更新单元格
        const currentStepData = stepsData.Steps[window.step];
        updateStepCell(title, currentStepData);
    });
    window.step++;
}

function startPlayback() {
    if (currentPlayback) {
        clearInterval(currentPlayback);
    }
    
    currentPlayback = setInterval(renderAllSteps, playbackSpeed);
}

function updateStepCell(title, stepData) {
    // 找到对应的迷宫容器
    const mazeContainer = document.querySelector(`.maze-box[data-title="${title}"]`);
    if (!mazeContainer) {
        console.error('找不到迷宫容器:', title);
        return;
    }

    // 找到step-layer
    const stepLayer = mazeContainer.querySelector('.step-layer');
    if (!stepLayer) {
        console.error('找不到step-layer:', title);
        return;
    }
    //console.log("stepData:", stepLayer, " stepLayer.style.gridTemplateColumns", stepLayer.style.gridTemplateColumns);
    // 计算在step-layer中的位置
    const size = parseInt(stepLayer.style.gridTemplateColumns.match(/repeat\((\d+), 8px\)/)[1]);
    const row = stepData.Pos[0];
    const col = stepData.Pos[1];
    const index = (row + 1) * (size) + (col + 1); // +1 是因为周围有墙
    //console.log("size:", size, "row:", row, "col:", col, "index:", index);

    // 获取对应的step-info单元格
    const stepInfo = stepLayer.children[index];
    if (!stepInfo) {
        console.error('找不到step-info单元格:', row, col);
        return;
    }

    // 根据类型添加新的点
    if (stepData.Type === "pop") { // 已检查的格子
        const dot = document.createElement('div');
        dot.className = 'step-info pop-dot';
        stepInfo.appendChild(dot);
    } else if (stepData.Type === "push") { // 路径
        const dot = document.createElement('div');
        dot.className = 'step-info push-dot';
        stepInfo.appendChild(dot);
    }
}
