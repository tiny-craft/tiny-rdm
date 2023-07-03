<script setup>
const props = defineProps({
    strokeWidth: {
        type: [Number, String],
        default: 3,
    },
})
</script>

<template>
    <svg viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path
            fill-rule="evenodd"
            clip-rule="evenodd"
            d="M20 5.91406H28V13.9141H43V21.9141H5V13.9141H20V5.91406Z"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            d="M8 40H40V22H8V40Z"
            fill="none"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linejoin="round"
        />
        <path
            d="M16 39.8976V33.9141"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            d="M24 39.8977V33.8977"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            d="M32 39.8976V33.9141"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            d="M12 40H36"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
    </svg>
</template>

<style lang="scss" scoped></style>
