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
            d="M20 6H6V20H20V6Z"
            fill="none"
            stroke="currentColor"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M20 28H6V42H20V28Z"
            fill="none"
            stroke="currentColor"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M42 6H28V20H42V6Z"
            fill="none"
            stroke="currentColor"
            stroke-linejoin="round" />
        <path :stroke-width="props.strokeWidth" d="M29 28V42" stroke="currentColor" stroke-linecap="round" />
        <path :stroke-width="props.strokeWidth" d="M41 28V42" stroke="currentColor" stroke-linecap="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
