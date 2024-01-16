const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
/**
 * convert byte value
 * @param {number} bytes
 * @param {number} decimals
 * @return {{unit: string, value: number}}
 */
export const convertBytes = (bytes, decimals = 2) => {
    if (bytes <= 0) {
        return {
            value: 0,
            unit: sizes[0],
        }
    }

    const k = 1024
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    const j = Math.min(i, sizes.length - 1)
    return {
        value: parseFloat((bytes / Math.pow(k, j)).toFixed(decimals)),
        unit: sizes[j],
    }
}

/**
 *
 * @param {number} bytes
 * @param {number} decimals
 * @return {string}
 */
export const formatBytes = (bytes, decimals = 2) => {
    const res = convertBytes(bytes, decimals)
    return res.value + ' ' + res.unit
}
