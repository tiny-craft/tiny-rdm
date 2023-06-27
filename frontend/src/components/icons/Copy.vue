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
            d="M13 38H41V16H30V4H13V38Z"
            fill="none"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            :stroke-width="props.strokeWidth"
            d="M30 4L41 16"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            :stroke-width="props.strokeWidth"
            d="M7 20V44H28"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path :stroke-width="props.strokeWidth" d="M19 20H23" stroke="currentColor" stroke-linecap="round" />
        <path :stroke-width="props.strokeWidth" d="M19 28H31" stroke="currentColor" stroke-linecap="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
