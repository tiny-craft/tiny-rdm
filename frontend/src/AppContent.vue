<script setup>
import ContentPane from './components/ContentPane.vue'
import NavigationPane from './components/NavigationPane.vue'
import { computed, nextTick, onMounted, provide, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { GetPreferences } from '../wailsjs/go/storage/PreferencesStorage.js'
import { get } from 'lodash'
import { useThemeVars } from 'naive-ui'
import NavMenu from './components/NavMenu.vue'

const themeVars = useThemeVars()

const data = reactive({
    asideWith: 300,
    hoverResize: false,
    resizing: false,
})

const preferences = ref({})
provide('preferences', preferences)
const i18n = useI18n()

onMounted(async () => {
    preferences.value = await GetPreferences()
    await nextTick(() => {
        i18n.locale.value = get(preferences.value, 'general.language', 'en')
    })
})

// TODO: apply font size to all elements
const getFontSize = computed(() => {
    return get(preferences.value, 'general.font_size', 'en')
})

const handleResize = (evt) => {
    if (data.resizing) {
        data.asideWith = Math.max(evt.clientX, 300)
    }
}

const stopResize = () => {
    data.resizing = false
    document.removeEventListener('mousemove', handleResize)
    document.removeEventListener('mouseup', stopResize)
    // TODO: Save sidebar x-position
}

const startResize = () => {
    data.resizing = true
    document.addEventListener('mousemove', handleResize)
    document.addEventListener('mouseup', stopResize)
}

const asideWidthVal = computed(() => {
    return data.asideWith + 'px'
})

const dragging = computed(() => {
    return data.hoverResize || data.resizing
})
</script>

<template>
    <!-- app content-->
    <div id="app-container" :class="{ dragging: dragging }" class="flex-box-h">
        <nav-menu />
        <div id="app-side" :style="{ width: asideWidthVal }" class="flex-box-h flex-item">
            <navigation-pane class="flex-item-expand"></navigation-pane>
            <div
                :class="{
                    'resize-divider-hover': data.hoverResize,
                    'resize-divider-drag': data.resizing,
                }"
                class="resize-divider"
                @mousedown="startResize"
                @mouseout="data.hoverResize = false"
                @mouseover="data.hoverResize = true"
            ></div>
        </div>
        <content-pane class="flex-item-expand" />
    </div>
</template>

<style lang="scss">
#app-container {
    height: 100%;
    overflow: hidden;
    border-top: var(--border-color) 1px solid;
    box-sizing: border-box;

    #app-toolbar {
        height: 40px;
        border-bottom: var(--border-color) 1px solid;
    }

    #app-side {
        //overflow: hidden;
        height: 100%;

        .resize-divider {
            //height: 100%;
            width: 2px;
            border-left-width: 5px;
            background-color: var(--border-color);
        }

        .resize-divider-hover {
            width: 5px;
        }

        .resize-divider-drag {
            //background-color: rgb(0, 105, 218);
            width: 5px;
            //background-color: var(--el-color-primary);
            background-color: v-bind('themeVars.primaryColor');
        }
    }
}

.dragging {
    cursor: col-resize !important;
}
</style>
