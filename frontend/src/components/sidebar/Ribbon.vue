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
import usePreferencesStore from 'stores/preferences.js'
import Record from '@/components/icons/Record.vue'
import { extraTheme } from '@/utils/extra_theme.js'
import useBrowserStore from 'stores/browser.js'

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

const iconSize = computed(() => Math.floor(props.width * 0.45))
const renderIcon = (icon) => {
    return () => h(NIcon, null, { default: () => h(icon, { strokeWidth: 3 }) })
}

const browserStore = useBrowserStore()
const i18n = useI18n()
const menuOptions = computed(() => {
    return [
        {
            label: i18n.t('ribbon.browser'),
            key: 'browser',
            icon: Database,
            show: browserStore.anyConnectionOpened,
        },
        {
            label: i18n.t('ribbon.server'),
            key: 'server',
            icon: Server,
        },
        {
            label: i18n.t('ribbon.log'),
            key: 'log',
            icon: Record,
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
        id="app-ribbon"
        :style="{
            width: props.width + 'px',
            minWidth: props.width + 'px',
        }"
        class="flex-box-v">
        <div class="ribbon-wrapper flex-box-v">
            <div
                v-for="(m, i) in menuOptions"
                v-show="m.show !== false"
                :key="i"
                :class="{ 'ribbon-item-active': props.value === m.key }"
                class="ribbon-item clickable"
                @click="emit('update:value', m.key)">
                <n-tooltip :delay="2" :show-arrow="false" placement="right">
                    <template #trigger>
                        <n-icon :size="iconSize">
                            <component :is="m.icon" :stroke-width="3.5" />
                        </n-icon>
                    </template>
                    {{ m.label }}
                </n-tooltip>
            </div>
        </div>
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
#app-ribbon {
    //height: 100vh;
    border-right: v-bind('exThemeVars.splitColor') solid 1px;
    background-color: v-bind('exThemeVars.ribbonColor');
    box-sizing: border-box;
    color: v-bind('themeVars.textColor2');

    .ribbon-wrapper {
        gap: 2px;
        margin-top: 5px;
        justify-content: center;
        align-items: center;
        box-sizing: border-box;
        padding-right: 3px;

        .ribbon-item {
            width: 100%;
            height: 100%;
            text-align: center;
            line-height: 1;
            color: v-bind('themeVars.textColor3');
            //border-left: 5px solid #000;
            border-radius: v-bind('themeVars.borderRadius');
            padding: 8px 0;
            position: relative;

            &:hover {
                background-color: rgba(0, 0, 0, 0.05);
                color: v-bind('themeVars.primaryColor');

                &:before {
                    position: absolute;
                    width: 3px;
                    left: 0;
                    top: 24%;
                    bottom: 24%;
                    border-radius: 9999px;
                    content: '';
                    background-color: v-bind('themeVars.primaryColor');
                }
            }
        }

        .ribbon-item-active {
            //background-color: v-bind('exThemeVars.ribbonActiveColor');
            color: v-bind('themeVars.primaryColor');

            &:hover {
                color: v-bind('themeVars.primaryColor') !important;
            }

            &:before {
                position: absolute;
                width: 3px;
                left: 0;
                top: 24%;
                bottom: 24%;
                border-radius: 9999px;
                content: '';
                background-color: v-bind('themeVars.primaryColor');
            }
        }
    }

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
