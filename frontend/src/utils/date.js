import dayjs from 'dayjs'
import { i18nGlobal } from '@/utils/i18n.js'

/**
 * convert seconds number to human-readable string
 * @param {number} duration duration in seconds
 * @return {string}
 */
export const toHumanReadable = (duration) => {
    const dur = dayjs.duration(duration, 'seconds')
    const days = dur.days()
    if (days > 0) {
        return days + i18nGlobal.t('common.unit_day') + ' ' + dur.format('HH:mm:ss')
    } else {
        return dur.format('HH:mm:ss')
    }
}
