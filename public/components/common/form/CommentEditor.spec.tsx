import React from "react"
import { fireEvent, render, screen, waitFor } from "@testing-library/react"
import { Fider, FiderContext } from "@fider/services"
import CommentEditor from "./CommentEditor"

const fider = Fider.initialize({
  settings: { environment: "development", oauth: [] },
  tenant: { allowedSchemes: "" },
  user: undefined,
})

const renderEditor = (props: Partial<React.ComponentProps<typeof CommentEditor>> = {}) =>
  render(
    <FiderContext.Provider value={fider}>
      <CommentEditor field="description" disabled={false} initialValue="" {...props} />
    </FiderContext.Provider>
  )

const editorImages = (container: HTMLElement) => Array.from(container.querySelectorAll<HTMLImageElement>(".ProseMirror img:not(.ProseMirror-separator)"))

describe("<CommentEditor /> unsaved images", () => {
  test("shows an image from a restored draft using its cached data", async () => {
    const onGetImageSrc = (bkey: string) => (bkey === "attachments/draft.png" ? "AAAA" : "")
    const { container } = renderEditor({ initialValue: "![](fider-image:attachments/draft.png)", onGetImageSrc })

    await waitFor(() => expect(editorImages(container)).toHaveLength(1))
    expect(editorImages(container)[0].getAttribute("src")).toEqual("data:image/jpeg;base64,AAAA")
  })

  test("keeps showing an image attached in this session after switching to markdown and back", async () => {
    const onImageUploaded = jest.fn()
    const { container } = renderEditor({ onImageUploaded })

    const fileInput = container.querySelector<HTMLInputElement>('input[type="file"]')
    expect(fileInput).not.toBeNull()
    const file = new File(["hello"], "pic.png", { type: "image/png" })
    fireEvent.change(fileInput as HTMLInputElement, { target: { files: [file] } })

    await waitFor(() => expect(editorImages(container)).toHaveLength(1))
    const dataUrl = editorImages(container)[0].getAttribute("src")
    expect(dataUrl).toMatch(/^data:image\/png;base64,/)

    fireEvent.click(screen.getByTitle("Markdown Mode"))
    const textarea = screen.getByTestId<HTMLTextAreaElement>("markdown-textarea")
    expect(textarea.value).toMatch(/^!\[\]\(fider-image:attachments\/[a-zA-Z0-9]+-pic\.png\)$/)

    fireEvent.click(screen.getByTitle("Rich Text Mode"))
    await waitFor(() => expect(editorImages(container)).toHaveLength(1))
    expect(editorImages(container)[0].getAttribute("src")).toEqual(dataUrl)
  })
})
