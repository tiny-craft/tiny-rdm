import { h } from 'vue'
import { NIcon } from 'naive-ui'

export function useRender() {
    return {
        /**
         *
         * @param {string|Object} icon
         * @param {{}} [props]
         * @return {*}
         */
        renderIcon: (icon, props = {}) => {
            return () => {
                return h(NIcon, null, {
                    default: () => h(icon, props),
                })
            }
        },
    }
}
