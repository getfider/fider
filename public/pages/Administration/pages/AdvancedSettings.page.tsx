import React from "react"

import { TextArea, Form, Button } from "@fider/components"
import { Failure, actions, Fider } from "@fider/services"
import { AdminBasePage } from "../components/AdminBasePage"

interface AdvancedSettingsPageProps {
  customCSS: string
  profanityWords: string
}

interface AdvancedSettingsPageState {
  customCSS: string
  profanityWords: string
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
      profanityWords: this.props.profanityWords,
    }
  }

  private setCustomCSS = (customCSS: string): void => {
    this.setState({ customCSS })
  }

  private setProfanityWords = (profanityWords: string): void => {
    this.setState({ profanityWords })
  }

  // Separate save function for profanity words.
  private handleSaveProfanityWords = async (): Promise<void> => {
    const result = await actions.updateProfanityWords(this.state.profanityWords)
    if (result.ok) {
      location.reload()
    } else {
      this.setState({ error: result.error })
    }
  }

  // Save function for custom CSS.
  private handleSaveCustomCSS = async (): Promise<void> => {
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
          {}
        </TextArea>

        <div className="field">
          <Button variant="primary" onClick={this.handleSaveCustomCSS}>
            Save Custom CSS
          </Button>
        </div>

        <TextArea
          field="profanityWords"
          label="Profanity Words (one per line)"
          disabled={!Fider.session.user.isAdministrator}
          minRows={5}
          value={this.state.profanityWords}
          onChange={this.setProfanityWords}
        >
          <p className="text-muted">
            Enter banned words, one per line. Any post or comment containing these words will be blocked.
          </p>
        </TextArea>

        <div className="field">
          <Button variant="primary" onClick={this.handleSaveProfanityWords}>
            Save Profanity Words
          </Button>
        </div>
      </Form>
    )
  }
}
