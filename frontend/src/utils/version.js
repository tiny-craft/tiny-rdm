import { get, isEmpty, map, size, split, trimStart } from 'lodash'

const toVerArr = (ver) => {
    const v = trimStart(ver, 'v')
    let vParts = split(v, '.')
    if (isEmpty(vParts)) {
        vParts = ['0']
    }
    return map(vParts, (v) => {
        let vNum = parseInt(v)
        return isNaN(vNum) ? 0 : vNum
    })
}

/**
 * compare two version strings
 * @param {string} v1
 * @param {string} v2
 * @return {number}
 */
export const compareVersion = (v1, v2) => {
    if (v1 !== v2) {
        const v1Nums = toVerArr(v1)
        const v2Nums = toVerArr(v2)
        const length = Math.max(size(v1Nums), size(v2Nums))

        for (let i = 0; i < length; i++) {
            const num1 = get(v1Nums, i, 0)
            const num2 = get(v2Nums, i, 0)
            if (num1 !== num2) {
                return num1 > num2 ? 1 : -1
            }
        }
    }
    return 0
}
