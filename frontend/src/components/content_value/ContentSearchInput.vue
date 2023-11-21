<script setup>
import { computed, reactive } from 'vue'
import { debounce, isEmpty, trim } from 'lodash'
import { NButton, NInput } from 'naive-ui'

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
const onInput = debounce(_onInput, 500, { leading: true, trailing: true })

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
        <n-input
            v-model:value="inputData.filter"
            :placeholder="$t('interface.filter')"
            clearable
            @clear="onClearFilter"
            @input="onInput">
            <template #prefix>
                <n-tooltip v-if="hasMatch">
                    <template #trigger>
                        <n-tag closable size="small" @close="onClearMatch">
                            {{ inputData.match }}
                        </n-tag>
                    </template>
                    {{ $t('interface.full_search_result', { pattern: inputData.match }) }}
                </n-tooltip>
            </template>
        </n-input>
        <n-button :disabled="hasMatch && !hasFilter" :focusable="false" @click="onFullSearch">
            {{ $t('interface.full_search') }}
        </n-button>
    </n-input-group>
</template>

<style lang="scss" scoped></style>
