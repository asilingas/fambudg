import { describe, it, expect } from "vitest"
import { formatCents, centsToInput, inputToCents } from "./format"

describe("formatCents", () => {
  it("formats positive cents as EUR", () => {
    expect(formatCents(1999)).toMatch(/19,99.*€/)
  })

  it("formats negative cents as EUR", () => {
    expect(formatCents(-4599)).toMatch(/-45,99.*€/)
  })

  it("formats zero", () => {
    expect(formatCents(0)).toMatch(/0,00.*€/)
  })
})

describe("centsToInput", () => {
  it("converts cents to decimal string", () => {
    expect(centsToInput(1999)).toBe("19.99")
  })
})

describe("inputToCents", () => {
  it("converts decimal string to cents", () => {
    expect(inputToCents("19.99")).toBe(1999)
  })

  it("handles whole numbers", () => {
    expect(inputToCents("20")).toBe(2000)
  })
})
