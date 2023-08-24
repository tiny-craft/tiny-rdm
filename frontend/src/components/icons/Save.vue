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
            d="M39.3 6H8.7C7.20883 6 6 7.20883 6 8.7V39.3C6 40.7912 7.20883 42 8.7 42H39.3C40.7912 42 42 40.7912 42 39.3V8.7C42 7.20883 40.7912 6 39.3 6Z"
            fill="none"
            stroke="currentColor"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M32 6V24H15V6H32Z"
            fill="none"
            stroke="currentColor"
            stroke-linejoin="round" />
        <path :stroke-width="props.strokeWidth" d="M26 13V17" stroke="currentColor" stroke-linecap="round" />
        <path :stroke-width="props.strokeWidth" d="M10.9971 6H35.9986" stroke="currentColor" stroke-linecap="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
