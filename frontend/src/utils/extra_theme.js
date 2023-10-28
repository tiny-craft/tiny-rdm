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
    titleColor: '#F2F2F2',
    sidebarColor: '#E9E9E9',
    splitColor: '#DADADA',
}

/**
 *
 * @type ExtraTheme
 */
export const extraDarkTheme = {
    titleColor: '#363636',
    sidebarColor: '#262626',
    splitColor: '#474747',
}

/**
 *
 * @param {boolean} dark
 * @return ExtraTheme
 */
export const extraTheme = (dark) => {
    return dark ? extraDarkTheme : extraLightTheme
}
