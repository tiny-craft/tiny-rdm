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
        heightMedium: '32px',
        lineHeight: 1.5,
        scrollbarWidth: '8px',
        tabColor: '#FFFFFF',
    },
    Button: {
        heightMedium: '32px',
        paddingSmall: '0 8px',
        paddingMedium: '0 12px',
    },
    Tag: {
        borderRadius: '4px',
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
    Card: {
        colorEmbedded: '#FAFAFA',
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
    DataTable: {
        thPaddingSmall: '6px 8px',
        tdPaddingSmall: '6px 8px',
    },
    Dropdown: {
        borderRadius: '5px',
        optionIconSizeMedium: '18px',
        padding: '6px 2px',
        optionColorHover: '#D33A31',
        optionTextColorHover: '#FFF',
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
        colorEmbedded: '#212121',
    },
    Dropdown: {
        color: '#272727',
    },
    Popover: {
        color: '#2C2C32',
    },
}

export const darkThemeOverrides = merge({}, themeOverrides, _darkThemeOverrides)
