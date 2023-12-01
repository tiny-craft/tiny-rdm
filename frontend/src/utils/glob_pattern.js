import { includes, isEmpty } from 'lodash'

const REDIS_GLOB_CHAR = ['?', '*', '[', ']', '{', '}']
export const isRedisGlob = (str) => {
    if (!isEmpty(str)) {
        for (const c of REDIS_GLOB_CHAR) {
            if (includes(str, c)) {
                return true
            }
        }
    }
    return false
}
