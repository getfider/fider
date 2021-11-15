import React, { useState } from "react"

import { Button, Form, Input, LegalFooter, TenantLogo } from "@fider/components"
import { actions, Failure } from "@fider/services"
import { t, Trans } from "@lingui/macro"
import { EmailVerificationKind } from "@fider/models"

import "./CompleteSignInProfile.page.scss"

interface CompleteSignInProfilePageProps {
  kind: EmailVerificationKind
  k: string
}

const CompleteSignInProfilePage = (props: CompleteSignInProfilePageProps) => {
  const [name, setName] = useState("")
  const [error, setError] = useState<Failure | undefined>()

  const submit = async () => {
    const result = await actions.completeProfile(props.kind, props.k, name)
    if (result.ok) {
      location.href = "/"
    } else if (result.error) {
      setError(result.error)
    }
  }

  return (
    <>
      <div id="p-complete-profile" className="page container w-max-6xl">
        <div>
          <div className="h-20 text-center mb-4">
            <TenantLogo size={100} />
          </div>

          <p className="text-title">
            <Trans id="modal.completeprofile.header">Complete your profile</Trans>
          </p>

          <p>
            <Trans id="modal.completeprofile.text">Because this is your first sign in, please enter your name.</Trans>
          </p>
          <Form error={error}>
            <Input
              field="name"
              onChange={setName}
              maxLength={100}
              placeholder={t({ id: "modal.completeprofile.name.placeholder", message: "Name" })}
              suffix={
                <Button type="submit" onClick={submit} variant="primary" disabled={name === ""}>
                  <Trans id="action.submit">Submit</Trans>
                </Button>
              }
            />
          </Form>

          <LegalFooter />
        </div>
      </div>
    </>
  )
}

export default CompleteSignInProfilePage
