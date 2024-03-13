<script setup>
import { computed, ref } from 'vue'
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
import { useRender } from '@/utils/render.js'
import wechatUrl from '@/assets/images/wechat_official.png'
import QRCode from '@/components/icons/QRCode.vue'
import Twitter from '@/components/icons/Twitter.vue'
import { trackEvent } from '@/utils/analytics.js'

const themeVars = useThemeVars()
const render = useRender()

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

const browserStore = useBrowserStore()
const i18n = useI18n()
const showWechat = ref(false)
const menuOptions = computed(() => {
    return [
        {
            label: 'ribbon.browser',
            key: 'browser',
            icon: Database,
            show: browserStore.anyConnectionOpened,
        },
        {
            label: 'ribbon.server',
            key: 'server',
            icon: Server,
        },
        {
            label: 'ribbon.log',
            key: 'log',
            icon: Record,
        },
    ]
})

const preferencesOptions = computed(() => {
    return [
        {
            label: 'menu.preferences',
            key: 'preferences',
            icon: Config,
        },
        // {
        //     label: 'menu.help',
        //     key: 'help',
        //     icon: Help,
        // },
        {
            label: 'menu.report_bug',
            key: 'report',
        },
        {
            label: 'menu.user_guide',
            key: 'help',
        },
        {
            label: 'menu.check_update',
            key: 'update',
        },
        {
            type: 'divider',
            key: 'd1',
        },
        {
            label: 'menu.about',
            key: 'about',
        },
    ]
})

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
        case 'report':
            BrowserOpenURL('https://github.com/tiny-craft/tiny-rdm/issues')
            break
        case 'help':
            if (prefStore.currentLanguage === 'zh') {
                BrowserOpenURL('https://redis.tinycraft.cc/zh/guide/')
            } else {
                BrowserOpenURL('https://redis.tinycraft.cc/guide/')
            }
            break
        case 'about':
            dialogStore.openAboutDialog()
            break
    }
}

const openWechatOfficial = () => {
    trackEvent('open', { target: 'wechat_official' })
    showWechat.value = true
}

const openX = () => {
    trackEvent('open', { target: 'x' })
    BrowserOpenURL('https://twitter.com/LykinHuang')
}

const openGithub = () => {
    trackEvent('open', { target: 'github' })
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
            <n-tooltip v-for="(m, i) in menuOptions" :key="i" :delay="2" :show-arrow="false" placement="right">
                <template #trigger>
                    <div
                        v-show="m.show !== false"
                        :class="{ 'ribbon-item-active': props.value === m.key }"
                        class="ribbon-item clickable"
                        @click="emit('update:value', m.key)">
                        <n-icon :size="iconSize">
                            <component :is="m.icon" :stroke-width="3.5" />
                        </n-icon>
                    </div>
                </template>
                {{ $t(m.label) }}
            </n-tooltip>
        </div>
        <div class="flex-item-expand"></div>
        <div class="nav-menu-item flex-box-v">
            <n-dropdown
                :options="preferencesOptions"
                :render-icon="({ icon }) => render.renderIcon(icon)"
                :render-label="({ label }) => render.renderLabel($t(label), { class: 'context-menu-item' })"
                content-class="nav-menu-button"
                trigger="click"
                @select="onSelectPreferenceMenu">
                <icon-button :icon="Config" :size="iconSize" :stroke-width="3" />
            </n-dropdown>
            <icon-button
                v-if="prefStore.currentLanguage === 'zh'"
                :icon="QRCode"
                :size="iconSize"
                :tooltip-delay="100"
                class="nav-menu-button"
                t-tooltip="ribbon.wechat_official"
                @click="openWechatOfficial" />
            <icon-button
                v-else
                :border="false"
                :icon="Twitter"
                :size="iconSize"
                :tooltip-delay="100"
                class="nav-menu-button"
                t-tooltip="ribbon.follow_x"
                @click="openX" />
            <icon-button
                :icon="Github"
                :size="iconSize"
                :tooltip-delay="100"
                class="nav-menu-button"
                t-tooltip="ribbon.github"
                @click="openGithub" />
        </div>

        <!-- wechat official modal -->
        <n-modal
            :show="showWechat"
            transform-origin="center"
            @close="showWechat = false"
            @mask-click="showWechat = false">
            <n-image :src="wechatUrl" :width="400" preview-disabled />
        </n-modal>
    </div>
</template>

<style lang="scss">
#app-ribbon {
    //height: 100vh;
    border-right: v-bind('exThemeVars.splitColor') solid 1px;
    background-color: v-bind('exThemeVars.ribbonColor');
    box-sizing: border-box;
    color: v-bind('themeVars.textColor2');
    --wails-draggable: drag;

    .ribbon-wrapper {
        gap: 2px;
        margin-top: 5px;
        justify-content: center;
        align-items: center;
        box-sizing: border-box;
        padding-right: 3px;
        --wails-draggable: none;

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
        padding: 10px 0 15px;
        gap: 20px;
        --wails-draggable: none;

        .nav-menu-button {
            margin-bottom: 6px;

            :hover {
                color: v-bind('themeVars.primaryColor');
            }
        }
    }
}
</style>
