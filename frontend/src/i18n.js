import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

const resources = {
  en: {
    translation: {
      "home.responseTitle": "Response",
      "drawer.skills": "Skills",
      "drawer.immunology": "Immunology",
      "drawer.chat": "AI Chat",
      "drawer.upload": "Upload",
      "drawer.settings": "Settings",
      "welcome": "Welcome to ESP Organizer",
      "loading": "Loading...",
      "error": "Error",
      "search": "Search",
      "query": "Query",
      "submit": "Submit",
      "clear": "Clear"
    }
  }
};

i18n
  .use(initReactI18next)
  .init({
    resources,
    lng: 'en',
    fallbackLng: 'en',
    interpolation: {
      escapeValue: false
    }
  });

export default i18n;
