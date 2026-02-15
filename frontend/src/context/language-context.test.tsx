import { describe, it, expect, beforeEach } from "vitest"
import { renderHook, act } from "@testing-library/react"
import { LanguageProvider, useLanguage } from "./language-context"
import type { ReactNode } from "react"

function wrapper({ children }: { children: ReactNode }) {
  return <LanguageProvider>{children}</LanguageProvider>
}

describe("useLanguage", () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it("defaults to English when no stored preference", () => {
    const { result } = renderHook(() => useLanguage(), { wrapper })
    expect(result.current.language).toBe("en")
  })

  it("uses stored language from localStorage", () => {
    localStorage.setItem("language", "lt")
    const { result } = renderHook(() => useLanguage(), { wrapper })
    expect(result.current.language).toBe("lt")
  })

  it("toggles from en to lt", () => {
    const { result } = renderHook(() => useLanguage(), { wrapper })
    act(() => result.current.toggleLanguage())
    expect(result.current.language).toBe("lt")
    expect(localStorage.getItem("language")).toBe("lt")
  })

  it("toggles from lt to en", () => {
    localStorage.setItem("language", "lt")
    const { result } = renderHook(() => useLanguage(), { wrapper })
    act(() => result.current.toggleLanguage())
    expect(result.current.language).toBe("en")
    expect(localStorage.getItem("language")).toBe("en")
  })

  it("t() returns English translation by default", () => {
    const { result } = renderHook(() => useLanguage(), { wrapper })
    expect(result.current.t("nav.dashboard")).toBe("Dashboard")
  })

  it("t() returns Lithuanian translation after toggle", () => {
    const { result } = renderHook(() => useLanguage(), { wrapper })
    act(() => result.current.toggleLanguage())
    expect(result.current.t("nav.dashboard")).toBe("Skydelis")
  })

  it("t() returns key if translation is missing", () => {
    const { result } = renderHook(() => useLanguage(), { wrapper })
    // Cast to any to test with a key that doesn't exist in the type
    const t = result.current.t as (key: string) => string
    expect(t("nonexistent.key")).toBe("nonexistent.key")
  })

  it("ignores invalid stored language values", () => {
    localStorage.setItem("language", "fr")
    const { result } = renderHook(() => useLanguage(), { wrapper })
    expect(result.current.language).toBe("en")
  })
})
