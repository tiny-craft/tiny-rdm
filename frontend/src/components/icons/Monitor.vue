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
            stroke-linecap="round"
            stroke-linejoin="round"
            width="36"
            x="6"
            y="6" />
        <path
            :stroke="props.inverse ? props.strokeColor : 'currentColor'"
            :stroke-width="props.strokeWidth"
            d="M12 25H15L19 14L22 36L27 23L31 29L34 25H37"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
