<template>
    <div class="module-header-wrapper">
        <div class="module-header">
            <div class="header-content">
                <h1 class="header-title">{{ title }}</h1>
                <p v-if="subtitle" class="header-subtitle">{{ subtitle }}</p>
            </div>
        </div>

        <!-- 提示框 -->
        <transition name="hint-fade">
            <div v-if="props.showHint && internalShowHint && hint" :class="['module-hint', `hint-${hint.type}`]"
                :role="hint.type === 'error' ? 'alert' : 'status'" aria-live="polite">
                <div class="hint-content">
                    <img v-if="hint.type === 'success'" src="/images/message-success.svg" alt="success"
                        class="hint-icon" />
                    <img v-else-if="hint.type === 'error'" src="/images/message-error.svg" alt="error"
                        class="hint-icon" />
                    <span class="hint-message">{{ hint.message }}</span>
                </div>
                <button class="hint-close" @click="closeHint" :aria-label="`关闭${hint.type}提示`">
                    ×
                </button>
            </div>
        </transition>
    </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'

const props = defineProps({
    title: {
        type: String,
        required: true
    },
    subtitle: {
        type: String,
        default: null
    },
    hint: {
        type: Object,
        default: null
    },
    showHint: {
        type: Boolean,
        default: true
    }
})

const internalShowHint = ref(true)
let hintTimeoutId = null

watch(
    () => props.hint,
    (newHint) => {
        if (newHint) {
            internalShowHint.value = true

            // 清除之前的超时
            if (hintTimeoutId) {
                clearTimeout(hintTimeoutId)
            }

            // 成功提示 5 秒后自动消退
            if (newHint.type === 'success') {
                hintTimeoutId = setTimeout(() => {
                    internalShowHint.value = false
                }, 5000)
            }
        }
    },
    { deep: true }
)

const closeHint = () => {
    internalShowHint.value = false
    if (hintTimeoutId) {
        clearTimeout(hintTimeoutId)
    }
}

const shouldShowHint = computed(() => {
    return props.showHint && internalShowHint.value
})
</script>

<style scoped>
.module-header-wrapper {
    flex-shrink: 0;
}

.module-header {
    padding: 24px;
    background-color: #ffffff;
    border-bottom: 1px solid #e1e4e8;
}

.header-content {
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.header-title {
    margin: 0;
    font-size: 24px;
    font-weight: 600;
    color: #24292e;
}

.header-subtitle {
    margin: 0;
    font-size: 14px;
    font-weight: 400;
    color: #586069;
}

/* 提示框样式 */
.module-hint {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 16px;
    border-radius: 4px;
    margin: 0;
    gap: 12px;
    border: 1px solid;
}

.hint-content {
    display: flex;
    align-items: center;
    gap: 12px;
    flex: 1;
}

.hint-icon {
    width: 20px;
    height: 20px;
    display: block;
    flex-shrink: 0;
}

.hint-message {
    font-size: 14px;
    line-height: 1.5;
}

.hint-close {
    background: transparent;
    border: none;
    cursor: pointer;
    font-size: 20px;
    padding: 0 4px;
    color: inherit;
    flex-shrink: 0;
    transition: opacity 0.2s;
}

.hint-close:hover {
    opacity: 0.7;
}

/* Success 提示 */
.hint-success {
    background-color: #f0f9ff;
    color: #0369a1;
    border-color: #bfdbfe;
}

.hint-success .hint-close {
    color: #0369a1;
}

/* Error 提示 */
.hint-error {
    background-color: #fef2f2;
    color: #991b1b;
    border-color: #fecaca;
}

.hint-error .hint-close {
    color: #991b1b;
}

/* Warning 提示 */
.hint-warning {
    background-color: #fffbeb;
    color: #92400e;
    border-color: #fde68a;
}

.hint-warning .hint-close {
    color: #92400e;
}

/* Info 提示 */
.hint-info {
    background-color: #f0f9ff;
    color: #0369a1;
    border-color: #bfdbfe;
}

.hint-info .hint-close {
    color: #0369a1;
}

/* 转场动画 */
.hint-fade-enter-active,
.hint-fade-leave-active {
    transition: opacity 0.3s ease;
}

.hint-fade-enter-from,
.hint-fade-leave-to {
    opacity: 0;
}
</style>
