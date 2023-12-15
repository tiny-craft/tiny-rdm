<script setup>
import { NInput } from 'naive-ui'

const props = defineProps({
    isEdit: Boolean,
    value: [String, Number],
    onUpdateValue: [Function, Array],
})

const emit = defineEmits(['update:value'])

const handleUpdateValue = (val) => {
    emit('update:value', val)
}
</script>

<template>
    <div style="min-height: 22px">
        <template v-if="props.isEdit">
            <n-input :value="props.value" @update:value="handleUpdateValue" />
        </template>
        <template v-else>
            {{ props.value }}
        </template>
    </div>
</template>
