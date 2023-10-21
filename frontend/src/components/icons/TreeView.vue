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
            d="M38 20H18V28H38V20Z"
            fill="none"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linejoin="round" />
        <path
            d="M32 6H18V14H32V6Z"
            fill="none"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linejoin="round" />
        <path
            d="M44 34H18V42H44V34Z"
            fill="none"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linejoin="round" />
        <path
            d="M17 10H5"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            d="M17 24H5"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            d="M17 38H5"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            d="M5 44V4"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
