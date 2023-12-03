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
        <path clip-rule="evenodd" d="M20 24H44H20Z" fill="none" fill-rule="evenodd" />
        <path
            :stroke-width="props.strokeWidth"
            d="M20 24H44"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path clip-rule="evenodd" d="M20 38H44H20Z" fill="none" fill-rule="evenodd" />
        <path
            :stroke-width="props.strokeWidth"
            d="M20 38H44"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path clip-rule="evenodd" d="M20 10H44H20Z" fill="none" fill-rule="evenodd" />
        <path
            :stroke-width="props.strokeWidth"
            d="M20 10H44"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <rect
            :stroke-width="props.strokeWidth"
            fill="none"
            height="8"
            stroke="currentColor"
            stroke-linejoin="round"
            width="8"
            x="4"
            y="34" />
        <rect
            :stroke-width="props.strokeWidth"
            fill="none"
            height="8"
            stroke="currentColor"
            stroke-linejoin="round"
            width="8"
            x="4"
            y="20" />
        <rect
            :stroke-width="props.strokeWidth"
            fill="none"
            height="8"
            stroke="currentColor"
            stroke-linejoin="round"
            width="8"
            x="4"
            y="6" />
    </svg>
</template>

<style lang="scss" scoped></style>
