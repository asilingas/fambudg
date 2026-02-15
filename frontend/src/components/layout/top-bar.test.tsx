import { render, screen, fireEvent } from "@testing-library/react"

const mockLogout = vi.fn()
const mockToggleTheme = vi.fn()
let mockTheme = "light"
let mockUser: { name: string; role: string } | null = { name: "Test User", role: "admin" }

vi.mock("@/context/auth-context", () => ({
  useAuth: () => ({ user: mockUser, logout: mockLogout }),
}))

vi.mock("@/hooks/use-theme", () => ({
  useTheme: () => ({ theme: mockTheme, toggleTheme: mockToggleTheme }),
}))

import { TopBar } from "./top-bar"

describe("TopBar", () => {
  beforeEach(() => {
    mockTheme = "light"
    mockUser = { name: "Test User", role: "admin" }
    vi.clearAllMocks()
  })

  it("renders user name and role", () => {
    render(<TopBar />)
    expect(screen.getByText("Test User")).toBeInTheDocument()
    expect(screen.getByText("admin")).toBeInTheDocument()
  })

  it("renders theme toggle button", () => {
    render(<TopBar />)
    expect(screen.getByLabelText("Toggle theme")).toBeInTheDocument()
  })

  it("calls toggleTheme on theme button click", () => {
    render(<TopBar />)
    fireEvent.click(screen.getByLabelText("Toggle theme"))
    expect(mockToggleTheme).toHaveBeenCalled()
  })

  it("calls logout on logout button click", () => {
    render(<TopBar />)
    fireEvent.click(screen.getByLabelText("Logout"))
    expect(mockLogout).toHaveBeenCalled()
  })

  it("returns null when no user", () => {
    mockUser = null
    const { container } = render(<TopBar />)
    expect(container.firstChild).toBeNull()
  })
})
