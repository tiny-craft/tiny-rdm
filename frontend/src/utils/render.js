import { h } from 'vue'
import { NIcon } from 'naive-ui'

export function useRender() {
    return {
        /**
         *
         * @param {string|Object} icon
         * @param {{}} [props]
         * @return {VNode}
         */
        renderIcon: (icon, props = {}) => {
            if (icon == null) {
                return undefined
            }
            return h(NIcon, null, {
                default: () => h(icon, props),
            })
        },

        /**
         *
         * @param {string} label
         * @param {{}} [props]
         * @return {VNode}
         */
        renderLabel: (label, props = {}) => {
            return h('div', props, label)
        },
    }
}
