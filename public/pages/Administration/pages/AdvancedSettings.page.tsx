import React from "react"

import { TextArea, Form, Button } from "@fider/components"
import { Failure, actions, Fider } from "@fider/services"
import { AdminBasePage } from "../components/AdminBasePage"

interface AdvancedSettingsPageProps {
  customCSS: string
  allowedSchemes: string
  licenseKey: string
  isCommercial: boolean
}

interface AdvancedSettingsPageState {
  customCSS: string
  allowedSchemes: string
  error?: Failure
  copied: boolean
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
      allowedSchemes: this.props.allowedSchemes,
      copied: false,
    }
  }

  private setCustomCSS = (customCSS: string): void => {
    this.setState({ customCSS })
  }

  private setAllowedSchemes = (allowedSchemes: string): void => {
    this.setState({ allowedSchemes })
  }

  private handleSave = async (): Promise<void> => {
    const result = await actions.updateTenantAdvancedSettings(this.state.customCSS, this.state.allowedSchemes)
    if (result.ok) {
      location.reload()
    } else {
      this.setState({ error: result.error })
    }
  }

  private copyLicenseKey = (): void => {
    navigator.clipboard.writeText(this.props.licenseKey)
    this.setState({ copied: true })
    setTimeout(() => this.setState({ copied: false }), 2000)
  }

  public content() {
    return (
      <Form error={this.state.error}>
        {this.props.isCommercial && this.props.licenseKey && (
          <div className="field">
            <label>Commercial License Key</label>
            <p className="text-muted">
              Use this key to run Fider self-hosted with commercial features (content moderation).
            </p>
            <div className="mt-2 mb-4">
              <div className="flex gap-2">
                <input
                  type="text"
                  id="licenseKey"
                  readOnly
                  value={this.props.licenseKey}
                  className="w-full font-mono text-sm bg-gray-50 border border-gray-300 rounded px-3 py-2"
                  style={{ fontFamily: "monospace" }}
                />
                <Button size="small" onClick={this.copyLicenseKey} className="whitespace-nowrap">
                  {this.state.copied ? "Copied!" : "Copy"}
                </Button>
              </div>
            </div>
          </div>
        )}

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

        {Fider.settings.allowAllowedSchemes && (
          <TextArea
            field="allowedSchemes"
            label="Allowed URL Schemes"
            disabled={!Fider.session.user.isAdministrator}
            minRows={3}
            value={this.state.allowedSchemes}
            onChange={this.setAllowedSchemes}
          >
            <p className="text-muted">
              By default, uncommon URL schemes are forbidden in links.
              <br />
              If you want to allow linking monero or bitcoin addresses, you should add <code>^monero:[48]</code> or <code>^bitcoin:(1|3|bc1)</code> here.
            </p>
            <p className="text-muted">
              These are regular expressions, one per line, matched against the link address. <code>^javascript</code> is always rejected.
            </p>
          </TextArea>
        )}

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
