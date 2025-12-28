<template>
    <div class="sidebar">
        <div class="sidebar-header">
            <div class="app-logo">
                <div class="logo-icon">
                    <img src="/logo-icon.png" alt="logo" class="logo-img" />
                </div>
                <span class="app-name">{{ appName }}</span>
            </div>
        </div>

        <nav class="sidebar-nav">
            <a v-for="item in menuItems" :key="item.id" :class="['nav-item', { active: currentRoute === item.id }]"
                @click="$emit('navigate', item.id)">
                <img :src="item.icon" :alt="item.label" class="nav-icon" />
                <span class="nav-text">{{ item.label }}</span>
            </a>
        </nav>
    </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { GetAppConfig } from '../../wailsjs/go/main/App'

const { t } = useI18n()

const props = defineProps({
    currentRoute: {
        type: String,
        default: 'hosts'
    }
})

defineEmits(['navigate'])

const appName = ref('LiSteward')

onMounted(async () => {
    try {
        const cfg = await GetAppConfig()
        if (cfg) {
            if (cfg.productName && cfg.productName !== '') appName.value = cfg.productName
            else if (cfg.productNameEn && cfg.productNameEn !== '') appName.value = cfg.productNameEn
        }
    } catch (e) {
        // ignore and keep default
    }
})

const menuItems = computed(() => [
    { id: 'hosts', icon: '/images/menu-hosts.svg', label: t('menu.hosts') },
    { id: 'node', icon: '/images/menu-node.svg', label: t('menu.node') },
    { id: 'settings', icon: '/images/menu-settings.svg', label: t('menu.settings') },
])
</script>

<style scoped>
.sidebar {
    width: 100%;
    height: 100%;
    background-color: #1a1c2c;
    display: flex;
    flex-direction: column;
    border-right: 1px solid #2a2d3a;
    box-shadow: none;
}

.sidebar-header {
    padding: 16px 16px;
    border-bottom: 1px solid #2a2d3a;
    background-color: #1a1c2c;
    flex-shrink: 0;
}

.app-logo {
    display: flex;
    align-items: center;
    gap: 8px;
}

.logo-icon {
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
}

.logo-icon .logo-img {
    width: 28px;
    height: 28px;
    object-fit: contain;
    filter: drop-shadow(0 0 4px rgba(52, 211, 153, 0.3));
}

.app-name {
    font-size: 16px;
    font-weight: 700;
    color: #34d399;
    white-space: nowrap;
}

.sidebar-nav {
    flex: 1;
    padding: 12px 8px;
    overflow-y: auto;
}

.nav-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 12px;
    margin-bottom: 2px;
    border-radius: 6px;
    color: #8b92a0;
    cursor: pointer;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    text-decoration: none;
    font-size: 13px;
}

.nav-item:hover {
    background-color: rgba(52, 211, 153, 0.1);
    color: #34d399;
    transform: translateX(4px);
}

.nav-item.active {
    background-color: rgba(52, 211, 153, 0.15);
    color: #34d399;
    border-left: 3px solid #34d399;
    padding-left: 13px;
}

.nav-icon {
    font-size: 18px;
    transition: transform 0.3s;
    flex-shrink: 0;
    width: 18px;
    height: 18px;
    object-fit: contain;
    display: flex;
    align-items: center;
    justify-content: center;
}

.nav-item:hover .nav-icon {
    transform: scale(1.1);
    filter: brightness(1.2);
}

.nav-text {
    font-size: 13px;
    font-weight: 500;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

/* 滚动条样式 */
.sidebar-nav::-webkit-scrollbar {
    width: 4px;
}

.sidebar-nav::-webkit-scrollbar-track {
    background: transparent;
}

.sidebar-nav::-webkit-scrollbar-thumb {
    background: #3a3d4a;
    border-radius: 2px;
}

.sidebar-nav::-webkit-scrollbar-thumb:hover {
    background: #4a4d5a;
}
</style>
