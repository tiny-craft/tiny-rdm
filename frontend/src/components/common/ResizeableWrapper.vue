<script setup>
import { useThemeVars } from 'naive-ui'
import { ref } from 'vue'

/**
 * Resizeable component wrapper
 */
const themeVars = useThemeVars()

const props = defineProps({
    size: {
        type: Number,
        default: 100,
    },
    minSize: {
        type: Number,
        default: 300,
    },
    maxSize: {
        type: Number,
        default: 0,
    },
    offset: {
        type: Number,
        default: 0,
    },
    disabled: {
        type: Boolean,
        default: false,
    },
    borderWidth: {
        type: Number,
        default: 4,
    },
})

const emit = defineEmits(['update:size'])

const resizing = ref(false)
const hover = ref(false)

const handleResize = (evt) => {
    if (resizing.value) {
        let size = evt.clientX - props.offset
        if (size < props.minSize) {
            size = props.minSize
        }
        if (props.maxSize > 0 && size > props.maxSize) {
            size = props.maxSize
        }
        emit('update:size', size)
    }
}

const stopResize = () => {
    resizing.value = false
    document.removeEventListener('mousemove', handleResize)
    document.removeEventListener('mouseup', stopResize)
}

const startResize = () => {
    if (props.disabled) {
        return
    }
    resizing.value = true
    document.addEventListener('mousemove', handleResize)
    document.addEventListener('mouseup', stopResize)
}

const handleMouseOver = () => {
    if (props.disabled) {
        return
    }
    hover.value = true
}
</script>

<template>
    <div :style="{ width: props.size + 'px' }" class="resize-wrapper flex-box-h">
        <slot></slot>
        <div
            :class="{
                'resize-divider-hover': hover,
                'resize-divider-drag': resizing,
                dragging: hover || resizing,
            }"
            :style="{ width: props.borderWidth + 'px', right: Math.floor(-props.borderWidth / 2) + 'px' }"
            class="resize-divider"
            @mousedown="startResize"
            @mouseout="hover = false"
            @mouseover="handleMouseOver" />
    </div>
</template>

<style lang="scss" scoped>
.resize-wrapper {
    position: relative;

    .resize-divider {
        position: absolute;
        top: 0;
        bottom: 0;
        transition: background-color 0.3s ease-in;
    }

    .resize-divider-hide {
        background-color: #0000;
    }

    .resize-divider-hover {
        background-color: v-bind('themeVars.borderColor');
    }

    .resize-divider-drag {
        background-color: v-bind('themeVars.primaryColor');
    }

    .dragging {
        cursor: col-resize !important;
    }
}
</style>
