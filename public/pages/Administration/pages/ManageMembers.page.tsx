import React from "react"
import { Input, Avatar, UserName, Icon, Dropdown, Button } from "@fider/components"
import { User, UserRole, UserStatus } from "@fider/models"
import { AdminBasePage } from "../components/AdminBasePage"
import IconSearch from "@fider/assets/images/heroicons-search.svg"
import IconX from "@fider/assets/images/heroicons-x.svg"
import IconDotsHorizontal from "@fider/assets/images/heroicons-dots-horizontal.svg"
import { actions, Fider } from "@fider/services"
import { HStack, VStack } from "@fider/components/layout"

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
}

const UserListItem = (props: UserListItemProps) => {
  const admin = props.user.role === UserRole.Administrator && <span>administrator</span>
  const collaborator = props.user.role === UserRole.Collaborator && <span>collaborator</span>
  const moderator = props.user.role === UserRole.Moderator && <span>moderator</span>
  const blocked = props.user.status === UserStatus.Blocked && <span className="text-red-700">blocked</span>
  const isVisitor = props.user.role === UserRole.Visitor

  const actionSelected = (actionName: string) => () => {
    props.onAction(actionName, props.user)
  }

  return (
    <HStack spacing={4}>
      <HStack spacing={4}>
        <Avatar user={props.user} />
        <VStack spacing={0}>
          <UserName user={props.user} showEmail={true} />
          <span className="text-muted">
            {admin} {moderator} {collaborator} {blocked}
          </span>
        </VStack>
      </HStack>
      {Fider.session.user.id !== props.user.id && Fider.session.user.isAdministrator && (
        <Dropdown renderHandle={<Icon sprite={IconDotsHorizontal} width="16" height="16" />}>
        {!blocked && (props.user.role !== UserRole.Administrator) && (
          <Dropdown.ListItem onClick={actionSelected("to-administrator")}>
            {props.user.role === UserRole.Moderator ? "Promote to Administrator" : "Promote to Administrator"}
          </Dropdown.ListItem>
        )}
        {!blocked && (props.user.role !== UserRole.Moderator) && (
          <Dropdown.ListItem onClick={actionSelected("to-moderator")}>
            {props.user.role === UserRole.Administrator ? "Demote to Moderator" : "Promote to Moderator"}
          </Dropdown.ListItem>
        )}
        {!blocked && (props.user.role !== UserRole.Collaborator) && (
          <Dropdown.ListItem onClick={actionSelected("to-collaborator")}>
            {props.user.role === UserRole.Administrator ? "Demote to Collaborator" : "Promote to Collaborator"}
          </Dropdown.ListItem>
        )}
        {!blocked && (props.user.role !== UserRole.Visitor) && (
          <Dropdown.ListItem onClick={actionSelected("to-visitor")}>
            Demote to Visitor
          </Dropdown.ListItem>
        )}
        {isVisitor && !blocked && (
          <Dropdown.ListItem onClick={actionSelected("block")}>Block User</Dropdown.ListItem>
        )}
        {isVisitor && !!blocked && (
          <Dropdown.ListItem onClick={actionSelected("unblock")}>Unblock User</Dropdown.ListItem>
        )}
      </Dropdown>
      
      )}
    </HStack>
  )
}

export default class ManageMembersPage extends AdminBasePage<ManageMembersPageProps, ManageMembersPageState> {
  public id = "p-admin-members"
  public name = "members"
  public title = "Members"
  public subtitle = "Manage your site administrators and collaborators"

  constructor(props: ManageMembersPageProps) {
    super(props)

    // Sort the users using our custom sort order (see sortByStaff below)
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
    return (
      user.name.toLowerCase().indexOf(query.toLowerCase()) >= 0 ||
      (user.email && user.email.toLowerCase().indexOf(query.toLowerCase()) >= 0) ||
      false
    )
  }

  private handleSearchFilterChanged = (query: string) => {
    const users = this.props.users
      .filter((x) => this.memberFilter(query, x))
      .sort(this.sortByStaff)
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

    if (actionName === "to-administrator") {
      await changeRole(UserRole.Administrator)
    } else if (actionName === "to-moderator") {
      await changeRole(UserRole.Moderator)
    } else if (actionName === "to-collaborator") {
      await changeRole(UserRole.Collaborator)
    } else if (actionName === "to-visitor") {
      await changeRole(UserRole.Visitor)
    } else if (actionName === "block") {
      await changeStatus(UserStatus.Blocked)
    } else if (actionName === "unblock") {
      await changeStatus(UserStatus.Active)
    }
  }

  private sortByStaff = (left: User, right: User) => {
    const rolePriority: { [key in UserRole]: number } = {
      [UserRole.Administrator]: 1,
      [UserRole.Moderator]: 2,
      [UserRole.Collaborator]: 3,
      [UserRole.Visitor]: 4,
    }

    if (rolePriority[left.role] === rolePriority[right.role]) {
      return left.name.localeCompare(right.name)
    }
    return rolePriority[left.role] - rolePriority[right.role]
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
        <div className="p-2">
          <VStack spacing={2} divide={true}>
            {this.state.visibleUsers.map((user) => (
              <UserListItem key={user.id} user={user} onAction={this.handleAction} />
            ))}
          </VStack>
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
            <strong>Moderators</strong> can moderate discussions and manage community interactions.
          </li>
          <li>
            <strong>Collaborators</strong> can edit and manage content, but not permissions and settings.
          </li>
          <li>
            <strong>Blocked</strong> users are unable to log into this site.
          </li>
        </ul>
      </>
    )
  }
}
