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
        <path
            d="M39 6H9C7.34315 6 6 7.34315 6 9V39C6 40.6569 7.34315 42 9 42H39C40.6569 42 42 40.6569 42 39V9C42 7.34315 40.6569 6 39 6Z"
            fill="none"
            stroke="currentColor"
            stroke-width="4" />
    </svg>
</template>

<style lang="scss" scoped></style>
