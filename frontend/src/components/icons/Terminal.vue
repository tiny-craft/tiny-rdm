<script setup>
const emit = defineEmits(['update:modelValue'])

const props = defineProps({
    inverse: {
        type: Boolean,
        default: false,
    },
    strokeWidth: {
        type: [Number, String],
        default: 3,
    },
})
</script>

<template>
    <svg v-if="props.inverse" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
        <rect
            x="4"
            y="8"
            width="40"
            height="32"
            rx="2"
            fill="currentColor"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linejoin="round" />
        <path
            d="M12 18L19 24L12 30"
            stroke="#FFF"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            d="M23 32H36"
            stroke="#FFF"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
    <svg v-else viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
        <rect
            x="4"
            y="8"
            width="40"
            height="32"
            rx="2"
            fill="none"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linejoin="round" />
        <path
            d="M12 18L19 24L12 30"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            d="M23 32H36"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
