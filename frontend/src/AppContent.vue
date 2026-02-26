<script setup>
import ContentPane from './components/content/ContentPane.vue'
import BrowserPane from './components/sidebar/BrowserPane.vue'
import { computed, onMounted, onUnmounted, reactive, ref, watchEffect } from 'vue'
import { debounce } from 'lodash'
import { useThemeVars } from 'naive-ui'
import Ribbon from './components/sidebar/Ribbon.vue'
import ConnectionPane from './components/sidebar/ConnectionPane.vue'
import ContentServerPane from './components/content/ContentServerPane.vue'
import useTabStore from './stores/tab.js'
import usePreferencesStore from './stores/preferences.js'
import ContentLogPane from './components/content/ContentLogPane.vue'
import ContentValueTab from '@/components/content/ContentValueTab.vue'
import ToolbarControlWidget from '@/components/common/ToolbarControlWidget.vue'
import { EventsOn, WindowIsFullscreen, WindowIsMaximised, WindowToggleMaximise } from 'wailsjs/runtime/runtime.js'
import { isMacOS, isWeb, isWindows } from '@/utils/platform.js'
import iconUrl from '@/assets/images/icon.png'
import ResizeableWrapper from '@/components/common/ResizeableWrapper.vue'
import { extraTheme } from '@/utils/extra_theme.js'

const themeVars = useThemeVars()

const props = defineProps({
    loading: Boolean,
})

const data = reactive({
    navMenuWidth: 50,
    toolbarHeight: 38,
})

const tabStore = useTabStore()
const prefStore = usePreferencesStore()
const logPaneRef = ref(null)
const exThemeVars = computed(() => {
    return extraTheme(prefStore.isDark)
})
// const preferences = ref({})
// provide('preferences', preferences)

const saveSidebarWidth = debounce(prefStore.savePreferences, 1000, { trailing: true })
const handleResize = () => {
    saveSidebarWidth()
}

watchEffect(() => {
    if (tabStore.nav === 'log') {
        logPaneRef.value?.refresh()
    }
})

const logoWrapperWidth = computed(() => {
    return `${data.navMenuWidth + prefStore.behavior.asideWidth - 4}px`
})

const logoPaddingLeft = ref(10)
const maximised = ref(false)
const hideRadius = ref(false)
const wrapperStyle = computed(() => {
    if (isWindows() || isWeb()) {
        return {}
    }
    return hideRadius.value
        ? {}
        : {
              border: `1px solid ${themeVars.value.borderColor}`,
              borderRadius: '10px',
          }
})
const spinStyle = computed(() => {
    if (isWindows() || isWeb()) {
        return {
            backgroundColor: themeVars.value.bodyColor,
        }
    }
    return hideRadius.value
        ? {
              backgroundColor: themeVars.value.bodyColor,
          }
        : {
              backgroundColor: themeVars.value.bodyColor,
              borderRadius: '10px',
          }
})

const onToggleFullscreen = (fullscreen) => {
    hideRadius.value = fullscreen
    if (fullscreen) {
        logoPaddingLeft.value = 10
    } else {
        logoPaddingLeft.value = isMacOS() ? 70 : 10
    }
}

const onToggleMaximize = (isMaximised) => {
    if (isMaximised) {
        maximised.value = true
        if (!isMacOS()) {
            hideRadius.value = true
        }
    } else {
        maximised.value = false
        if (!isMacOS()) {
            hideRadius.value = false
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
    window.addEventListener('keydown', onKeyShortcut)
})

onUnmounted(() => {
    window.removeEventListener('keydown', onKeyShortcut)
})

const onKeyShortcut = (e) => {
    const isCtrlOn = isMacOS() ? e.metaKey : e.ctrlKey
    switch (e.key) {
        case 'w':
            if (isCtrlOn) {
                // close current tab
                const tabStore = useTabStore()
                const currentTab = tabStore.currentTab
                if (currentTab != null) {
                    tabStore.closeTab(currentTab.name)
                }
            }
            break
    }
}
</script>

<template>
    <!-- app content-->
    <n-spin :show="props.loading" :style="spinStyle" :theme-overrides="{ opacitySpinning: 0 }">
        <div id="app-content-wrapper" :style="wrapperStyle" class="flex-box-v">
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
                        width: logoWrapperWidth,
                        minWidth: logoWrapperWidth,
                        paddingLeft: `${logoPaddingLeft}px`,
                    }">
                    <n-space :size="3" :wrap="false" :wrap-item="false" align="center">
                        <n-avatar :size="32" :src="iconUrl" color="#0000" style="min-width: 32px" />
                        <div style="min-width: 68px; white-space: nowrap; font-weight: 800">Tiny RDM</div>
                        <transition name="fade">
                            <n-text v-if="tabStore.nav === 'browser'" class="ellipsis" strong style="font-size: 13px">
                                - {{ tabStore.currentTabName }}
                            </n-text>
                        </transition>
                    </n-space>
                </div>
                <!-- browser tabs -->
                <div v-show="tabStore.nav === 'browser'" class="app-toolbar-tab flex-item-expand">
                    <content-value-tab />
                </div>
                <div class="flex-item-expand" style="min-width: 15px"></div>
                <!-- simulate window control buttons -->
                <toolbar-control-widget
                    v-if="!isMacOS() && !isWeb()"
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
                <ribbon v-model:value="tabStore.nav" :width="data.navMenuWidth" />
                <!-- browser page -->
                <div v-show="tabStore.nav === 'browser'" class="content-area flex-box-h flex-item-expand">
                    <resizeable-wrapper
                        v-model:size="prefStore.behavior.asideWidth"
                        :min-size="300"
                        :offset="data.navMenuWidth"
                        class="flex-item"
                        @update:size="handleResize">
                        <browser-pane
                            v-for="t in tabStore.tabs"
                            v-show="tabStore.currentTabName === t.name"
                            :key="t.name"
                            :db="t.db"
                            :server="t.name"
                            class="app-side flex-item-expand" />
                    </resizeable-wrapper>
                    <content-pane
                        v-for="t in tabStore.tabs"
                        v-show="tabStore.currentTabName === t.name"
                        :key="t.name"
                        :server="t.name"
                        class="flex-item-expand" />
                </div>

                <!-- server list page -->
                <div v-show="tabStore.nav === 'server'" class="content-area flex-box-h flex-item-expand">
                    <resizeable-wrapper
                        v-model:size="prefStore.behavior.asideWidth"
                        :min-size="300"
                        :offset="data.navMenuWidth"
                        class="flex-item"
                        @update:size="handleResize">
                        <connection-pane class="app-side flex-item-expand" />
                    </resizeable-wrapper>
                    <content-server-pane class="flex-item-expand" />
                </div>

                <!-- log page -->
                <div v-show="tabStore.nav === 'log'" class="content-area flex-box-h flex-item-expand">
                    <content-log-pane ref="logPaneRef" class="flex-item-expand" />
                </div>
            </div>
        </div>
    </n-spin>
</template>

<style lang="scss" scoped>
#app-content-wrapper {
    width: 100vw;
    height: 100vh;
    height: 100dvh;
    overflow: hidden;
    box-sizing: border-box;
    background-color: v-bind('themeVars.bodyColor');
    color: v-bind('themeVars.textColorBase');

    #app-toolbar {
        background-color: v-bind('exThemeVars.titleColor');
        border-bottom: 1px solid v-bind('exThemeVars.splitColor');

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
        overflow: auto;
    }

    #app-content {
        height: calc(100% - 60px);

        .content-area {
            overflow: hidden;
        }
    }

    .app-side {
        //overflow: hidden;
        height: 100%;
        background-color: v-bind('exThemeVars.sidebarColor');
        border-right: 1px solid v-bind('exThemeVars.splitColor');
    }
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
