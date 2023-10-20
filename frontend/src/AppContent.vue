<script setup>
import ContentPane from './components/content/ContentPane.vue'
import BrowserPane from './components/sidebar/BrowserPane.vue'
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { debounce, get } from 'lodash'
import { useThemeVars } from 'naive-ui'
import NavMenu from './components/sidebar/NavMenu.vue'
import ConnectionPane from './components/sidebar/ConnectionPane.vue'
import ContentServerPane from './components/content/ContentServerPane.vue'
import useTabStore from './stores/tab.js'
import usePreferencesStore from './stores/preferences.js'
import useConnectionStore from './stores/connections.js'
import ContentLogPane from './components/content/ContentLogPane.vue'
import ContentValueTab from '@/components/content/ContentValueTab.vue'
import ToolbarControlWidget from '@/components/common/ToolbarControlWidget.vue'
import { EventsOn, WindowIsFullscreen, WindowIsMaximised, WindowToggleMaximise } from 'wailsjs/runtime/runtime.js'
import { isMacOS } from '@/utils/platform.js'
import iconUrl from '@/assets/images/icon.png'

const themeVars = useThemeVars()

const props = defineProps({
    loading: Boolean,
})

const data = reactive({
    navMenuWidth: 60,
    hoverResize: false,
    resizing: false,
    toolbarHeight: 45,
})

const tabStore = useTabStore()
const prefStore = usePreferencesStore()
const connectionStore = useConnectionStore()
const logPaneRef = ref(null)
// const preferences = ref({})
// provide('preferences', preferences)

const saveSidebarWidth = debounce(prefStore.savePreferences, 1000, { trailing: true })
const handleResize = (evt) => {
    if (data.resizing) {
        prefStore.setAsideWidth(Math.max(evt.clientX - data.navMenuWidth, 300))
        saveSidebarWidth()
    }
}

const stopResize = () => {
    data.resizing = false
    document.removeEventListener('mousemove', handleResize)
    document.removeEventListener('mouseup', stopResize)
}

const startResize = () => {
    data.resizing = true
    document.addEventListener('mousemove', handleResize)
    document.addEventListener('mouseup', stopResize)
}

const asideWidthVal = computed(() => {
    return prefStore.behavior.asideWidth + 'px'
})

const dragging = computed(() => {
    return data.hoverResize || data.resizing
})

watch(
    () => tabStore.nav,
    (nav) => {
        if (nav === 'log') {
            logPaneRef.value?.refresh()
        }
    },
)

const border = computed(() => {
    const color = isMacOS() && false ? '#0000' : themeVars.value.borderColor
    return `1px solid ${color}`
})

const borderRadius = ref(10)
const logoPaddingLeft = ref(10)
const maximised = ref(false)
const toggleWindowRadius = (on) => {
    borderRadius.value = on ? 10 : 0
}

const onToggleFullscreen = (fullscreen) => {
    if (fullscreen) {
        logoPaddingLeft.value = 10
        toggleWindowRadius(false)
    } else {
        logoPaddingLeft.value = isMacOS() ? 70 : 10
        toggleWindowRadius(true)
    }
}

const onToggleMaximize = (isMaximised) => {
    if (isMaximised) {
        maximised.value = true
        if (!isMacOS()) {
            toggleWindowRadius(false)
        }
    } else {
        maximised.value = false
        if (!isMacOS()) {
            toggleWindowRadius(true)
        }
    }
}

EventsOn('window_changed', (info) => {
    const { fullscreen, maximised } = info
    onToggleFullscreen(fullscreen === true)
    onToggleMaximize(maximised)
})

onMounted(async () => {
    const fullscreen = await WindowIsFullscreen()
    onToggleFullscreen(fullscreen === true)
    const maximised = await WindowIsMaximised()
    onToggleMaximize(maximised)
})
</script>

<template>
    <!-- app content-->
    <n-spin
        :show="props.loading"
        :style="{ backgroundColor: themeVars.bodyColor, borderRadius: `${borderRadius}px`, border }"
        :theme-overrides="{ opacitySpinning: 0 }">
        <div
            id="app-content-wrapper"
            :style="{
                backgroundColor: themeVars.bodyColor,
                color: themeVars.textColorBase,
            }"
            class="flex-box-v">
            <!-- title bar -->
            <div
                id="app-toolbar"
                :style="{ height: data.toolbarHeight + 'px' }"
                class="flex-box-h"
                style="--wails-draggable: drag"
                @dblclick="WindowToggleMaximise">
                <!-- title -->
                <div
                    id="app-toolbar-title"
                    :style="{
                        width: `${data.navMenuWidth + prefStore.behavior.asideWidth - 4}px`,
                        paddingLeft: `${logoPaddingLeft}px`,
                    }">
                    <n-space :size="3" :wrap="false" :wrap-item="false" align="center">
                        <n-avatar :size="35" :src="iconUrl" color="#0000" style="min-width: 35px" />
                        <div style="min-width: 68px; font-weight: 800">Tiny RDM</div>
                        <transition name="fade">
                            <n-text v-if="tabStore.nav === 'browser'" class="ellipsis" strong style="font-size: 13px">
                                - {{ get(tabStore.currentTab, 'name') }}
                            </n-text>
                        </transition>
                    </n-space>
                </div>
                <div
                    :class="{
                        'resize-divider-hover': data.hoverResize,
                        'resize-divider-drag': data.resizing,
                    }"
                    class="resize-divider resize-divider-hide"
                    @mousedown="startResize"
                    @mouseout="data.hoverResize = false"
                    @mouseover="data.hoverResize = true" />
                <!-- browser tabs -->
                <div v-show="tabStore.nav === 'browser'" class="app-toolbar-tab flex-item-expand">
                    <content-value-tab />
                </div>
                <div class="flex-item-expand"></div>
                <!-- simulate window control buttons -->
                <toolbar-control-widget
                    v-if="!isMacOS()"
                    :maximised="maximised"
                    :size="data.toolbarHeight"
                    style="align-self: flex-start" />
            </div>

            <!-- content -->
            <div
                id="app-content"
                :style="prefStore.generalFont"
                class="flex-box-h flex-item-expand"
                style="--wails-draggable: none">
                <nav-menu v-model:value="tabStore.nav" :width="data.navMenuWidth" />
                <!-- browser page -->
                <div v-show="tabStore.nav === 'browser'" :class="{ dragging }" class="flex-box-h flex-item-expand">
                    <div id="app-side" :style="{ width: asideWidthVal }" class="flex-box-h flex-item">
                        <browser-pane
                            v-for="t in tabStore.tabs"
                            v-show="get(tabStore.currentTab, 'name') === t.name"
                            :key="t.name"
                            class="flex-item-expand" />
                        <div
                            :class="{
                                'resize-divider-hover': data.hoverResize,
                                'resize-divider-drag': data.resizing,
                            }"
                            class="resize-divider"
                            @mousedown="startResize"
                            @mouseout="data.hoverResize = false"
                            @mouseover="data.hoverResize = true" />
                    </div>
                    <content-pane class="flex-item-expand" />
                </div>

                <!-- server list page -->
                <div v-show="tabStore.nav === 'server'" :class="{ dragging }" class="flex-box-h flex-item-expand">
                    <div id="app-side" :style="{ width: asideWidthVal }" class="flex-box-h flex-item">
                        <connection-pane class="flex-item-expand" />
                        <div
                            :class="{
                                'resize-divider-hover': data.hoverResize,
                                'resize-divider-drag': data.resizing,
                            }"
                            class="resize-divider"
                            @mousedown="startResize"
                            @mouseout="data.hoverResize = false"
                            @mouseover="data.hoverResize = true" />
                    </div>
                    <content-server-pane class="flex-item-expand" />
                </div>

                <!-- log page -->
                <div v-show="tabStore.nav === 'log'" class="flex-box-h flex-item-expand">
                    <content-log-pane ref="logPaneRef" class="flex-item-expand" />
                </div>
            </div>
        </div>
    </n-spin>
</template>

<style lang="scss" scoped>
#app-content-wrapper {
    width: calc(100vw - 2px);
    height: calc(100vh - 2px);
    overflow: hidden;
    box-sizing: border-box;
    border-radius: 10px;

    #app-toolbar {
        background-color: v-bind('themeVars.tabColor');
        border-bottom: 1px solid v-bind('themeVars.borderColor');

        &-title {
            padding-left: 10px;
            padding-right: 10px;
            box-sizing: border-box;
            align-self: center;
            align-items: baseline;
        }
    }

    .app-toolbar-tab {
        align-self: flex-end;
        margin-bottom: -1px;
        margin-left: 3px;
    }

    #app-content {
        height: calc(100% - 60px);
    }

    #app-side {
        //overflow: hidden;
        height: 100%;
        background-color: v-bind('themeVars.tabColor');
    }
}

.resize-divider {
    width: 3px;
    border-right: 1px solid v-bind('themeVars.borderColor');
}

.resize-divider-hide {
    background-color: #0000;
    border-right-color: #0000;
}

.resize-divider-hover {
    background-color: v-bind('themeVars.borderColor');
    border-right-color: v-bind('themeVars.borderColor');
}

.resize-divider-drag {
    background-color: v-bind('themeVars.primaryColor');
    border-right-color: v-bind('themeVars.primaryColor');
}

.dragging {
    cursor: col-resize !important;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
}

.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.3s ease;
}
</style>
