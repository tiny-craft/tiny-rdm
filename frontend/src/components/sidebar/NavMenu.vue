<script setup>
import { computed, h } from 'vue'
import { NIcon, useThemeVars } from 'naive-ui'
import Database from '@/components/icons/Database.vue'
import { useI18n } from 'vue-i18n'
import Server from '@/components/icons/Server.vue'
import IconButton from '@/components/common/IconButton.vue'
import Config from '@/components/icons/Config.vue'
import useDialogStore from 'stores/dialog.js'
import Github from '@/components/icons/Github.vue'
import { BrowserOpenURL } from 'wailsjs/runtime/runtime.js'
import useConnectionStore from 'stores/connections.js'
import usePreferencesStore from 'stores/preferences.js'
import Record from '@/components/icons/Record.vue'
import { extraTheme } from '@/utils/extra_theme.js'

const themeVars = useThemeVars()

const props = defineProps({
    value: {
        type: String,
        default: 'server',
    },
    width: {
        type: Number,
        default: 60,
    },
})

const emit = defineEmits(['update:value'])

const iconSize = computed(() => Math.floor(props.width * 0.4))
const renderIcon = (icon) => {
    return () => h(NIcon, null, { default: () => h(icon, { strokeWidth: 3 }) })
}

const connectionStore = useConnectionStore()
const i18n = useI18n()
const menuOptions = computed(() => {
    return [
        {
            label: i18n.t('ribbon.browser'),
            key: 'browser',
            icon: renderIcon(Database),
            show: connectionStore.anyConnectionOpened,
        },
        {
            label: i18n.t('ribbon.server'),
            key: 'server',
            icon: renderIcon(Server),
        },
        {
            label: i18n.t('ribbon.log'),
            key: 'log',
            icon: renderIcon(Record),
        },
    ]
})

const preferencesOptions = computed(() => {
    return [
        {
            label: i18n.t('menu.preferences'),
            key: 'preferences',
            icon: renderIcon(Config),
        },
        // {
        //     label: i18n.t('menu.help'),
        //     key: 'help',
        //     icon: renderIcon(Help),
        // },
        {
            label: i18n.t('menu.check_update'),
            key: 'update',
        },
        {
            type: 'divider',
            key: 'd1',
        },
        {
            label: i18n.t('menu.about'),
            key: 'about',
        },
    ]
})

const renderContextLabel = (option) => {
    return h('div', { class: 'context-menu-item' }, option.label)
}

const dialogStore = useDialogStore()
const prefStore = usePreferencesStore()
const onSelectPreferenceMenu = (key) => {
    switch (key) {
        case 'preferences':
            dialogStore.openPreferencesDialog()
            break
        case 'update':
            prefStore.checkForUpdate(true)
            break
        case 'about':
            dialogStore.openAboutDialog()
            break
    }
}

const openGithub = () => {
    BrowserOpenURL('https://github.com/tiny-craft/tiny-rdm')
}

const exThemeVars = computed(() => {
    return extraTheme(prefStore.isDark)
})
</script>

<template>
    <div
        id="app-nav-menu"
        :style="{
            width: props.width + 'px',
        }"
        class="flex-box-v">
        <n-menu
            :collapsed="true"
            :collapsed-icon-size="iconSize"
            :collapsed-width="props.width"
            :options="menuOptions"
            :value="props.value"
            @update:value="(val) => emit('update:value', val)" />
        <div class="flex-item-expand"></div>
        <div class="nav-menu-item flex-box-v">
            <n-dropdown
                :options="preferencesOptions"
                :render-label="renderContextLabel"
                trigger="click"
                @select="onSelectPreferenceMenu">
                <icon-button :icon="Config" :size="iconSize" :stroke-width="3" class="nav-menu-button" />
            </n-dropdown>
            <icon-button :icon="Github" :size="iconSize" class="nav-menu-button" @click="openGithub" />
        </div>
    </div>
</template>

<style lang="scss">
#app-nav-menu {
    //height: 100vh;
    border-right: v-bind('exThemeVars.splitColor') solid 1px;
    background-color: v-bind('exThemeVars.sidebarColor');
    box-sizing: border-box;

    .nav-menu-item {
        align-items: center;
        padding: 10px 0;
        gap: 15px;

        .nav-menu-button {
            margin-bottom: 6px;

            :hover {
                color: v-bind('themeVars.primaryColor');
            }
        }
    }
}
</style>
