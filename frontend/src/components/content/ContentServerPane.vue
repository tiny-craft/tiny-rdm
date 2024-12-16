<script setup>
import { computed } from 'vue'
import AddLink from '@/components/icons/AddLink.vue'
import useDialogStore from 'stores/dialog.js'
import { NButton, useThemeVars } from 'naive-ui'
import { BrowserOpenURL } from 'wailsjs/runtime/runtime.js'
import { find, includes, isEmpty } from 'lodash'
import usePreferencesStore from 'stores/preferences.js'

const themeVars = useThemeVars()
const dialogStore = useDialogStore()
const prefStore = usePreferencesStore()

const onOpenSponsor = (link) => {
    BrowserOpenURL(link)
}

const sponsorAd = computed(() => {
    try {
        const content = localStorage.getItem('sponsor_ad')
        const ads = JSON.parse(content)
        const ad = find(ads, ({ region }) => {
            return isEmpty(region) || includes(region, prefStore.currentLanguage)
        })
        return ad || null
    } catch {
        return null
    }
})
</script>

<template>
    <div class="content-container flex-box-v">
        <!-- TODO: replace icon to app icon -->
        <n-empty :description="$t('interface.empty_server_content')">
            <template #extra>
                <n-button :focusable="false" @click="dialogStore.openNewDialog()">
                    <template #icon>
                        <n-icon :component="AddLink" size="18" />
                    </template>
                    {{ $t('interface.new_conn') }}
                </n-button>
            </template>
        </n-empty>

        <n-button v-if="sponsorAd != null" class="sponsor-ad" style="" text @click="onOpenSponsor(sponsorAd.link)">
            {{ sponsorAd.name }}
        </n-button>
    </div>
</template>

<style lang="scss" scoped>
@use '@/styles/content';

.content-container {
    justify-content: center;
    padding: 5px;
    box-sizing: border-box;

    & > .sponsor-ad {
        text-align: center;
        margin-top: 20px;
        vertical-align: bottom;
        color: v-bind('themeVars.textColor3');
    }
}

.color-preset-item {
    width: 24px;
    height: 24px;
    margin-right: 2px;
    border: white 3px solid;
    cursor: pointer;

    &_selected,
    &:hover {
        border-color: #cdd0d6;
    }
}
</style>
