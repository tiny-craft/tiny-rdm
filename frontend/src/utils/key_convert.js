import { join, map } from 'lodash'

/**
 * converted binary data in strings to hex format
 * @param {string|number[]} key
 * @return {string}
 */
export function decodeRedisKey(key) {
    if (key instanceof Array) {
        // char array, convert to hex string
        return join(
            map(key, (k) => {
                if (k >= 32 && k <= 126) {
                    return String.fromCharCode(k)
                }
                return '\\x' + k.toString(16).toUpperCase().padStart(2, '0')
            }),
            '',
        )
    }

    return key
}

/**
 * convert char code array to string
 * @param {string|number[]} key
 * @return {string}
 */
export function nativeRedisKey(key) {
    if (key instanceof Array) {
        return map(key, (c) => String.fromCharCode(c)).join('')
    }
    return key
}
