<script setup>
const props = defineProps({
    fillColor: {
        type: String,
        default: '#f2c55c',
    },
})
</script>

<template>
    <svg viewBox="0 0 512 512" xmlns="http://www.w3.org/2000/svg">
        <path
            :fill="props.fillColor"
            d="M218.1 167.17c0 13 0 25.6 4.1 37.4c-43.1 50.6-156.9 184.3-167.5 194.5a20.17 20.17 0 0 0-6.7 15c0 8.5 5.2 16.7 9.6 21.3c6.6 6.9 34.8 33 40 28c15.4-15 18.5-19 24.8-25.2c9.5-9.3-1-28.3 2.3-36s6.8-9.2 12.5-10.4s15.8 2.9 23.7 3c8.3.1 12.8-3.4 19-9.2c5-4.6 8.6-8.9 8.7-15.6c.2-9-12.8-20.9-3.1-30.4s23.7 6.2 34 5s22.8-15.5 24.1-21.6s-11.7-21.8-9.7-30.7c.7-3 6.8-10 11.4-11s25 6.9 29.6 5.9c5.6-1.2 12.1-7.1 17.4-10.4c15.5 6.7 29.6 9.4 47.7 9.4c68.5 0 124-53.4 124-119.2S408.5 48 340 48s-121.9 53.37-121.9 119.17zM400 144a32 32 0 1 1-32-32a32 32 0 0 1 32 32z"
            stroke="currentColor"
            stroke-linejoin="round"
            stroke-width="32"
        ></path>
    </svg>
    <!--  <svg viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">-->
    <!--    <path-->
    <!--        d="m38.3682,9.27721c2.5423,2.4953 3.5456,6.1544 2.6289,9.5881c-0.9166,3.4337 -3.6127,6.1156 -7.0646,7.0274c-3.4519,0.9119 -7.1303,-0.0861 -9.63875,-2.6151c-3.77167,-3.8845 -3.71773,-10.0592 0.12124,-13.878c3.83901,-3.8187 10.04631,-3.8724 13.95141,-0.1206l0.0018,-0.0018z"-->
    <!--        :fill="props.fillColor" stroke="currentColor" :stroke-width="props.strokeWidth" stroke-linejoin="round"/>-->
    <!--    <path d="m6.5,40.5l17,-17" stroke="currentColor" :stroke-width="props.strokeWidth" stroke-linecap="round"-->
    <!--          stroke-linejoin="round"/>-->
    <!--    <path d="m10.5,36.8l5.4285,5.4l6.3334,-6.3l-5.4286,-5.4l-6.3333,6.3z" :fill="props.fillColor"-->
    <!--          stroke="currentColor" :stroke-width="props.strokeWidth" stroke-linejoin="round"/>-->
    <!--  </svg>-->
</template>

<style lang="scss" scoped></style>
