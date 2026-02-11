import React from "react"
import MailSentIllustration from "@fider/assets/images/undraw_mail-sent.svg"

import { LegalFooter, TenantLogo, Icon } from "@fider/components"
import { basePath } from "@fider/services"
import { Trans } from "@lingui/react/macro"

import "./LoginEmailSent.page.scss"

const LoginEmailSentPage = ({ email }: { email: string }) => {
  return (
    <>
      <div id="p-email-sent" className="page container w-max-6xl bg-gray-100">
        <div className="flex flex-y justify-center flex-items-center full-height py-4">
          <div className="text-center mb-8">
            <a href={`${basePath()}/`}>
              <TenantLogo size={50} />
            </a>
          </div>

          <div className="box shadow-sm text-center w-full">
            <Icon sprite={MailSentIllustration} height="120" className="mb-4" />

            <p className="text-xl text-center mb-4 text-gray-800">
              <Trans id="signin.message.emailsent">
                We have just sent a confirmation link to <b>{email}</b>. Click the link and you’ll be signed in.
              </Trans>
            </p>

            <LegalFooter />
          </div>
        </div>
      </div>
    </>
  )
}
export default LoginEmailSentPage
