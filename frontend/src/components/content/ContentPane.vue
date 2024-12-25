<script setup>
import { computed, nextTick, ref, watch } from 'vue'
import { find, map, toUpper } from 'lodash'
import useTabStore from 'stores/tab.js'
import ContentServerStatus from '@/components/content_value/ContentServerStatus.vue'
import Status from '@/components/icons/Status.vue'
import { useThemeVars } from 'naive-ui'
import { BrowserTabType } from '@/consts/browser_tab_type.js'
import Terminal from '@/components/icons/Terminal.vue'
import Log from '@/components/icons/Log.vue'
import Detail from '@/components/icons/Detail.vue'
import ContentValueWrapper from '@/components/content_value/ContentValueWrapper.vue'
import ContentCli from '@/components/content_value/ContentCli.vue'
import Monitor from '@/components/icons/Monitor.vue'
import ContentSlog from '@/components/content_value/ContentSlog.vue'
import ContentMonitor from '@/components/content_value/ContentMonitor.vue'
import { decodeRedisKey } from '@/utils/key_convert.js'
import ContentPubsub from '@/components/content_value/ContentPubsub.vue'
import Subscribe from '@/components/icons/Subscribe.vue'

const themeVars = useThemeVars()

/**
 * @typedef {Object} ServerStatusItem
 * @property {string} name
 * @property {Object} info
 * @property {boolean} autoRefresh
 * @property {boolean} loading loading status for refresh
 * @property {boolean} autoLoading loading status for auto refresh
 */

const props = defineProps({
    server: String,
})

const tabStore = useTabStore()
const tab = computed(() =>
    map(tabStore.tabs, (item) => ({
        key: item.name,
        label: item.title,
    })),
)

const tabContent = computed(() => {
    const tab = find(tabStore.tabs, { name: props.server })
    if (tab == null) {
        return {}
    }
    return {
        name: tab.name,
        subTab: tab.subTab,
        type: toUpper(tab.type),
        db: tab.db,
        keyPath: tab.keyCode != null ? decodeRedisKey(tab.keyCode) : tab.key,
        keyCode: tab.keyCode,
        ttl: tab.ttl,
        value: tab.value,
        size: tab.size || 0,
        length: tab.length || 0,
        decode: tab.decode,
        format: tab.format,
        matchPattern: tab.matchPattern || '',
        end: tab.end === true,
        loading: tab.loading === true,
    }
})

const isBlankValue = computed(() => {
    return tabContent.value?.keyPath == null
})

const selectedSubTab = computed(() => {
    const { subTab = BrowserTabType.Status } = tabStore.currentTab || {}
    return subTab
})

// BUG: naive-ui tabs will set the bottom line to '0px' after switch to another page and back again
// watch parent tabs' changing and call 'syncBarPosition' manually
const tabsRef = ref(null)
const cliRef = ref(null)
watch(
    () => tabContent.value?.name,
    (name) => {
        if (name === props.server) {
            nextTick().then(() => {
                tabsRef.value?.syncBarPosition()
                cliRef.value?.resizeTerm()
            })
        }
    },
)
</script>

<template>
    <div class="content-container flex-box-v">
        <n-tabs
            ref="tabsRef"
            :tabs-padding="5"
            :theme-overrides="{
                tabFontWeightActive: 'normal',
                tabGapSmallLine: '10px',
                tabGapMediumLine: '10px',
                tabGapLargeLine: '10px',
            }"
            :value="selectedSubTab"
            class="content-sub-tab"
            :default-value="BrowserTabType.Status.toString()"
            pane-class="content-sub-tab-pane"
            placement="top"
            tab-style="padding-left: 10px; padding-right: 10px;"
            type="line"
            @update:value="tabStore.switchSubTab">
            <!-- server status pane -->
            <n-tab-pane :name="BrowserTabType.Status.toString()" display-directive="show:lazy">
                <template #tab>
                    <n-space :size="5" :wrap-item="false" align="center" inline justify="center">
                        <n-icon size="16">
                            <status
                                :inverse="selectedSubTab === BrowserTabType.Status.toString()"
                                :stroke-color="themeVars.tabColor"
                                stroke-width="4" />
                        </n-icon>
                        <span>{{ $t('interface.sub_tab.status') }}</span>
                    </n-space>
                </template>
                <content-server-status
                    :pause="selectedSubTab !== BrowserTabType.Status.toString()"
                    :server="props.server" />
            </n-tab-pane>

            <!-- key detail pane -->
            <n-tab-pane :name="BrowserTabType.KeyDetail.toString()" display-directive="show:lazy">
                <template #tab>
                    <n-space :size="5" :wrap-item="false" align="center" inline justify="center">
                        <n-icon size="16">
                            <detail
                                :inverse="selectedSubTab === BrowserTabType.KeyDetail.toString()"
                                :stroke-color="themeVars.tabColor"
                                stroke-width="4" />
                        </n-icon>
                        <span>{{ $t('interface.sub_tab.key_detail') }}</span>
                    </n-space>
                </template>
                <content-value-wrapper :blank="isBlankValue" :content="tabContent" />
            </n-tab-pane>

            <!-- cli pane -->
            <n-tab-pane :name="BrowserTabType.Cli.toString()" display-directive="show:lazy">
                <template #tab>
                    <n-space :size="5" :wrap-item="false" align="center" inline justify="center">
                        <n-icon size="16">
                            <terminal
                                :inverse="selectedSubTab === BrowserTabType.Cli.toString()"
                                :stroke-color="themeVars.tabColor"
                                stroke-width="4" />
                        </n-icon>
                        <span>{{ $t('interface.sub_tab.cli') }}</span>
                    </n-space>
                </template>
                <content-cli ref="cliRef" :name="props.server" />
            </n-tab-pane>

            <!-- slow log pane -->
            <n-tab-pane :name="BrowserTabType.SlowLog.toString()" display-directive="show:lazy">
                <template #tab>
                    <n-space :size="5" :wrap-item="false" align="center" inline justify="center">
                        <n-icon size="16">
                            <log
                                :inverse="selectedSubTab === BrowserTabType.SlowLog.toString()"
                                :stroke-color="themeVars.tabColor"
                                stroke-width="4" />
                        </n-icon>
                        <span>{{ $t('interface.sub_tab.slow_log') }}</span>
                    </n-space>
                </template>
                <content-slog :server="props.server" />
            </n-tab-pane>

            <!-- command monitor pane -->
            <n-tab-pane :name="BrowserTabType.CmdMonitor.toString()" display-directive="show:lazy">
                <template #tab>
                    <n-space :size="5" :wrap-item="false" align="center" inline justify="center">
                        <n-icon size="16">
                            <monitor
                                :inverse="selectedSubTab === BrowserTabType.CmdMonitor.toString()"
                                :stroke-color="themeVars.tabColor"
                                stroke-width="4" />
                        </n-icon>
                        <span>{{ $t('interface.sub_tab.cmd_monitor') }}</span>
                    </n-space>
                </template>
                <content-monitor :server="props.server" />
            </n-tab-pane>

            <!-- pub/sub message pane -->
            <n-tab-pane :name="BrowserTabType.PubMessage.toString()" display-directive="show:lazy">
                <template #tab>
                    <n-space :size="5" :wrap-item="false" align="center" inline justify="center">
                        <n-icon size="16">
                            <subscribe
                                :inverse="selectedSubTab === BrowserTabType.PubMessage.toString()"
                                :stroke-color="themeVars.tabColor"
                                stroke-width="4" />
                        </n-icon>
                        <span>{{ $t('interface.sub_tab.pub_message') }}</span>
                    </n-space>
                </template>
                <content-pubsub :server="props.server" />
            </n-tab-pane>
        </n-tabs>
    </div>
</template>

<style lang="scss" scoped>
@use '@/styles/content';

.content-container {
    //padding: 5px 5px 0;
    //padding-top: 0;
    box-sizing: border-box;
    background-color: v-bind('themeVars.tabColor');
}
</style>

<style lang="scss">
.content-sub-tab {
    background-color: v-bind('themeVars.tabColor');
    height: 100%;
}

.content-sub-tab-pane {
    padding: 0 !important;
    height: 100%;
    background-color: v-bind('themeVars.bodyColor');
    overflow: hidden;
}

.n-tabs .n-tabs-bar {
    transition: none !important;
}
</style>
