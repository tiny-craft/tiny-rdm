/**
 * @typedef ExtraTheme
 * @property {string} titleColor
 * @property {string} sidebarColor
 * @property {string} splitColor
 */

/**
 *
 * @type ExtraTheme
 */
export const extraLightTheme = {
    titleColor: '#F0F0F4',
    sidebarColor: '#F6F6F6',
    splitColor: '#E0E0E6',
}

/**
 *
 * @type ExtraTheme
 */
export const extraDarkTheme = {
    titleColor: '#202020',
    sidebarColor: '#202124',
    splitColor: '#323138',
}

/**
 *
 * @param {boolean} dark
 * @return ExtraTheme
 */
export const extraTheme = (dark) => {
    return dark ? extraDarkTheme : extraLightTheme
}
