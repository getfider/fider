import React, { useState } from "react"
import { OAuthConfig, OAuthConfigStatus, ImageUpload } from "@fider/models"
import { Failure, actions } from "@fider/services"
import { Form, Button, Input, SocialSignInButton, Field, ImageUploader, Toggle, hasError } from "@fider/components"
import { useFider } from "@fider/hooks"
import { HStack, VStack } from "@fider/components/layout"

interface OAuthFormProps {
  config?: OAuthConfig
  onCancel: () => void
  cantDisable: boolean
}

export const OAuthForm: React.FC<OAuthFormProps> = (props) => {
  const fider = useFider()
  const [provider] = useState((props.config && props.config.provider) || "")
  const [displayName, setDisplayName] = useState((props.config && props.config.displayName) || "")
  const [enabled, setEnabled] = useState((props.config && props.config.status === OAuthConfigStatus.Enabled) || false)
  const [isTrusted, setTrusted] = useState((props.config && props.config.isTrusted) || false)
  const [clientID, setClientID] = useState((props.config && props.config.clientID) || "")
  const [clientSecret, setClientSecret] = useState((props.config && props.config.clientSecret) || "")
  const [clientSecretEnabled, setClientSecretEnabled] = useState(!props.config)
  const [authorizeURL, setAuthorizeURL] = useState((props.config && props.config.authorizeURL) || "")
  const [tokenURL, setTokenURL] = useState((props.config && props.config.tokenURL) || "")
  const [profileURL, setProfileURL] = useState((props.config && props.config.profileURL) || "")
  const [scope, setScope] = useState((props.config && props.config.scope) || "")
  const [jsonUserIDPath, setJSONUserIDPath] = useState((props.config && props.config.jsonUserIDPath) || "")
  const [jsonUserNamePath, setJSONUserNamePath] = useState((props.config && props.config.jsonUserNamePath) || "")
  const [jsonUserEmailPath, setJSONUserEmailPath] = useState((props.config && props.config.jsonUserEmailPath) || "")
  const [jsonUserRolesPath, setJSONUserRolesPath] = useState((props.config && props.config.jsonUserRolesPath) || "")
  const [allowedRoles, setAllowedRoles] = useState((props.config && props.config.allowedRoles) || "")
  const [logo, setLogo] = useState<ImageUpload | undefined>()
  const [logoURL, setLogoURL] = useState<string | undefined>()
  const [logoBlobKey, setLogoBlobKey] = useState((props.config && props.config.logoBlobKey) || "")
  const [error, setError] = useState<Failure | undefined>()
  const [showAdvanced, setShowAdvanced] = useState(!!(jsonUserRolesPath || allowedRoles))

  const handleSave = async () => {
    const result = await actions.saveOAuthConfig({
      provider,
      status: enabled ? OAuthConfigStatus.Enabled : OAuthConfigStatus.Disabled,
      isTrusted,
      displayName,
      clientID,
      clientSecret: clientSecretEnabled ? clientSecret : "",
      authorizeURL,
      tokenURL,
      profileURL,
      scope,
      jsonUserIDPath,
      jsonUserNamePath,
      jsonUserEmailPath,
      jsonUserRolesPath,
      allowedRoles,
      logo,
    })
    if (result.ok) {
      location.reload()
    } else {
      setError(result.error)
      // Don't leave an error hidden inside the collapsed section
      if (hasError("jsonUserRolesPath", result.error) || hasError("allowedRoles", result.error)) {
        setShowAdvanced(true)
      }
    }
  }

  const handleLogoChange = (newLogo: ImageUpload, instanceID: string, previewURL: string) => {
    setLogo(newLogo)
    setLogoURL(previewURL)
    setLogoBlobKey("")
  }

  const handleCancel = async () => {
    props.onCancel()
  }

  const enableClientSecret = () => {
    setClientSecret("")
    setClientSecretEnabled(true)
  }

  const isAdmin = fider.session.user.isAdministrator
  const siteAccess = fider.session.tenant.isPrivate ? "private" : fider.session.tenant.isSignUpDisabled ? "public with sign-ups disabled" : "public"
  const title = props.config ? `OAuth Provider: ${props.config.displayName}` : "New OAuth Provider"

  return (
    <>
      <h2 className="text-display mb-6">{title}</h2>
      <Form error={error}>
        <Input field="displayName" label="Display Name" maxLength={50} value={displayName} disabled={!isAdmin} onChange={setDisplayName} />

        <ImageUploader label="Logo" field="logo" bkey={logoBlobKey} disabled={!isAdmin} onChange={handleLogoChange}>
          <p className="text-muted mt-1">Square JPG, GIF or PNG, smaller than 50KB and at least 24x24 pixels.</p>
        </ImageUploader>

        <Field label="Button Preview">
          {/* The preview has no URL, so clicking it would reload the page and lose the form */}
          <div className="grid w-max-4xl" onClickCapture={(e) => e.preventDefault()}>
            <SocialSignInButton option={{ displayName: displayName || "Button", provider, logoBlobKey, logoURL }} />
          </div>
        </Field>

        <h3 className="text-title mt-8 mb-1">Credentials</h3>
        <p className="text-muted mb-4">From the OAuth app you registered with your provider.</p>

        <Input field="clientID" label="Client ID" maxLength={100} value={clientID} disabled={!isAdmin} onChange={setClientID} />

        <Input field="clientSecret" label="Client Secret" maxLength={500} value={clientSecret} disabled={!clientSecretEnabled} onChange={setClientSecret}>
          {!clientSecretEnabled && (
            <p className="text-muted mt-1">
              Hidden for security reasons.{" "}
              <Button variant="link" size="no-padding" onClick={enableClientSecret}>
                Change
              </Button>
            </p>
          )}
        </Input>

        <Input field="authorizeURL" label="Authorize URL" maxLength={300} value={authorizeURL} disabled={!isAdmin} onChange={setAuthorizeURL} />

        <Input field="tokenURL" label="Token URL" maxLength={300} value={tokenURL} disabled={!isAdmin} onChange={setTokenURL} />

        <Input field="scope" label="Scope" maxLength={100} value={scope} disabled={!isAdmin} onChange={setScope}>
          <p className="text-muted mt-1">
            Separate scopes with a space. Only request what&apos;s needed to read the user&apos;s <strong>id</strong>, <strong>name</strong> and{" "}
            <strong>email</strong>.
          </p>
        </Input>

        <h3 className="text-title mt-8 mb-1">User Profile</h3>
        <p className="text-muted mb-4">
          How Fider reads the user&apos;s details after they sign in.{" "}
          <a rel="noopener" className="text-link" target="_blank" href="https://docs.fider.io/configuring-oauth#configuring-the-json-paths">
            Learn about JSON paths
          </a>
          .
        </p>

        <Input field="profileURL" label="Profile API URL" maxLength={300} value={profileURL} disabled={!isAdmin} onChange={setProfileURL}>
          <p className="text-muted mt-1">Leave empty to read the user&apos;s details from the access token.</p>
        </Input>

        <div className="grid grid-cols-3 gap-x-4">
          <Input field="jsonUserIDPath" label="ID JSON Path" maxLength={100} value={jsonUserIDPath} disabled={!isAdmin} onChange={setJSONUserIDPath}>
            <p className="text-muted mt-1">Required. Must be unique.</p>
          </Input>
          <Input field="jsonUserNamePath" label="Name JSON Path" maxLength={100} value={jsonUserNamePath} disabled={!isAdmin} onChange={setJSONUserNamePath}>
            <p className="text-muted mt-1">Recommended.</p>
          </Input>
          <Input
            field="jsonUserEmailPath"
            label="Email JSON Path"
            maxLength={100}
            value={jsonUserEmailPath}
            disabled={!isAdmin}
            onChange={setJSONUserEmailPath}
          >
            <p className="text-muted mt-1">Recommended.</p>
          </Input>
        </div>

        <h3 className="text-title mt-8 mb-4">Access</h3>

        <VStack spacing={6}>
          <HStack justify="between" align="start" spacing={4}>
            <div>
              <strong>Trusted source</strong>
              <p className="text-muted mt-1">
                Users from this provider can join without an invite, even when the site is private or sign-ups are disabled. Use it for company identity
                providers like Okta or Microsoft Entra, not public ones like Facebook. This site is currently <strong>{siteAccess}</strong>.
              </p>
            </div>
            <Toggle field="isTrusted" ariaLabel="Trusted source" active={isTrusted} onToggle={setTrusted} />
          </HStack>

          <HStack justify="between" align="start" spacing={4}>
            <div>
              <strong>Enabled</strong>
              <p className="text-muted mt-1">
                {enabled
                  ? "Shown on the sign in page. Test it before enabling: the Test button appears once you've saved."
                  : "Hidden from the sign in page. Users can't sign in with this provider."}
              </p>
              {props.cantDisable && <p className="text-muted mt-1">Enable email authentication before disabling your last OAuth provider.</p>}
            </div>
            <Toggle field="status" ariaLabel="Enabled" disabled={props.cantDisable} active={enabled} onToggle={setEnabled} />
          </HStack>
        </VStack>

        <details className="mt-8" open={showAdvanced} onToggle={(e) => setShowAdvanced(e.currentTarget.open)}>
          <summary className="text-title clickable mb-4">Advanced</summary>
          <Input
            field="jsonUserRolesPath"
            label="Roles JSON Path"
            maxLength={100}
            value={jsonUserRolesPath}
            disabled={!isAdmin}
            onChange={setJSONUserRolesPath}
          >
            <p className="text-muted mt-1">Where to find the user&apos;s roles in their profile.</p>
          </Input>
          <Input field="allowedRoles" label="Allowed Roles" maxLength={500} value={allowedRoles} disabled={!isAdmin} onChange={setAllowedRoles}>
            <p className="text-muted mt-1">
              Comma-separated, e.g. <strong>ROLE_ADMIN,ROLE_TEACHER</strong>. Only enforced when a Roles JSON Path is set. Leave empty to allow everyone.
            </p>
          </Input>
        </details>

        <HStack className="mt-8">
          <Button variant="primary" onClick={handleSave}>
            Save
          </Button>
          <Button variant="tertiary" onClick={handleCancel}>
            Cancel
          </Button>
        </HStack>
      </Form>
    </>
  )
}
