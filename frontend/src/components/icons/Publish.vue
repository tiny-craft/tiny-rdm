<script setup>
const props = defineProps({
    strokeWidth: {
        type: [Number, String],
        default: 3,
    },
})
</script>

<template>
    <svg fill="none" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
        <path
            :stroke-width="props.strokeWidth"
            d="M43 5L29.7 43L22.1 25.9L5 18.3L43 5Z"
            stroke="currentColor"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M43.0001 5L22.1001 25.9"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
