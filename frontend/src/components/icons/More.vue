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
            d="M24 44C35.0457 44 44 35.0457 44 24C44 12.9543 35.0457 4 24 4C12.9543 4 4 12.9543 4 24C4 35.0457 12.9543 44 24 44Z"
            fill="none"
            stroke="currentColor"
            :stroke-width="strokeWidth"
            stroke-linejoin="round" />
        <circle cx="14" cy="24" r="3" fill="currentColor" />
        <circle cx="24" cy="24" r="3" fill="currentColor" />
        <circle cx="34" cy="24" r="3" fill="currentColor" />
    </svg>
</template>

<style lang="scss" scoped></style>
