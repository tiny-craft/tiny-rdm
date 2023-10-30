<script setup>
const props = defineProps({
    inverse: {
        type: Boolean,
        default: false,
    },
    strokeWidth: {
        type: [Number, String],
        default: 3,
    },
    strokeColor: {
        type: String,
        default: '#FFF',
    },
})
</script>

<template>
    <svg fill="none" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
        <rect
            :fill="props.inverse ? 'currentColor' : 'none'"
            :stroke-width="props.strokeWidth"
            height="34"
            stroke="currentColor"
            stroke-linejoin="round"
            width="28"
            x="13"
            y="10" />
        <path
            d="M35 10V4H8C7.44772 4 7 4.44772 7 5V38H13"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="3" />
        <path
            :stroke="props.inverse ? props.strokeColor : 'currentColor'"
            :stroke-width="props.strokeWidth"
            d="M21 22H33"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke="props.inverse ? props.strokeColor : 'currentColor'"
            :stroke-width="props.strokeWidth"
            d="M21 30H33"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
