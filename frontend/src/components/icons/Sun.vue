<script setup>
const props = defineProps({
    strokeWidth: {
        type: [Number, String],
        default: 2,
    },
})
</script>

<template>
    <svg
        :stroke-width="props.strokeWidth"
        fill="none"
        stroke="currentColor"
        stroke-linecap="round"
        stroke-linejoin="round"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg">
        <circle cx="12" cy="12" r="5" />
        <path
            d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42" />
    </svg>
</template>

<style lang="scss" scoped></style>
