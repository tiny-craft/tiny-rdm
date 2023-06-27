<script setup>
const props = defineProps({
    strokeWidth: {
        type: [Number, String],
        default: 3,
    },
    fillColor: {
        type: String,
        default: '#f2c55c',
    },
})
</script>

<template>
    <svg fill="none" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
        <path
            :fill="props.fillColor"
            d="M4 11.9143L24 19L44 11.9143L24 5L4 11.9143Z"
            stroke="currentColor"
            stroke-linejoin="round"
            stroke-width="3"
        />
        <path
            :stroke-width="props.strokeWidth"
            d="M4 20L24 27L44 20"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            :stroke-width="props.strokeWidth"
            d="M4 28L24 35L44 28"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            :stroke-width="props.strokeWidth"
            d="M4 36L24 43L44 36"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
    </svg>
</template>

<style lang="scss" scoped></style>
