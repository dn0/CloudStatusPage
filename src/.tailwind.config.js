const defaultTheme = require('tailwindcss/defaultTheme')
const colors = require('tailwindcss/colors')

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    'pkg/mon/web/templates/*.templ',
    'pkg/mon/web/views/**/*.templ',
  ],
  darkMode: 'media',
  safelist: [
    'htmx-request', // For #loading-main
  ],
  theme: {
    colors: {
      transparent: 'transparent',
      current: 'currentColor',
      black: colors.black,
      white: colors.white,
      gray: colors.gray,
      red: colors.red,
      green: colors.green,
      blue: colors.sky,
      amber: colors.amber,
      emerald: colors.emerald,
    },
    extend: {
      fontFamily: {
        'sans': ['Inter', ...defaultTheme.fontFamily.sans],
      },
      fontSize: {
        ss: ['13px', '18px'],
      },
      colors: {
        'drk': '#18181b',
        'lgt': '#fafafa',
      },
      rotate: {
        '40': '40deg',
      },
      translate: {
        'w14-x': '37.5px',
        'h24-y': '15.5px',
        'h32-y': '30px',
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ]
}
