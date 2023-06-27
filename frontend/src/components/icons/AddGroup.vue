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
            d="M5 8C5 6.89543 5.89543 6 7 6H19L24 12H41C42.1046 12 43 12.8954 43 14V40C43 41.1046 42.1046 42 41 42H7C5.89543 42 5 41.1046 5 40V8Z"
            fill="none"
            stroke="currentColor"
            stroke-linejoin="round"
        />
        <path :stroke-width="props.strokeWidth" d="M18 27H30" stroke="currentColor" stroke-linecap="round" />
        <path :stroke-width="props.strokeWidth" d="M24 21L24 33" stroke="currentColor" stroke-linecap="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
