<script setup>
const props = defineProps({
    open: {
        type: Boolean,
        default: false,
    },
    strokeWidth: {
        type: [Number, String],
        default: 3,
    },
    fillColor: {
        type: String,
        default: '#ffce78',
    },
})
</script>

<template>
    <svg v-if="props.open" fill="none" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
        <path
            :stroke-width="props.strokeWidth"
            d="M4 9V41L9 21H39.5V15C39.5 13.8954 38.6046 13 37.5 13H24L19 7H6C4.89543 7 4 7.89543 4 9Z"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :fill="props.fillColor || 'currentColor'"
            :stroke-width="props.strokeWidth"
            d="M40 41L44 21H8.8125L4 41H40Z"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
    <svg v-else fill="none" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
        <path
            :fill="props.fillColor || 'currentColor'"
            :stroke-width="props.strokeWidth"
            d="M5 8C5 6.89543 5.89543 6 7 6H19L24 12H41C42.1046 12 43 12.8954 43 14V40C43 41.1046 42.1046 42 41 42H7C5.89543 42 5 41.1046 5 40V8Z"
            stroke="currentColor"
            stroke-linejoin="round" />
        <path :stroke-width="props.strokeWidth" d="M43 22H5" stroke="currentColor" stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M5 16V28"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M43 16V28"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
