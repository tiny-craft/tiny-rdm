<script setup>
import { computed, reactive } from 'vue'
import { debounce, isEmpty, trim } from 'lodash'
import { NButton, NInput } from 'naive-ui'
import IconButton from '@/components/common/IconButton.vue'
import Help from '@/components/icons/Help.vue'

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
})

const emit = defineEmits(['filterChanged', 'matchChanged'])

/**
 *
 * @type {UnwrapNestedRefs<{filter: string, match: string}>}
 */
const inputData = reactive({
    match: '',
    filter: '',
})

const hasMatch = computed(() => {
    return !isEmpty(trim(inputData.match))
})

const hasFilter = computed(() => {
    return !isEmpty(trim(inputData.filter))
})

const onFullSearch = () => {
    inputData.filter = trim(inputData.filter)
    if (!isEmpty(inputData.filter)) {
        inputData.match = inputData.filter
        inputData.filter = ''
        emit('matchChanged', inputData.match, inputData.filter)
    }
}

const _onInput = () => {
    emit('filterChanged', inputData.filter)
}
const onInput = debounce(_onInput, props.debounceWait, { leading: true, trailing: true })

const onKeyup = (evt) => {
    if (evt.key === 'Enter') {
        onFullSearch()
    }
}

const onClearFilter = () => {
    inputData.filter = ''
    onClearMatch()
}

const onClearMatch = () => {
    const changed = !isEmpty(inputData.match)
    inputData.match = ''
    if (changed) {
        emit('matchChanged', inputData.match, inputData.filter)
    } else {
        emit('filterChanged', inputData.filter)
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
            clearable
            @clear="onClearFilter"
            @input="onInput"
            @keyup.enter="onKeyup">
            <template #prefix>
                <slot name="prefix" />
                <n-tooltip v-if="hasMatch">
                    <template #trigger>
                        <n-tag closable size="small" @close="onClearMatch">
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
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :component="Help" />
                        </template>
                        <div class="text-block" style="max-width: 600px">
                            {{ $t('dialogue.filter.filter_pattern_tip') }}
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
            border
            small
            stroke-width="4"
            t-tooltip="interface.full_search"
            @click="onFullSearch" />
        <n-button v-else :disabled="hasMatch && !hasFilter" :focusable="false" @click="onFullSearch">
            {{ $t('interface.full_search') }}
        </n-button>
        <slot name="append" />
    </n-input-group>
</template>

<style lang="scss" scoped></style>
