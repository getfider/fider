import React from "react"

import { TextArea, Form, Button } from "@fider/components"
import { Failure, actions, Fider } from "@fider/services"
import { AdminBasePage } from "../components/AdminBasePage"

interface AdvancedSettingsPageProps {
  customCSS: string
}

interface AdvancedSettingsPageState {
  customCSS: string
  error?: Failure
}

export default class AdvancedSettingsPage extends AdminBasePage<AdvancedSettingsPageProps, AdvancedSettingsPageState> {
  public id = "p-admin-advanced"
  public name = "advanced"
  public title = "Advanced"
  public subtitle = "Manage your site settings"

  constructor(props: AdvancedSettingsPageProps) {
    super(props)

    this.state = {
      customCSS: this.props.customCSS,
    }
  }

  private setCustomCSS = (customCSS: string): void => {
    this.setState({ customCSS })
  }

  private handleSave = async (): Promise<void> => {
    const result = await actions.updateTenantAdvancedSettings(this.state.customCSS)
    if (result.ok) {
      location.reload()
    } else {
      this.setState({ error: result.error })
    }
  }

  public content() {
    return (
      <Form error={this.state.error}>
        <TextArea
          field="customCSS"
          label="Custom CSS"
          disabled={!Fider.session.user.isAdministrator}
          minRows={10}
          value={this.state.customCSS}
          onChange={this.setCustomCSS}
        >
          <p className="text-muted">
            Custom CSS allows you to change the look and feel of Fider and apply your own branding.
            <br />
            This is a powerful and flexible feature, but requires basic understanding of <a href="https://developer.mozilla.org/en-US/docs/Learn/CSS">CSS</a>.
          </p>
          <p className="text-muted">
            Custom CSS might break the design of your site as Fider evolves. You can minimize conflict by following these recommendations:
          </p>
          <ul className="text-muted">
            <li>
              <strong>Avoid nested selectors</strong>: Fider might change the structure of the HTML at any time. It&apos;s likely that such changes would
              invalidate some rules.
            </li>
            <li>
              <strong>Keep it simple</strong>: Customize only the essential.
            </li>
          </ul>
        </TextArea>

        {Fider.session.user.isAdministrator && (
          <div className="field">
            <Button variant="primary" onClick={this.handleSave}>
              Save
            </Button>
          </div>
        )}
      </Form>
    )
  }
}
