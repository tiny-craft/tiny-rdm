let inited = false

/**
 * load umami analytics module
 * @param {boolean} allowTrack
 * @return {Promise<void>}
 */
export const loadModule = async (allowTrack = true) => {
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
}

const enable = () => {
    return inited && umami
}

export const trackEvent = async (event, data) => {
    if (enable() || event === 'startup') {
        umami.track(({ website, language }) => ({
            language,
            website,
            name: event,
            data,
        }))
    }
}
