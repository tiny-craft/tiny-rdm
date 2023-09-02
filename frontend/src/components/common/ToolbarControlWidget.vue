<script setup>
import WindowMin from '@/components/icons/WindowMin.vue'
import WindowMax from '@/components/icons/WindowMax.vue'
import WindowClose from '@/components/icons/WindowClose.vue'
import { computed } from 'vue'
import { useThemeVars } from 'naive-ui'
import { Quit, WindowMinimise, WindowToggleMaximise } from 'wailsjs/runtime/runtime.js'

const themeVars = useThemeVars()
const props = defineProps({
    size: {
        type: Number,
        default: 35,
    },
})

const buttonSize = computed(() => {
    return props.size + 'px'
})

const handleMinimise = async () => {
    WindowMinimise()
}

const handleMaximise = () => {
    WindowToggleMaximise()
}

const handleClose = () => {
    Quit()
}
</script>

<template>
    <n-space :wrap-item="false" align="center" justify="center" :size="0">
        <div class="btn-wrapper" @click="handleMinimise">
            <window-min />
        </div>
        <div class="btn-wrapper" @click="handleMaximise">
            <window-max />
        </div>
        <div class="btn-wrapper" @click="handleClose">
            <window-close />
        </div>
    </n-space>
</template>

<style scoped lang="scss">
.btn-wrapper {
    width: v-bind('buttonSize');
    height: v-bind('buttonSize');
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    --wails-draggable: none;

    &:hover {
        cursor: pointer;
    }

    &:not(:last-child) {
        &:hover {
            background-color: v-bind('themeVars.closeColorHover');
        }

        &:active {
            background-color: v-bind('themeVars.closeColorPressed');
        }
    }

    &:last-child {
        &:hover {
            background-color: v-bind('themeVars.primaryColorHover');
        }

        &:active {
            background-color: v-bind('themeVars.primaryColorPressed');
        }
    }
}
</style>
