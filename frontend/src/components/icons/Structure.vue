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
            :stroke-width="strokeWidth"
            d="M38 20H18V28H38V20Z"
            fill="none"
            stroke="currentColor"
            stroke-linejoin="round"
        />
        <path
            :stroke-width="strokeWidth"
            d="M32 6H18V14H32V6Z"
            fill="none"
            stroke="currentColor"
            stroke-linejoin="round"
        />
        <path
            :stroke-width="strokeWidth"
            d="M44 34H18V42H44V34Z"
            fill="none"
            stroke="currentColor"
            stroke-linejoin="round"
        />
        <path
            :stroke-width="strokeWidth"
            d="M17 10H5"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            :stroke-width="strokeWidth"
            d="M17 24H5"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            :stroke-width="strokeWidth"
            d="M17 38H5"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            :stroke-width="strokeWidth"
            d="M5 44V4"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
    </svg>
</template>

<style lang="scss" scoped></style>
