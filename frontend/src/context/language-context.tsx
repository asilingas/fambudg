import { createContext, useContext, useState, useEffect, useCallback, type ReactNode } from "react"
import { translations, type Language, type TranslationKey } from "@/lib/i18n"

interface LanguageContextValue {
  language: Language
  toggleLanguage: () => void
  t: (key: TranslationKey) => string
}

const LanguageContext = createContext<LanguageContextValue | null>(null)

export function LanguageProvider({ children }: { children: ReactNode }) {
  const [language, setLanguage] = useState<Language>(() => {
    const stored = localStorage.getItem("language")
    if (stored === "en" || stored === "lt") return stored
    return "en"
  })

  useEffect(() => {
    localStorage.setItem("language", language)
  }, [language])

  const toggleLanguage = useCallback(() => {
    setLanguage((prev) => (prev === "en" ? "lt" : "en"))
  }, [])

  const t = useCallback(
    (key: TranslationKey): string => {
      return translations[language][key] ?? key
    },
    [language],
  )

  return (
    <LanguageContext.Provider value={{ language, toggleLanguage, t }}>
      {children}
    </LanguageContext.Provider>
  )
}

export function useLanguage() {
  const context = useContext(LanguageContext)
  if (!context) {
    throw new Error("useLanguage must be used within a LanguageProvider")
  }
  return context
}
