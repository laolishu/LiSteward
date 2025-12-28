<template>
    <div class="main-layout">
        <!-- 左侧导航栏 -->
        <div class="sidebar-wrapper" :style="{ width: sidebarWidth + 'px' }">
            <Sidebar @navigate="handleNavigate" :collapsed="sidebarCollapsed" :current-route="props.currentRoute" />
        </div>

        <!-- 拖动分隔符（包含中间收起/展开按钮） -->
        <div class="divider" :class="{ dragging: isDragging }" @mousedown="startDragging">
            <div class="divider-line"></div>
            <button class="divider-toggle" @click.stop="toggleSidebar" @mousedown.stop>
                <span v-if="sidebarCollapsed">▶</span>
                <span v-else>◀</span>
            </button>
        </div>

        <!-- 右侧内容区 -->
        <div class="content-wrapper">
            <!-- 模块内容区 -->
            <div class="content-main">
                <slot />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import Sidebar from '../components/Sidebar.vue'

// 接收currentRoute作为prop
const props = defineProps({
    currentRoute: {
        type: String,
        default: 'hosts'
    }
})

// 定义事件
const emit = defineEmits(['navigate'])

// 状态管理
const sidebarWidth = ref(250) // 侧栏宽度，px
const sidebarCollapsed = ref(false) // 侧栏是否折叠
const isDragging = ref(false) // 是否正在拖动分隔符
const minSidebarWidth = 180 // 最小侧栏宽度
const maxSidebarWidth = 320 // 最大侧栏宽度

// 方法：切换侧栏折叠状态
function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
    if (sidebarCollapsed.value) {
        sidebarWidth.value = 60 // 折叠后只显示图标
    } else {
        sidebarWidth.value = 250 // 展开到默认宽度
    }
}

// 方法：处理导航事件
function handleNavigate(moduleName: string) {
    // 向上级父组件发射事件
    emit('navigate', moduleName)
}

// 方法：开始拖动分隔符
function startDragging(e: MouseEvent) {
    isDragging.value = true
    const startX = e.clientX
    const startWidth = sidebarWidth.value

    function handleMouseMove(moveEvent: MouseEvent) {
        const delta = moveEvent.clientX - startX
        let newWidth = startWidth + delta

        // 限制宽度范围
        if (newWidth < minSidebarWidth) {
            newWidth = minSidebarWidth
        } else if (newWidth > maxSidebarWidth) {
            newWidth = maxSidebarWidth
        }

        sidebarWidth.value = newWidth
    }

    function handleMouseUp() {
        isDragging.value = false
        document.removeEventListener('mousemove', handleMouseMove)
        document.removeEventListener('mouseup', handleMouseUp)
    }

    document.addEventListener('mousemove', handleMouseMove)
    document.addEventListener('mouseup', handleMouseUp)
}

// 生命周期
onMounted(() => {
    // 可以从localStorage恢复侧栏宽度
    const savedWidth = localStorage.getItem('sidebar-width')
    if (savedWidth) {
        const width = parseInt(savedWidth, 10)
        // 确保宽度在允许的范围内
        if (width >= minSidebarWidth && width <= maxSidebarWidth) {
            sidebarWidth.value = width
        } else {
            // 如果超出范围，重置为默认值并保存
            localStorage.setItem('sidebar-width', '250')
        }
    }
})

onUnmounted(() => {
    // 保存侧栏宽度到localStorage
    localStorage.setItem('sidebar-width', sidebarWidth.value.toString())
})
</script>

<style scoped>
/* 主容器 */
.main-layout {
    display: flex;
    height: 100vh;
    width: 100vw;
    background-color: #f6f8fa;
    overflow: hidden;

    /* 设置最小宽高 */
    min-width: 1200px;
    min-height: 780px;
}

/* 左侧导航栏容器 */
.sidebar-wrapper {
    flex-shrink: 0;
    height: 100%;
    min-height: 0;
    background-color: #ffffff;
    border-right: 1px solid #e1e4e8;
    overflow-y: auto;
    overflow-x: hidden;
    transition: width 0.3s ease;
    padding-bottom: 12px;

    /* 自定义滚动条 */
    scrollbar-width: thin;
    scrollbar-color: #d0d7de #f6f8fa;
}

.sidebar-wrapper::-webkit-scrollbar {
    width: 6px;
}

.sidebar-wrapper::-webkit-scrollbar-track {
    background: #f6f8fa;
}

.sidebar-wrapper::-webkit-scrollbar-thumb {
    background: #d0d7de;
    border-radius: 3px;
}

.sidebar-wrapper::-webkit-scrollbar-thumb:hover {
    background: #c9d1d9;
}

/* 拖动分隔符（包含中间收起/展开按钮） */
.divider {
    flex-shrink: 0;
    width: 24px;
    /* 为按钮留出空间 */
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    user-select: none;
}

.divider-line {
    position: absolute;
    left: 50%;
    transform: translateX(-50%);
    width: 4px;
    height: 100%;
    background-color: #e1e4e8;
    cursor: col-resize;
    transition: background-color 0.2s ease;
}

.divider:hover .divider-line,
.divider.dragging .divider-line {
    background-color: #34d399;
}

.divider-toggle {
    position: relative;
    z-index: 5;
    width: 28px;
    height: 28px;
    border-radius: 999px;
    border: 1px solid rgba(0, 0, 0, 0.06);
    background: #ffffff;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 1px 2px rgba(16, 24, 40, 0.04);
    cursor: pointer;
    font-size: 12px;
    color: #6b7280;
    transition: transform 0.12s ease, background-color 0.12s ease, color 0.12s ease;
}

.divider-toggle:hover {
    transform: translateX(-1px);
    background-color: #f6f8fa;
    color: #10b981;
}

.divider-toggle:active {
    transform: translateX(0);
}

/* 右侧内容区 */
.content-wrapper {
    flex: 1;
    display: flex;
    flex-direction: column;
    height: 100%;
    min-height: 0;
    overflow: hidden;
    background-color: #f6f8fa;
}

/* 模块内容区 */
.content-main {
    flex: 1;
    overflow: auto;
    background-color: #f6f8fa;
}

.content-main::-webkit-scrollbar {
    width: 8px;
}

.content-main::-webkit-scrollbar-track {
    background: transparent;
}

.content-main::-webkit-scrollbar-thumb {
    background: #d0d7de;
    border-radius: 4px;
}

.content-main::-webkit-scrollbar-thumb:hover {
    background: #c9d1d9;
}

/* 响应式设计 */
@media (max-width: 1400px) {
    .main-layout {
        min-width: 1000px;
    }
}

@media (max-width: 1000px) {
    .main-layout {
        min-width: 800px;
    }

    .sidebar-wrapper {
        max-width: 300px;
    }
}

@media (max-width: 800px) {
    .main-layout {
        min-width: 600px;
        min-height: 600px;
    }

    .sidebar-wrapper {
        max-width: 250px;
    }
}
</style>
