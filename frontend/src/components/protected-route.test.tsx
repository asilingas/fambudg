import { describe, it, expect, vi } from "vitest"
import { screen } from "@testing-library/react"
import { render } from "@testing-library/react"
import { MemoryRouter, Routes, Route } from "react-router-dom"
import { ProtectedRoute } from "./protected-route"

const mockUseAuth = vi.fn()

vi.mock("@/context/auth-context", () => ({
  useAuth: (...args: unknown[]) => mockUseAuth(...args),
}))

describe("ProtectedRoute", () => {
  it("redirects to /login when not authenticated", () => {
    mockUseAuth.mockReturnValue({ user: null, loading: false })

    render(
      <MemoryRouter initialEntries={["/"]}>
        <Routes>
          <Route path="/login" element={<div>Login Page</div>} />
          <Route element={<ProtectedRoute />}>
            <Route path="/" element={<div>Dashboard</div>} />
          </Route>
        </Routes>
      </MemoryRouter>,
    )

    expect(screen.getByText("Login Page")).toBeInTheDocument()
  })

  it("shows nothing while loading", () => {
    mockUseAuth.mockReturnValue({ user: null, loading: true })

    const { container } = render(
      <MemoryRouter initialEntries={["/"]}>
        <Routes>
          <Route path="/login" element={<div>Login Page</div>} />
          <Route element={<ProtectedRoute />}>
            <Route path="/" element={<div>Dashboard</div>} />
          </Route>
        </Routes>
      </MemoryRouter>,
    )

    expect(container.textContent).toBe("")
  })

  it("renders children when authenticated", () => {
    mockUseAuth.mockReturnValue({
      user: { id: "1", email: "admin@test.com", name: "Admin", role: "admin" },
      loading: false,
    })

    render(
      <MemoryRouter initialEntries={["/"]}>
        <Routes>
          <Route path="/login" element={<div>Login Page</div>} />
          <Route element={<ProtectedRoute />}>
            <Route path="/" element={<div>Dashboard</div>} />
          </Route>
        </Routes>
      </MemoryRouter>,
    )

    expect(screen.getByText("Dashboard")).toBeInTheDocument()
  })

  it("redirects when role is not allowed", () => {
    mockUseAuth.mockReturnValue({
      user: { id: "2", email: "child@test.com", name: "Kid", role: "child" },
      loading: false,
    })

    render(
      <MemoryRouter initialEntries={["/admin"]}>
        <Routes>
          <Route path="/" element={<div>Home</div>} />
          <Route element={<ProtectedRoute allowedRoles={["admin"]} />}>
            <Route path="/admin" element={<div>Admin Only</div>} />
          </Route>
        </Routes>
      </MemoryRouter>,
    )

    expect(screen.getByText("Home")).toBeInTheDocument()
  })
})
