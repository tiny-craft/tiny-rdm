<script setup>
const props = defineProps({
    strokeWidth: {
        type: [Number, String],
        default: 2,
    },
})
</script>

<template>
    <svg :stroke-width="props.strokeWidth" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
    </svg>
</template>

<style lang="scss" scoped></style>
