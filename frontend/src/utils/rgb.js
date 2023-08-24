import { padStart, size, startsWith } from 'lodash'

/**
 * @typedef {Object} RGB
 * @property {number} r
 * @property {number} g
 * @property {number} b
 */

/**
 * parse hex color to rgb object
 * @param hex
 * @return {RGB}
 */
export function parseHexColor(hex) {
    if (size(hex) < 6) {
        return { r: 0, g: 0, b: 0 }
    }
    if (startsWith(hex, '#')) {
        hex = hex.slice(1)
    }
    const bigint = parseInt(hex, 16)
    const r = (bigint >> 16) & 255
    const g = (bigint >> 8) & 255
    const b = bigint & 255
    return { r, g, b }
}

/**
 * do gamma correction with an RGB object
 * @param {RGB} rgb
 * @param {Number} gamma
 * @return {RGB}
 */
export function hexGammaCorrection(rgb, gamma) {
    if (typeof rgb !== 'object') {
        return { r: 0, g: 0, b: 0 }
    }
    return {
        r: Math.max(0, Math.min(255, Math.round(rgb.r * gamma))),
        g: Math.max(0, Math.min(255, Math.round(rgb.g * gamma))),
        b: Math.max(0, Math.min(255, Math.round(rgb.b * gamma))),
    }
}

/**
 * RGB object to hex color string
 * @param {RGB} rgb
 * @return {string}
 */
export function toHexColor(rgb) {
    return (
        '#' +
        padStart(rgb.r.toString(16), 2, '0') +
        padStart(rgb.g.toString(16), 2, '0') +
        padStart(rgb.b.toString(16), 2, '0')
    )
}
