/**
 * HTTP API adapter layer - replaces Wails RPC bindings for web mode.
 * All functions match the original Wails-generated function signatures.
 */

const API_BASE = '/api'

async function post(path, body = {}) {
    const resp = await fetch(`${API_BASE}${path}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'same-origin',
        body: JSON.stringify(body),
    })
    if (resp.status === 401) {
        window.dispatchEvent(new Event('rdm:unauthorized'))
        return { success: false, msg: 'unauthorized' }
    }
    return resp.json()
}

async function get(path, params = {}) {
    const query = new URLSearchParams(params).toString()
    const url = query ? `${API_BASE}${path}?${query}` : `${API_BASE}${path}`
    const resp = await fetch(url, { credentials: 'same-origin' })
    if (resp.status === 401) {
        window.dispatchEvent(new Event('rdm:unauthorized'))
        return { success: false, msg: 'unauthorized' }
    }
    return resp.json()
}

async function del(path, params = {}) {
    const query = new URLSearchParams(params).toString()
    const url = query ? `${API_BASE}${path}?${query}` : `${API_BASE}${path}`
    const resp = await fetch(url, { method: 'DELETE', credentials: 'same-origin' })
    if (resp.status === 401) {
        window.dispatchEvent(new Event('rdm:unauthorized'))
        return { success: false, msg: 'unauthorized' }
    }
    return resp.json()
}

// ==================== Connection Service ====================

export function ListConnection() {
    return get('/connection/list')
}

export function GetConnection(name) {
    return get('/connection/get', { name })
}

export function SaveConnection(name, param) {
    return post('/connection/save', { name, param })
}

export function SaveSortedConnection(conns) {
    return post('/connection/save-sorted', { conns })
}

export function TestConnection(param) {
    return post('/connection/test', param)
}

export function DeleteConnection(name) {
    return del('/connection/delete', { name })
}

export function CreateGroup(name) {
    return post('/connection/group/create', { name })
}

export function RenameGroup(name, newName) {
    return post('/connection/group/rename', { name, newName })
}

export function DeleteGroup(name, includeConn) {
    return del('/connection/group/delete', { name, includeConn })
}

export function SaveLastDB(name, db) {
    return post('/connection/save-last-db', { name, db })
}

export function SaveRefreshInterval(name, interval) {
    return post('/connection/save-refresh-interval', { name, interval })
}

export async function ExportConnections() {
    // Web mode: trigger browser download of connections zip
    try {
        const resp = await fetch(`${API_BASE}/connection/export-download`, {
            credentials: 'same-origin',
        })
        if (resp.status === 401) {
            window.dispatchEvent(new Event('rdm:unauthorized'))
            return { success: false, msg: 'unauthorized' }
        }
        if (!resp.ok) {
            const err = await resp.json().catch(() => ({}))
            return { success: false, msg: err.msg || 'export failed' }
        }
        const blob = await resp.blob()
        const url = URL.createObjectURL(blob)
        const a = document.createElement('a')
        const disposition = resp.headers.get('Content-Disposition') || ''
        const match = disposition.match(/filename=(.+)/)
        a.download = match ? match[1] : 'connections.zip'
        a.href = url
        a.click()
        URL.revokeObjectURL(url)
        return { success: true, data: { path: '' } }
    } catch {
        return { success: false, msg: 'export failed' }
    }
}

export async function ImportConnections() {
    // Web mode: open file picker, upload zip to backend
    return new Promise((resolve) => {
        const input = document.createElement('input')
        input.type = 'file'
        input.accept = '.zip'
        input.onchange = async () => {
            if (input.files && input.files[0]) {
                const formData = new FormData()
                formData.append('file', input.files[0])
                try {
                    const resp = await fetch(`${API_BASE}/connection/import-upload`, {
                        method: 'POST',
                        credentials: 'same-origin',
                        body: formData,
                    })
                    if (resp.status === 401) {
                        window.dispatchEvent(new Event('rdm:unauthorized'))
                        resolve({ success: false, msg: 'unauthorized' })
                        return
                    }
                    resolve(await resp.json())
                } catch {
                    resolve({ success: false, msg: 'import failed' })
                }
            } else {
                resolve({ success: false, msg: '' })
            }
        }
        // User cancelled file picker
        input.addEventListener('cancel', () => resolve({ success: false, msg: '' }))
        input.click()
    })
}

export function ParseConnectURL(url) {
    return post('/connection/parse-url', { url })
}

export function ListSentinelMasters(param) {
    return post('/connection/list-sentinel-masters', param)
}

// ==================== Browser Service ====================

export function OpenConnection(name) {
    return post('/browser/open-connection', { name })
}

export function CloseConnection(name) {
    return post('/browser/close-connection', { name })
}

export function OpenDatabase(server, db) {
    return post('/browser/open-database', { server, db })
}

export function ServerInfo(name) {
    return post('/browser/server-info', { name })
}

export function LoadNextKeys(server, db, match, keyType, exactMatch) {
    return post('/browser/load-next-keys', { server, db, match, keyType, exactMatch })
}

export function LoadNextAllKeys(server, db, match, keyType, exactMatch) {
    return post('/browser/load-next-all-keys', { server, db, match, keyType, exactMatch })
}

export function LoadAllKeys(server, db, match, keyType, exactMatch) {
    return post('/browser/load-all-keys', { server, db, match, keyType, exactMatch })
}

export function GetKeyType(param) {
    return post('/browser/get-key-type', param)
}

export function GetKeySummary(param) {
    return post('/browser/get-key-summary', param)
}

export function GetKeyDetail(param) {
    return post('/browser/get-key-detail', param)
}

export function ConvertValue(value, decode, format) {
    return post('/browser/convert-value', { value, decode, format })
}

export function SetKeyValue(param) {
    return post('/browser/set-key-value', param)
}

export function GetHashValue(param) {
    return post('/browser/get-hash-value', param)
}

export function SetHashValue(param) {
    return post('/browser/set-hash-value', param)
}

export function AddHashField(server, db, key, action, fieldItems) {
    return post('/browser/add-hash-field', { server, db, key, action, fieldItems })
}

export function AddListItem(server, db, key, action, items) {
    return post('/browser/add-list-item', { server, db, key, action, items })
}

export function SetListItem(param) {
    return post('/browser/set-list-item', param)
}

export function SetSetItem(server, db, key, remove, members) {
    return post('/browser/set-set-item', { server, db, key, remove, members })
}

export function UpdateSetItem(param) {
    return post('/browser/update-set-item', param)
}

export function UpdateZSetValue(param) {
    return post('/browser/update-zset-value', param)
}

export function AddZSetValue(server, db, key, action, valueScore) {
    return post('/browser/add-zset-value', { server, db, key, action, valueScore })
}

export function AddStreamValue(server, db, key, id, fieldItems) {
    return post('/browser/add-stream-value', { server, db, key, id, fieldItems })
}

export function RemoveStreamValues(server, db, key, ids) {
    return post('/browser/remove-stream-values', { server, db, key, ids })
}

export function SetKeyTTL(server, db, key, ttl) {
    return post('/browser/set-key-ttl', { server, db, key, ttl })
}

export function BatchSetTTL(server, db, keys, ttl, serialNo) {
    return post('/browser/batch-set-ttl', { server, db, keys, ttl, serialNo })
}

export function DeleteKey(server, db, key, async) {
    return post('/browser/delete-key', { server, db, key, async })
}

export function DeleteKeys(server, db, keys, serialNo) {
    return post('/browser/delete-keys', { server, db, keys, serialNo })
}

export function DeleteKeysByPattern(server, db, pattern) {
    return post('/browser/delete-keys-by-pattern', { server, db, pattern })
}

export function RenameKey(server, db, key, newKey) {
    return post('/browser/rename-key', { server, db, key, newKey })
}

export function ExportKey(server, db, keys, path, includeExpire) {
    return post('/browser/export-key', { server, db, keys, path, includeExpire })
}

export function ImportCSV(server, db, path, conflict, ttl) {
    return post('/browser/import-csv', { server, db, path, conflict, ttl })
}

export function FlushDB(server, db, async) {
    return post('/browser/flush-db', { server, db, async })
}

export function GetSlowLogs(server, db, num) {
    return post('/browser/get-slow-logs', { server, db, num })
}

export function GetClientList(server, db) {
    return post('/browser/get-client-list', { server, db })
}

export function GetCmdHistory() {
    return post('/browser/get-cmd-history')
}

export function CleanCmdHistory() {
    return post('/browser/clean-cmd-history')
}

// ==================== CLI Service ====================

export function StartCli(server, db) {
    return post('/cli/start', { server, db })
}

export function CloseCli(server) {
    return post('/cli/close', { server })
}

// ==================== Monitor Service ====================

export function StartMonitor(server) {
    return post('/monitor/start', { server })
}

export function StopMonitor(server) {
    return post('/monitor/stop', { server })
}

export function ExportLog(logs) {
    return post('/monitor/export-log', { logs })
}

// ==================== Pubsub Service ====================

export function Publish(server, channel, payload) {
    return post('/pubsub/publish', { server, channel, payload })
}

export function StartSubscribe(server) {
    return post('/pubsub/subscribe', { server })
}

export function StopSubscribe(server) {
    return post('/pubsub/unsubscribe', { server })
}

// ==================== Preferences Service ====================

export function GetPreferences() {
    return get('/preferences/get')
}

export function SetPreferences(pf) {
    return post('/preferences/set', pf)
}

export function UpdatePreferences(value) {
    return post('/preferences/update', value)
}

export function RestorePreferences() {
    return post('/preferences/restore')
}

// Common fonts to probe when Local Font Access API is unavailable
const CANDIDATE_FONTS = [
    // Sans-serif
    'Arial', 'Helvetica', 'Helvetica Neue', 'Verdana', 'Geneva', 'Tahoma',
    'Trebuchet MS', 'Lucida Grande', 'Lucida Sans Unicode', 'Segoe UI',
    'Roboto', 'Noto Sans', 'Open Sans', 'Lato', 'Source Sans Pro',
    // Serif
    'Times New Roman', 'Georgia', 'Palatino', 'Book Antiqua', 'Cambria',
    'Noto Serif',
    // Monospace
    'Courier New', 'Consolas', 'Monaco', 'Menlo', 'DejaVu Sans Mono',
    'Fira Code', 'JetBrains Mono', 'Source Code Pro', 'Cascadia Code',
    // CJK
    'Microsoft YaHei', 'PingFang SC', 'PingFang TC', 'Hiragino Sans GB',
    'Noto Sans SC', 'Noto Sans TC', 'Noto Sans JP', 'Noto Sans KR',
    'Source Han Sans SC', 'Source Han Sans TC', 'WenQuanYi Micro Hei',
    'Yu Gothic', 'Meiryo', 'Malgun Gothic',
]

async function queryBrowserFonts() {
    await document.fonts.ready
    return CANDIDATE_FONTS.filter((f) => document.fonts.check(`16px "${f}"`)).map((name) => ({ name, path: '' }))
}

export async function GetFontList() {
    try {
        const fonts = await queryBrowserFonts()
        return { success: true, data: { fonts } }
    } catch (_) {
        return { success: true, data: { fonts: [] } }
    }
}

export function GetBuildInDecoder() {
    return get('/preferences/buildin-decoder')
}

export function GetAppVersion() {
    return get('/preferences/version')
}

export function CheckForUpdate() {
    return get('/preferences/check-update')
}

// ==================== System Service ====================

// Alias used in App.vue
export function Info() {
    return get('/system/info')
}

// Web replacement for native file dialog
export async function SelectFile(title, ext) {
    return new Promise((resolve) => {
        const input = document.createElement('input')
        input.type = 'file'
        if (ext && Array.isArray(ext) && ext.length > 0) {
            input.accept = ext.map((e) => '.' + e.replace(/^\./, '')).join(',')
        }
        input.onchange = async () => {
            if (input.files && input.files[0]) {
                const formData = new FormData()
                formData.append('file', input.files[0])
                try {
                    const resp = await fetch('/api/system/select-file', {
                        method: 'POST',
                        credentials: 'same-origin',
                        body: formData,
                    })
                    if (resp.status === 401) {
                        window.dispatchEvent(new Event('rdm:unauthorized'))
                        resolve({ success: false, msg: 'unauthorized' })
                        return
                    }
                    resolve(await resp.json())
                } catch {
                    resolve({ success: false, msg: 'upload failed' })
                }
            } else {
                resolve({ success: false, msg: '' })
            }
        }
        input.addEventListener('cancel', () => resolve({ success: false, msg: '' }))
        input.click()
    })
}

export async function SaveFile(title, defaultName, ext) {
    // In web mode, file save dialogs are not applicable
    // The backend ExportLog etc. will handle download differently
    return { success: true, data: { path: '' } }
}

// ==================== Auth Service ====================

export async function Login(username, password) {
    return await post('/auth/login', { username, password })
}

export async function Logout() {
    return await post('/auth/logout')
}
