/**
 * plugins/vuetify.ts
 *
 * Framework documentation: https://vuetifyjs.com`
 */

// Styles
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'

// Composables
import { createVuetify } from 'vuetify'

// https://vuetifyjs.com/en/introduction/why-vuetify/#feature-guides
export default createVuetify({
  theme: {
    defaultTheme: 'blackOrange',
    themes: {
      blackOrange: {
        dark: true,
        colors: {
          background: '#000000',
          surface: '#000000',
          'surface-bright': '#111111',
          'surface-light': '#1a1a1a',
          'surface-variant': '#0a0a0a',
          'on-surface-variant': '#66ccff',
          primary: '#66ccff',
          'primary-darken-1': '#3399cc',
          secondary: '#99ddff',
          'secondary-darken-1': '#3388bb',
          error: '#ff4444',
          info: '#66ccff',
          success: '#66ccff',
          warning: '#ffcc00',
          'on-background': '#66ccff',
          'on-surface': '#66ccff',
          'on-primary': '#000000',
          'on-secondary': '#000000',
          'on-error': '#000000',
          'on-info': '#000000',
          'on-success': '#000000',
          'on-warning': '#000000',
        }
      }
    }
  },
})
