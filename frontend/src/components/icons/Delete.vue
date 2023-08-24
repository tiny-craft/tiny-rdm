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
            d="M15 12L16.2 5H31.8L33 12"
            stroke="currentColor"
            stroke-linejoin="round" />
        <path :stroke-width="props.strokeWidth" d="M6 12H42" stroke="currentColor" stroke-linecap="round" />
        <path
            :stroke-width="props.strokeWidth"
            clip-rule="evenodd"
            d="M37 12L35 43H13L11 12H37Z"
            fill="none"
            fill-rule="evenodd"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path :stroke-width="props.strokeWidth" d="M20 22V34" stroke="currentColor" stroke-linecap="round" />
        <path :stroke-width="props.strokeWidth" d="M28 22V34" stroke="currentColor" stroke-linecap="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
