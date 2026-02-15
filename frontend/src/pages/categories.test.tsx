import { describe, it, expect, vi, beforeEach } from "vitest"
import { screen, waitFor } from "@testing-library/react"
import { renderWithProviders } from "@/test/test-utils"
import CategoriesPage from "./categories"

const mockUseAuth = vi.fn()

vi.mock("@/context/auth-context", async () => {
  const actual = await vi.importActual<typeof import("@/context/auth-context")>(
    "@/context/auth-context",
  )
  return {
    ...actual,
    useAuth: (...args: unknown[]) => mockUseAuth(...args),
  }
})

vi.mock("@/lib/api", () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
    interceptors: {
      request: { use: vi.fn() },
      response: { use: vi.fn() },
    },
  },
}))

import api from "@/lib/api"

const mockedApi = vi.mocked(api)

const categories = [
  { id: "c1", name: "Groceries", type: "expense", sortOrder: 1 },
  { id: "c2", name: "Salary", type: "income", sortOrder: 2 },
]

beforeEach(() => {
  vi.clearAllMocks()
  mockedApi.get.mockImplementation((url: string) => {
    if (url === "/categories") return Promise.resolve({ data: categories })
    return Promise.reject(new Error("not mocked"))
  })
})

describe("CategoriesPage", () => {
  it("renders categories table", async () => {
    mockUseAuth.mockReturnValue({
      user: { id: "u1", role: "admin", name: "Admin", email: "a@test.com" },
      loading: false,
    })

    renderWithProviders(<CategoriesPage />)

    await waitFor(() => {
      expect(screen.getByText("Groceries")).toBeInTheDocument()
    })
    expect(screen.getByText("Salary")).toBeInTheDocument()
  })

  it("shows add button for admin", async () => {
    mockUseAuth.mockReturnValue({
      user: { id: "u1", role: "admin", name: "Admin", email: "a@test.com" },
      loading: false,
    })

    renderWithProviders(<CategoriesPage />)

    await waitFor(() => {
      expect(screen.getByRole("button", { name: /add category/i })).toBeInTheDocument()
    })
  })

  it("shows add button for member", async () => {
    mockUseAuth.mockReturnValue({
      user: { id: "u2", role: "member", name: "Member", email: "m@test.com" },
      loading: false,
    })

    renderWithProviders(<CategoriesPage />)

    await waitFor(() => {
      expect(screen.getByRole("button", { name: /add category/i })).toBeInTheDocument()
    })
  })

  it("hides add button for child", async () => {
    mockUseAuth.mockReturnValue({
      user: { id: "u3", role: "child", name: "Kid", email: "k@test.com" },
      loading: false,
    })

    renderWithProviders(<CategoriesPage />)

    await waitFor(() => {
      expect(screen.getByText("Groceries")).toBeInTheDocument()
    })
    expect(screen.queryByRole("button", { name: /add category/i })).not.toBeInTheDocument()
  })

  it("hides edit/delete buttons for non-admin", async () => {
    mockUseAuth.mockReturnValue({
      user: { id: "u2", role: "member", name: "Member", email: "m@test.com" },
      loading: false,
    })

    renderWithProviders(<CategoriesPage />)

    await waitFor(() => {
      expect(screen.getByText("Groceries")).toBeInTheDocument()
    })
    // Member shouldn't see edit/delete icons
    const editButtons = screen.queryAllByRole("button", { name: "" })
    // No icon buttons with Pencil/Trash2 should exist
    expect(editButtons.length).toBe(0)
  })
})
