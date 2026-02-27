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
        <circle cx="12" cy="12" r="10" />
        <path d="M12 2a10 10 0 0 1 0 20V2" />
    </svg>
</template>

<style lang="scss" scoped></style>
