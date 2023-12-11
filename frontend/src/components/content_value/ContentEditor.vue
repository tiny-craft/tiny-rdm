<script setup>
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import * as monaco from 'monaco-editor'
import usePreferencesStore from 'stores/preferences.js'
import { useThemeVars } from 'naive-ui'

const props = defineProps({
    content: {
        type: String,
    },
    language: {
        type: String,
        default: 'json',
    },
    readonly: {
        type: String,
    },
    loading: {
        type: Boolean,
    },
    showLineNum: {
        type: Boolean,
        default: true,
    },
    border: {
        type: Boolean,
        default: false,
    },
})

const emit = defineEmits(['reset', 'input', 'save'])

const themeVars = useThemeVars()
/** @type {HTMLElement|null} */
const editorRef = ref(null)
/** @type monaco.editor.IStandaloneCodeEditor */
let editorNode = null

const destroyEditor = () => {
    if (editorNode != null && editorNode.dispose != null) {
        const model = editorNode.getModel()
        if (model != null) {
            model.dispose()
        }
        editorNode.dispose()
        editorNode = null
    }
}

const readonlyValue = computed(() => {
    return props.readonly || props.loading
})

const pref = usePreferencesStore()
onMounted(async () => {
    if (editorRef.value != null) {
        const { fontSize, fontFamily = undefined } = pref.generalFont
        editorNode = monaco.editor.create(editorRef.value, {
            // value: props.content,
            theme: pref.isDark ? 'rdm-dark' : 'rdm-light',
            language: props.language,
            lineNumbers: props.showLineNum ? 'on' : 'off',
            readOnly: readonlyValue.value,
            colorDecorators: true,
            accessibilitySupport: 'off',
            wordWrap: 'on',
            tabSize: 2,
            folding: true,
            fontFamily,
            fontSize,
            scrollBeyondLastLine: false,
            automaticLayout: true,
            scrollbar: {
                useShadows: false,
                verticalScrollbarSize: '10px',
            },
            // formatOnType: true,
            contextmenu: false,
            lineNumbersMinChars: 2,
            lineDecorationsWidth: 0,
            minimap: {
                enabled: false,
            },
            selectionHighlight: false,
            renderLineHighlight: 'gutter',
        })

        // add shortcut for save
        editorNode.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyS, (event) => {
            emit('save')
        })

        // editorNode.onDidChangeModelLanguageConfiguration(() => {
        //     editorNode?.getAction('editor.action.formatDocument')?.run()
        // })

        if (editorNode.onDidChangeModelContent) {
            editorNode.onDidChangeModelContent(() => {
                emit('input', editorNode.getValue())
            })
        }
    }
})

watch(
    () => props.content,
    async (content) => {
        if (editorNode != null) {
            editorNode.setValue(content)
            await nextTick(() => emit('reset', content))
        }
    },
)

watch(
    () => readonlyValue.value,
    (readOnly) => {
        if (editorNode != null) {
            editorNode.updateOptions({
                readOnly,
            })
        }
    },
)

watch(
    () => props.language,
    (language) => {
        if (editorNode != null) {
            const model = editorNode.getModel()
            if (model != null) {
                monaco.editor.setModelLanguage(model, language)
            }
        }
    },
)

watch(
    () => pref.isDark,
    (dark) => {
        if (editorNode != null) {
            editorNode.updateOptions({
                theme: dark ? 'rdm-dark' : 'rdm-light',
            })
        }
    },
)

onUnmounted(() => {
    destroyEditor()
})
</script>

<template>
    <div :class="{ 'editor-border': props.border === true }" style="position: relative">
        <div ref="editorRef" class="editor-inst" />
    </div>
</template>

<style lang="scss" scoped>
.editor-border {
    border: 1px solid v-bind('themeVars.borderColor');
    border-radius: v-bind('themeVars.borderRadius');
    padding: 3px;
    box-sizing: border-box;
}

.editor-inst {
    position: absolute;
    top: 2px;
    bottom: 2px;
    left: 2px;
    right: 2px;
}
</style>
