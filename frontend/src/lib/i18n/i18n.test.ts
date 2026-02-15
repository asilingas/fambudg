import { describe, it, expect } from "vitest"
import { en } from "./en"
import { lt } from "./lt"

describe("i18n translation files", () => {
  const enKeys = Object.keys(en).sort()
  const ltKeys = Object.keys(lt).sort()

  it("both languages have the same number of keys", () => {
    expect(enKeys.length).toBe(ltKeys.length)
  })

  it("both languages have identical keys", () => {
    expect(enKeys).toEqual(ltKeys)
  })

  it("no empty values in English translations", () => {
    for (const [key, value] of Object.entries(en)) {
      expect(value, `en["${key}"] is empty`).not.toBe("")
    }
  })

  it("no empty values in Lithuanian translations", () => {
    for (const [key, value] of Object.entries(lt)) {
      expect(value, `lt["${key}"] is empty`).not.toBe("")
    }
  })

  it("English keys missing in Lithuanian", () => {
    const missingInLt = enKeys.filter((k) => !ltKeys.includes(k))
    expect(missingInLt, `Keys in EN but not LT: ${missingInLt.join(", ")}`).toEqual([])
  })

  it("Lithuanian keys missing in English", () => {
    const missingInEn = ltKeys.filter((k) => !enKeys.includes(k))
    expect(missingInEn, `Keys in LT but not EN: ${missingInEn.join(", ")}`).toEqual([])
  })
})
