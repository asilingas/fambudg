import { renderHook, act } from "@testing-library/react"
import { useTheme } from "./use-theme"

describe("useTheme", () => {
  beforeEach(() => {
    localStorage.clear()
    document.documentElement.classList.remove("light", "dark")
  })

  it("defaults to light when no stored preference and no system dark mode", () => {
    window.matchMedia = vi.fn().mockImplementation(() => ({
      matches: false,
    }))

    const { result } = renderHook(() => useTheme())
    expect(result.current.theme).toBe("light")
  })

  it("uses stored theme from localStorage", () => {
    localStorage.setItem("theme", "dark")
    window.matchMedia = vi.fn().mockImplementation(() => ({
      matches: false,
    }))

    const { result } = renderHook(() => useTheme())
    expect(result.current.theme).toBe("dark")
  })

  it("toggles from light to dark", () => {
    window.matchMedia = vi.fn().mockImplementation(() => ({
      matches: false,
    }))

    const { result } = renderHook(() => useTheme())
    act(() => result.current.toggleTheme())
    expect(result.current.theme).toBe("dark")
    expect(localStorage.getItem("theme")).toBe("dark")
  })

  it("toggles from dark to light", () => {
    localStorage.setItem("theme", "dark")
    window.matchMedia = vi.fn().mockImplementation(() => ({
      matches: false,
    }))

    const { result } = renderHook(() => useTheme())
    act(() => result.current.toggleTheme())
    expect(result.current.theme).toBe("light")
    expect(localStorage.getItem("theme")).toBe("light")
  })

  it("adds theme class to document element", () => {
    window.matchMedia = vi.fn().mockImplementation(() => ({
      matches: false,
    }))

    renderHook(() => useTheme())
    expect(document.documentElement.classList.contains("light")).toBe(true)
  })
})
