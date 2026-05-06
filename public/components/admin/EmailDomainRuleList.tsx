import React, { useState } from "react"
import { Trans } from "@lingui/macro"
import { Button, Input } from "@fider/components"
import { addEmailDomainRule, deleteEmailDomainRule, EmailDomainRule } from "@fider/services/actions/disposable"
import { notify } from "@fider/services"

import "./EmailDomainRuleList.scss"

interface Props {
  ruleType: "deny" | "allow"
  rules: EmailDomainRule[]
  onChange: () => void
}

export const EmailDomainRuleList: React.FC<Props> = ({ ruleType, rules, onChange }) => {
  const [input, setInput] = useState("")
  const [error, setError] = useState<string | undefined>()
  const [busy, setBusy] = useState(false)

  const submit = async () => {
    const domain = input.trim().toLowerCase()
    if (!domain) return
    setBusy(true)
    setError(undefined)
    const result = await addEmailDomainRule(domain, ruleType)
    setBusy(false)
    if (!result.ok) {
      const fieldErr = result.error?.errors?.find((e) => e.field === "domain")
      setError(fieldErr?.message ?? "Could not add domain.")
      return
    }
    setInput("")
    onChange()
  }

  const remove = async (id: number) => {
    const result = await deleteEmailDomainRule(id)
    if (!result.ok) {
      notify.error("Could not remove domain.")
      return
    }
    onChange()
  }

  const placeholder = ruleType === "deny" ? "e.g. mailinator.com" : "e.g. company.com"

  return (
    <div className="c-email-domain-rule-list">
      <div className="c-email-domain-rule-list__chips">
        {rules.length === 0 && (
          <span className="text-muted">
            <Trans id="admin.disposable.empty">No domains added yet.</Trans>
          </span>
        )}
        {rules.map((r) => (
          <span key={r.id} className="c-email-domain-rule-list__chip">
            <span className="c-email-domain-rule-list__domain">{r.domain}</span>
            <button type="button" aria-label="Remove" className="c-email-domain-rule-list__remove" onClick={() => remove(r.id)}>
              ×
            </button>
          </span>
        ))}
      </div>
      <div className="c-email-domain-rule-list__add">
        <Input field="domain" placeholder={placeholder} value={input} onChange={setInput} disabled={busy} />
        <Button onClick={submit} disabled={busy || !input.trim()}>
          <Trans id="admin.disposable.add">Add</Trans>
        </Button>
      </div>
      {error && <div className="c-email-domain-rule-list__error">{error}</div>}
    </div>
  )
}
