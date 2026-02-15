export function formatCents(cents: number): string {
  const dollars = cents / 100
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
  }).format(dollars)
}

export function centsToInput(cents: number): string {
  return (cents / 100).toFixed(2)
}

export function inputToCents(value: string): number {
  return Math.round(parseFloat(value) * 100)
}
