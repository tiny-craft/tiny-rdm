let inited = false

/**
 * load umami analytics module
 * @param {boolean} allowTrack
 * @return {Promise<void>}
 */
export const loadModule = async (allowTrack = true) => {
    try {
        await new Promise((resolve, reject) => {
            const script = document.createElement('script')
            script.setAttribute('src', 'https://analytics.tinycraft.cc/script.js')
            script.setAttribute('data-website-id', 'ad6de51d-1e27-44a5-958d-319679c56aec')
            script.setAttribute('data-cache', 'true')
            script.setAttribute('data-auto-track', allowTrack !== false ? 'true' : 'false')
            script.onload = () => {
                inited = true
                resolve()
            }
            script.onerror = () => {
                inited = false
                reject()
            }
            document.body.appendChild(script)
        })
    } catch {
        // Script blocked by CSP or network error — silently ignore
    }
}

const enable = () => {
    return inited && typeof umami !== 'undefined'
}

export const trackEvent = async (event, data) => {
    if (!enable()) {
        return
    }
    try {
        umami.track(({ website, language }) => ({
            language,
            website,
            name: event,
            data,
        }))
    } catch {
        // umami not available — silently ignore
    }
}
