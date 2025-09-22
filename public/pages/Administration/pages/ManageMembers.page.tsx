import React, { useState, useEffect, useCallback } from "react"
import { Input, Avatar, Icon, Dropdown, Button } from "@fider/components"
import { User, UserRole, UserStatus } from "@fider/models"
import IconSearch from "@fider/assets/images/heroicons-search.svg"
import IconX from "@fider/assets/images/heroicons-x.svg"
import IconDotsHorizontal from "@fider/assets/images/heroicons-dots-horizontal.svg"
import HeroIconFilter from "@fider/assets/images/heroicons-filter.svg"
import { actions, Fider } from "@fider/services"
import { AdminPageContainer } from "../components/AdminBasePage"
import { HStack, VStack } from "@fider/components/layout"

interface ManageMembersPageProps {
  users: User[]
}

interface UserListItemProps {
  user: User
  onAction: (actionName: string, user: User) => Promise<void>
}

interface UserListItemExtendedProps extends UserListItemProps {
  isLast?: boolean
}

const UserListItem = (props: UserListItemExtendedProps) => {
  const admin = props.user.role === UserRole.Administrator && <span className="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded">administrator</span>
  const collaborator = props.user.role === UserRole.Collaborator && <span className="text-xs bg-green-100 text-green-800 px-2 py-1 rounded">collaborator</span>
  const blocked = props.user.status === UserStatus.Blocked && <span className="text-xs bg-red-100 text-red-800 px-2 py-1 rounded">blocked</span>
  const isVisitor = props.user.role === UserRole.Visitor

  const actionSelected = (actionName: string) => () => {
    props.onAction(actionName, props.user)
  }

  return (
    <div
      className={`border-b border-gray-200 grid gap-4 py-4 px-4 flex-items-center bg-white hover:bg-gray-100 ${props.isLast ? "rounded-md-b" : ""}`}
      style={{ gridTemplateColumns: "minmax(200px, 1fr) minmax(280px, 2fr) minmax(120px, 150px) 100px" }}
    >
      <HStack>
        <Avatar user={props.user} />
        <div className="text-subtitle">{props.user.name}</div>
      </HStack>

      <div className="text-muted text-ellipsis" title={props.user.email}>
        {props.user.email || "No email"}
      </div>

      <div>
        {admin} {collaborator} {blocked}
        {isVisitor && !blocked && <span className="text-xs text-gray-600">visitor</span>}
      </div>

      <div className="flex justify-end relative">
        {Fider.session.user.id !== props.user.id && Fider.session.user.isAdministrator && (
          <div className="relative z-10">
            <Dropdown renderHandle={<Icon sprite={IconDotsHorizontal} width="16" height="16" />}>
              {!blocked && (!!collaborator || isVisitor) && (
                <Dropdown.ListItem onClick={actionSelected("to-administrator")}>Promote to Administrator</Dropdown.ListItem>
              )}
              {!blocked && (!!admin || isVisitor) && <Dropdown.ListItem onClick={actionSelected("to-collaborator")}>Promote to Collaborator</Dropdown.ListItem>}
              {!blocked && (!!collaborator || !!admin) && <Dropdown.ListItem onClick={actionSelected("to-visitor")}>Demote to Visitor</Dropdown.ListItem>}
              {isVisitor && !blocked && <Dropdown.ListItem onClick={actionSelected("block")}>Block User</Dropdown.ListItem>}
              {isVisitor && !!blocked && <Dropdown.ListItem onClick={actionSelected("unblock")}>Unblock User</Dropdown.ListItem>}
            </Dropdown>
          </div>
        )}
      </div>
    </div>
  )
}

export default function ManageMembersPage(props: ManageMembersPageProps) {
  const [query, setQuery] = useState("")
  const [roleFilter, setRoleFilter] = useState<UserRole | "all">("all")
  const [users, setUsers] = useState<User[]>([])
  const [visibleUsers, setVisibleUsers] = useState<User[]>([])
  const [searchTimeoutId, setSearchTimeoutId] = useState<number | undefined>(undefined)

  // Initialize state from URL parameters and props
  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search)
    const initialQuery = urlParams.get("query") || ""
    const initialRoleFilter = (urlParams.get("roles") as UserRole) || "all"

    setQuery(initialQuery)
    setRoleFilter(initialRoleFilter)

    const sortedUsers = props.users.sort(sortByStaff)
    setUsers(sortedUsers)
    setVisibleUsers(sortedUsers.slice(0, 10))
  }, [props.users])

  const sortByStaff = (left: User, right: User) => {
    if (right.role === left.role) {
      if (left.name < right.name) {
        return -1
      } else if (left.name > right.name) {
        return 1
      }
      return 0
    }

    if (right.role !== UserRole.Visitor) {
      return 1
    }
    return -1
  }

  const reloadUsers = useCallback(async (searchQuery: string, roleFilterValue: UserRole | "all") => {
    const params = new URLSearchParams()
    if (searchQuery) {
      params.append("query", searchQuery)
    }
    if (roleFilterValue !== "all") {
      params.append("roles", roleFilterValue.toString())
    }

    const response = await fetch(`/api/v1/users${params.toString() ? "?" + params.toString() : ""}`)
    if (response.ok) {
      const usersData = await response.json()
      const sortedUsers = usersData.sort(sortByStaff)
      setUsers(sortedUsers)
      setVisibleUsers(sortedUsers.slice(0, 10))
    }
  }, [])

  const handleSearchFilterChanged = useCallback(
    (newQuery: string) => {
      setQuery(newQuery)

      // Debounce the API call for search
      if (searchTimeoutId) {
        clearTimeout(searchTimeoutId)
      }

      const timeoutId = window.setTimeout(() => {
        reloadUsers(newQuery, roleFilter)
      }, 300)

      setSearchTimeoutId(timeoutId)
    },
    [roleFilter, reloadUsers, searchTimeoutId]
  )

  const handleRoleFilterChanged = useCallback(
    (newRoleFilter: UserRole | "all") => {
      setRoleFilter(newRoleFilter)
      reloadUsers(query, newRoleFilter)
    },
    [query, reloadUsers]
  )

  const clearSearch = useCallback(() => {
    if (searchTimeoutId) {
      clearTimeout(searchTimeoutId)
    }
    setQuery("")
    reloadUsers("", roleFilter)
  }, [roleFilter, reloadUsers, searchTimeoutId])

  const showMore = useCallback(() => {
    setVisibleUsers(users.slice(0, visibleUsers.length + 10))
  }, [users, visibleUsers.length])

  const handleAction = useCallback(
    async (actionName: string, user: User) => {
      const changeRole = async (role: UserRole) => {
        const result = await actions.changeUserRole(user.id, role)
        if (result.ok) {
          user.role = role
          // Update the user in current state without full reload
          const updatedUsers = users.map((u) => (u.id === user.id ? user : u))
          const sortedUsers = updatedUsers.sort(sortByStaff)
          setUsers(sortedUsers)
          setVisibleUsers(sortedUsers.slice(0, visibleUsers.length))
        }
      }

      const changeStatus = async (status: UserStatus) => {
        const action = status === UserStatus.Blocked ? actions.blockUser : actions.unblockUser
        const result = await action(user.id)
        if (result.ok) {
          user.status = status
          // Update the user in current state without full reload
          const updatedUsers = users.map((u) => (u.id === user.id ? user : u))
          setUsers(updatedUsers)
          setVisibleUsers(updatedUsers.slice(0, visibleUsers.length))
        }
      }

      if (actionName === "to-collaborator") {
        await changeRole(UserRole.Collaborator)
      } else if (actionName === "to-visitor") {
        await changeRole(UserRole.Visitor)
      } else if (actionName === "to-administrator") {
        await changeRole(UserRole.Administrator)
      } else if (actionName === "block") {
        await changeStatus(UserStatus.Blocked)
      } else if (actionName === "unblock") {
        await changeStatus(UserStatus.Active)
      }
    },
    [users, visibleUsers.length]
  )

  return (
    <AdminPageContainer id="p-admin-members" name="members" title="Members" subtitle="Manage your site administrators and collaborators">
      <div className="flex gap-4 flex-items-center mb-4">
        <div className="flex-grow">
          <Input
            field="query"
            icon={query ? IconX : IconSearch}
            onIconClick={query ? clearSearch : undefined}
            placeholder="Search for users by name / email ..."
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
          <Dropdown.ListItem onClick={() => handleRoleFilterChanged(UserRole.Administrator)}>
            <span className={roleFilter === UserRole.Administrator ? "text-semibold" : ""}>Administrators</span>
          </Dropdown.ListItem>
          <Dropdown.ListItem onClick={() => handleRoleFilterChanged(UserRole.Collaborator)}>
            <span className={roleFilter === UserRole.Collaborator ? "text-semibold" : ""}>Collaborators</span>
          </Dropdown.ListItem>
          <Dropdown.ListItem onClick={() => handleRoleFilterChanged(UserRole.Visitor)}>
            <span className={roleFilter === UserRole.Visitor ? "text-semibold" : ""}>Visitors</span>
          </Dropdown.ListItem>
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
          {visibleUsers.map((user, index) => (
            <UserListItem key={user.id} user={user} onAction={handleAction} isLast={index === visibleUsers.length - 1} />
          ))}
        </div>
      </VStack>

      <p className="text-muted pt-4">
        {!query && roleFilter === "all" && (
          <>
            Showing {visibleUsers.length} of {users.length} registered users.
          </>
        )}
        {(query || roleFilter !== "all") && (
          <>
            Showing {visibleUsers.length} of {users.length} users
            {query && <> matching &apos;{query}&apos;</>}
            {roleFilter !== "all" && <> with role &apos;{roleFilter}&apos;</>}.
          </>
        )}
        {visibleUsers.length < users.length && (
          <Button variant="tertiary" onClick={showMore}>
            view more
          </Button>
        )}
      </p>

      <ul className="text-muted">
        <li>
          <strong>Administrators</strong> have full access to edit and manage content, permissions and all site settings.
        </li>
        <li>
          <strong>Collaborators</strong> can edit and manage content, but not permissions and settings.
        </li>
        <li>
          <strong>Blocked</strong> users are unable to sign into this site.
        </li>
      </ul>
    </AdminPageContainer>
  )
}
