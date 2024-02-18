<script setup>
const props = defineProps({
    value: String,
})

const emit = defineEmits(['update:value'])
</script>

<template>
    <n-form-item :label="$t('common.value')">
        <n-input
            :rows="6"
            :value="props.value"
            placeholder=""
            type="textarea"
            @input="(val) => emit('update:value', val)" />
    </n-form-item>
</template>

<style lang="scss" scoped></style>
