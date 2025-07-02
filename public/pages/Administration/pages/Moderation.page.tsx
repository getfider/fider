import React from "react"
import { Button, Loader } from "@fider/components"
import { actions, notify, Fider } from "@fider/services"
import { AdminBasePage } from "@fider/pages/Administration/components/AdminBasePage"
import { HStack, VStack } from "@fider/components/layout"

interface ModerationItem {
  id: number
  type: "post" | "comment"
  title?: string
  content: string
  user: {
    id: number
    name: string
    email: string
  }
  createdAt: string
  postId?: number
  postNumber?: number
  postSlug?: string
  postTitle?: string
}

interface ModerationPageState {
  items: ModerationItem[]
  loading: boolean
  selectedItems: Set<string>
  selectAll: boolean
}

export default class ModerationPage extends AdminBasePage<any, ModerationPageState> {
  public id = "p-admin-moderation"
  public name = "moderation"
  public title = "Moderation"
  public subtitle = "Review and approve content"

  constructor(props: any) {
    super(props)

    this.state = {
      items: [],
      loading: true,
      selectedItems: new Set(),
      selectAll: false,
    }
  }

  public async componentDidMount() {
    await this.loadModerationItems()
  }

  private loadModerationItems = async () => {
    this.setState({ loading: true })
    try {
      const response = await fetch("/_api/admin/moderation/items")
      if (response.ok) {
        const data = await response.json()
        this.setState({ items: data.items || [], loading: false })
      } else {
        notify.error("Failed to load moderation items")
        this.setState({ loading: false })
      }
    } catch (error) {
      notify.error("Failed to load moderation items")
      this.setState({ loading: false })
    }
  }

  private handleSelectItem = (itemKey: string, checked: boolean) => {
    const selectedItems = new Set(this.state.selectedItems)
    if (checked) {
      selectedItems.add(itemKey)
    } else {
      selectedItems.delete(itemKey)
    }
    
    const selectAll = selectedItems.size === this.state.items.length
    this.setState({ selectedItems, selectAll })
  }

  private handleSelectAll = (checked: boolean) => {
    if (checked) {
      const selectedItems = new Set(this.state.items.map(item => `${item.type}-${item.id}`))
      this.setState({ selectedItems, selectAll: true })
    } else {
      this.setState({ selectedItems: new Set(), selectAll: false })
    }
  }

  private bulkApprove = async () => {
    const { selectedItems } = this.state
    if (selectedItems.size === 0) {
      notify.error("Please select items to approve")
      return
    }

    const postIDs: number[] = []
    const commentIDs: number[] = []

    selectedItems.forEach(itemKey => {
      const [type, id] = itemKey.split('-')
      if (type === 'post') {
        postIDs.push(parseInt(id))
      } else if (type === 'comment') {
        commentIDs.push(parseInt(id))
      }
    })

    try {
      const response = await actions.bulkApproveItems(postIDs, commentIDs)
      if (response.ok) {
        notify.success(`Approved ${selectedItems.size} item(s) successfully`)
        await this.loadModerationItems()
        this.setState({ selectedItems: new Set(), selectAll: false })
      } else {
        notify.error("Failed to approve items")
      }
    } catch (error) {
      notify.error("Failed to approve items")
    }
  }

  private bulkDecline = async () => {
    const { selectedItems } = this.state
    if (selectedItems.size === 0) {
      notify.error("Please select items to decline")
      return
    }

    if (!confirm(`Are you sure you want to decline ${selectedItems.size} item(s)? This action cannot be undone.`)) {
      return
    }

    const postIDs: number[] = []
    const commentIDs: number[] = []

    selectedItems.forEach(itemKey => {
      const [type, id] = itemKey.split('-')
      if (type === 'post') {
        postIDs.push(parseInt(id))
      } else if (type === 'comment') {
        commentIDs.push(parseInt(id))
      }
    })

    try {
      const response = await actions.bulkDeclineItems(postIDs, commentIDs)
      if (response.ok) {
        notify.success(`Declined ${selectedItems.size} item(s) successfully`)
        await this.loadModerationItems()
        this.setState({ selectedItems: new Set(), selectAll: false })
      } else {
        notify.error("Failed to decline items")
      }
    } catch (error) {
      notify.error("Failed to decline items")
    }
  }

  private approveItem = async (item: ModerationItem) => {
    try {
      let response
      if (item.type === 'post') {
        response = await actions.approvePost(item.id)
      } else {
        response = await actions.approveComment(item.id)
      }

      if (response.ok) {
        notify.success(`${item.type === 'post' ? 'Post' : 'Comment'} approved successfully`)
        await this.loadModerationItems()
      } else {
        notify.error(`Failed to approve ${item.type}`)
      }
    } catch (error) {
      notify.error(`Failed to approve ${item.type}`)
    }
  }

  private declineItem = async (item: ModerationItem) => {
    if (!confirm(`Are you sure you want to decline this ${item.type}? This action cannot be undone.`)) {
      return
    }

    try {
      let response
      if (item.type === 'post') {
        response = await actions.declinePost(item.id)
      } else {
        response = await actions.declineComment(item.id)
      }

      if (response.ok) {
        notify.success(`${item.type === 'post' ? 'Post' : 'Comment'} declined successfully`)
        await this.loadModerationItems()
      } else {
        notify.error(`Failed to decline ${item.type}`)
      }
    } catch (error) {
      notify.error(`Failed to decline ${item.type}`)
    }
  }

  private formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  public content() {
    const { items, loading, selectedItems, selectAll } = this.state

    if (!Fider.session.tenant.isModerationEnabled) {
      return (
        <div className="p-4 bg-blue-50 rounded border-l-4 border-blue-500">
          <p className="text-blue-800">
            Content moderation is currently disabled for this site. You can enable it in the{" "}
            <a href="/admin/privacy" className="text-blue-600 underline">Privacy Settings</a>.
          </p>
        </div>
      )
    }

    if (loading) {
      return <Loader />
    }

    if (items.length === 0) {
      return (
        <div className="text-center py-8">
          <h3 className="text-lg font-medium text-gray-900 mb-2">No pending items</h3>
          <p className="text-gray-600">All content has been moderated.</p>
        </div>
      )
    }

    return (
      <VStack spacing={4}>
        <div className="flex justify-between items-center">
          <h3 className="text-lg font-medium">
            {items.length} item(s) awaiting moderation
          </h3>
          <HStack spacing={2}>
            <Button 
              variant="primary" 
              size="small" 
              onClick={this.bulkApprove}
              disabled={selectedItems.size === 0}
            >
              Approve Selected ({selectedItems.size})
            </Button>
            <Button 
              variant="danger" 
              size="small" 
              onClick={this.bulkDecline}
              disabled={selectedItems.size === 0}
            >
              Decline Selected ({selectedItems.size})
            </Button>
          </HStack>
        </div>

        <div className="bg-white shadow rounded-lg overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  <input
                    type="checkbox"
                    checked={selectAll}
                    onChange={(e) => this.handleSelectAll(e.target.checked)}
                    className="rounded border-gray-300"
                  />
                </th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Type
                </th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Content
                </th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Author
                </th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Created
                </th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {items.map((item) => {
                const itemKey = `${item.type}-${item.id}`
                const isSelected = selectedItems.has(itemKey)
                
                return (
                  <tr key={itemKey} className={isSelected ? "bg-blue-50" : ""}>
                    <td className="px-4 py-4 whitespace-nowrap">
                      <input
                        type="checkbox"
                        checked={isSelected}
                        onChange={(e) => this.handleSelectItem(itemKey, e.target.checked)}
                        className="rounded border-gray-300"
                      />
                    </td>
                    <td className="px-4 py-4 whitespace-nowrap">
                      <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                        item.type === 'post' 
                          ? 'bg-green-100 text-green-800' 
                          : 'bg-blue-100 text-blue-800'
                      }`}>
                        {item.type === 'post' ? 'Post' : 'Comment'}
                      </span>
                    </td>
                    <td className="px-4 py-4">
                      <div className="max-w-xs">
                        <div className="font-medium text-gray-900 truncate">
                          {item.title || (item.type === 'comment' ? 'Comment' : 'Post')}
                        </div>
                        <div className="text-sm text-gray-500 truncate">
                          {item.content}
                        </div>
                        {item.type === 'comment' && item.postTitle && (
                          <div className="text-xs text-gray-400 truncate">
                            on: {item.postTitle}
                          </div>
                        )}
                      </div>
                    </td>
                    <td className="px-4 py-4 whitespace-nowrap text-sm text-gray-900">
                      {item.user.name}
                    </td>
                    <td className="px-4 py-4 whitespace-nowrap text-sm text-gray-500">
                      {this.formatDate(item.createdAt)}
                    </td>
                    <td className="px-4 py-4 whitespace-nowrap text-sm">
                      <HStack spacing={1}>
                        <Button
                          variant="primary"
                          size="small"
                          onClick={() => this.approveItem(item)}
                        >
                          Approve
                        </Button>
                        <Button
                          variant="danger"
                          size="small"
                          onClick={() => this.declineItem(item)}
                        >
                          Decline
                        </Button>
                        {item.type === 'post' && item.postNumber && item.postSlug && (
                          <a
                            href={`/posts/${item.postNumber}/${item.postSlug}`}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-blue-600 hover:text-blue-800 text-xs"
                          >
                            View
                          </a>
                        )}
                        {item.type === 'comment' && item.postNumber && item.postSlug && (
                          <a
                            href={`/posts/${item.postNumber}/${item.postSlug}#comment-${item.id}`}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-blue-600 hover:text-blue-800 text-xs"
                          >
                            View
                          </a>
                        )}
                      </HStack>
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>
      </VStack>
    )
  }
}