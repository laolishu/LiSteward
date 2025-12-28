<template>
    <div class="settings-container">
        <ModuleHeader :title="$t('node.title')" />

        <div class="settings-content">
            <div class="settings-section">
                <div class="section-title">
                    <h2>{{ $t('node.title') }}</h2>
                </div>

                <div class="settings-item">
                    <div class="item-label">{{ $t('node.nvmVersion') }}</div>
                    <div class="item-control">
                        <div v-if="nvmVersion" class="nvm-version-display">{{ nvmVersion }}</div>
                        <button v-else class="btn-open" @click="installNvm">{{ $t('node.installNvm') }}</button>
                    </div>
                </div>

                <div class="settings-item">
                    <div class="item-label">{{ $t('node.nodeRepository') }}</div>
                    <div class="item-control">
                        <div class="nvm-version-display" style="flex: 0 0 70%">{{ nvmRoot }}</div>
                        <button class="btn-install" @click="installNewVersion" :title="$t('node.installButton')"
                            aria-label="install-node-version">
                            {{ $t('node.installButton') }}
                        </button>
                    </div>
                </div>

                <div class="section-title" style="margin-top:8px">
                    <h2>{{ $t('node.versions') }}</h2>
                    <button class="btn-refresh" @click="refresh" :title="$t('hosts.reload')" aria-label="refresh">
                        <img src="/images/refresh.svg" alt="refresh" width="18" height="18" />
                    </button>
                </div>

                <div v-if="loading" class="message-success">{{ $t('hosts.loading') }}</div>

                <div class="code-editor">
                    <table class="entries-table">
                        <thead>
                            <tr>
                                <th class="col-index">序号</th>
                                <th class="col-version">版本</th>
                                <th class="col-actions">操作</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(v, idx) in versions" :key="v.version">
                                <td class="col-index">{{ idx + 1 }}</td>
                                <td class="col-version">
                                    {{ v.version }}
                                    <span v-if="v.isCurrent" class="current-badge">{{ $t('node.current') }}</span>
                                </td>
                                <td class="col-actions">
                                    <button v-if="!v.isCurrent" class="action-icon" @click="useVersion(v.version)"
                                        :title="$t('node.use')" aria-label="use">
                                        <img src="/images/use.svg" alt="use" width="16" height="16" />
                                    </button>
                                    <button class="action-icon" @click="uninstallVersion(v.version)"
                                        :title="$t('node.uninstall')" aria-label="uninstall">
                                        <img src="/images/action-delete.svg" alt="uninstall" width="16" height="16" />
                                    </button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>

            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import ModuleHeader from '../components/ModuleHeader.vue'
import { DetectNvm, InstallNvm, ListNodeVersions, GetNvmDir, SetNvmDir, InstallNodeVersion, UseNodeVersion, UninstallNodeVersion, GetNvmVersion, GetNvmRoot, OpenNodeAvailableList } from '../../wailsjs/go/main/App'

const { t } = useI18n()

const nvmDir = ref('')
const versions = ref([])
const loading = ref(false)
const nvmVersion = ref('')
const nvmRoot = ref('')

async function loadNvmVersion() {
    try {
        const version = await GetNvmVersion()
        nvmVersion.value = version
    } catch (e) {
        nvmVersion.value = ''
        console.error('Failed to get NVM version:', e)
    }
}

async function loadNvmRoot() {
    try {
        const root = await GetNvmRoot()
        nvmRoot.value = root
    } catch (e) {
        nvmRoot.value = ''
        console.error('Failed to get NVM root:', e)
    }
}

async function detect() {
    try {
        const [ok, dir] = await DetectNvm()
        if (ok) {
            nvmDir.value = dir || (await GetNvmDir()) || nvmDir.value
            await loadNvmVersion()
            await refresh()
        } else {
            // nothing
        }
    } catch (e) {
        console.error(e)
    }
}

async function installNvm() {
    try {
        const msg = await InstallNvm()
        alert(msg)
    } catch (e) {
        console.error(e)
    }
}

async function refresh() {
    loading.value = true
    try {
        const list = await ListNodeVersions()
        versions.value = list || []
    } finally {
        loading.value = false
    }
}

async function useVersion(v) {
    await UseNodeVersion(v)
    alert('Switched to ' + v)
    await refresh()
}

async function uninstallVersion(v) {
    await UninstallNodeVersion(v)
    await refresh()
}

async function installNewVersion() {
    try {
        await OpenNodeAvailableList()
    } catch (e) {
        console.error('Failed to open node available list:', e)
    }
}

// on mount
(async () => {
    try {
        nvmDir.value = await GetNvmDir() || ''
        await loadNvmVersion()
        await loadNvmRoot()
        await refresh()
    } catch (e) {
        // ignore
    }
})()
</script>

<style scoped>
/* Reuse Settings-like styles for consistent look */
.settings-container {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    background-color: #f6f8fa;
    overflow-y: auto
}

.settings-header {
    display: none;
}

.settings-content {
    flex: 1;
    padding: 24px;
    overflow-y: auto
}

.settings-section {
    background-color: #ffffff;
    border-radius: 8px;
    padding: 24px;
    margin-bottom: 20px;
    border: 1px solid #e1e4e8
}

.section-title h2 {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
    color: #24292e
}

.section-title {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 16px;
}

.btn-refresh {
    background: transparent;
    border: none;
    cursor: pointer;
    padding: 4px 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
    transition: background 0.2s;
    flex-shrink: 0;
}

.btn-refresh:hover {
    background: rgba(16, 185, 129, 0.1);
}

.btn-refresh img {
    display: block;
}

.settings-item {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 16px
}

.item-label {
    flex-shrink: 0;
    width: 120px;
    font-size: 14px;
    color: #586069;
    font-weight: 500
}

.item-control {
    flex: 1;
    display: flex;
    gap: 12px;
    align-items: center
}

.nvm-version-display {
    font-size: 14px;
    color: #24292e;
    font-weight: 500;
    padding: 8px 12px;
    background-color: #f6f8fa;
    border-radius: 4px;
    border: 1px solid #e1e4e8;
}

.control-right {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    gap: 12px
}

.btn-open {
    padding: 8px 12px;
    background: #1B7C34;
    color: #fff;
    border: none;
    border-radius: 6px;
    cursor: pointer
}

.btn-open:hover {
    background: #2DA544
}

.btn-install {
    flex: 0 0 auto;
    padding: 8px 16px;
    background-color: #1f6feb;
    color: white;
    border: none;
    border-radius: 6px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: background-color 0.2s;
}

.btn-install:hover {
    background-color: #388bfd;
}

.btn-install:active {
    background-color: #1f6feb;
}

.language-select {
    width: 100%;
    max-width: 420px;
    padding: 8px 12px;
    border: 1px solid #e1e4e8;
    border-radius: 6px;
    font-size: 14px;
    color: #24292e;
    background-color: #fff
}

.language-select:hover {
    border-color: #34d399
}

.language-select:focus {
    outline: none;
    border-color: #34d399;
    box-shadow: 0 0 0 3px rgba(52, 211, 153, 0.1)
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
    gap: 8px
}

.text-muted {
    color: #666;
    font-size: 13px
}

.current-badge {
    display: inline-block;
    margin-left: 8px;
    padding: 2px 6px;
    background-color: #10B981;
    color: white;
    border-radius: 3px;
    font-size: 11px;
    font-weight: 600;
}

/* Table styles similar to HostsManager entries table */
.entries-table {
    width: 100%;
    border-collapse: collapse;
}

.entries-table thead th {
    text-align: left;
    padding: 10px;
    font-size: 13px;
    color: #586069;
    border-bottom: 1px solid #e1e4e8;
    background: #fff;
}

.entries-table tbody td {
    padding: 10px;
    border-bottom: 1px solid #f1f3f5;
    vertical-align: middle;
}

.entries-table .col-actions {
    width: 120px;
}

.entries-table .col-index {
    width: 60px
}

.entries-table .col-version {
    width: 200px
}

.action-icon {
    background: transparent;
    border: none;
    cursor: pointer;
    padding: 6px 8px;
    font-size: 16px;
    border-radius: 6px;
}

.action-icon:hover {
    background: rgba(52, 211, 153, 0.12);
}

.action-icon img {
    display: block
}
</style>
