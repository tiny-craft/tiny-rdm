import { includes, isEmpty, toUpper, trim } from 'lodash'

/**
 * join execute path and arguments into a command string
 * @param {string} path
 * @param {string[]} args
 * @param {string} [emptyContent]
 * @return {string}
 */
export const joinCommand = (path, args = [], emptyContent = '-') => {
    let cmd = ''
    path = trim(path)
    if (!isEmpty(path)) {
        let containValuePlaceholder = false
        cmd = includes(path, ' ') ? `"${path}"` : path
        if (args) {
            for (let part of args) {
                part = trim(part)
                if (isEmpty(part)) {
                    continue
                }
                if (includes(part, ' ')) {
                    cmd += ' "' + part + '"'
                } else {
                    if (toUpper(part) === '{VALUE}') {
                        part = '{VALUE}'
                        containValuePlaceholder = true
                    }
                    cmd += ' ' + part
                }
            }
        }
        if (!containValuePlaceholder) {
            cmd += ' {VALUE}'
        }
    }
    return cmd || emptyContent
}
