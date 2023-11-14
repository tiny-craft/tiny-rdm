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
            height="36"
            rx="3"
            stroke="currentColor"
            stroke-linejoin="round"
            width="36"
            x="6"
            y="6" />
        <rect
            :fill="props.inverse ? props.strokeColor : 'none'"
            :stroke="props.inverse ? props.strokeColor : 'currentColor'"
            :stroke-width="props.strokeWidth"
            height="8"
            stroke-linejoin="round"
            width="8"
            x="13"
            y="13" />
        <path
            :stroke="props.inverse ? props.strokeColor : 'currentColor'"
            :stroke-width="props.strokeWidth"
            d="M27 13L35 13"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke="props.inverse ? props.strokeColor : 'currentColor'"
            :stroke-width="props.strokeWidth"
            d="M27 20L35 20"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke="props.inverse ? props.strokeColor : 'currentColor'"
            :stroke-width="props.strokeWidth"
            d="M13 28L35 28"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke="props.inverse ? props.strokeColor : 'currentColor'"
            :stroke-width="props.strokeWidth"
            d="M13 35H35"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
