import React, { useState } from "react"

import { Button, Form, Input, TenantLogo } from "@fider/components"
import { actions, Failure } from "@fider/services"
import { i18n } from "@lingui/core"
import { Trans } from "@lingui/react/macro"

import { EmailVerificationKind } from "@fider/models"

import "./CompleteSignInProfile.page.scss"

interface CompleteSignInProfilePageProps {
  kind: EmailVerificationKind
  k: string
  c?: string
}

const CompleteSignInProfilePage = (props: CompleteSignInProfilePageProps) => {
  const [name, setName] = useState("")
  const [error, setError] = useState<Failure | undefined>()

  const submit = async () => {
    const result = await actions.completeProfile(props.kind, props.k, name)
    if (result.ok) {
      if (props.c !== undefined) {
        location.href = "/?c=" + props.c
      } else {
        location.href = "/"
      }
    } else if (result.error) {
      setError(result.error)
    }
  }

  return (
    <>
      <div id="p-complete-profile" className="page container w-max-6xl bg-gray-100">
        <div className="flex flex-y justify-center full-height py-4">
          <div className="text-center mb-8">
            <a href="/">
              <TenantLogo size={50} />
            </a>
          </div>

          <div className="box shadow-sm text-center w-full">
            <p className="text-title text-center">
              <Trans id="modal.completeprofile.header">Complete your profile</Trans>
            </p>

            <p>
              <Trans id="modal.completeprofile.text">Because this is your first sign in, please enter your name.</Trans>
            </p>
            <Form error={error} className="mb-4">
              <Input
                field="name"
                onChange={setName}
                maxLength={100}
                placeholder={i18n._("modal.completeprofile.name.placeholder", { message: "Name" })}
                suffix={
                  <Button type="submit" onClick={submit} variant="primary" disabled={name === ""}>
                    <Trans id="action.submit">Submit</Trans>
                  </Button>
                }
              />
            </Form>
          </div>
        </div>
      </div>
    </>
  )
}

export default CompleteSignInProfilePage
