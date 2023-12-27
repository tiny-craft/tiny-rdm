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
            d="M6 24V42H42V24"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M33 15L24 6L15 15"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M24 6V32"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
