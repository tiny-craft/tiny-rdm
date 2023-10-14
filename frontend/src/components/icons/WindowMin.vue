<script setup>
const props = defineProps({
    size: {
        type: [Number, String],
        default: 14,
    },
})
</script>

<template>
    <svg :width="props.size" :height="props.size" viewBox="0 0 48 48" fill="none">
        <path d="M8 24L40 24" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
