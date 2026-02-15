import { describe, it, expect, vi, beforeEach } from "vitest"
import { screen, waitFor } from "@testing-library/react"
import { renderWithProviders } from "@/test/test-utils"
import ImportExportPage from "./import-export"

vi.mock("@/lib/api", () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    interceptors: {
      request: { use: vi.fn() },
      response: { use: vi.fn() },
    },
  },
}))

import api from "@/lib/api"

const mockedApi = vi.mocked(api)

beforeEach(() => {
  vi.clearAllMocks()
  localStorage.setItem("token", "test-token")
  mockedApi.get.mockImplementation((url: string) => {
    if (url === "/auth/me") {
      return Promise.resolve({
        data: { id: "u1", email: "admin@test.com", name: "Admin", role: "admin" },
      })
    }
    return Promise.reject(new Error("not mocked"))
  })
})

describe("ImportExportPage", () => {
  it("renders export and import sections", async () => {
    renderWithProviders(<ImportExportPage />)

    await waitFor(() => {
      expect(screen.getByText("Export Transactions")).toBeInTheDocument()
    })
    expect(screen.getByText("Import Transactions")).toBeInTheDocument()
  })

  it("shows export button", async () => {
    renderWithProviders(<ImportExportPage />)

    await waitFor(() => {
      expect(screen.getByRole("button", { name: /export csv/i })).toBeInTheDocument()
    })
  })

  it("shows import button and file input", async () => {
    renderWithProviders(<ImportExportPage />)

    await waitFor(() => {
      expect(screen.getByRole("button", { name: /import csv/i })).toBeInTheDocument()
    })
    expect(screen.getByLabelText(/csv file/i)).toBeInTheDocument()
  })
})
