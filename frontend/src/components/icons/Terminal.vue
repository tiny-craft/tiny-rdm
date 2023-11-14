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
            height="32"
            rx="2"
            stroke="currentColor"
            stroke-linejoin="round"
            width="40"
            x="4"
            y="8" />
        <path
            :stroke="props.inverse ? props.strokeColor : 'currentColor'"
            :stroke-width="props.strokeWidth"
            d="M12 18L19 24L12 30"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke="props.inverse ? props.strokeColor : 'currentColor'"
            :stroke-width="props.strokeWidth"
            d="M23 32H36"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
