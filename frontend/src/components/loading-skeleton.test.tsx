import { render, screen } from "@testing-library/react"
import { DashboardSkeleton, PageSkeleton, TableSkeleton, CardSkeleton } from "./loading-skeleton"

describe("Loading Skeletons", () => {
  it("renders DashboardSkeleton with skeleton elements", () => {
    const { container } = render(<DashboardSkeleton />)
    const skeletons = container.querySelectorAll("[data-slot='skeleton']")
    expect(skeletons.length).toBeGreaterThan(0)
  })

  it("renders PageSkeleton with skeleton elements", () => {
    const { container } = render(<PageSkeleton />)
    const skeletons = container.querySelectorAll("[data-slot='skeleton']")
    expect(skeletons.length).toBeGreaterThan(0)
  })

  it("renders TableSkeleton with configurable rows", () => {
    const { container } = render(<TableSkeleton rows={3} />)
    // 1 header skeleton + 3 row skeletons = 4 total
    const skeletons = container.querySelectorAll("[data-slot='skeleton']")
    expect(skeletons.length).toBe(4)
  })

  it("renders CardSkeleton", () => {
    const { container } = render(<CardSkeleton />)
    const skeletons = container.querySelectorAll("[data-slot='skeleton']")
    expect(skeletons.length).toBe(2)
  })
})
