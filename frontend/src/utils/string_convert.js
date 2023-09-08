import { map, padStart } from 'lodash'

/**
 * convert string to json
 * @param str
 * @return {string}
 */
export const toJsonText = (str) => {
    try {
        const jsonObj = JSON.parse(str)
        return JSON.stringify(jsonObj, null, 2)
    } catch (e) {
        return str
    }
}

/**
 * convert string from base64
 * @param str
 * @return {string}
 */
export const fromBase64 = (str) => {
    try {
        return atob(str)
    } catch (e) {
        return str
    }
}

/**
 * convert string from base64 to json
 * @param str
 * @return {string}
 */
export const fromBase64Json = (str) => {
    try {
        const text = atob(str)
        const jsonObj = JSON.parse(text)
        return JSON.stringify(jsonObj, null, 2)
    } catch (e) {
        return str
    }
}

/**
 * convert string to hex string
 * @param str
 * @return {string}
 */
export const toHex = (str) => {
    const hexArr = map(str, (char) => {
        const charCode = char.charCodeAt(0)
        return charCode.toString(16)
    })
    return hexArr.join(' ')
}

/**
 * convert string to binary string
 * @param str
 * @return {string}
 */
export const toBinary = (str) => {
    const codeUnits = map(str, (char) => {
        let code = char.charCodeAt(0).toString(2)
        code = padStart(code, 8, '0')
        return code
    })
    return codeUnits.join(' ')
}
