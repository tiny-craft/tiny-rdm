import { merge } from 'lodash'

/**
 *
 * @type import('naive-ui').GlobalThemeOverrides
 */
export const themeOverrides = {
    common: {
        primaryColor: '#D33A31',
        primaryColorHover: '#FF6B6B',
        primaryColorPressed: '#D5271C',
        primaryColorSuppl: '#FF6B6B',
        borderRadius: '4px',
        borderRadiusSmall: '3px',
        lineHeight: 1.5,
        scrollbarWidth: '8px',
        tabColor: '#FFFFFF',
    },
    Button: {
        heightMedium: '32px',
    },
    Tag: {
        // borderRadius: '3px'
        heightLarge: '32px',
    },
    Input: {
        heightMedium: '32px',
    },
    Tabs: {
        tabGapSmallCard: '2px',
        tabGapMediumCard: '2px',
        tabGapLargeCard: '2px',
        tabFontWeightActive: 450,
    },
    Form: {
        labelFontSizeTopSmall: '12px',
        labelFontSizeTopMedium: '13px',
        labelFontSizeTopLarge: '13px',
        labelHeightSmall: '18px',
        labelHeightMedium: '18px',
        labelHeightLarge: '18px',
        labelPaddingVertical: '0 0 5px 2px',
        feedbackHeightSmall: '18px',
        feedbackHeightMedium: '18px',
        feedbackHeightLarge: '20px',
        feedbackFontSizeSmall: '11px',
        feedbackFontSizeMedium: '12px',
        feedbackFontSizeLarge: '12px',
        labelTextColor: 'rgb(113,120,128)',
        labelFontWeight: '450',
    },
    Radio: {
        buttonColorActive: '#D13B37',
        buttonTextColorActive: '#FFF',
    },
    Message: {
        margin: '0 0 38px 0',
    },
}

/**
 *
 * @type import('naive-ui').GlobalThemeOverrides
 */
const _darkThemeOverrides = {
    common: {
        bodyColor: '#1E1E1E',
        tabColor: '#1E1E1E',
        borderColor: '#515151',
    },
    Tree: {
        nodeTextColor: '#CECED0',
    },
    Card: {
        colorEmbedded: '#18181C',
    },
    Dropdown: {
        color: '#272727',
    },
}

export const darkThemeOverrides = merge({}, themeOverrides, _darkThemeOverrides)
