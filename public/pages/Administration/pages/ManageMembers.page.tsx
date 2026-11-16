import React, { useState, useEffect, useCallback } from "react"
import { Input, Avatar, Icon, Dropdown, Pagination } from "@fider/components"
import { User, UserRole, UserStatus } from "@fider/models"
import IconSearch from "@fider/assets/images/heroicons-search.svg"
import IconX from "@fider/assets/images/heroicons-x.svg"
import IconDotsHorizontal from "@fider/assets/images/heroicons-dots-horizontal.svg"
import IconCheck from "@fider/assets/images/heroicons-check.svg"
import HeroIconFilter from "@fider/assets/images/heroicons-filter.svg"
import { actions, Fider, Result } from "@fider/services"
import { AdminPageContainer } from "../components/AdminBasePage"
import { HStack, VStack } from "@fider/components/layout"

interface ManageMembersPageProps {
  users: User[]
  totalPages: number
}

type Level = "administrator" | "collaborator" | "trusted" | "member" | "blocked"

const levels: { level: Level; label: string; plural: string; badgeClass: string }[] = [
  { level: "administrator", label: "Administrator", plural: "Administrators", badgeClass: "bg-blue-100 text-blue-800" },
  { level: "collaborator", label: "Collaborator", plural: "Collaborators", badgeClass: "bg-green-100 text-green-800" },
  { level: "trusted", label: "Trusted member", plural: "Trusted members", badgeClass: "bg-green-100 text-green-800" },
  { level: "member", label: "Member", plural: "Members", badgeClass: "text-gray-600" },
  { level: "blocked", label: "Blocked", plural: "Blocked", badgeClass: "bg-red-100 text-red-800" },
]

const getLevel = (user: User): Level => {
  if (user.status === UserStatus.Blocked) return "blocked"
  if (user.role === UserRole.Administrator) return "administrator"
  if (user.role === UserRole.Collaborator) return "collaborator"
  return user.isTrusted ? "trusted" : "member"
}

interface UserListItemProps {
  user: User
  onChangeLevel: (user: User, level: Level) => Promise<void>
  isLast?: boolean
}

const UserListItem = (props: UserListItemProps) => {
  const current = getLevel(props.user)
  const badge = levels.find((l) => l.level === current)

  return (
    <div
      className={`border-b border-gray-200 grid gap-4 py-4 px-4 flex-items-center bg-white hover ${props.isLast ? "rounded-md-b" : ""}`}
      style={{ gridTemplateColumns: "minmax(200px, 1fr) minmax(280px, 2fr) minmax(120px, 150px) 100px" }}
    >
      <HStack>
        <Avatar user={props.user} />
        <div className="text-subtitle">{props.user.name}</div>
      </HStack>

      <div className="text-muted nowrap" title={props.user.email}>
        {props.user.email || "No email"}
      </div>

      <div>{badge && <span className={`text-xs px-2 py-1 rounded ${badge.badgeClass}`}>{badge.label.toLowerCase()}</span>}</div>

      <div className="flex justify-end relative">
        {Fider.session.user.id !== props.user.id && Fider.session.user.isAdministrator && (
          <div className="relative z-10">
            <Dropdown renderHandle={<Icon sprite={IconDotsHorizontal} width="16" height="16" />}>
              {levels.map(({ level, label }) => (
                <Dropdown.ListItem key={level} onClick={level === current ? undefined : () => props.onChangeLevel(props.user, level)}>
                  {level === current ? <Icon sprite={IconCheck} className="mr-2" width="16" height="16" /> : <span className="w-4 mr-2" />}
                  {label}
                </Dropdown.ListItem>
              ))}
            </Dropdown>
          </div>
        )}
      </div>
    </div>
  )
}

export default function ManageMembersPage(props: ManageMembersPageProps) {
  const [query, setQuery] = useState("")
  const [roleFilter, setRoleFilter] = useState<Level | "all">("all")
  const [users, setUsers] = useState<User[]>(props.users)
  const [currentPage, setCurrentPage] = useState(1)
  const [totalPages, setTotalPages] = useState(props.totalPages)
  const [searchTimeoutId, setSearchTimeoutId] = useState<number | undefined>(undefined)
  const pageSize = 10

  // Initialize state from URL parameters and load first page
  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search)
    const initialQuery = urlParams.get("query") || ""
    const initialRoleFilter = (urlParams.get("roles") as Level) || "all"
    const initialPage = parseInt(urlParams.get("page") || "1")

    setQuery(initialQuery)
    setRoleFilter(initialRoleFilter)
    setCurrentPage(initialPage)
  }, [])

  const reloadUsers = useCallback(
    async (searchQuery: string, roleFilterValue: Level | "all", page = 1) => {
      const params = new URLSearchParams()
      if (searchQuery) {
        params.append("query", searchQuery)
      }
      if (roleFilterValue !== "all") {
        params.append("roles", roleFilterValue)
      }
      params.append("page", page.toString())
      params.append("limit", pageSize.toString())

      const response = await fetch(`/api/v1/users?${params.toString()}`)
      if (response.ok) {
        const data = await response.json()
        setUsers(data.users)
        setTotalPages(data.totalPages)
        setCurrentPage(page)
      }
    },
    [pageSize]
  )

  const handleSearchFilterChanged = useCallback(
    (newQuery: string) => {
      setQuery(newQuery)

      // Debounce the API call for search
      if (searchTimeoutId) {
        clearTimeout(searchTimeoutId)
      }

      const timeoutId = window.setTimeout(() => {
        reloadUsers(newQuery, roleFilter, 1) // Reset to page 1 when searching
      }, 300)

      setSearchTimeoutId(timeoutId)
    },
    [roleFilter, reloadUsers, searchTimeoutId]
  )

  const handleRoleFilterChanged = useCallback(
    (newRoleFilter: Level | "all") => {
      setRoleFilter(newRoleFilter)
      reloadUsers(query, newRoleFilter, 1) // Reset to page 1 when changing filter
    },
    [query, reloadUsers]
  )

  const clearSearch = useCallback(() => {
    if (searchTimeoutId) {
      clearTimeout(searchTimeoutId)
    }
    setQuery("")
    reloadUsers("", roleFilter, 1)
  }, [roleFilter, reloadUsers, searchTimeoutId])

  const handlePageChange = useCallback(
    (page: number) => {
      reloadUsers(query, roleFilter, page)
    },
    [query, roleFilter, reloadUsers]
  )

  const handleChangeLevel = useCallback(
    async (user: User, level: Level) => {
      // Blocked users are demoted to member so they drop out of staff lists
      const role = level === "administrator" ? UserRole.Administrator : level === "collaborator" ? UserRole.Collaborator : UserRole.Visitor
      const status = level === "blocked" ? UserStatus.Blocked : UserStatus.Active
      const isTrusted = level === "member" ? false : level === "trusted" ? true : user.isTrusted

      const steps: (() => Promise<Result>)[] = []
      if (user.role !== role) steps.push(() => actions.changeUserRole(user.id, role))
      if (user.status !== status) steps.push(() => (status === UserStatus.Blocked ? actions.blockUser(user.id) : actions.unblockUser(user.id)))
      if (user.isTrusted !== isTrusted) steps.push(() => (isTrusted ? actions.trustUser(user.id) : actions.untrustUser(user.id)))

      for (const step of steps) {
        const result = await step()
        if (!result.ok) {
          // Earlier steps may have succeeded, so show the user's real state
          await reloadUsers(query, roleFilter, currentPage)
          return
        }
      }

      setUsers((prev) => prev.map((u) => (u.id === user.id ? { ...user, role, status, isTrusted } : u)))
    },
    [query, roleFilter, currentPage, reloadUsers]
  )

  return (
    <AdminPageContainer id="p-admin-members" name="users" title="Members" subtitle="Manage your site administrators and collaborators">
      <div className="flex gap-4 flex-items-center mb-4">
        <div className="flex-grow">
          <Input
            field="query"
            icon={query ? IconX : IconSearch}
            onIconClick={query ? clearSearch : undefined}
            placeholder="Search by name / email ..."
            value={query}
            onChange={handleSearchFilterChanged}
          />
        </div>
        <Dropdown
          renderHandle={
            <div className="flex flex-items-center h-10 text-medium text-xs rounded-md uppercase border border-gray-400 text-gray-800 p-2 px-3 hover">
              <Icon sprite={HeroIconFilter} className="h-5 pr-1" />
              Role
              {roleFilter !== "all" && <div className="bg-gray-200 inline-block rounded-full px-2 py-1 w-min-4 text-2xs text-center ml-2">1</div>}
            </div>
          }
        >
          <Dropdown.ListItem onClick={() => handleRoleFilterChanged("all")}>
            <span className={roleFilter === "all" ? "text-semibold" : ""}>All Roles</span>
          </Dropdown.ListItem>
          {levels.map(({ level, plural }) => (
            <Dropdown.ListItem key={level} onClick={() => handleRoleFilterChanged(level)}>
              <span className={roleFilter === level ? "text-semibold" : ""}>{plural}</span>
            </Dropdown.ListItem>
          ))}
        </Dropdown>
      </div>

      <VStack className="rounded-md border border-gray-200 relative">
        <div
          className="grid rounded-md-t gap-4 py-3 px-4 bg-gray-100 text-category"
          style={{ gridTemplateColumns: "minmax(200px, 1fr) minmax(280px, 2fr) minmax(120px, 150px) 100px" }}
        >
          <div>Name</div>
          <div>Email</div>
          <div>Role</div>
        </div>
        <div>
          {users.map((user, index) => (
            <UserListItem key={user.id} user={user} onChangeLevel={handleChangeLevel} isLast={index === users.length - 1} />
          ))}
        </div>
      </VStack>

      <div className="pt-4">
        <Pagination currentPage={currentPage} totalPages={totalPages} onPageChange={handlePageChange} />
      </div>

      <ul className="text-muted">
        <li>
          <strong>Administrators</strong> have full access to edit and manage content, permissions and all site settings.
        </li>
        <li>
          <strong>Collaborators</strong> can edit and manage content, but not permissions and settings.
        </li>
        <li>
          <strong>Trusted members</strong> won&apos;t need to have their content moderated if moderation is enabled.
        </li>
        <li>
          <strong>Blocked</strong> users are unable to sign into this site.
        </li>
      </ul>
    </AdminPageContainer>
  )
}
