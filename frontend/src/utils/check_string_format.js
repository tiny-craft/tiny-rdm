import { endsWith, isEmpty, size, startsWith } from 'lodash'

export const IsRedisKey = (str, separator) => {
    if (isEmpty(separator)) {
        separator = ':'
    }
}

/**
 * check string is json
 * @param str
 * @returns {boolean}
 * @constructor
 */
export const IsJson = (str) => {
    if (size(str) >= 2) {
        if (startsWith(str, '{') && endsWith(str, '}')) {
            return true
        }
        if (startsWith(str, '[') && endsWith(str, ']')) {
            return true
        }
    }
    return false
}
