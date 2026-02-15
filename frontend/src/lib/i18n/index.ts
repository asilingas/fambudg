import { en } from "./en"
import { lt } from "./lt"

export type TranslationKey = keyof typeof en
export type Language = "en" | "lt"

export const translations: Record<Language, Record<TranslationKey, string>> = {
  en,
  lt,
}
