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
            d="M7 42H43"
            stroke="currentColorr"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            :stroke-width="props.strokeWidth"
            d="M11 26.7199V34H18.3172L39 13.3081L31.6951 6L11 26.7199Z"
            fill="none"
            stroke="currentColor"
            stroke-linejoin="round"
        />
    </svg>
</template>

<style lang="scss" scoped></style>
