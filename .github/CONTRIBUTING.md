## Tiny RDM Contribute Guide

### Multi-language Contributions

#### Adding New Language

1. New file: Add a new JSON file in the [frontend/src/langs](../frontend/src/langs/), with the file naming format is "
   {language}-{region}.json", e.g. English is "en-us.json", simplified Chinese is "zh-cn.json". Highly recommended to duplicate the [en-us.json](../frontend/src/langs/en-us.json) file and rename it.
2. Fill content: Refer to [en-us.json](../frontend/src/langs/en-us.json), or duplicate the file and modify the language content.
3. Update codes: Edit[frontend/src/langs/index.js](.../frontend/src/langs/index.js), import the new language data inside.
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
4. Submit review once there are no issues with the translation context in the application. (learn how to submit)

### Code Submission`(To be completed)`

#### Pull Request Title
The format of PR's title like "<type>: <description>"
- type: PR type
- description: PR description

PR type list below:

| type     | description                                        |
|----------|----------------------------------------------------|
| revert   | Revert a commit                                    |
| feat     | New features                                       |
| perf     | Performance improvements                           |
| fix      | Fix any bugs                                       |
| style    | Style updates                                      |
| docs     | Document updates                                   |
| refactor | Code refactors                                     |
| chore    | Some chores                                        |
| ci       | Automation process configuration or script updates |
