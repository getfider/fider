import { PostStatus, Status, StatusKind } from "./post"
import { Tenant } from "./identity"

// resolveStatus returns the tenant's Status entry for the given slug if the
// tenant has a status catalogue and a matching active row; otherwise null.
// All callers should pair this with a fallback to the hardcoded PostStatus
// enum so pre-migration tenants (or tenants whose admin hasn't seeded a
// matching row) still render correctly.
export const resolveStatus = (tenant: Tenant | null | undefined, slug: string): Status | null => {
  if (!tenant || !tenant.statuses || tenant.statuses.length === 0) return null
  for (const s of tenant.statuses) {
    if (s.slug === slug && s.isActive) return s
  }
  return null
}

// isClosedKind reports whether a kind is one of the closed semantic kinds.
// Useful for filter UIs that group active vs. closed statuses.
export const isClosedKind = (kind: StatusKind): boolean =>
  kind === "closed-completed" || kind === "closed-declined" || kind === "duplicate"

// statusListFor returns the runtime status list to render. Prefers the
// tenant catalogue; falls back to PostStatus.All for tenants without one.
export const statusListFor = (tenant: Tenant | null | undefined): Array<{ value: string; label: string; closed: boolean; filterable: boolean }> => {
  if (tenant && tenant.statuses && tenant.statuses.length > 0) {
    return tenant.statuses
      .filter((s) => s.isActive)
      .sort((a, b) => a.sortOrder - b.sortOrder)
      .map((s) => ({ value: s.slug, label: s.label, closed: isClosedKind(s.kind), filterable: s.filterable }))
  }
  return PostStatus.All.map((p) => ({ value: p.value, label: p.title, closed: p.closed, filterable: p.filterable }))
}

// The 6 built-in slugs have entries in locale/<lang>/client.json under
// enum.poststatus.<slug>. Custom slugs don't, so the i18n lookup returns
// the literal key string. Pages should consult this set to decide whether
// to call i18n._() for translation or fall straight through to the
// tenant-defined label.
const BUILT_IN_STATUS_SLUGS = new Set<string>(["open", "planned", "started", "completed", "declined", "duplicate"])

// statusLabel resolves the user-visible label for a status row in the
// runtime list. For built-in slugs it consults the i18n catalog so locale
// overrides still apply; for custom slugs it returns the tenant-defined
// label verbatim (the locale catalog has no entry to translate from).
export const statusLabel = (item: { value: string; label: string }, translate: (id: string, fallback: string) => string): string => {
  if (BUILT_IN_STATUS_SLUGS.has(item.value)) {
    return translate(`enum.poststatus.${item.value}`, item.label)
  }
  return item.label
}
