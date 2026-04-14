import React from "react"
import { Button, Icon } from "@fider/components"

import { Tag } from "@fider/models"
import { actions, Failure, Fider } from "@fider/services"
import { ImportTagsResult } from "@fider/services/actions/tag"
import { AdminBasePage } from "../components/AdminBasePage"
import { TagFormState, TagForm } from "../components/TagForm"
import { TagListItem } from "../components/TagListItem"
import { VStack } from "@fider/components/layout"
import IconDownload from "@fider/assets/images/heroicons-download.svg"
import IconUpload from "@fider/assets/images/heroicons-upload.svg"
import IconPlus from "@fider/assets/images/heroicons-plus.svg"

interface ManageTagsPageProps {
  tags: Tag[]
}

interface ManageTagsPageState {
  isAdding: boolean
  allTags: Tag[]
  deleting?: number
  editing?: number
  importStatus?: "idle" | "loading" | "success" | "error"
  importResult?: ImportTagsResult
  importError?: string
}

const tagSorter = (t1: Tag, t2: Tag) => {
  if (t1.name < t2.name) {
    return -1
  } else if (t1.name > t2.name) {
    return 1
  }
  return 0
}

export default class ManageTagsPage extends AdminBasePage<ManageTagsPageProps, ManageTagsPageState> {
  public id = "p-admin-tags"
  public name = "tags"
  public title = "Tags"
  public subtitle = "Manage your site tags"

  constructor(props: ManageTagsPageProps) {
    super(props)
    this.state = {
      isAdding: false,
      allTags: this.props.tags,
      importStatus: "idle",
    }
  }

  private addNew = async () => {
    this.setState({
      isAdding: true,
      deleting: undefined,
      editing: undefined,
    })
  }

  private cancelAdd = () => {
    this.setState({ isAdding: false })
  }

  private saveNewTag = async (data: TagFormState): Promise<Failure | undefined> => {
    const result = await actions.createTag(data.name, data.color, data.isPublic)
    if (result.ok) {
      this.setState({
        isAdding: false,
        allTags: this.state.allTags.concat(result.data).sort(tagSorter),
      })
    } else {
      return result.error
    }
  }

  private handleTagDeleted = (tag: Tag) => {
    const idx = this.state.allTags.indexOf(tag)
    this.setState({
      allTags: this.state.allTags.splice(idx, 1) && this.state.allTags,
    })
  }

  private handleTagEdited = () => {
    this.setState({
      allTags: this.state.allTags.sort(tagSorter),
    })
  }

  private getTagList(filter: (tag: Tag) => boolean) {
    return this.state.allTags.filter(filter).map((t) => {
      return <TagListItem key={t.id} tag={t} onTagDeleted={this.handleTagDeleted} onTagEdited={this.handleTagEdited} />
    })
  }

  private handleImportFile = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return

    this.setState({ importStatus: "loading", importResult: undefined, importError: undefined })

    try {
      const text = await file.text()
      let parsed: Array<{ name: string; color: string; isPublic: boolean }>

      try {
        parsed = JSON.parse(text)
      } catch {
        this.setState({ importStatus: "error", importError: "Invalid JSON file." })
        return
      }

      if (!Array.isArray(parsed)) {
        this.setState({ importStatus: "error", importError: "JSON must be an array of tags." })
        return
      }

      const result = await actions.importTags(parsed)
      if (result.ok) {
        // Reload newly-created tags by merging them in (simplest: re-fetch via reload)
        this.setState({ importStatus: "success", importResult: result.data })
      } else {
        this.setState({ importStatus: "error", importError: "Import failed. Please check the file format." })
      }
    } catch {
      this.setState({ importStatus: "error", importError: "An unexpected error occurred." })
    }

    // Reset the file input so the same file can be selected again
    e.target.value = ""
  }

  public content() {
    const publicTaglist = this.getTagList((t) => t.isPublic)
    const privateTagList = this.getTagList((t) => !t.isPublic)

    const form =
      Fider.session.user.isAdministrator &&
      (this.state.isAdding ? (
        <TagForm onSave={this.saveNewTag} onCancel={this.cancelAdd} />
      ) : (
        <Button variant="secondary" onClick={this.addNew}>
          <Icon sprite={IconPlus} />
          <span>Add new</span>
        </Button>
      ))

    const { importStatus, importResult, importError } = this.state

    return (
      <VStack spacing={8}>
        <div>
          <h2 className="text-display">Public Tags</h2>
          <p className="text-muted">These tags are visible to all visitors.</p>
          <VStack spacing={4} divide={true}>
            {publicTaglist.length === 0 ? <p className="text-muted">There aren&apos;t any public tags yet.</p> : publicTaglist}
          </VStack>
        </div>
        <div>
          <h2 className="text-display">Private Tags</h2>
          <p className="text-muted">These tags are only visible for members of this site.</p>
          <VStack spacing={4} divide={true}>
            {privateTagList.length === 0 ? <p className="text-muted">There aren&apos;t any private tags yet.</p> : privateTagList}
          </VStack>
        </div>
        <div>{form}</div>

        {Fider.session.user.isAdministrator && (
          <div className="mt-4">
            <h2 className="text-display">Import &amp; Export Tags</h2>
            <p className="text-muted">
              Download all your tags as a JSON file, or upload a <code>tags.json</code> file to restore or add tags. Existing tags with the same name will be
              skipped.
            </p>
            <div className="flex gap-2">
              <Button variant="secondary" href="/admin/export/tags.json">
                <Icon sprite={IconDownload} />
                <span>Export tags.json</span>
              </Button>
              <label>
                <input
                  type="file"
                  accept=".json,application/json"
                  className="hidden"
                  onChange={this.handleImportFile}
                  disabled={importStatus === "loading"}
                  id="tags-import-input"
                />
                <Button variant="secondary" onClick={() => document.getElementById("tags-import-input")?.click()} disabled={importStatus === "loading"}>
                  <Icon sprite={IconUpload} />
                  <span>{importStatus === "loading" ? "Importing…" : "Import tags.json"}</span>
                </Button>
              </label>
            </div>

            {importStatus === "success" && importResult && (
              <div className="mt-2 text-green-700">
                ✓ Import complete — <strong>{importResult.created}</strong> tag(s) created, <strong>{importResult.skipped}</strong> skipped.
                {importResult.errors && importResult.errors.length > 0 && (
                  <ul className="mt-1 text-red-600 text-sm list-disc list-inside">
                    {importResult.errors.map((e, i) => (
                      <li key={i}>{e}</li>
                    ))}
                  </ul>
                )}
              </div>
            )}
            {importStatus === "error" && <div className="mt-2 text-red-600">✗ {importError}</div>}
          </div>
        )}
      </VStack>
    )
  }
}
