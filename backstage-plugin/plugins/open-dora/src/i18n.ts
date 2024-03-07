// /* istanbul ignore file */
import i18next from 'i18next';
import { initReactI18next } from 'react-i18next';

// Import all translation files
import translationEnglish from './locales/en/translation.json';
import translationDutch from './locales/nl/translation.json';

// Using translation
const resources = {
  en: {
    translation: translationEnglish,
  },
  nl: {
    translation: translationDutch,
  },
};

i18next.use(initReactI18next).init({
  resources,
  lng: 'en', // default language
});

export default i18next;