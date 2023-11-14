<script setup>
const props = defineProps({
    size: {
        type: [Number, String],
        default: 14,
    },
    strokeWidth: {
        type: [Number, String],
        default: 3,
    },
})
</script>

<template>
    <svg :height="props.size" :width="props.size" fill="none" viewBox="0 0 48 48">
        <path
            :stroke-width="props.strokeWidth"
            d="M8 8L40 40"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="4" />
        <path
            :stroke-width="props.strokeWidth"
            d="M8 40L40 8"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="4" />
    </svg>
</template>

<style lang="scss" scoped></style>
