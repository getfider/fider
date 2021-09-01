import React from "react"

// defines DOM related expect methods
import "@testing-library/jest-dom/extend-expect"

// Mock for LinguiJS so we don't need to setup i18n on each test
jest.mock("@lingui/react", () => ({
  Trans: function TransMock({ children }: { children: React.ReactNode }) {
    return <>{children}</>
  },

  t: function tMock(id: string): string {
    return id
  },

  Plural: function PluralMock({ value, one, other }: { value: number; one: React.ReactNode; other: React.ReactNode }) {
    return <>{value > 1 ? other : one}</>
  },
}))
