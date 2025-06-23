import React, { useState } from "react"

import { Button, ButtonClickEvent, TextArea, Form, Input, ImageUploader, Select } from "@fider/components"
import { AdminPageContainer } from "../components/AdminBasePage"
import { actions, Failure, Fider } from "@fider/services"
import { ImageUpload } from "@fider/models"
import { useFider } from "@fider/hooks"
import locales from "@locale/locales"

const GeneralSettingsPage = () => {
  const fider = useFider()
  const [title, setTitle] = useState<string>(fider.session.tenant.name)
  const [welcomeMessage, setWelcomeMessage] = useState<string>(fider.session.tenant.welcomeMessage)
  const [invitation, setInvitation] = useState<string>(fider.session.tenant.invitation)
  const [logo, setLogo] = useState<ImageUpload | undefined>(undefined)
  const [cname, setCNAME] = useState<string>(fider.session.tenant.cname)
  const [locale, setLocale] = useState<string>(fider.session.tenant.locale)
  const [error, setError] = useState<Failure | undefined>(undefined)

  const handleSave = async (e: ButtonClickEvent) => {
    const result = await actions.updateTenantSettings({ title, cname, welcomeMessage, invitation, logo, locale })
    if (result.ok) {
      e.preventEnable()
      location.href = `/`
    } else if (result.error) {
      setError(result.error)
    }
  }

  const dnsInstructions = (): JSX.Element => {
    const isApex = cname.split(".").length <= 2
    const recordType = isApex ? "ALIAS" : "CNAME"
    return (
      <>
        <strong>{cname}</strong> {recordType}{" "}
        <strong>
          {fider.session.tenant.subdomain}
          {fider.settings.domain}
        </strong>
      </>
    )
  }

  return (
    <AdminPageContainer id="p-admin-general" name="general" title="General" subtitle="Manage your site settings">
      <Form error={error}>
        <Input field="title" label="Your Fider board's title" maxLength={60} value={title} disabled={!fider.session.user.isAdministrator} onChange={setTitle}>
          <p className="text-muted">Keep it short and snappy. Your product / service name is usually best.</p>
        </Input>

        <TextArea
          field="welcomeMessage"
          label="Welcome Message"
          value={welcomeMessage}
          disabled={!fider.session.user.isAdministrator}
          onChange={setWelcomeMessage}
        >
          <p className="text-muted">
            The message is shown on this site&apos;s home page. Use it to help visitors understand what this space is about and the importance of their
            feedback.
          </p>
        </TextArea>

        <Input
          field="invitation"
          label="Invitation"
          maxLength={60}
          value={invitation}
          disabled={!fider.session.user.isAdministrator}
          placeholder="Enter your suggestion here..."
          onChange={setInvitation}
        >
          <p className="text-muted">Placeholder text in the suggestion&apos;s box. It should invite your visitors into sharing their feedback.</p>
        </Input>

        <ImageUploader label="Your Logo" field="logo" bkey={fider.session.tenant.logoBlobKey} disabled={!fider.session.user.isAdministrator} onChange={setLogo}>
          <p className="text-muted">JPG, GIF or PNG smaller than 100KB, minimum size 200x200 pixels.</p>
        </ImageUploader>

        {!Fider.isSingleHostMode() && (
          <Input
            field="cname"
            label="Custom Domain"
            maxLength={100}
            placeholder="feedback.yourcompany.com"
            value={cname}
            disabled={!fider.session.user.isAdministrator}
            onChange={setCNAME}
          >
            <div className="text-muted">
              {cname ? (
                [
                  <p key={0}>Enter the following record into your DNS zone records:</p>,
                  <p key={1}>{dnsInstructions()}</p>,
                  <p key={2}>Please note that it may take up to 72 hours for the change to take effect worldwide due to DNS propagation.</p>,
                ]
              ) : (
                <p>
                  Use custom domains to access Fider via your own domain name <code>feedback.yourcompany.com</code>
                </p>
              )}
            </div>
          </Input>
        )}

        <Select
          label="Locale"
          field="locale"
          defaultValue={locale}
          options={Object.entries(locales).map(([k, v]) => ({
            value: k,
            label: v.text,
          }))}
          onChange={(o) => setLocale(o?.value || "en")}
        >
          {locale !== "en" && (
            <>
              <p className="text-muted">
                This language is translated by the Open Source community. If you find a mistake or would like to improve its quality, you can find the
                translations on{" "}
                <a className="text-link" target="_blank" rel="noopener" href="https://github.com/getfider/fider/tree/main/locale">
                  GitHub
                </a>{" "}
                and contribute with your own translations.
              </p>
              <p className="text-muted">Only public pages are translated. Internal and/or administrative pages will remain in English.</p>
            </>
          )}
        </Select>

        <div className="field">
          <Button disabled={!fider.session.user.isAdministrator} variant="primary" onClick={handleSave}>
            Save
          </Button>
        </div>
      </Form>
    </AdminPageContainer>
  )
}

export default GeneralSettingsPage
