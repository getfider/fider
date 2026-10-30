import React from "react"
import { act, fireEvent, render, screen, waitFor } from "@testing-library/react"
import { actions, FiderContext } from "@fider/services"
import { fiderMock, setupModalRoot } from "@fider/services/testing"
import { ShareFeedback } from "./ShareFeedback"

jest.mock("@fider/services/actions", () => {
  const actual = jest.requireActual("@fider/services/actions")
  return { ...actual, createPost: jest.fn() }
})

const mockCreatePost = actions.createPost as jest.Mock

type InputProps = {
  children?: React.ReactNode
  label?: string
  value?: string
  onChange?: (value: string) => void
}

type FormProps = {
  children?: React.ReactNode
  error?: { errors?: Array<{ field?: string; message?: string }> }
}

jest.mock("@fider/components/common", () => ({
  Button: (props: React.ButtonHTMLAttributes<HTMLButtonElement> & { variant?: string }) => {
    const buttonProps = { ...props }
    delete buttonProps.variant
    return <button {...buttonProps} />
  },
  CloseIcon: () => null,
  Form: ({ children, error }: FormProps) => (
    <form>
      {children}
      {error && <div data-testid="form-error">{error.errors?.map((item) => item.message).join(" ")}</div>}
    </form>
  ),
  Input: ({ children, label, value, onChange }: InputProps) => (
    <label>
      {label}
      <input aria-label={label} value={value} onChange={(event) => onChange?.(event.currentTarget.value)} />
      {children}
    </label>
  ),
  LegalFooter: () => <div data-testid="legal-footer" />,
  Modal: {
    Window: ({ children, isOpen }: { children?: React.ReactNode; isOpen: boolean }) => (isOpen ? <div role="dialog">{children}</div> : null),
    Header: ({ children }: { children?: React.ReactNode }) => <>{children}</>,
    Content: ({ children }: { children?: React.ReactNode }) => <>{children}</>,
  },
}))

jest.mock("@fider/components/common/SignInControl", () => ({
  SignInControl: ({ onCodeVerified }: { onCodeVerified?: () => void }) => (
    <button data-testid="signin-control" onClick={onCodeVerified}>
      Verify code
    </button>
  ),
}))

jest.mock("@fider/components/common/TagsSelect", () => ({ TagsSelect: () => null }))
jest.mock("@fider/components/common/form/CommentEditor", () => () => null)
jest.mock("@fider/hooks/useAttachments", () => ({
  useAttachments: () => ({
    attachments: [],
    handleImageUploaded: jest.fn(),
    getImageSrc: jest.fn(),
    clearAttachments: jest.fn(),
  }),
}))
jest.mock("./SimilarPosts", () => ({ SimilarPosts: () => null }))
function MockTrans({ children, id, message }: { children?: React.ReactNode; id?: string; message?: string }) {
  return <>{children || message || (id === "newpost.modal.submit" ? "Submit your idea" : id)}</>
}

jest.mock("@lingui/react/macro", () => ({ Trans: MockTrans }))
jest.mock("@lingui/react", () => ({ Trans: MockTrans }))
jest.mock("@lingui/core", () => ({
  i18n: {
    _: ({ message, values }: { message: string; values?: Record<string, number> }) =>
      message.replace(/\{(\w+)\}/g, (_, key: string) => String(values?.[key] ?? `{${key}}`)),
  },
}))
describe("<ShareFeedback />", () => {
  beforeEach(() => {
    setupModalRoot()
    localStorage.clear()
    sessionStorage.clear()
    jest.clearAllMocks()
    mockCreatePost.mockReset().mockResolvedValue({
      ok: false,
      error: { errors: [{ field: "title", message: "Title is too short" }] },
    })
  })

  const renderShareFeedback = (isAuthenticated = false) => {
    const fider = isAuthenticated ? fiderMock.authenticated() : fiderMock.notAuthenticated()

    return render(
      <FiderContext.Provider value={fider}>
        <ShareFeedback isOpen={true} placeholder="" onClose={jest.fn()} tags={[]} />
      </FiderContext.Provider>
    )
  }

  test("explains the title minimum and counts collapsed whitespace", () => {
    renderShareFeedback(true)
    const title = screen.getByLabelText("Give your idea a title")

    fireEvent.change(title, { target: { value: "  one   two " } })

    expect(screen.getByText("Title must be at least 10 characters (3 more needed).")).toBeInTheDocument()
    expect(screen.queryByRole("button", { name: "Submit your idea" })).not.toBeInTheDocument()

    fireEvent.change(title, { target: { value: "1234567890" } })

    expect(screen.queryByText(/Title must be at least/)).not.toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Submit your idea" })).toBeInTheDocument()
  })

  test("keeps the authenticated flow after post creation fails following code verification", async () => {
    renderShareFeedback()
    const title = screen.getByLabelText("Give your idea a title")

    fireEvent.change(title, { target: { value: "too short" } })
    expect(screen.getByTestId("signin-control")).toBeInTheDocument()

    await act(async () => {
      fireEvent.click(screen.getByTestId("signin-control"))
    })

    await waitFor(() => expect(mockCreatePost).toHaveBeenCalledWith("too short", "", [], []))
    expect(screen.queryByTestId("signin-control")).not.toBeInTheDocument()
    await waitFor(() => expect(screen.getByTestId("form-error")).toHaveTextContent("Title is too short"))
    expect(screen.queryByTestId("legal-footer")).not.toBeInTheDocument()
  })
})
