import React from "react"
import { Button, ButtonClickEvent, Form, Input, Select, SelectOption, Toggle } from "@fider/components"
import { HStack, VStack } from "@fider/components/layout"
import { Status, StatusKind } from "@fider/models"
import { actions, Failure } from "@fider/services"
import { i18n } from "@lingui/core"
import { Trans } from "@lingui/react/macro"
import { AdminBasePage } from "../components/AdminBasePage"

interface ManageStatusesPageProps {
  statuses: Status[]
}

interface ManageStatusesPageState {
  statuses: Status[]
  isAdding: boolean
  editingId: number | null
  draftSlug: string
  draftLabel: string
  draftKind: StatusKind
  draftColor: string
  draftIcon: string
  draftShowOnHome: boolean
  draftFilterable: boolean
  draftSortOrder: number
  error?: Failure
  busy: boolean
}

const kindOptions = (): SelectOption[] => [
  { value: "open", label: i18n._({ id: "admin.statuses.kind.open", message: "Open — default initial state" }) },
  { value: "active", label: i18n._({ id: "admin.statuses.kind.active", message: "Active — accepted, in progress" }) },
  { value: "closed-completed", label: i18n._({ id: "admin.statuses.kind.closedcompleted", message: "Closed (completed) — done, positive resolution" }) },
  { value: "closed-declined", label: i18n._({ id: "admin.statuses.kind.closeddeclined", message: "Closed (declined) — closed without action" }) },
  { value: "duplicate", label: i18n._({ id: "admin.statuses.kind.duplicate", message: "Duplicate — merged into another post" }) },
]

const colorOptions = (): SelectOption[] => [
  { value: "blue", label: i18n._({ id: "admin.statuses.color.blue", message: "Blue" }) },
  { value: "green", label: i18n._({ id: "admin.statuses.color.green", message: "Green" }) },
  { value: "yellow", label: i18n._({ id: "admin.statuses.color.yellow", message: "Yellow" }) },
  { value: "red", label: i18n._({ id: "admin.statuses.color.red", message: "Red" }) },
  { value: "gray", label: i18n._({ id: "admin.statuses.color.gray", message: "Gray" }) },
]

const slugify = (s: string): string =>
  s
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "")
    .slice(0, 50)

export default class ManageStatusesPage extends AdminBasePage<ManageStatusesPageProps, ManageStatusesPageState> {
  public id = "p-admin-statuses"
  public name = "statuses"
  public title = i18n._({ id: "admin.statuses.page.title", message: "Statuses" })
  public subtitle = i18n._({ id: "admin.statuses.page.subtitle", message: "Customize the list of post statuses your admins can apply" })

  constructor(props: ManageStatusesPageProps) {
    super(props)
    this.state = {
      statuses: props.statuses ?? [],
      isAdding: false,
      editingId: null,
      draftSlug: "",
      draftLabel: "",
      draftKind: "open",
      draftColor: "blue",
      draftIcon: "lightbulb",
      draftShowOnHome: true,
      draftFilterable: true,
      draftSortOrder: 100,
      busy: false,
    }
  }

  private openAdd = () => {
    this.setState({
      isAdding: true,
      editingId: null,
      draftSlug: "",
      draftLabel: "",
      draftKind: "open",
      draftColor: "blue",
      draftIcon: "lightbulb",
      draftShowOnHome: true,
      draftFilterable: true,
      draftSortOrder: 100,
      error: undefined,
    })
  }

  private openEdit = (status: Status) => {
    this.setState({
      isAdding: true,
      editingId: status.id,
      draftSlug: status.slug,
      draftLabel: status.label,
      draftKind: status.kind,
      draftColor: status.color,
      draftIcon: status.icon,
      draftShowOnHome: status.showOnHome,
      draftFilterable: status.filterable,
      draftSortOrder: status.sortOrder,
      error: undefined,
    })
  }

  private cancelAdd = () => {
    this.setState({ isAdding: false, editingId: null, error: undefined })
  }

  private updateLabel = (value: string) => {
    const auto = !this.state.draftSlug || this.state.draftSlug === slugify(this.state.draftLabel)
    this.setState({
      draftLabel: value,
      draftSlug: auto ? slugify(value) : this.state.draftSlug,
    })
  }

  private save = async (e: ButtonClickEvent) => {
    e.preventEnable()
    this.setState({ busy: true, error: undefined })

    if (this.state.editingId !== null) {
      const existing = this.state.statuses.find((s) => s.id === this.state.editingId)
      const result = await actions.updateStatus(this.state.editingId, {
        label: this.state.draftLabel,
        color: this.state.draftColor,
        icon: this.state.draftIcon,
        showOnHome: this.state.draftShowOnHome,
        filterable: this.state.draftFilterable,
        sortOrder: this.state.draftSortOrder,
        isActive: existing ? existing.isActive : true,
      })
      if (result.ok) {
        this.setState({
          statuses: this.state.statuses
            .map((s) =>
              s.id === this.state.editingId
                ? {
                    ...s,
                    label: this.state.draftLabel,
                    color: this.state.draftColor,
                    icon: this.state.draftIcon,
                    showOnHome: this.state.draftShowOnHome,
                    filterable: this.state.draftFilterable,
                    sortOrder: this.state.draftSortOrder,
                  }
                : s
            )
            .sort((a, b) => a.sortOrder - b.sortOrder),
          isAdding: false,
          editingId: null,
          busy: false,
        })
      } else {
        this.setState({ busy: false, error: result.error })
      }
      return
    }

    const result = await actions.createStatus({
      slug: this.state.draftSlug,
      label: this.state.draftLabel,
      kind: this.state.draftKind,
      color: this.state.draftColor,
      icon: this.state.draftIcon,
      showOnHome: this.state.draftShowOnHome,
      filterable: this.state.draftFilterable,
      sortOrder: this.state.draftSortOrder,
    })
    if (result.ok) {
      this.setState({
        statuses: [...this.state.statuses, result.data].sort((a, b) => a.sortOrder - b.sortOrder),
        isAdding: false,
        busy: false,
      })
    } else {
      this.setState({ busy: false, error: result.error })
    }
  }

  private toggleActive = async (status: Status) => {
    const result = await actions.updateStatus(status.id, {
      label: status.label,
      color: status.color,
      icon: status.icon,
      showOnHome: status.showOnHome,
      filterable: status.filterable,
      sortOrder: status.sortOrder,
      isActive: !status.isActive,
    })
    if (result.ok) {
      this.setState({
        statuses: this.state.statuses.map((s) => (s.id === status.id ? { ...s, isActive: !s.isActive } : s)),
      })
    }
  }

  private remove = async (status: Status) => {
    if (status.isSystem) return
    const prompt = i18n._({
      id: "admin.statuses.delete.confirm",
      values: { label: status.label },
      message: 'Delete the "{label}" status? Posts already using it must be reassigned first.',
    })
    if (!window.confirm(prompt)) return
    const result = await actions.deleteStatus(status.id)
    if (result.ok) {
      this.setState({ statuses: this.state.statuses.filter((s) => s.id !== status.id) })
    } else {
      window.alert(
        result.error?.errors?.[0]?.message ??
          i18n._({ id: "admin.statuses.delete.failed", message: "Could not delete this status." })
      )
    }
  }

  public content() {
    return (
      <VStack spacing={4}>
        <p className="text-sm text-muted">
          <Trans id="admin.statuses.help">
            The 6 built-in statuses are seeded for every site and can be renamed, recolored, or deactivated but not deleted.
            Add custom statuses for your own workflow — e.g. <em>Under Review</em> between Open and Planned.
          </Trans>
        </p>

        <table className="w-full text-sm">
          <thead>
            <tr className="text-left text-muted">
              <th className="py-2"><Trans id="admin.statuses.table.label">Label</Trans></th>
              <th><Trans id="admin.statuses.table.slug">Slug</Trans></th>
              <th><Trans id="admin.statuses.table.kind">Kind</Trans></th>
              <th><Trans id="admin.statuses.table.color">Color</Trans></th>
              <th><Trans id="admin.statuses.table.home">Home</Trans></th>
              <th><Trans id="admin.statuses.table.filter">Filter</Trans></th>
              <th><Trans id="admin.statuses.table.active">Active</Trans></th>
              <th />
            </tr>
          </thead>
          <tbody>
            {this.state.statuses.map((s) => (
              <tr key={s.id} className="border-t">
                <td className="py-2">
                  {s.label}{" "}
                  {s.isSystem && (
                    <span className="text-xs text-muted ml-1">
                      <Trans id="admin.statuses.table.system">(system)</Trans>
                    </span>
                  )}
                </td>
                <td><code>{s.slug}</code></td>
                <td>{s.kind}</td>
                <td>{s.color}</td>
                <td>{s.showOnHome ? i18n._({ id: "admin.statuses.table.yes", message: "yes" }) : "—"}</td>
                <td>{s.filterable ? i18n._({ id: "admin.statuses.table.yes", message: "yes" }) : "—"}</td>
                <td>
                  <Toggle active={s.isActive} onToggle={() => this.toggleActive(s)} />
                </td>
                <td>
                  <HStack spacing={2}>
                    <Button variant="tertiary" size="small" onClick={() => this.openEdit(s)}>
                      <Trans id="admin.statuses.action.edit">Edit</Trans>
                    </Button>
                    {!s.isSystem && (
                      <Button variant="danger" size="small" onClick={() => this.remove(s)}>
                        <Trans id="admin.statuses.action.delete">Delete</Trans>
                      </Button>
                    )}
                  </HStack>
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {this.state.isAdding ? (
          <Form error={this.state.error}>
            <VStack spacing={4}>
              <Input
                field="label"
                label={i18n._({ id: "admin.statuses.form.label", message: "Label" })}
                value={this.state.draftLabel}
                onChange={this.updateLabel}
              />
              {this.state.editingId === null ? (
                <>
                  <Input
                    field="slug"
                    label={i18n._({ id: "admin.statuses.form.slug", message: "Slug (URL-safe)" })}
                    value={this.state.draftSlug}
                    onChange={(v) => this.setState({ draftSlug: slugify(v) })}
                  />
                  <Select
                    field="kind"
                    label={i18n._({ id: "admin.statuses.form.kind", message: "Semantic kind" })}
                    defaultValue={this.state.draftKind}
                    options={kindOptions()}
                    onChange={(opt) => this.setState({ draftKind: (opt?.value ?? "open") as StatusKind })}
                  />
                </>
              ) : (
                <p className="text-sm text-muted">
                  <Trans id="admin.statuses.form.lockednote">
                    Slug <code>{this.state.draftSlug}</code> and semantic kind <code>{this.state.draftKind}</code> are fixed once a status is created.
                  </Trans>
                </p>
              )}
              <Select
                field="color"
                label={i18n._({ id: "admin.statuses.form.color", message: "Color" })}
                defaultValue={this.state.draftColor}
                options={colorOptions()}
                onChange={(opt) => this.setState({ draftColor: opt?.value ?? "blue" })}
              />
              <HStack spacing={4}>
                <Toggle active={this.state.draftShowOnHome} onToggle={(v) => this.setState({ draftShowOnHome: v })} />
                <span className="text-sm">
                  <Trans id="admin.statuses.form.showonhome">Show posts in this status on the home page</Trans>
                </span>
              </HStack>
              <HStack spacing={4}>
                <Toggle active={this.state.draftFilterable} onToggle={(v) => this.setState({ draftFilterable: v })} />
                <span className="text-sm">
                  <Trans id="admin.statuses.form.filterable">Include this status in the home-page filter</Trans>
                </span>
              </HStack>
              <HStack spacing={2}>
                <Button variant="primary" onClick={this.save} disabled={this.state.busy}>
                  {this.state.editingId === null ? (
                    <Trans id="admin.statuses.form.save">Save status</Trans>
                  ) : (
                    <Trans id="admin.statuses.form.update">Update status</Trans>
                  )}
                </Button>
                <Button variant="tertiary" onClick={this.cancelAdd}>
                  <Trans id="action.cancel">Cancel</Trans>
                </Button>
              </HStack>
            </VStack>
          </Form>
        ) : (
          <Button variant="primary" onClick={this.openAdd}>
            <Trans id="admin.statuses.add">Add a custom status</Trans>
          </Button>
        )}
      </VStack>
    )
  }
}
