<script setup>
import iconUrl from '@/assets/images/icon.png'
import useDialog from 'stores/dialog.js'
import { useThemeVars } from 'naive-ui'
import { BrowserOpenURL } from 'wailsjs/runtime/runtime.js'
import { GetAppVersion } from 'wailsjs/go/services/preferencesService.js'
import { ref, onMounted } from 'vue'

const themeVars = useThemeVars()
const dialogStore = useDialog()
const version = ref('')

onMounted(() => {
    GetAppVersion().then(({ data }) => {
        version.value = data.version
    })
})

const onOpenSource = () => {
    BrowserOpenURL('https://github.com/tiny-craft/tiny-rdm')
}

const onOpenWebsite = () => {
    BrowserOpenURL('https://redis.tinycraft.cc/')
}
</script>

<template>
    <n-modal v-model:show="dialogStore.aboutDialogVisible" :show-icon="false" preset="dialog" transform-origin="center">
        <n-space vertical align="center" :wrap-item="false" :wrap="false" :size="10">
            <n-avatar :size="120" color="#0000" :src="iconUrl"></n-avatar>
            <div class="about-app-title">Tiny RDM</div>
            <n-text>{{ version }}</n-text>
            <n-space align="center" :wrap-item="false" :wrap="false" :size="5">
                <n-text class="about-link" @click="onOpenSource">{{ $t('dialogue.about.source') }}</n-text>
                <n-divider vertical />
                <n-text class="about-link" @click="onOpenWebsite">{{ $t('dialogue.about.website') }}</n-text>
            </n-space>
            <div class="about-copyright" :style="{ color: themeVars.textColor3 }">
                Copyright © 2023 Tinycraft.cc All rights reserved
            </div>
        </n-space>
    </n-modal>
</template>

<style scoped lang="scss">
.about-app-title {
    font-weight: bold;
    font-size: 18px;
    margin: 5px;
}

.about-link {
    cursor: pointer;

    &:hover {
        text-decoration: underline;
    }
}

.about-copyright {
    font-size: 12px;
}
</style>
