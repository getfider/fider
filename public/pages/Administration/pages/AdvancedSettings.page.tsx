import "./AdvancedSettings.page.scss"

import React from "react"

import { TextArea, Form, Button } from "@fider/components"
import { Failure, actions, Fider } from "@fider/services"
import { FaStar } from "react-icons/fa"
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
  public icon = FaStar
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
          <p className="info">
            Custom CSS allows you to change the look and feel of Fider so that you can apply your own branding.
            <br />
            This is a powerful and flexible feature, but requires basic understanding of <a href="https://developer.mozilla.org/en-US/docs/Learn/CSS">CSS</a>.
          </p>
          <p className="info">
            Custom CSS might break the design of your site as Fider evolves. By doing this, you&apos;re taking this risk, and you will need to fix issues if
            they arise. <br /> You can minimize some issues by following these recommendations:
          </p>
          <ul className="info">
            <li>
              <strong>Avoid nested selectors</strong>: Fider might change the structure of the HTML at any time, and it&apos;s likely that such changes would
              invalidate some rules.
            </li>
            <li>
              <strong>Keep it short</strong>: Customize only the essential. Avoid changing the style or structure of the entire site.
            </li>
          </ul>
        </TextArea>

        {Fider.session.user.isAdministrator && (
          <div className="field">
            <Button color="positive" onClick={this.handleSave}>
              Save
            </Button>
          </div>
        )}
      </Form>
    )
  }
}
