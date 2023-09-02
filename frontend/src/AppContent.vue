<script setup>
import ContentPane from './components/content/ContentPane.vue'
import BrowserPane from './components/sidebar/BrowserPane.vue'
import { computed, reactive, ref, watch } from 'vue'
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
import { WindowToggleMaximise } from 'wailsjs/runtime/runtime.js'
import { isMacOS } from '@/utils/platform.js'
import iconUrl from '@/assets/images/icon.png'

const themeVars = useThemeVars()

const data = reactive({
    initializing: false,
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

const saveWidth = debounce(prefStore.savePreferences, 1000, { trailing: true })
const handleResize = (evt) => {
    if (data.resizing) {
        prefStore.setAsideWidth(Math.max(evt.clientX - data.navMenuWidth, 300))
        saveWidth()
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
    return prefStore.general.asideWidth + 'px'
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
</script>

<template>
    <!-- app content-->
    <div id="app-content-wrapper" class="flex-box-v">
        <!-- title bar -->
        <div
            id="app-toolbar"
            class="flex-box-h"
            style="--wails-draggable: drag"
            :style="{ height: data.toolbarHeight + 'px' }"
            @dblclick="WindowToggleMaximise">
            <!-- title -->
            <div
                id="app-toolbar-title"
                :style="{
                    width: `${data.navMenuWidth + prefStore.general.asideWidth - 4}px`,
                    paddingLeft: isMacOS() ? '70px' : '10px',
                }">
                <n-space align="center" :wrap-item="false" :wrap="false" :size="3">
                    <n-avatar :src="iconUrl" color="#0000" :size="35" style="min-width: 35px" />
                    <div style="min-width: 68px; font-weight: 800">Tiny RDM</div>
                    <transition name="fade">
                        <n-text v-if="tabStore.nav === 'browser'" strong class="ellipsis" style="font-size: 13px">
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
            <toolbar-control-widget v-if="!isMacOS()" :size="data.toolbarHeight" style="align-self: flex-start" />
        </div>

        <!-- content -->
        <div id="app-content" :style="prefStore.generalFont" class="flex-box-h flex-item-expand">
            <nav-menu v-model:value="tabStore.nav" :width="data.navMenuWidth" />
            <!-- browser page-->
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
</template>

<style lang="scss" scoped>
#app-content-wrapper {
    height: 100%;
    overflow: hidden;
    box-sizing: border-box;

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
