import { Environment } from 'wailsjs/runtime/runtime.js'

let os = ''

export async function loadEnvironment() {
    const env = await Environment()
    os = env.platform
}

export function isMacOS() {
    return os === 'darwin'
}

export function isWindows() {
    return os === 'windows'
}

export function isWeb() {
    return os === 'web'
}
