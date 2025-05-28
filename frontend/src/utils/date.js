import { i18nGlobal } from '@/utils/i18n.js'
import { padStart } from 'lodash'

/**
 * convert seconds number to human-readable string
 * @param {number} duration duration in seconds
 * @return {string}
 */
export const toHumanReadable = (duration) => {
    const days = Math.floor(duration / 86400)
    const hours = Math.floor((duration % 86400) / 3600)
    const minutes = Math.floor((duration % 3600) / 60)
    const seconds = duration % 60
    const time = `${padStart(hours, 2, '0')}:${padStart(minutes, 2, '0')}:${padStart(seconds, 2, '0')}`
    if (days > 0) {
        return days + i18nGlobal.t('common.unit_day') + ' ' + time
    } else {
        return time
    }
}
