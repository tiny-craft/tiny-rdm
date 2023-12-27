<script setup>
import { SelectFile } from 'wailsjs/go/services/systemService.js'
import { get, isEmpty } from 'lodash'

const props = defineProps({
    value: String,
    placeholder: String,
    disabled: Boolean,
    ext: String,
})

const emit = defineEmits(['update:value'])

const handleSelectFile = async () => {
    const { success, data } = await SelectFile('', isEmpty(props.ext) ? null : [props.ext])
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
        <n-input v-model:value="props.value" :disabled="props.disabled" :placeholder="placeholder" clearable />
        <n-button :disabled="props.disabled" :focusable="false" @click="handleSelectFile">...</n-button>
    </n-input-group>
</template>

<style lang="scss" scoped></style>
