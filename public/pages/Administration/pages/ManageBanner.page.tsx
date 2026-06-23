import React from "react"
import { Button, Form, Select, SelectOption, TextArea, Toggle } from "@fider/components"
import { HStack, VStack } from "@fider/components/layout"
import { actions, Failure, notify } from "@fider/services"
import { i18n } from "@lingui/core"
import { Trans } from "@lingui/react/macro"
import { AdminBasePage } from "../components/AdminBasePage"

interface ManageBannerPageProps {
  siteBannerEnabled: boolean
  siteBannerMessage: string
  siteBannerVariant: string
}

interface ManageBannerPageState {
  enabled: boolean
  message: string
  variant: string
  error?: Failure
  busy: boolean
}

const VARIANTS = (): SelectOption[] => [
  { value: "info", label: i18n._({ id: "admin.banner.variant.info", message: "Info — neutral announcement" }) },
  { value: "success", label: i18n._({ id: "admin.banner.variant.success", message: "Success — green confirmation" }) },
  { value: "warning", label: i18n._({ id: "admin.banner.variant.warning", message: "Warning — scheduled maintenance" }) },
  { value: "danger", label: i18n._({ id: "admin.banner.variant.danger", message: "Danger — incident in progress" }) },
  { value: "brand", label: i18n._({ id: "admin.banner.variant.brand", message: "Brand — tenant primary color" }) },
]

const MAX_LEN = 500

export default class ManageBannerPage extends AdminBasePage<ManageBannerPageProps, ManageBannerPageState> {
  public id = "p-admin-banner"
  public name = "banner"
  public title = i18n._({ id: "admin.banner.page.title", message: "Site Banner" })
  public subtitle = i18n._({ id: "admin.banner.page.subtitle", message: "Show a site-wide notice above the header for maintenance, releases, or incidents" })

  constructor(props: ManageBannerPageProps) {
    super(props)
    this.state = {
      enabled: props.siteBannerEnabled,
      message: props.siteBannerMessage,
      variant: props.siteBannerVariant || "info",
      busy: false,
    }
  }

  // The Button auto-re-enables after onClick resolves unless we call
  // event.preventEnable() — we want it to re-enable here, so the event is
  // intentionally untouched.
  private save = async () => {
    this.setState({ busy: true, error: undefined })
    const result = await actions.updateTenantSiteBanner({
      enabled: this.state.enabled,
      message: this.state.message,
      variant: this.state.variant,
    })
    if (result.ok) {
      this.setState({ busy: false })
      notify.success(i18n._({ id: "admin.banner.saved", message: "Banner settings saved." }))
    } else {
      this.setState({ busy: false, error: result.error })
    }
  }

  public content() {
    const remaining = Math.max(0, MAX_LEN - this.state.message.length)
    return (
      <Form error={this.state.error}>
        <VStack spacing={4}>
          {this.state.enabled && this.state.message.trim() !== "" && (
            <div>
              <p className="text-xs text-muted mb-1">
                <Trans id="admin.banner.preview">Preview</Trans>
              </p>
              <div
                className={`site-banner site-banner--${this.state.variant}`}
                style={{
                  padding: "10px 24px",
                  borderRadius: 6,
                  fontWeight: 600,
                  lineHeight: 1.45,
                  textAlign: "center",
                  whiteSpace: "pre-wrap",
                }}
              >
                {this.state.message}
              </div>
            </div>
          )}

          <HStack spacing={4}>
            <Toggle active={this.state.enabled} onToggle={(v) => this.setState({ enabled: v })} />
            <span className="text-sm">
              <Trans id="admin.banner.enable">Show site-wide banner</Trans>
            </span>
          </HStack>

          <Select
            field="variant"
            label={i18n._({ id: "admin.banner.form.variant", message: "Color variant" })}
            defaultValue={this.state.variant}
            options={VARIANTS()}
            onChange={(opt) => this.setState({ variant: opt?.value ?? "info" })}
          />

          <TextArea
            field="message"
            label={i18n._({ id: "admin.banner.form.message", message: "Message" })}
            value={this.state.message}
            onChange={(v) => this.setState({ message: v.slice(0, MAX_LEN) })}
            minRows={3}
          />
          <p className="text-xs text-muted">
            <Trans id="admin.banner.charsleft">{remaining} characters remaining</Trans>
          </p>

          <HStack spacing={2}>
            <Button variant="primary" onClick={this.save} disabled={this.state.busy}>
              <Trans id="admin.banner.save">Save banner</Trans>
            </Button>
          </HStack>
        </VStack>
      </Form>
    )
  }
}
