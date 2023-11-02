## Tiny RDM 代码贡献指南

### 多国语言贡献 {#language}

#### 增加新的语言
1. 创建文件：在[frontend/src/langs](../frontend/src/langs/)目录下新增语言配置JSON文件，文件名格式为“{语言}-{地区}.json”，如英文为“en-us.json”，简体中文为“zh-cn.json”，建议直接复制[en-us.json](../frontend/src/langs/en-us.json)文件进行改名。
2. 填充内容：参考[en-us.json](../frontend/src/langs/en-us.json)，或者直接克隆一份文件，对语言部分内容进行修改。
3. 代码修改：在[frontend/src/langs/index.js](.../frontend/src/langs/index.js)文件内导入新增的语言数据
    ```javascript
    import en from './en-us'
    // import your new localize file 'zh-cn' here
    import zh from './zh-cn'
    
    export const lang = {
        en,
        // export new language data 'zh' here
        zh,
    }
    ```
4. 检查应用中对应翻译语境无问题后，可提交审核（[查看如何提交](#pull_request)）

### 代码提交`(待完善)` {#pull_request}

#### PR提交规范
PR提交格式为“<type>: <description>”
- type: 提交类型
- description: 提交内容描述

其中提交类型如下：

| 提交类型     | 类型描述         |
|----------|--------------|
| revert   | 回退某个commit提交 |
| feat     | 新功能/新特性      |
| perf     | 功能、体验等方面的优化  |
| fix      | 修复问题         |
| style    | 样式相关修改       |
| docs     | 文档更新         |
| refactor | 代码重构         |
| chore    | 杂项修改         |
| ci       | 自动化流程配置或脚本修改 |

