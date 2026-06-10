// Vuetify 3 setup with a professional Green + White theme.
import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import { createVuetify } from 'vuetify'

const rentacarTheme = {
  dark: false,
  colors: {
    background: '#F5F7F5',
    surface: '#FFFFFF',
    primary: '#1B5E20', // deep green
    'primary-darken-1': '#0F3D12',
    secondary: '#43A047',
    accent: '#66BB6A',
    error: '#C62828',
    info: '#1565C0',
    success: '#2E7D32',
    warning: '#EF6C00',
  },
}

export default createVuetify({
  theme: {
    defaultTheme: 'rentacarTheme',
    themes: { rentacarTheme },
  },
  defaults: {
    VCard: { rounded: 'lg' },
    VBtn: { rounded: 'md' },
    VTextField: { variant: 'outlined', density: 'comfortable' },
    VSelect: { variant: 'outlined', density: 'comfortable' },
    VTextarea: { variant: 'outlined', density: 'comfortable' },
  },
})
