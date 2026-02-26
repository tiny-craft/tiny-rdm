/**
 * Web-mode stubs for wailsjs/runtime/runtime.js
 * Replaces Wails desktop runtime functions with browser equivalents.
 */

import { offWsEvent, onWsEvent, reconnectWebSocket, sendWsMessage, waitForWebSocket } from '@/utils/websocket.js'

// Don't auto-connect — wait for explicit call after login
// connectWebSocket()

// ==================== Events ====================

export function EventsOn(event, callback) {
    onWsEvent(event, callback)
}

export function EventsOnce(event, callback) {
    const wrapper = (...args) => {
        offWsEvent(event, wrapper)
        callback(...args)
    }
    onWsEvent(event, wrapper)
}

export function EventsEmit(event, ...data) {
    sendWsMessage({ event, data: data.length === 1 ? data[0] : data })
}

export function EventsOff(event) {
    offWsEvent(event)
}

// ==================== Clipboard ====================

export async function ClipboardGetText() {
    try {
        return await navigator.clipboard.readText()
    } catch {
        // clipboard.readText() requires HTTPS + user permission grant
        // Throw so callers can show a meaningful error instead of silent empty string
        throw new Error('clipboard permission denied')
    }
}

export async function ClipboardSetText(text) {
    try {
        await navigator.clipboard.writeText(text)
    } catch {
        // fallback
        const ta = document.createElement('textarea')
        ta.value = text
        ta.style.position = 'fixed'
        ta.style.left = '-9999px'
        document.body.appendChild(ta)
        ta.select()
        document.execCommand('copy')
        document.body.removeChild(ta)
    }
}

// ==================== Browser ====================

export function BrowserOpenURL(url) {
    window.open(url, '_blank')
}

// ==================== Window Management (no-ops for web) ====================

export function WindowMinimise() {}
export function WindowMaximise() {}
export function WindowToggleMaximise() {}
export function WindowIsMaximised() { return false }
export function WindowIsFullscreen() { return false }
export function WindowSetDarkTheme() {}
export function WindowSetLightTheme() {}
export function Quit() {}

// ==================== Environment ====================

export async function Environment() {
    return {
        buildType: 'production',
        platform: 'web',
        arch: 'web',
    }
}

// ==================== WebSocket Management ====================

export { reconnectWebSocket as ReconnectWebSocket }
export { waitForWebSocket as WaitForWebSocket }
