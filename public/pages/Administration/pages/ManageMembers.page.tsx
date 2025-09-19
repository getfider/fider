import React from "react"
import { Input, Avatar, Icon, Dropdown, Button } from "@fider/components"
import { User, UserRole, UserStatus } from "@fider/models"
import { AdminBasePage } from "../components/AdminBasePage"
import IconSearch from "@fider/assets/images/heroicons-search.svg"
import IconX from "@fider/assets/images/heroicons-x.svg"
import IconDotsHorizontal from "@fider/assets/images/heroicons-dots-horizontal.svg"
import { actions, Fider } from "@fider/services"
import { HStack } from "@fider/components/layout"

interface ManageMembersPageState {
  query: string
  users: User[]
  visibleUsers: User[]
}

interface ManageMembersPageProps {
  users: User[]
}

interface UserListItemProps {
  user: User
  onAction: (actionName: string, user: User) => Promise<void>
  isEven: boolean
}

const UserListItem = (props: UserListItemProps) => {
  const admin = props.user.role === UserRole.Administrator && <span className="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded">administrator</span>
  const collaborator = props.user.role === UserRole.Collaborator && <span className="text-xs bg-green-100 text-green-800 px-2 py-1 rounded">collaborator</span>
  const blocked = props.user.status === UserStatus.Blocked && <span className="text-xs bg-red-100 text-red-800 px-2 py-1 rounded">blocked</span>
  const isVisitor = props.user.role === UserRole.Visitor

  const actionSelected = (actionName: string) => () => {
    props.onAction(actionName, props.user)
  }

  return (
    <div
      className={`grid gap-4 py-4 px-4 flex-items-center ${props.isEven ? "bg-gray-50" : "bg-white"} hover:bg-gray-100`}
      style={{ gridTemplateColumns: "2fr 2fr 1fr 100px" }}
    >
      <HStack>
        <Avatar user={props.user} />
        <div className="text-subtitle">{props.user.name}</div>
      </HStack>

      <div className="text-muted">{props.user.email || "No email"}</div>

      <div>
        {admin} {collaborator} {blocked}
        {isVisitor && !blocked && <span className="text-xs text-gray-600">visitor</span>}
      </div>

      <div className="flex justify-end">
        {Fider.session.user.id !== props.user.id && Fider.session.user.isAdministrator && (
          <Dropdown renderHandle={<Icon sprite={IconDotsHorizontal} width="16" height="16" />}>
            {!blocked && (!!collaborator || isVisitor) && (
              <Dropdown.ListItem onClick={actionSelected("to-administrator")}>Promote to Administrator</Dropdown.ListItem>
            )}
            {!blocked && (!!admin || isVisitor) && <Dropdown.ListItem onClick={actionSelected("to-collaborator")}>Promote to Collaborator</Dropdown.ListItem>}
            {!blocked && (!!collaborator || !!admin) && <Dropdown.ListItem onClick={actionSelected("to-visitor")}>Demote to Visitor</Dropdown.ListItem>}
            {isVisitor && !blocked && <Dropdown.ListItem onClick={actionSelected("block")}>Block User</Dropdown.ListItem>}
            {isVisitor && !!blocked && <Dropdown.ListItem onClick={actionSelected("unblock")}>Unblock User</Dropdown.ListItem>}
          </Dropdown>
        )}
      </div>
    </div>
  )
}

export default class ManageMembersPage extends AdminBasePage<ManageMembersPageProps, ManageMembersPageState> {
  public id = "p-admin-members"
  public name = "members"
  public title = "Members"
  public subtitle = "Manage your site administrators and collaborators"

  constructor(props: ManageMembersPageProps) {
    super(props)

    const users = this.props.users.sort(this.sortByStaff)

    this.state = {
      query: "",
      users,
      visibleUsers: users.slice(0, 10),
    }
  }

  private showMore = (): void => {
    this.setState({
      visibleUsers: this.state.users.slice(0, this.state.visibleUsers.length + 10),
    })
  }

  private clearSearch = () => {
    this.handleSearchFilterChanged("")
  }

  private memberFilter = (query: string, user: User): boolean => {
    return user.name.toLowerCase().indexOf(query.toLowerCase()) >= 0 || (user.email && user.email.toLowerCase().indexOf(query.toLowerCase()) >= 0) || false
  }

  private handleSearchFilterChanged = (query: string) => {
    const users = this.props.users.filter((x) => this.memberFilter(query, x)).sort(this.sortByStaff)
    this.setState({ query, users, visibleUsers: users.slice(0, 10) })
  }

  private handleAction = async (actionName: string, user: User) => {
    const changeRole = async (role: UserRole) => {
      const result = await actions.changeUserRole(user.id, role)
      if (result.ok) {
        user.role = role
      }
      this.handleSearchFilterChanged(this.state.query)
    }

    const changeStatus = async (status: UserStatus) => {
      const action = status === UserStatus.Blocked ? actions.blockUser : actions.unblockUser
      const result = await action(user.id)
      if (result.ok) {
        user.status = status
      }
      this.forceUpdate()
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
  }

  private sortByStaff = (left: User, right: User) => {
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

  public content() {
    return (
      <>
        <Input
          field="query"
          icon={this.state.query ? IconX : IconSearch}
          onIconClick={this.state.query ? this.clearSearch : undefined}
          placeholder="Search for users by name / email ..."
          value={this.state.query}
          onChange={this.handleSearchFilterChanged}
        />
        <div className="shadow rounded border overflow-hidden">
          <div className="grid gap-4 py-3 px-4 bg-gray-100 text-category border-b" style={{ gridTemplateColumns: "2fr 2fr 1fr 100px" }}>
            <div>User</div>
            <div>Email</div>
            <div>Role</div>
            <div className="text-center">Actions</div>
          </div>
          <div>
            {this.state.visibleUsers.map((user, index) => (
              <UserListItem key={user.id} user={user} onAction={this.handleAction} isEven={index % 2 === 0} />
            ))}
          </div>
        </div>
        <p className="text-muted pt-4">
          {!this.state.query && (
            <>
              Showing {this.state.visibleUsers.length} of {this.state.users.length} registered users.
            </>
          )}
          {this.state.query && (
            <>
              Showing {this.state.visibleUsers.length} of {this.state.users.length} users matching &apos;{this.state.query}&apos;.
            </>
          )}
          {this.state.visibleUsers.length < this.state.users.length && (
            <Button variant="tertiary" onClick={this.showMore}>
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
      </>
    )
  }
}
