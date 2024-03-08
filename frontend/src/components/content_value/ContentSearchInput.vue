<script setup>
import { computed, nextTick, reactive } from 'vue'
import { debounce, isEmpty, trim } from 'lodash'
import { NButton, NInput } from 'naive-ui'
import IconButton from '@/components/common/IconButton.vue'
import SpellCheck from '@/components/icons/SpellCheck.vue'

const props = defineProps({
    fullSearchIcon: {
        type: [String, Object],
        default: null,
    },
    debounceWait: {
        type: Number,
        default: 500,
    },
    small: {
        type: Boolean,
        default: false,
    },
    useGlob: {
        type: Boolean,
        default: false,
    },
    exact: {
        type: Boolean,
        default: false,
    },
})

const emit = defineEmits(['filterChanged', 'matchChanged', 'exactChanged'])

/**
 *
 * @type {UnwrapNestedRefs<{filter: string, match: string, exact: boolean}>}
 */
const inputData = reactive({
    match: '',
    filter: '',
    exact: false,
})

const hasMatch = computed(() => {
    return !isEmpty(trim(inputData.match))
})

const hasFilter = computed(() => {
    return !isEmpty(trim(inputData.filter))
})

const onExactChecked = () => {
    // update search search result
    if (hasMatch.value) {
        nextTick(() => onForceFullSearch())
    }
}

const onFullSearch = () => {
    inputData.filter = trim(inputData.filter)
    if (!isEmpty(inputData.filter)) {
        inputData.match = inputData.filter
        inputData.filter = ''
        emit('matchChanged', inputData.match, inputData.filter, inputData.exact)
    }
}

const onForceFullSearch = () => {
    inputData.filter = trim(inputData.filter)
    emit('matchChanged', inputData.match, inputData.filter, inputData.exact)
}

const _onInput = () => {
    emit('filterChanged', inputData.filter, inputData.exact)
}
const onInput = debounce(_onInput, props.debounceWait, { leading: true, trailing: true })

const onClearFilter = () => {
    inputData.filter = ''
    onClearMatch()
}

const onUpdateMatch = () => {
    inputData.filter = inputData.match
    onClearMatch()
}

const onClearMatch = () => {
    const changed = !isEmpty(inputData.match)
    inputData.match = ''
    if (changed) {
        emit('matchChanged', inputData.match, inputData.filter, inputData.exact)
    } else {
        emit('filterChanged', inputData.filter, inputData.exact)
    }
}

defineExpose({
    reset: onClearFilter,
})
</script>

<template>
    <n-input-group>
        <slot name="prepend" />
        <n-input
            v-model:value="inputData.filter"
            :placeholder="$t('interface.filter')"
            :size="props.small ? 'small' : ''"
            :theme-overrides="{ paddingSmall: '0 3px', paddingMedium: '0 6px' }"
            clearable
            @clear="onClearFilter"
            @input="onInput"
            @keyup.enter="onFullSearch">
            <template #prefix>
                <slot name="prefix" />
                <n-tooltip v-if="hasMatch" placement="bottom">
                    <template #trigger>
                        <n-tag closable size="small" @close="onClearMatch" @dblclick="onUpdateMatch">
                            {{ inputData.match }}
                        </n-tag>
                    </template>
                    {{
                        $t('interface.full_search_result', {
                            pattern: props.useGlob ? inputData.match : '*' + inputData.match + '*',
                        })
                    }}
                </n-tooltip>
            </template>
            <template #suffix>
                <template v-if="props.useGlob">
                    <n-tooltip placement="bottom" trigger="hover">
                        <template #trigger>
                            <n-tag
                                v-model:checked="inputData.exact"
                                :checkable="true"
                                :type="props.exact ? 'primary' : 'default'"
                                size="small"
                                strong
                                style="padding: 0 5px"
                                @updateChecked="onExactChecked">
                                <n-icon :size="14">
                                    <spell-check :stroke-width="2" />
                                </n-icon>
                            </n-tag>
                        </template>
                        <div class="text-block" style="max-width: 600px">
                            {{ $t('dialogue.filter.exact_match_tip') }}
                        </div>
                    </n-tooltip>
                </template>
            </template>
        </n-input>

        <icon-button
            v-if="props.fullSearchIcon"
            :disabled="hasMatch && !hasFilter"
            :icon="props.fullSearchIcon"
            :size="small ? 16 : 20"
            :tooltip-delay="1"
            border
            small
            stroke-width="4"
            @click="onFullSearch">
            <template #tooltip>
                <div class="text-block" style="max-width: 600px">
                    {{ $t('dialogue.filter.filter_pattern_tip') }}
                </div>
            </template>
        </icon-button>
        <n-button v-else :disabled="hasMatch && !hasFilter" :focusable="false" @click="onFullSearch">
            {{ $t('interface.full_search') }}
        </n-button>
        <slot name="append" />
    </n-input-group>
</template>

<style lang="scss" scoped>
//:deep(.n-input__prefix) {
//    max-width: 50%;
//}
//:deep(.n-tag__content) {
//    overflow: hidden;
//}
</style>
