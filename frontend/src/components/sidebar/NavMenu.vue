<script setup>
import { computed, h } from 'vue'
import { NIcon, useThemeVars } from 'naive-ui'
import ToggleDb from '../icons/ToggleDb.vue'
import { useI18n } from 'vue-i18n'
import ToggleServer from '../icons/ToggleServer.vue'
import IconButton from '../common/IconButton.vue'
import Config from '../icons/Config.vue'
import useDialogStore from '../../stores/dialog.js'
import Github from '../icons/Github.vue'
import { BrowserOpenURL } from '../../../wailsjs/runtime/runtime.js'
import Log from '../icons/Log.vue'
import useConnectionStore from '../../stores/connections.js'
import Help from '../icons/Help.vue'

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
    return () => h(NIcon, null, { default: () => h(icon) })
}

const connectionStore = useConnectionStore()
const i18n = useI18n()
const menuOptions = computed(() => {
    return [
        {
            label: i18n.t('browser'),
            key: 'browser',
            icon: renderIcon(ToggleDb),
            show: connectionStore.anyConnectionOpened,
        },
        {
            label: i18n.t('server'),
            key: 'server',
            icon: renderIcon(ToggleServer),
        },
        {
            label: i18n.t('log'),
            key: 'log',
            icon: renderIcon(Log),
        },
    ]
})

const preferencesOptions = computed(() => {
    return [
        {
            label: i18n.t('preferences'),
            key: 'preferences',
            icon: renderIcon(Config),
        },
        {
            label: i18n.t('help'),
            key: 'help',
            icon: renderIcon(Help),
        },
        {
            label: i18n.t('about'),
            key: 'about',
        },
        {
            label: i18n.t('check_update'),
            key: 'update',
        },
    ]
})

const renderContextLabel = (option) => {
    return h('div', { class: 'context-menu-item' }, option.label)
}

const dialogStore = useDialogStore()
const onSelectPreferenceMenu = (key) => {
    switch (key) {
        case 'preferences':
            dialogStore.openPreferencesDialog()
            break
        case 'update':
            break
    }
}

const openGithub = () => {
    BrowserOpenURL('https://github.com/tiny-craft/tiny-rdm')
}
</script>

<template>
    <div
        id="app-nav-menu"
        :style="{
            width: props.width + 'px',
        }"
        class="flex-box-v"
    >
        <n-menu
            :collapsed-width="props.width"
            :value="props.value"
            :collapsed="true"
            :collapsed-icon-size="iconSize"
            @update:value="(val) => emit('update:value', val)"
            :options="menuOptions"
        />
        <div class="flex-item-expand"></div>
        <div class="nav-menu-item flex-box-v">
            <n-dropdown
                :animated="false"
                :keyboard="false"
                :options="preferencesOptions"
                :render-label="renderContextLabel"
                trigger="click"
                @select="onSelectPreferenceMenu"
            >
                <icon-button :icon="Config" :size="iconSize" class="nav-menu-button" />
            </n-dropdown>
            <icon-button :icon="Github" :size="iconSize" class="nav-menu-button" @click="openGithub" />
        </div>
    </div>
</template>

<style lang="scss">
#app-nav-menu {
    height: 100vh;
    //border-right: v-bind('themeVars.borderColor') solid 1px;

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
