<template>
    <div class="hosts-manager">
        <ModuleHeader :title="$t('hosts.title')" :hint="messageHint" />

        <div class="module-layout">
            <!-- 主内容：左侧 -->
            <div class="right-panel">
                <div class="editor-toolbar">
                    <button @click="toggleSidebar" class="btn-icon" :title="t('profile.scheme')">
                        <img src="/images/toolbar-menu.svg" alt="menu" />
                    </button>
                    <div class="search-box">
                        <img src="/images/toolbar-search.svg" alt="search" class="search-icon" />
                        <input v-model="searchQuery" type="text" :placeholder="t('hosts.search')"
                            class="search-input" />
                    </div>
                    <select v-model="filterStatus" class="filter-select">
                        <option value="all">{{ t('hosts.all') }}</option>
                        <option value="enabled">{{ t('hosts.enabled') }}</option>
                        <option value="disabled">{{ t('hosts.disabled') }}</option>
                    </select>
                    <!-- 刷新按钮 -->
                    <button @click="reloadHosts" class="btn-icon" :title="t('hosts.reload')">
                        <img src="/images/refresh.svg" alt="refresh" />
                    </button>
                    <!-- 应用按钮 -->
                    <button @click="applyChanges" :disabled="!hasUnsavedChanges || isSaving" class="btn-icon"
                        :title="t('hosts.apply')">
                        <img src="/images/toolbar-apply.svg" alt="apply" />
                    </button>
                </div>

                <div class="editor-container">
                    <div class="add-entry-bar">
                        <input v-model="newEntry.domain" type="text" :placeholder="t('hosts.domain')"
                            class="entry-input domain-input" />
                        <input v-model="newEntry.ip" type="text" :placeholder="t('hosts.ip')"
                            class="entry-input ip-input" />
                        <input v-model="newEntry.alternate_ip" type="text" :placeholder="t('hosts.alternateIp')"
                            class="entry-input alternate-ip-input" />
                        <input v-model="newEntry.description" type="text" :placeholder="t('hosts.description')"
                            class="entry-input description-input" />
                        <button @click="addEntry" class="btn-add-entry">{{ t('hosts.addEntry') }}</button>
                    </div>

                    <div class="code-editor">
                        <div v-if="loading" class="loading-state">
                            <div class="spinner"></div>
                            <p>{{ t('hosts.loading') }}</p>
                        </div>
                        <div v-else-if="filteredEntries.length === 0" class="empty-state">
                            <p>{{ t('hosts.empty') }}</p>
                            <small>{{ t('hosts.clickAdd') }}</small>
                        </div>
                        <table v-else class="entries-table">
                            <thead>
                                <tr>
                                    <th class="col-index">{{ t('hosts.index') }}</th>
                                    <th class="col-status">{{ t('hosts.status') }}</th>
                                    <th class="col-domain">{{ t('hosts.domain') }}</th>
                                    <th class="col-ip">IP</th>
                                    <th class="col-alternate-ip">{{ t('hosts.alternateIp') }}</th>
                                    <th class="col-comment">{{ t('hosts.remark') }}</th>
                                    <th class="col-actions">{{ t('hosts.operations') }}</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="(entry, index) in filteredEntries" :key="index"
                                    :class="{ disabled: !entry.enabled }">
                                    <td class="col-index">{{ index + 1 }}</td>
                                    <td class="col-status">
                                        <button @click="toggleEntry(getOriginalIndex(entry))"
                                            :class="['toggle-indicator', { enabled: entry.enabled }]"
                                            :title="entry.enabled ? t('hosts.clickDisable') : t('hosts.clickEnable')">
                                            <img v-if="entry.enabled" src="/images/message-success.svg" alt="enabled"
                                                class="status-img" />
                                            <img v-else src="/images/message-error.svg" alt="disabled"
                                                class="status-img" />
                                        </button>
                                    </td>
                                    <td class="col-domain">{{ entry.domain }}</td>
                                    <td class="col-ip">{{ entry.ip }}</td>
                                    <td class="col-alternate-ip">{{ entry.alternate_ip || '-' }}</td>
                                    <td class="col-comment">
                                        <span v-if="entry.description">{{ entry.description }}</span>
                                        <span v-else-if="entry.comment"># {{ entry.comment }}</span>
                                    </td>
                                    <td class="col-actions">
                                        <button v-if="entry.alternate_ip" @click="swapIPs(getOriginalIndex(entry))"
                                            class="action-btn swap-btn"
                                            :title="`${t('hosts.swap')} (${entry.ip} ⇄ ${entry.alternate_ip})`">
                                            <img src="/images/action-swap.svg" alt="swap" />
                                        </button>
                                        <button @click="editEntry(getOriginalIndex(entry))" class="action-btn"
                                            :title="t('hosts.edit')">
                                            <img src="/images/action-edit.svg" alt="edit" />
                                        </button>
                                        <button @click="deleteEntry(getOriginalIndex(entry))" class="action-btn"
                                            :title="t('hosts.delete')">
                                            <img src="/images/action-delete.svg" alt="delete" />
                                        </button>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>

                    <!-- 底部状态栏 -->
                    <div class="status-bar">
                        <div class="status-left">
                            <span v-if="hasUnsavedChanges" class="status-item unsaved-item">
                                ● {{ t('hosts.unsaved') }} ({{ entries.length }} {{ t('hosts.items') }})
                            </span>
                            <span v-else class="status-item saved-item">
                                <img src="/images/message-success.svg" alt="saved"
                                    style="width: 14px; height: 14px; display: inline-block; margin-right: 4px;" />
                                {{ t('hosts.saved') }}
                            </span>
                        </div>
                        <div class="status-right">
                            <span class="status-item">{{ t('hosts.total') }} {{ entries.length }} {{ t('hosts.entries')
                                }} | {{ t('hosts.enabled') }} {{entries.filter(e =>
                                    e.enabled).length}} {{ t('hosts.entries') }}</span>
                        </div>
                    </div>
                </div>
            </div>

            <!-- 侧栏（右侧抽屉，推挤式） -->
            <div :class="['profile-sidebar', { collapsed: !sidebarOpen }]" :style="{ width: profilePanelWidth + 'px' }">
                <div class="panel-header">
                    <h3>{{ t('profile.scheme') }}</h3>
                    <button @click="showSaveProfileDialog = true" class="btn-add">+</button>
                </div>

                <div class="profile-list">
                    <div v-for="profile in profiles" :key="profile.name"
                        :class="['profile-card', { active: selectedProfile === profile.name }]"
                        @click="selectProfile(profile.name)">
                        <div class="profile-info">
                            <div class="profile-name">{{ profile.name }}</div>
                            <div class="profile-meta">
                                <span>{{ profile.entries?.length || 0 }} {{ t('hosts.items') }}</span>
                                <span>{{ formatDate(profile.updated_at) }}</span>
                            </div>
                        </div>
                        <div class="profile-actions">
                            <button @click.stop="openRenameProfileDialog(profile.name)" class="btn-edit"
                                :title="t('profile.rename')">
                                <img src="/images/action-edit.svg" alt="edit" />
                            </button>
                            <button @click.stop="deleteProfileConfirm(profile.name)" class="btn-delete"
                                :title="t('hosts.delete')">
                                <img src="/images/action-delete.svg" alt="delete" />
                            </button>
                        </div>
                    </div>
                </div>

                <div class="panel-section">
                    <div class="section-header">
                        <h3>{{ t('profile.backups') }}</h3>
                    </div>
                    <div class="backup-list">
                        <div v-for="backup in (backups || []).slice(0, 5)" :key="backup.id" class="backup-item">
                            <div class="backup-info">
                                <div class="backup-time">{{ formatDate(backup.timestamp) }}</div>
                                <div class="backup-size">{{ formatSize(backup.size) }}</div>
                            </div>
                            <button @click="restoreBackup(backup.id)" class="btn-restore">{{ t('profile.restore')
                            }}</button>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div v-if="showSaveProfileDialog" class="modal-overlay" @click.self="showSaveProfileDialog = false">
            <div class="modal-dialog">
                <div class="modal-header">
                    <h3>{{ t('profile.saveAs') }}</h3>
                    <button @click="showSaveProfileDialog = false" class="btn-close">×</button>
                </div>
                <div class="modal-body">
                    <input v-model="profileName" type="text" :placeholder="t('profile.inputName')" class="modal-input"
                        @keyup.enter="saveProfile" />
                </div>
                <div class="modal-footer">
                    <button @click="showSaveProfileDialog = false" class="btn btn-secondary">{{ t('dialog.cancel')
                        }}</button>
                    <button @click="saveProfile" class="btn btn-primary">{{ t('dialog.save') }}</button>
                </div>
            </div>
        </div>

        <!-- 编辑条目对话框 -->
        <div v-if="editingEntry !== null" class="modal-overlay" @click.self="cancelEdit">
            <div class="modal-dialog">
                <div class="modal-header">
                    <h3>{{ t('dialog.editEntry') }}</h3>
                    <button @click="cancelEdit" class="btn-close">×</button>
                </div>
                <div class="modal-body">
                    <div class="form-group"><label>{{ t('hosts.domain') }}</label><input v-model="editingEntry.domain"
                            type="text" class="modal-input" /></div>
                    <div class="form-group"><label>{{ t('hosts.ip') }}</label><input v-model="editingEntry.ip"
                            type="text" class="modal-input" /></div>
                    <div class="form-group"><label>{{ t('hosts.alternateIp') }}</label><input
                            v-model="editingEntry.alternate_ip" type="text" class="modal-input"
                            :placeholder="t('profile.example')" /></div>
                    <div class="form-group"><label>{{ t('hosts.description') }}</label><input
                            v-model="editingEntry.description" type="text" class="modal-input"
                            :placeholder="t('profile.exampleEnv')" /></div>
                </div>
                <div class="modal-footer">
                    <button @click="cancelEdit" class="btn btn-secondary">{{ t('dialog.cancel') }}</button>
                    <button @click="saveEdit" class="btn btn-primary">{{ t('dialog.save') }}</button>
                </div>
            </div>
        </div>

        <!-- 重命名方案对话框 -->
        <div v-if="showRenameProfileDialog" class="modal-overlay" @click.self="cancelRenameProfile">
            <div class="modal-dialog">
                <div class="modal-header">
                    <h3>{{ t('profile.rename') }}</h3>
                    <button @click="cancelRenameProfile" class="btn-close">×</button>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <input v-model="newProfileName" type="text" class="modal-input" @keyup.enter="renameProfile"
                            autofocus />
                    </div>
                </div>
                <div class="modal-footer">
                    <button @click="cancelRenameProfile" class="btn btn-secondary">{{ t('dialog.cancel') }}</button>
                    <button @click="renameProfile" class="btn btn-primary">{{ t('dialog.confirm') }}</button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import ModuleHeader from '../components/ModuleHeader.vue'
import {
    GetHostsEntries, SaveHostsEntries, ValidateHostEntry,
    CreateBackup, ListBackups, RestoreBackup, DeleteBackup,
    SaveProfile, LoadProfile, ListProfiles, DeleteProfile, ApplyProfile,
    SaveHostsEntriesWithReadOnlyRestore, RenameProfile, SwapIPs
} from '../../wailsjs/go/main/App'

const { t } = useI18n()

// 状态
const entries = ref([])
const initialEntries = ref([]) // 保存初始状态，用于刷新时对比
const profiles = ref([])
const backups = ref([])
const loading = ref(false)
const message = ref({ text: '', type: '' })
const messageHint = ref(null)
const searchQuery = ref('')
const filterStatus = ref('all')
const selectedProfile = ref('')
const showSaveProfileDialog = ref(false)
const profileName = ref('')
const sidebarOpen = ref(false)
const profilePanelWidth = ref(280)
const isDragging = ref(false)
const hasUnsavedChanges = ref(false) // 缓存标志：是否有未保存的更改
const isSaving = ref(false) // 保存中状态

// 重命名对话框状态
const showRenameProfileDialog = ref(false)
const renamingProfile = ref('') // 当前正在重命名的方案名称
const newProfileName = ref('') // 输入框中的新名称

const toggleSidebar = () => {
    sidebarOpen.value = !sidebarOpen.value
}

// 新条目
const newEntry = ref({
    ip: '',
    domain: '',
    comment: '',
    alternate_ip: '',
    description: '',
    enabled: true
})

// 编辑中的条目
const editingEntry = ref(null)
const editingIndex = ref(-1)

// 过滤后的条目
const filteredEntries = computed(() => {
    let result = entries.value

    if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        result = result.filter(e =>
            e.ip.toLowerCase().includes(query) ||
            e.domain.toLowerCase().includes(query)
        )
    }

    if (filterStatus.value === 'enabled') {
        result = result.filter(e => e.enabled)
    } else if (filterStatus.value === 'disabled') {
        result = result.filter(e => !e.enabled)
    }

    return result
})

// 选择方案
const selectProfile = async (name) => {
    // 检查是否有未保存的更改
    if (hasUnsavedChanges.value) {
        if (!confirm(t('dialog.switchProfileUnsaved'))) {
            return
        }
    }

    // 再次确认切换方案
    if (!confirm(t('dialog.switchProfileConfirm').replace('{name}', name))) {
        return
    }

    try {
        await ApplyProfile(name)
        await loadHosts()
        selectedProfile.value = name
        showMessage(`方案 "${name}" 已切换`, 'success')
    } catch (error) {
        showMessage('切换方案失败: ' + error, 'error')
    }
}

// 删除方案确认
const deleteProfileConfirm = async (name) => {
    if (!confirm(t('dialog.deleteProfileConfirm').replace('{name}', name))) {
        return
    }

    try {
        await DeleteProfile(name)
        showMessage(`方案 "${name}" 已删除`, 'success')
        loadProfiles()
    } catch (error) {
        showMessage('删除方案失败: ' + error, 'error')
    }
}

// 打开重命名对话框
const openRenameProfileDialog = (profileName) => {
    renamingProfile.value = profileName
    newProfileName.value = profileName // 预填当前名称
    showRenameProfileDialog.value = true
}

// 重命名方案
const renameProfile = async () => {
    const trimmedName = newProfileName.value.trim()

    // 验证新名称非空
    if (!trimmedName) {
        showMessage('方案名称不能为空', 'error')
        return
    }

    // 检查新名称与原名称是否相同
    if (trimmedName === renamingProfile.value) {
        showMessage('新名称与原名称相同，无需修改', 'success')
        cancelRenameProfile()
        return
    }

    try {
        await RenameProfile(renamingProfile.value, trimmedName)
        showMessage(`方案已重命名为 '${trimmedName}'`, 'success')
        cancelRenameProfile()
        loadProfiles()
    } catch (error) {
        showMessage('重命名方案失败: ' + error, 'error')
    }
}

// 取消重命名
const cancelRenameProfile = () => {
    showRenameProfileDialog.value = false
    renamingProfile.value = ''
    newProfileName.value = ''
}

// 验证 IP 地址格式
const isValidIP = (ip) => {
    if (!ip || ip.trim() === '') return true // 空字符串表示可选字段

    // IPv4 和 IPv6 正则表达式
    const ipv4Regex = /^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/
    const ipv6Regex = /^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4})$/

    return ipv4Regex.test(ip) || ipv6Regex.test(ip)
}

// 获取原始索引
const getOriginalIndex = (entry) => {
    return entries.value.findIndex(e =>
        e.ip === entry.ip && e.domain === entry.domain && e.comment === entry.comment && e.alternate_ip === entry.alternate_ip
    )
}

// 加载 Hosts 文件
const loadHosts = async () => {
    loading.value = true
    try {
        entries.value = await GetHostsEntries()
        initialEntries.value = JSON.parse(JSON.stringify(entries.value)) // 保存初始副本
        hasUnsavedChanges.value = false
        showMessage('加载成功', 'success')
    } catch (error) {
        showMessage('加载失败: ' + error, 'error')
    } finally {
        loading.value = false
    }
}

// 应用缓存中的更改到 Hosts 文件
const applyChanges = async () => {
    if (!hasUnsavedChanges.value) {
        showMessage('没有未保存的更改', 'error')
        return
    }

    if (!confirm(t('dialog.applySaveConfirm'))) {
        return
    }

    isSaving.value = true
    try {
        await SaveHostsEntriesWithReadOnlyRestore(entries.value)
        initialEntries.value = JSON.parse(JSON.stringify(entries.value)) // 更新初始副本
        hasUnsavedChanges.value = false
        showMessage('更改已保存到 Hosts 文件', 'success')
        loadBackups() // 刷新备份列表
    } catch (error) {
        showMessage('保存失败: ' + error, 'error')
    } finally {
        isSaving.value = false
    }
}

// 刷新：从 Hosts 文件重新加载，丢弃未保存的更改
const reloadHosts = async () => {
    if (hasUnsavedChanges.value) {
        const unsavedCount = entries.value.length - initialEntries.value.length
        if (!confirm(t('dialog.reloadWithoutSave'))) {
            return
        }
    }

    await loadHosts()
    showMessage('已从 Hosts 文件重新加载', 'success')
}

// 保存 Hosts 文件（旧方法，保留以兼容）
const saveHosts = async () => {
    if (!confirm(t('dialog.saveHostsConfirm'))) {
        return
    }

    loading.value = true
    try {
        await SaveHostsEntries(entries.value)
        showMessage('保存成功！已自动创建备份', 'success')
        loadBackups()
    } catch (error) {
        showMessage('保存失败: ' + error, 'error')
    } finally {
        loading.value = false
    }
}

// 添加条目
const addEntry = async () => {
    if (!newEntry.value.ip.trim() || !newEntry.value.domain.trim()) {
        showMessage('请输入 IP 和域名', 'error')
        return
    }

    // 验证 IP 格式
    if (!isValidIP(newEntry.value.ip)) {
        showMessage('主 IP 格式无效', 'error')
        return
    }

    // 验证备用 IP 格式（如果提供）
    if (newEntry.value.alternate_ip && !isValidIP(newEntry.value.alternate_ip)) {
        showMessage('备用 IP 格式无效', 'error')
        return
    }

    try {
        await ValidateHostEntry(newEntry.value)
        entries.value.push({ ...newEntry.value })

        // 设置缓存标志
        hasUnsavedChanges.value = true

        newEntry.value = { ip: '', domain: '', comment: '', alternate_ip: '', description: '', enabled: true }
        showMessage('条目已添加到缓存，点击"应用"保存到文件', 'success')
    } catch (error) {
        showMessage('添加失败: ' + (error.message || error), 'error')
    }
}

// 切换启用/禁用
const toggleEntry = (index) => {
    entries.value[index].enabled = !entries.value[index].enabled

    // 设置缓存标志
    hasUnsavedChanges.value = true
}

// 编辑条目
const editEntry = (index) => {
    editingIndex.value = index
    editingEntry.value = { ...entries.value[index] }
}

// 保存编辑
const saveEdit = async () => {
    // 验证 IP 格式
    if (!isValidIP(editingEntry.value.ip)) {
        showMessage('主 IP 格式无效', 'error')
        return
    }

    // 验证备用 IP 格式（如果提供）
    if (editingEntry.value.alternate_ip && !isValidIP(editingEntry.value.alternate_ip)) {
        showMessage('备用 IP 格式无效', 'error')
        return
    }

    try {
        await ValidateHostEntry(editingEntry.value)
        entries.value[editingIndex.value] = { ...editingEntry.value }

        // 设置缓存标志
        hasUnsavedChanges.value = true

        cancelEdit()
        showMessage('修改已保存到缓存，点击"应用"保存到文件', 'success')
    } catch (error) {
        showMessage('保存失败: ' + (error.message || error), 'error')
    }
}

// 取消编辑
const cancelEdit = () => {
    editingEntry.value = null
    editingIndex.value = -1
}

// 删除条目
const deleteEntry = async (index) => {
    if (!confirm(t('dialog.deleteEntryConfirm'))) {
        return
    }

    try {
        entries.value.splice(index, 1)

        // 设置缓存标志
        hasUnsavedChanges.value = true

        showMessage('条目已删除（缓存中），点击"应用"保存到文件', 'success')
    } catch (error) {
        showMessage('删除失败: ' + (error.message || error), 'error')
    }
}

// 交换 IP
const swapIPs = async (index) => {
    try {
        const swappedEntries = await SwapIPs(entries.value, index)

        // 更新本地条目
        entries.value = swappedEntries

        // 设置缓存标志
        hasUnsavedChanges.value = true

        showMessage(`已交换 ${entries.value[index].domain} 的 IP：${entries.value[index].ip} ⇄ ${entries.value[index].alternate_ip}`, 'success')
    } catch (error) {
        showMessage('交换 IP 失败: ' + (error.message || error), 'error')
    }
}

// 保存方案
const saveProfile = async () => {
    if (!profileName.value.trim()) {
        showMessage('请输入方案名称', 'error')
        return
    }

    try {
        await SaveProfile(profileName.value, entries.value)
        showMessage(`方案 "${profileName.value}" 已保存`, 'success')
        showSaveProfileDialog.value = false
        profileName.value = ''
        loadProfiles()
    } catch (error) {
        showMessage('保存方案失败: ' + error, 'error')
    }
}

// 加载方案列表
const loadProfiles = async () => {
    try {
        profiles.value = await ListProfiles()
    } catch (error) {
        console.error('加载方案列表失败:', error)
    }
}

// 创建备份
const createBackup = async () => {
    try {
        await CreateBackup()
        showMessage('备份创建成功', 'success')
        loadBackups()
    } catch (error) {
        showMessage('创建备份失败: ' + error, 'error')
    }
}

// 加载备份列表
const loadBackups = async () => {
    try {
        backups.value = await ListBackups()
    } catch (error) {
        console.error('加载备份列表失败:', error)
    }
}

// 恢复备份
const restoreBackup = async (backupID) => {
    if (!confirm(t('dialog.restoreBackupConfirm'))) {
        return
    }

    try {
        await RestoreBackup(backupID)
        showMessage('备份已恢复', 'success')
        loadHosts()
    } catch (error) {
        showMessage('恢复备份失败: ' + error, 'error')
    }
}

// 显示消息
const showMessage = (text, type = 'success', duration = 2000) => {
    messageHint.value = { message: text, type }

    // 成功提示 2 秒后自动消退
    if (type === 'success') {
        setTimeout(() => {
            messageHint.value = null
        }, duration)
    }
    // 错误提示不自动消退，需用户手动关闭
}

// 格式化日期
const formatDate = (timestamp) => {
    if (!timestamp) return '-'
    const date = new Date(timestamp)
    const now = new Date()
    const diff = now - date

    if (diff < 60000) return '刚刚'
    if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
    if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`
    if (diff < 604800000) return `${Math.floor(diff / 86400000)}天前`

    return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
}

// 格式化文件大小
const formatSize = (bytes) => {
    if (bytes < 1024) return bytes + ' B'
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
    return (bytes / 1024 / 1024).toFixed(1) + ' MB'
}

// 拖动分隔符
const startDragging = () => {
    isDragging.value = true
    const startX = event.clientX
    const startWidth = profilePanelWidth.value

    const handleMouseMove = (e) => {
        const delta = e.clientX - startX
        const newWidth = Math.max(200, Math.min(500, startWidth + delta))
        profilePanelWidth.value = newWidth
    }

    const handleMouseUp = () => {
        isDragging.value = false
        document.removeEventListener('mousemove', handleMouseMove)
        document.removeEventListener('mouseup', handleMouseUp)
    }

    document.addEventListener('mousemove', handleMouseMove)
    document.addEventListener('mouseup', handleMouseUp)
}

// 初始化
onMounted(() => {
    loadHosts()
    loadProfiles()
    loadBackups()
})
</script>

<style scoped>
/* 页面容器 */
.hosts-manager {
    height: 100%;
    display: flex;
    flex-direction: column;
    background: #F6F8FA;
    color: #24292E;
}

/* 顶部工具栏 - 隐藏，由 MainLayout 提供 */
.toolbar {
    display: none;
}

/* 按钮样式 */
.btn {
    padding: 8px 16px;
    border: none;
    border-radius: 6px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    font-size: 14px;
}

.btn-primary {
    background: #1B7C34;
    color: white;
}

.btn-primary:hover:not(:disabled) {
    background: #2DA544;
    transform: translateY(-2px);
}

.btn-primary:disabled {
    background: #D0D7DE;
    color: #6A737D;
    cursor: not-allowed;
}

.btn-secondary {
    background: #F6F8FA;
    color: #24292E;
    border: 1px solid #D0D7DE;
}

.btn-secondary:hover {
    background: #EEEFF2;
    color: #24292E;
}

.btn-icon {
    padding: 6px 12px;
    background: transparent;
    border: 1px solid #D0D7DE;
    color: #586069;
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
}

.btn-icon img {
    width: 16px;
    height: 16px;
    object-fit: contain;
}

.btn-icon:hover {
    background: #F3F4F6;
    color: #24292E;
}

.btn-icon:hover img {
    filter: brightness(0.8);
}

/* 主布局 - 改为模块布局 */
.module-layout {
    display: flex;
    flex-direction: row;
    flex: 1;
    overflow: hidden;
    position: relative;
    width: 100%;
}

/* 侧栏切换按钮 - 隐藏，由 MainLayout 提供 */
.sidebar-toggle {
    display: none;
}

/* 右侧面板 - 主内容区 */
.right-panel {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
}

/* 方案侧栏 - 右侧抽屉 */
.profile-sidebar {
    background: #FFFFFF;
    border-right: 1px solid #E1E4E8;
    border-left: none;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    transition: width 0.25s ease, opacity 0.25s ease;
    position: relative;
    flex-shrink: 0;
    box-shadow: 4px 0 12px rgba(0, 0, 0, 0.04);
}

.profile-sidebar.collapsed {
    width: 0 !important;
    min-width: 0 !important;
    padding: 0 !important;
    margin: 0 !important;
    opacity: 0;
    overflow: hidden;
    border: none !important;
    box-shadow: none !important;
}

.profile-sidebar:not(.collapsed) {
    opacity: 1;
}

.panel-header {
    padding: 16px;
    border-bottom: 1px solid #E1E4E8;
    display: flex;
    justify-content: space-between;
    align-items: center;
    position: sticky;
    top: 0;
    background: #FFFFFF;
}

.panel-header h3 {
    margin: 0;
    font-size: 14px;
    font-weight: 600;
    color: #24292E;
}

.btn-add {
    width: 28px;
    height: 28px;
    background: #1B7C34;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
}

.btn-add:hover {
    background: #2DA544;
}

/* 方案列表 */
.profile-list {
    flex: 1;
    padding: 8px;
    overflow-y: auto;
}

.profile-card {
    padding: 12px;
    margin-bottom: 8px;
    background: #F9F9F9;
    border: 1px solid #E1E4E8;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.profile-card:hover {
    background: #F3F4F6;
    border-color: #34D399;
    box-shadow: 0 2px 8px rgba(52, 211, 153, 0.12);
}

.profile-card.active {
    background: #E8F5E9;
    border-color: #34D399;
    box-shadow: 0 2px 8px rgba(52, 211, 153, 0.15);
}

.profile-info {
    flex: 1;
    min-width: 0;
}

.profile-name {
    font-size: 13px;
    font-weight: 600;
    color: #24292E;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.profile-meta {
    display: flex;
    gap: 8px;
    font-size: 11px;
    color: #6A737D;
    margin-top: 4px;
}

.profile-actions {
    display: flex;
    gap: 4px;
    opacity: 0;
    transition: opacity 0.15s;
}

.profile-card:hover .profile-actions {
    opacity: 1;
}

.btn-edit {
    width: 24px;
    height: 24px;
    background: transparent;
    color: #D0D7DE;
    border: none;
    cursor: pointer;
    border-radius: 3px;
    transition: all 0.2s;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0;
}

.btn-edit img {
    width: 16px;
    height: 16px;
    object-fit: contain;
    filter: brightness(0.5);
}

.btn-edit:hover {
    background: #EEF;
    color: #0969DA;
}

.btn-edit:hover img {
    filter: brightness(1);
}

.btn-delete {
    width: 24px;
    height: 24px;
    background: transparent;
    color: #D0D7DE;
    border: none;
    cursor: pointer;
    border-radius: 3px;
    transition: all 0.2s;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0;
}

.btn-delete img {
    width: 16px;
    height: 16px;
    object-fit: contain;
    filter: brightness(0.5);
}

.btn-delete:hover {
    background: #FEE;
    color: #CB2431;
}

.btn-delete:hover img {
    filter: brightness(1);
}

/* 备份部分 */
.panel-section {
    border-top: 1px solid #E1E4E8;
    padding: 16px 8px;
}

.section-header {
    padding: 0 8px 8px 8px;
}

.section-header h3 {
    margin: 0;
    font-size: 12px;
    font-weight: 600;
    color: #586069;
    text-transform: uppercase;
}

.backup-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.backup-item {
    padding: 10px;
    background: #F9F9F9;
    border: 1px solid #E1E4E8;
    border-radius: 4px;
    font-size: 11px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    transition: all 0.2s;
}

.backup-item:hover {
    background: #F3F4F6;
    border-color: #D0D7DE;
}

.backup-info {
    flex: 1;
}

.backup-time {
    color: #24292E;
    font-weight: 500;
}

.backup-size {
    color: #6A737D;
    font-size: 10px;
    margin-top: 2px;
}

.btn-restore {
    padding: 4px 8px;
    background: #1B7C34;
    color: white;
    border: none;
    border-radius: 3px;
    font-size: 10px;
    cursor: pointer;
    flex-shrink: 0;
}

.btn-restore:hover {
    background: #2DA544;
}

/* 可拖动分隔符 */
.divider {
    width: 4px;
    background: #E1E4E8;
    cursor: col-resize;
    transition: background 0.2s;
    position: relative;
}

.divider:hover {
    background: #D0D7DE;
}

/* 右侧面板 */
.right-panel {
    flex: 1;
    display: flex;
    flex-direction: column;
    background: #FFFFFF;
    overflow: hidden;
    transition: margin-right 0.25s ease;
}

/* 编辑器工具栏 */
.editor-toolbar {
    padding: 12px 16px;
    background: linear-gradient(135deg, #FFFFFF 0%, #F9F9F9 100%);
    border-bottom: 1px solid #E1E4E8;
    display: flex;
    gap: 12px;
    align-items: center;
    flex-shrink: 0;
}

/* 未保存更改指示 */
@keyframes pulse {

    0%,
    100% {
        opacity: 1;
    }

    50% {
        opacity: 0.6;
    }
}

.search-box {
    flex: 1;
    display: flex;
    align-items: center;
    background: #F6F8FA;
    border: 1px solid #D0D7DE;
    border-radius: 6px;
    padding: 0 8px;
}

.search-icon {
    margin-right: 4px;
    font-size: 14px;
    width: 16px;
    height: 16px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    object-fit: contain;
}

.search-input {
    flex: 1;
    background: transparent;
    border: none;
    color: #24292E;
    padding: 8px;
    font-size: 13px;
    outline: none;
}

.search-input::placeholder {
    color: #6A737D;
}

.filter-select {
    padding: 8px 12px;
    background: #F6F8FA;
    border: 1px solid #D0D7DE;
    color: #24292E;
    border-radius: 6px;
    font-size: 13px;
    cursor: pointer;
}

.filter-select:focus {
    outline: none;
    border-color: #34D399;
}

/* 编辑器容器 */
.editor-container {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
}

/* 底部状态栏 */
.status-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 16px;
    background: #F6F8FA;
    border-top: 1px solid #E1E4E8;
    font-size: 12px;
    color: #586069;
    flex-shrink: 0;
    margin-top: auto;
}

.status-left {
    display: flex;
    gap: 16px;
}

.status-right {
    display: flex;
    gap: 16px;
}

.status-item {
    display: flex;
    align-items: center;
    white-space: nowrap;
}

.unsaved-item {
    color: #D1242F;
    font-weight: 500;
    animation: pulse 1.5s ease-in-out infinite;
}

.saved-item {
    color: #1B7C34;
}

/* 添加条目栏 */
.add-entry-bar {
    display: flex;
    gap: 8px;
    padding: 12px 16px;
    background: #F3F4F6;
    border-bottom: 1px solid #E1E4E8;
    flex-shrink: 0;
    flex-wrap: wrap;
}

.entry-input {
    flex: 1;
    padding: 8px 12px;
    background: #FFFFFF;
    border: 1px solid #D0D7DE;
    color: #24292E;
    border-radius: 4px;
    font-size: 13px;
}

.entry-input:focus {
    outline: none;
    border-color: #34D399;
    background: #FFFFFF;
    box-shadow: 0 0 0 2px rgba(52, 211, 153, 0.1);
}

.entry-input::placeholder {
    color: #6A737D;
}

.btn-add-entry {
    padding: 8px 12px;
    background: #1B7C34;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-weight: 500;
    font-size: 12px;
    white-space: nowrap;
}

.btn-add-entry:hover {
    background: #2DA544;
}

/* 代码编辑器区域 */
.code-editor {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow-y: auto;
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace;
    font-size: 13px;
    line-height: 1.5;
    background: #FAFBFC;
}

.loading-state,
.empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 120px;
    color: #6A737D;
}

.spinner {
    width: 24px;
    height: 24px;
    border: 3px solid #E1E4E8;
    border-top-color: #34D399;
    border-radius: 50%;
    animation: spin 1s linear infinite;
}

@keyframes spin {
    to {
        transform: rotate(360deg);
    }
}

.empty-state p {
    margin: 0 0 8px 0;
}

.empty-state small {
    font-size: 12px;
}

/* 表格样式 */
.entries-table {
    width: 100%;
    border-collapse: collapse;
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace;
    font-size: 13px;
    flex-shrink: 0;
}

.entries-table thead {
    background: #F3F4F6;
    position: sticky;
    top: 0;
    z-index: 10;
}

.entries-table th {
    padding: 8px 6px;
    text-align: left;
    font-weight: 600;
    color: #24292E;
    font-size: 12px;
    border-bottom: 2px solid #34D399;
}

.entries-table td {
    padding: 6px;
    border-bottom: 1px solid #E1E4E8;
    vertical-align: middle;
}

.entries-table tbody tr {
    transition: background-color 0.15s;
}

.entries-table tbody tr:hover {
    background: #E8F5E9;
}

.entries-table tbody tr.disabled {
    opacity: 0.6;
}

/* 列宽 */
.col-index {
    width: 40px;
    text-align: center;
}

.col-status {
    width: 45px;
    text-align: center;
}

.col-domain {
    width: 200px;
}

.col-ip {
    width: 110px;
    text-align: left;
}

.col-alternate-ip {
    width: 110px;
    text-align: left;
    color: #586069;
    font-size: 13px;
}

.col-comment {
    flex: 1;
}

.col-actions {
    width: 80px;
    text-align: right;
}

.entries-table tbody td.col-actions {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 4px;
}

.line-number {
    display: inline-block;
    width: 48px;
    text-align: right;
    color: #6A737D;
    user-select: none;
    padding-right: 12px;
    border-right: 2px solid #E1E4E8;
    margin-right: 12px;
}

.toggle-indicator {
    width: 24px;
    height: 24px;
    background: transparent;
    border: none;
    color: #CB2431;
    cursor: pointer;
    font-size: 16px;
    padding: 0;
    flex-shrink: 0;
    transition: all 0.2s;
    font-weight: bold;
    display: flex;
    align-items: center;
    justify-content: center;
}

.toggle-indicator img,
.status-img {
    width: 16px;
    height: 16px;
    object-fit: contain;
}

.toggle-indicator:hover {
    color: #D1242F;
    transform: scale(1.2);
}

.toggle-indicator:hover .status-img {
    filter: brightness(1.1);
}

.toggle-indicator.enabled {
    color: #1B7C34;
}

.entry-ip {
    color: #1B7C34;
    font-weight: 500;
}

.entry-domain {
    color: #0969DA;
    word-break: break-all;
}

.entry-comment {
    color: #6A737D;
    word-break: break-all;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.line-actions {
    display: flex;
    gap: 8px;
    opacity: 0;
    transition: opacity 0.15s;
    flex-shrink: 0;
    color: #24292E;
    font-size: 12px;
    font-weight: 600;
}

.code-line:hover .line-actions {
    opacity: 1;
}

.code-line.header-row .line-actions {
    opacity: 1;
    color: #24292E;
}

.action-btn {
    width: 28px;
    height: 28px;
    background: transparent;
    border: none;
    color: #6A737D;
    cursor: pointer;
    border-radius: 4px;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0;
}

.action-btn img {
    width: 16px;
    height: 16px;
    object-fit: contain;
}

.action-btn:hover {
    background: #F0F0F0;
}

.action-btn:hover img {
    filter: brightness(0.8);
}

.swap-btn img {
    filter: brightness(1);
}

.swap-btn {
    color: #0366D6;
}

.swap-btn:hover {
    background: #D4E4FF;
    color: #0366D6;
}

/* 模态框 */
.modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.3);
    backdrop-filter: blur(4px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1001;
    animation: fadeIn 0.2s ease-out;
}

@keyframes fadeIn {
    from {
        opacity: 0;
    }

    to {
        opacity: 1;
    }
}

.modal-dialog {
    background: #FFFFFF;
    border: 1px solid #D0D7DE;
    border-radius: 8px;
    max-width: 400px;
    width: 90%;
    animation: scaleIn 0.3s ease-out;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

@keyframes scaleIn {
    from {
        transform: scale(0.9);
        opacity: 0;
    }

    to {
        transform: scale(1);
        opacity: 1;
    }
}

.modal-header {
    padding: 16px;
    border-bottom: 1px solid #E1E4E8;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.modal-header h3 {
    margin: 0;
    font-size: 16px;
    color: #24292E;
}

.btn-close {
    width: 28px;
    height: 28px;
    background: transparent;
    border: none;
    color: #6A737D;
    cursor: pointer;
    font-size: 20px;
}

.btn-close:hover {
    color: #24292E;
    background: #F0F0F0;
}

.modal-body {
    padding: 16px;
}

.form-group {
    margin-bottom: 16px;
}

.form-group:last-child {
    margin-bottom: 0;
}

.form-group label {
    display: block;
    font-size: 13px;
    font-weight: 500;
    color: #24292E;
    margin-bottom: 6px;
}

.modal-input {
    width: 100%;
    padding: 8px 12px;
    background: #F6F8FA;
    border: 1px solid #D0D7DE;
    color: #24292E;
    border-radius: 6px;
    font-size: 13px;
    box-sizing: border-box;
}

.modal-input:focus {
    outline: none;
    border-color: #34D399;
    background: #FFFFFF;
    box-shadow: 0 0 0 2px rgba(52, 211, 153, 0.1);
}

.modal-input::placeholder {
    color: #6A737D;
}

.modal-footer {
    padding: 12px 16px;
    border-top: 1px solid #E1E4E8;
    display: flex;
    gap: 8px;
    justify-content: flex-end;
}

/* 自定义滚动条 */
::-webkit-scrollbar {
    width: 8px;
    height: 8px;
}

::-webkit-scrollbar-track {
    background: #F6F8FA;
}

::-webkit-scrollbar-thumb {
    background: #D0D7DE;
    border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
    background: #B8BEC4;
}

/* 响应式设计 */
@media (max-width: 1200px) {
    .col-domain {
        width: 150px;
    }

    .col-ip {
        width: 100px;
    }

    .col-alternate-ip {
        width: 100px;
    }

    .entry-input {
        min-width: 100px;
    }
}

@media (max-width: 768px) {
    .module-layout {
        flex-direction: column;
    }

    .profile-sidebar {
        width: 100% !important;
        max-height: 300px;
        border-right: none;
        border-bottom: 1px solid #E1E4E8;
    }

    .profile-sidebar.collapsed {
        max-height: 0 !important;
        min-height: 0 !important;
        padding: 0 !important;
        margin: 0 !important;
    }

    .add-entry-bar {
        flex-wrap: wrap;
    }

    .entry-input {
        min-width: 80px;
        flex-grow: 1;
    }

    .col-domain,
    .col-ip,
    .col-alternate-ip,
    .col-comment {
        flex: 1;
        width: auto !important;
        text-align: left !important;
    }

    .col-actions {
        width: 70px;
    }
}
</style>
