/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './app/components/**/*.{vue,js,ts}',
    './app/layouts/**/*.vue',
    './app/pages/**/*.vue',
    './app/composables/**/*.{js,ts}',
    './app/plugins/**/*.{js,ts}',
    './app/app.vue',
    './app/error.vue',
  ],
  theme: {
    extend: {
      colors: {
        // Wheelzie brand palette — coral/red accent on a light surface.
        primary: {
          50: '#fff1ef',
          100: '#ffe0db',
          200: '#ffc5bc',
          300: '#ff9d8e',
          400: '#ff6a54',
          500: '#f5452b',
          600: '#e22d18',
          700: '#bd2113',
          800: '#9c1f15',
          900: '#812018',
        },
        ink: {
          DEFAULT: '#1f2430',
          soft: '#5b6472',
          muted: '#8b93a3',
        },
        surface: {
          DEFAULT: '#ffffff',
          muted: '#f6f7fb',
          border: '#eceef3',
        },
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'Avenir', 'Helvetica', 'Arial', 'sans-serif'],
      },
      boxShadow: {
        card: '0 1px 2px rgba(16, 24, 40, 0.04), 0 4px 16px rgba(16, 24, 40, 0.06)',
      },
      borderRadius: {
        xl: '1rem',
        '2xl': '1.25rem',
      },
    },
  },
  plugins: [],
}
