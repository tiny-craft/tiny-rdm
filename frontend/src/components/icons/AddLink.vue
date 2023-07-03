<script setup>
const props = defineProps({
    strokeWidth: {
        type: [Number, String],
        default: 3,
    },
})
</script>

<template>
    <svg viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path
            d="M24.0605 10L24.0239 38"
            stroke="currentColor"
            :stroke-width="strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            d="M10 24L38 24"
            stroke="currentColor"
            :stroke-width="strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
    </svg>
</template>

<style lang="scss" scoped></style>
