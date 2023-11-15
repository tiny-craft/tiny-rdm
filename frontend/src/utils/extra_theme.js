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
    ribbonColor: '#F9F9F9',
    ribbonActiveColor: '#E3E3E3',
    sidebarColor: '#F2F2F2',
    splitColor: '#DADADA',
}

/**
 *
 * @type ExtraTheme
 */
export const extraDarkTheme = {
    titleColor: '#262626',
    ribbonColor: '#2C2C2C',
    ribbonActiveColor: '#363636',
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
