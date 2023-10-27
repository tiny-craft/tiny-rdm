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
    <svg v-if="props.inverse" fill="none" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
        <rect
            :stroke-width="props.strokeWidth"
            fill="currentColor"
            height="32"
            rx="2"
            stroke="currentColor"
            stroke-linejoin="round"
            width="40"
            x="4"
            y="8" />
        <path
            :stroke-width="props.strokeWidth"
            d="M12 18L19 24L12 30"
            stroke="#FFF"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M23 32H36"
            stroke="#FFF"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
    <svg v-else fill="none" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
        <rect
            :stroke-width="props.strokeWidth"
            fill="none"
            height="32"
            rx="2"
            stroke="currentColor"
            stroke-linejoin="round"
            width="40"
            x="4"
            y="8" />
        <path
            :stroke-width="props.strokeWidth"
            d="M12 18L19 24L12 30"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M23 32H36"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
