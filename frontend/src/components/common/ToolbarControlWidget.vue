<script setup>
import WindowMin from '@/components/icons/WindowMin.vue'
import WindowMax from '@/components/icons/WindowMax.vue'
import WindowClose from '@/components/icons/WindowClose.vue'
import { computed } from 'vue'
import { useThemeVars } from 'naive-ui'
import { Quit, WindowMinimise, WindowToggleMaximise } from 'wailsjs/runtime/runtime.js'
import WindowRestore from '@/components/icons/WindowRestore.vue'

const themeVars = useThemeVars()
const props = defineProps({
    size: {
        type: Number,
        default: 35,
    },
    maximised: {
        type: Boolean,
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
    <n-space :size="0" :wrap-item="false" align="center" justify="center">
        <n-tooltip :delay="1000" :show-arrow="false">
            {{ $t('menu.minimise') }}
            <template #trigger>
                <div class="btn-wrapper" @click="handleMinimise">
                    <window-min />
                </div>
            </template>
        </n-tooltip>
        <n-tooltip v-if="maximised" :delay="1000" :show-arrow="false">
            {{ $t('menu.restore') }}
            <template #trigger>
                <div class="btn-wrapper" @click="handleMaximise">
                    <window-restore />
                </div>
            </template>
        </n-tooltip>
        <n-tooltip v-else :delay="1000" :show-arrow="false">
            {{ $t('menu.maximise') }}
            <template #trigger>
                <div class="btn-wrapper" @click="handleMaximise">
                    <window-max />
                </div>
            </template>
        </n-tooltip>
        <n-tooltip :delay="1000" :show-arrow="false">
            {{ $t('menu.close') }}
            <template #trigger>
                <div class="btn-wrapper" @click="handleClose">
                    <window-close />
                </div>
            </template>
        </n-tooltip>
    </n-space>
</template>

<style lang="scss" scoped>
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
