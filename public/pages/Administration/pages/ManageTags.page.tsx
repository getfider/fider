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
  if (t1.isPublic !== t2.isPublic) {
    return t1.isPublic ? -1 : 1
  }
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

  private sortedTags(): Tag[] {
    return this.state.allTags.slice().sort(tagSorter)
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
    const tags = this.sortedTags()
    const gridTemplateColumns = "minmax(200px, 1fr) minmax(100px, 150px) 200px"
    const canAdd = Fider.session.user.isAdministrator
    const lastTagIsLast = !canAdd

    const addRow = canAdd && (
      <div className="py-3 px-4 bg-white rounded-md-b">
        {this.state.isAdding ? (
          <TagForm onSave={this.saveNewTag} onCancel={this.cancelAdd} />
        ) : (
          <Button variant="tertiary" onClick={this.addNew}>
            <Icon sprite={IconPlus} />
            <span>Add new tag</span>
          </Button>
        )}
      </div>
    )

    const { importStatus, importResult, importError } = this.state

    return (
      <VStack spacing={8}>
        <VStack className="rounded-md border border-gray-200 relative">
          <div className="grid rounded-md-t gap-4 py-3 px-4 bg-gray-100 text-category" style={{ gridTemplateColumns }}>
            <div>Tag</div>
            <div>Visibility</div>
            <div></div>
          </div>
          <div>
            {tags.length === 0 ? (
              <div className={`py-4 px-4 bg-white text-muted ${lastTagIsLast ? "rounded-md-b" : "border-b border-gray-200"}`}>
                There aren&apos;t any tags yet.
              </div>
            ) : (
              tags.map((tag, index) => (
                <TagListItem
                  key={tag.id}
                  tag={tag}
                  gridTemplateColumns={gridTemplateColumns}
                  isLast={lastTagIsLast && index === tags.length - 1}
                  onTagDeleted={this.handleTagDeleted}
                  onTagEdited={this.handleTagEdited}
                />
              ))
            )}
          </div>
          {addRow}
        </VStack>

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
