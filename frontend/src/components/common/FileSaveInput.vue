<script setup>
import { SaveFile } from 'wailsjs/go/services/systemService.js'
import { get } from 'lodash'

const props = defineProps({
    value: String,
    placeholder: String,
    disabled: Boolean,
    defaultPath: String,
})

const emit = defineEmits(['update:value'])

const onInput = (val) => {
    emit('update:value', val)
}

const onClear = () => {
    emit('update:value', '')
}

const handleSaveFile = async () => {
    const { success, data } = await SaveFile(null, props.defaultPath, ['csv'])
    if (success) {
        const path = get(data, 'path', '')
        emit('update:value', path)
    } else {
        emit('update:value', '')
    }
}
</script>

<template>
    <n-input-group>
        <n-input
            :disabled="props.disabled"
            :placeholder="placeholder"
            :value="props.value"
            clearable
            @clear="onClear"
            @input="onInput" />
        <n-button :disabled="props.disabled" :focusable="false" @click="handleSaveFile">...</n-button>
    </n-input-group>
</template>

<style lang="scss" scoped></style>
