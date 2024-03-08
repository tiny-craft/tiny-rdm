<script setup>
const props = defineProps({
    strokeWidth: {
        type: [Number, String],
        default: 3,
    },
})
</script>

<template>
    <svg
        :stroke-width="props.strokeWidth"
        class="lucide lucide-spell-check"
        fill="none"
        stroke="currentColor"
        stroke-linecap="round"
        stroke-linejoin="round"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg">
        <path d="m6 16 6-12 6 12" />
        <path d="M8 12h8" />
        <path d="m16 20 2 2 4-4" />
    </svg>
</template>

<style lang="scss" scoped></style>
