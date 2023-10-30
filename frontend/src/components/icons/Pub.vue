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
        <path
            :stroke-width="props.strokeWidth"
            d="M33 38H22V30H36V22H44V38H39L36 41L33 38Z"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :fill="props.inverse ? 'currentColor' : 'none'"
            :stroke-width="props.strokeWidth"
            d="M4 6H36V30H17L13 34L9 30H4V6Z"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke="props.inverse ? props.strokeColor : 'currentColor'"
            :stroke-width="props.strokeWidth"
            d="M12 22H18"
            stroke-linecap="round" />
        <path
            :stroke="props.inverse ? props.strokeColor : 'currentColor'"
            :stroke-width="props.strokeWidth"
            d="M12 14H24"
            stroke-linecap="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
