<script setup lang="ts">
import { ref } from 'vue'
import MainLayout from './layouts/MainLayout.vue'
import HostsManager from './views/HostsManager.vue'
import Settings from './views/Settings.vue'
import NodeManager from './views/NodeManager.vue'

// 路由状态
const currentModule = ref('hosts')
const mainLayoutRef = ref<InstanceType<typeof MainLayout> | null>(null)

// 处理模块导航
const handleNavigate = (moduleName: string) => {
  currentModule.value = moduleName
}

// 处理同步状态更新
const updateSyncStatus = (syncing: boolean, error: boolean = false) => {
  mainLayoutRef.value?.setSyncStatus(syncing, error)
}

// 暴露给子组件使用
defineExpose({
  updateSyncStatus,
})
</script>

<template>
  <MainLayout ref="mainLayoutRef" @navigate="handleNavigate" :current-route="currentModule">
    <!-- 根据当前模块显示不同的内容 -->
    <HostsManager v-if="currentModule === 'hosts'" @navigate="handleNavigate" />
    <Settings v-else-if="currentModule === 'settings'" @navigate="handleNavigate" />
    <NodeManager v-else-if="currentModule === 'node'" @navigate="handleNavigate" />

    <div v-else class="module-placeholder">
      <h2>{{ currentModule }} 模块</h2>
      <p>功能开发中...</p>
    </div>
  </MainLayout>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html,
body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  overflow: hidden;
  height: 100%;
  width: 100%;
}

#app {
  height: 100vh;
  width: 100vw;
  overflow: hidden;
}

.module-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #666;
  background-color: #f6f8fa;
}

.module-placeholder h2 {
  font-size: 24px;
  margin-bottom: 12px;
  color: #24292e;
}

.module-placeholder p {
  font-size: 16px;
  color: #57606a;
}
</style>
