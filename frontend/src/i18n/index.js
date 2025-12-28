import { createI18n } from 'vue-i18n'
import zh from './locales/zh-CN.json'
import en from './locales/en-US.json'

export function initI18n() {
  // 从 localStorage 读取语言偏好
  const savedLocale = localStorage.getItem('app-language')
  const locale = savedLocale || 'zh-CN'
  
  const i18n = createI18n({
    legacy: false,          // 使用 Composition API 模式
    locale: locale,
    fallbackLocale: 'zh-CN',
    messages: {
      'zh-CN': zh,
      'en-US': en
    }
  })
  
  return i18n
}

export function setLanguage(locale) {
  localStorage.setItem('app-language', locale)
}

export function getLanguage() {
  return localStorage.getItem('app-language') || 'zh-CN'
}
