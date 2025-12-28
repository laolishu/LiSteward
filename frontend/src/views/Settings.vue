<template>
    <div class="settings-container">
        <ModuleHeader :title="$t('settings.title')" :hint="languageChangedHint" />

        <div class="settings-content">
            <!-- 通用设置 -->
            <div class="settings-section">
                <div class="section-title">
                    <h2>{{ $t('settings.general') }}</h2>
                </div>

                <div class="settings-item">
                    <div class="item-label">{{ $t('settings.currentVersion') }}</div>
                    <div class="item-control control-right">
                        <div class="version-display">{{ appInfo?.version || '-' }}</div>
                    </div>
                </div>

                <div class="settings-item">
                    <div class="item-label">{{ $t('settings.contactAuthor') }}</div>
                    <div class="item-control control-right">
                        <a :href="appInfo?.email ? `mailto:${appInfo.email}` : '#'" class="author-email">{{
                            appInfo?.author || '-' }} &lt;{{ appInfo?.email || '-' }}&gt;</a>
                    </div>
                </div>

                <div class="settings-item">
                    <div class="item-label">{{ $t('settings.languageSetting') }}</div>
                    <div class="item-control control-right">
                        <select v-model="selectedLanguage" @change="changeLanguage" class="language-select">
                            <option value="zh-CN">中文 (Chinese)</option>
                            <option value="en-US">English</option>
                        </select>
                    </div>
                </div>

                <div class="settings-item">
                    <div class="item-label">{{ $t('settings.autoStart') }}</div>
                    <div class="item-control control-right">
                        <input type="checkbox" v-model="autoStart" @change="toggleAutoStart" />
                    </div>
                </div>





            </div>

            <!-- 关于部分 已移除 -->

            <!-- 赞赏二维码部分 -->
            <div class="settings-section sponsor-section" v-if="appConfig?.sponsorQrCode">
                <div class="section-title">
                    <h2>{{ $t('settings.sponsor') }}</h2>
                </div>
                <p class="sponsor-tip">{{ $t('settings.sponsorTip') }}</p>

                <!-- 图片加载成功 -->
                <div class="qr-code-container" v-if="qrCodeLoaded">
                    <img :src="appConfig.sponsorQrCode" alt="WeChat Sponsor QR" class="qr-code" />
                </div>

                <!-- 图片加载中 (初始状态) -->
                <div class="qr-code-container loading" v-else-if="qrCodeErrorRetried === 0">
                    <div class="qr-code qr-loading">加载中...</div>
                </div>

                <!-- 图片加载失败 -->
                <div class="sponsor-error-content" v-else-if="qrCodeErrorRetried > 0">
                    <p class="sponsor-error-tip">
                        赞赏二维码加载失败，请检查网络连接或访问官方网站。
                    </p>
                    <button @click="retryLoadQrCode" class="btn-retry">重试</button>
                </div>

                <!-- 隐藏的img标签用于加载 -->
                <img :src="appConfig.sponsorQrCode" alt="" style="display:none" @load="onQrCodeLoad"
                    @error="onQrCodeLoadError" />
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import ModuleHeader from '../components/ModuleHeader.vue'
import { GetUserLanguagePreference, SetUserLanguagePreference, GetAppConfig, GetAutoStart, SetAutoStart, GetAppInfo } from '../../wailsjs/go/main/App'

const { locale, t } = useI18n()
const selectedLanguage = ref('zh-CN')
const languageChangedHint = ref(null)
const appConfig = ref(null)
const appInfo = ref(null)
const autoStart = ref(false)
const qrCodeLoaded = ref(false)
const qrCodeErrorRetried = ref(0)
// 简单消息提示（当前页面使用 window.alert 来显示简短信息）
const showMessage = (text, type) => {
    // type 可用于后续拓展，目前使用 alert
    try {
        window.alert(text)
    } catch (e) {
        console.log(text)
    }
}

onMounted(async () => {
    // 获取当前用户语言偏好
    try {
        const lang = await GetUserLanguagePreference()
        selectedLanguage.value = lang || 'zh-CN'
        locale.value = lang || 'zh-CN'
    } catch (e) {
        console.error('Failed to get language preference:', e)
    }

    // 获取应用配置
    try {
        appConfig.value = await GetAppConfig()
    } catch (e) {
        console.error('Failed to get app config:', e)
    }

    // 获取应用版本和信息
    try {
        appInfo.value = await GetAppInfo()
    } catch (e) {
        console.error('Failed to get app info:', e)
    }

    // 获取用户开机自启动设置
    try {
        const as = await GetAutoStart()
        autoStart.value = !!as
    } catch (e) {
        console.error('Failed to get auto start preference:', e)
    }
})

const changeLanguage = async () => {
    try {
        // 切换 i18n 语言
        locale.value = selectedLanguage.value

        // 保存到后端
        await SetUserLanguagePreference(selectedLanguage.value)

        // 保存到 localStorage
        localStorage.setItem('app-language', selectedLanguage.value)

        // 显示成功提示
        languageChangedHint.value = {
            type: 'success',
            message: t('settings.languageChange')
        }
        setTimeout(() => {
            languageChangedHint.value = null
        }, 2000)
    } catch (e) {
        console.error('Failed to change language:', e)
        languageChangedHint.value = null
    }
}

// 打开程序目录
const openAppDir = async () => {
    try {
        // 直接使用 runtime bridge，生成的 wailsjs 可能未包含新方法，使用 window 对象调用
        if (window && window['go'] && window['go']['main'] && window['go']['main']['App'] && window['go']['main']['App']['OpenAppDir']) {
            await window['go']['main']['App']['OpenAppDir']()
            showMessage('已尝试打开程序目录', 'success')
        } else {
            showMessage('后端方法不可用，请重建前端绑定', 'error')
        }
    } catch (e) {
        console.error('打开目录失败:', e)
        showMessage('打开目录失败: ' + e, 'error')
    }
}

const toggleAutoStart = async () => {
    try {
        await SetAutoStart(autoStart.value)
        showMessage('设置已保存', 'success')
    } catch (e) {
        console.error('Failed to set auto start:', e)
        showMessage('保存失败: ' + e, 'error')
    }
}

// 处理二维码加载成功
const onQrCodeLoad = () => {
    console.log('Sponsor QR code loaded successfully')
    qrCodeLoaded.value = true
}

// 处理二维码加载失败
const onQrCodeLoadError = () => {
    console.warn('Failed to load sponsor QR code:', appConfig.value?.sponsorQrCode)
    qrCodeLoaded.value = false
    qrCodeErrorRetried.value++
}

// 重试加载二维码
const retryLoadQrCode = () => {
    qrCodeLoaded.value = false
    qrCodeErrorRetried.value = 0
    // 触发隐藏的img标签重新加载
    if (appConfig.value?.sponsorQrCode) {
        const img = new Image()
        img.onload = onQrCodeLoad
        img.onerror = onQrCodeLoadError
        img.src = appConfig.value.sponsorQrCode + '?t=' + Date.now() // 添加时间戳避免缓存
    }
}


</script>

<style scoped>
.settings-container {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    background-color: #f6f8fa;
    overflow-y: auto;
}

.settings-header {
    display: none;
}

.settings-content {
    flex: 1;
    padding: 24px;
    overflow-y: auto;
}

.settings-section {
    background-color: #ffffff;
    border-radius: 8px;
    padding: 24px;
    margin-bottom: 20px;
    border: 1px solid #e1e4e8;
}

.section-title h2 {
    margin: 0 0 16px 0;
    font-size: 16px;
    font-weight: 600;
    color: #24292e;
}

.settings-item {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 16px;
}

.item-label {
    flex-shrink: 0;
    width: 120px;
    font-size: 14px;
    color: #586069;
    font-weight: 500;
}

.item-control {
    flex: 1;
}

.control-right {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    gap: 12px;
}

.btn-open {
    padding: 8px 12px;
    background: #1B7C34;
    color: #fff;
    border: none;
    border-radius: 6px;
    cursor: pointer;
}

.btn-open:hover {
    background: #2DA544;
}

.author-email {
    color: #0366d6;
    text-decoration: none;
}

.author-email:hover {
    text-decoration: underline;
}

.language-select {
    width: 200px;
    padding: 8px 12px;
    border: 1px solid #e1e4e8;
    border-radius: 6px;
    font-size: 14px;
    color: #24292e;
    background-color: #ffffff;
    cursor: pointer;
    transition: border-color 0.2s;
}

.language-select:hover {
    border-color: #34d399;
}

.language-select:focus {
    outline: none;
    border-color: #34d399;
    box-shadow: 0 0 0 3px rgba(52, 211, 153, 0.1);
}

.message-success {
    padding: 12px 16px;
    margin-top: 12px;
    background-color: #f0fdf4;
    border: 1px solid #dcfce7;
    border-radius: 6px;
    color: #166534;
    font-size: 14px;
    display: flex;
    align-items: center;
    gap: 8px;
}

.message-error {
    padding: 12px 16px;
    margin-top: 12px;
    background-color: #fef2f2;
    border: 1px solid #fee2e2;
    border-radius: 6px;
    color: #991b1b;
    font-size: 14px;
    display: flex;
    align-items: center;
    gap: 8px;
}



.about-item {
    margin-bottom: 16px;
    padding: 12px 0;
}

.app-info {
    margin-bottom: 12px;
}

.app-info h3 {
    margin: 0 0 8px 0;
    font-size: 18px;
    font-weight: 600;
    color: #24292e;
}

.version {
    margin: 0;
    font-size: 14px;
    color: #666;
}

.email {
    font-size: 12px;
    color: #999;
    margin: 4px 0 0 0;
}

.external-link {
    color: #0366d6;
    text-decoration: none;
    word-break: break-all;
}

.external-link:hover {
    text-decoration: underline;
}

.sponsor-section {
    background-color: #f9f9f9;
    text-align: center;
}

.sponsor-tip {
    margin: 0 0 12px 0;
    font-size: 14px;
    color: #666;
}

.qr-code-container {
    display: flex;
    justify-content: center;
}

.qr-code {
    width: 200px;
    height: 200px;
    border: 1px solid #ddd;
    border-radius: 4px;
}

.sponsor-error {
    background-color: #fef2f2;
    border-color: #fee2e2;
}

.sponsor-error-tip {
    margin: 0 0 16px 0;
    font-size: 14px;
    color: #991b1b;
    padding: 12px;
    background-color: #fff5f5;
    border: 1px solid #fecaca;
    border-radius: 4px;
}

.btn-retry {
    padding: 8px 16px;
    background-color: #0366d6;
    color: #ffffff;
    border: none;
    border-radius: 6px;
    font-size: 14px;
    cursor: pointer;
    transition: background-color 0.2s;
}

.btn-retry:hover {
    background-color: #0256c7;
}

.btn-retry:active {
    background-color: #0256c7;
    transform: scale(0.98);
}

.qr-loading {
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: #f5f5f5;
    color: #666;
    font-size: 14px;
}

.sponsor-error-content {
    text-align: center;
}


.message-fade-enter-from,
.message-fade-leave-to {
    opacity: 0;
}
</style>
