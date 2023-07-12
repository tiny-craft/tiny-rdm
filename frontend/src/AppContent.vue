<script setup>
import ContentPane from './components/content/ContentPane.vue'
import BrowserPane from './components/sidebar/BrowserPane.vue'
import { computed, nextTick, onBeforeMount, onMounted, reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { debounce, get } from 'lodash'
import { useThemeVars } from 'naive-ui'
import NavMenu from './components/sidebar/NavMenu.vue'
import ConnectionPane from './components/sidebar/ConnectionPane.vue'
import ContentServerPane from './components/content/ContentServerPane.vue'
import useTabStore from './stores/tab.js'
import usePreferencesStore from './stores/preferences.js'
import useConnectionStore from './stores/connections.js'

const themeVars = useThemeVars()

const data = reactive({
    initializing: false,
    navMenuWidth: 60,
    hoverResize: false,
    resizing: false,
})

const tabStore = useTabStore()
const prefStore = usePreferencesStore()
const connectionStore = useConnectionStore()
// const preferences = ref({})
// provide('preferences', preferences)
const i18n = useI18n()

onBeforeMount(async () => {
    try {
        data.initializing = true
        await prefStore.loadPreferences()
        i18n.locale.value = get(prefStore.general, 'language', 'en')
        await prefStore.loadFontList()
        await connectionStore.initConnections()
    } finally {
        data.initializing = false
    }
})

const saveWidth = debounce(prefStore.savePreferences, 1000, { trailing: true })
const handleResize = (evt) => {
    if (data.resizing) {
        prefStore.setNavWidth(Math.max(evt.clientX - data.navMenuWidth, 300))
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
    return prefStore.general.navMenuWidth + 'px'
})

const dragging = computed(() => {
    return data.hoverResize || data.resizing
})
</script>

<template>
    <!-- app content-->
    <!--    <div id="app-container"></div>-->
    <n-spin :show="data.initializing" :theme-overrides="{ opacitySpinning: 0 }">
        <template #description> {{ $t('launching') }} </template>
        <div id="app-container" :class="{ dragging }" class="flex-box-h" :style="prefStore.generalFont">
            <nav-menu v-model:value="tabStore.nav" :width="data.navMenuWidth" />
            <!-- structure page-->
            <div v-show="tabStore.nav === 'structure'" class="flex-box-h flex-item-expand">
                <div id="app-side" :style="{ width: asideWidthVal }" class="flex-box-h flex-item">
                    <browser-pane
                        v-for="t in tabStore.tabs"
                        v-show="get(tabStore.currentTab, 'name') === t.name"
                        :key="t.name"
                        class="flex-item-expand"
                    />
                    <div
                        :class="{
                            'resize-divider-hover': data.hoverResize,
                            'resize-divider-drag': data.resizing,
                        }"
                        class="resize-divider"
                        @mousedown="startResize"
                        @mouseout="data.hoverResize = false"
                        @mouseover="data.hoverResize = true"
                    />
                </div>
                <content-pane class="flex-item-expand" />
            </div>

            <!-- server list page -->
            <div v-show="tabStore.nav === 'server'" class="flex-box-h flex-item-expand">
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
                        @mouseover="data.hoverResize = true"
                    />
                </div>
                <content-server-pane class="flex-item-expand" />
            </div>

            <!-- log page -->
            <div v-show="tabStore.nav === 'log'">display log</div>
        </div>
    </n-spin>
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
