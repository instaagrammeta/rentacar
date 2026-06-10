// Vue I18n configuration. The application UI is entirely in Russian.
import { createI18n } from 'vue-i18n'
import ru from '@/i18n/ru'

export default createI18n({
  legacy: false,
  globalInjection: true,
  locale: 'ru',
  fallbackLocale: 'ru',
  messages: { ru },
})
