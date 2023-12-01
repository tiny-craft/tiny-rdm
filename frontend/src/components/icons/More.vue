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
        <circle :r="props.strokeWidth" cx="12" cy="24" fill="currentColor" />
        <circle :r="props.strokeWidth" cx="24" cy="24" fill="currentColor" />
        <circle :r="props.strokeWidth" cx="36" cy="24" fill="currentColor" />
    </svg>
</template>

<style lang="scss" scoped></style>
