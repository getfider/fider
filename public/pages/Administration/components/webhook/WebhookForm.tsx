import "./WebhookForm.scss"

import React, { useEffect, useState } from "react"
import { Button, Field, Form, Input, Loader, Message, Select, SelectOption, TextArea, Toggle } from "@fider/components"
import { actions, Failure } from "@fider/services"
import { HStack, VStack } from "@fider/components/layout"
import { Webhook, WebhookData, WebhookPreviewResult, WebhookStatus, WebhookType } from "@fider/models"
import { HoverInfo } from "@fider/components/common/HoverInfo"
import { WebhookTemplateInfoModal } from "@fider/pages/Administration/components/webhook/WebhookTemplateInfoModal"

interface WebhookFormProps {
  webhook?: Webhook
  onSave: (data: WebhookData) => Promise<Failure | undefined>
  onCancel: () => void
}

interface HttpHeaderProps {
  header?: string
  value?: string
  onEdit: (header: string, value: string) => void
  onRemove?: (header: string) => void
  allHeaders?: string[]
}

const HttpHeader = (props: HttpHeaderProps) => {
  const [header, setHeader] = useState(props.header || "")
  const [value, setValue] = useState(props.value || "")
  const [editing, setEditing] = useState(props.onRemove === undefined)

  const duplicate = props.allHeaders && props.allHeaders.includes(header)

  const suffix = editing ? (
    props.onRemove === undefined ? (
      <Button
        variant="primary"
        onClick={() => {
          props.onEdit(header, value)
          setHeader("")
          setValue("")
        }}
        disabled={duplicate || header.length === 0 || value.length === 0}
      >
        Add
      </Button>
    ) : (
      <Button
        variant="primary"
        onClick={() => {
          setEditing(false)
          props.onEdit(header, value)
        }}
        disabled={value.length === 0}
      >
        Save
      </Button>
    )
  ) : (
    <>
      <Button variant="secondary" onClick={() => setEditing(true)}>
        Edit
      </Button>
      <Button variant="danger" onClick={() => props.onRemove && props.onRemove(header)}>
        Remove
      </Button>
    </>
  )

  return (
    <HStack justify="full" spacing={4}>
      <Input
        field={`header-${header}`}
        value={header}
        onChange={setHeader}
        placeholder="Header"
        disabled={props.onRemove !== undefined}
        suffix={duplicate ? "Duplicate" : undefined}
      />
      <Input field={`value-${header}`} value={value} onChange={setValue} placeholder="Value" disabled={!editing} suffix={suffix} />
    </HStack>
  )
}

export const WebhookForm = (props: WebhookFormProps) => {
  const [name, setName] = useState(props.webhook?.name || "")
  const [type, _setType] = useState(props.webhook?.type || WebhookType.NEW_POST)
  const [status, _setStatus] = useState(props.webhook?.status || WebhookStatus.DISABLED)
  const [url, setUrl] = useState(props.webhook?.url || "")
  const [content, setContent] = useState(props.webhook?.content || "")
  const [httpMethod, setHttpMethod] = useState(props.webhook?.http_method || "POST")
  const [httpHeaders, _setHttpHeaders] = useState(props.webhook?.http_headers || {})
  const [typing, setTyping] = useState<NodeJS.Timeout | undefined>()
  const [preview, setPreview] = useState<WebhookPreviewResult | null | undefined>()
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [error, setError] = useState<Failure | undefined>()

  const calculatePreview = () => {
    actions
      .previewWebhook(type, url, content)
      .then(
        (result) => (result.ok ? result.data : null),
        () => null
      )
      .then(setPreview)
  }

  useEffect(calculatePreview, [])
  useEffect(() => {
    if (typing) clearTimeout(typing)
    setPreview(undefined)
    setTyping(
      setTimeout(() => {
        calculatePreview()
        setTyping(undefined)
      }, 2_000)
    )
  }, [url, content])

  const handleSave = async () => {
    const error = await props.onSave({ name, type, status, url, content, http_method: httpMethod, http_headers: httpHeaders })
    if (error) {
      setError(error)
    }
  }
  const handleCancel = () => props.onCancel()

  const setType = (option?: SelectOption) => _setType(option?.value as WebhookType)
  const setStatus = (active: boolean) => _setStatus(active ? WebhookStatus.ENABLED : WebhookStatus.DISABLED)

  const setHttpHeader = (header: string, value: string) => {
    _setHttpHeaders((headers) => ({
      ...headers,
      [header]: value,
    }))
  }

  const removeHttpHeader = (header: string) => {
    _setHttpHeaders((headers) => {
      const { [header]: _, ...remaining } = headers
      return remaining
    })
  }

  const showModal = () => setIsModalOpen(true)
  const hideModal = () => setIsModalOpen(false)

  const allHeaders = Object.keys(httpHeaders)
  const title = props.webhook ? `Webhook #${props.webhook.id}: ${props.webhook.name}` : "New webhook"
  return (
    <>
      {status === WebhookStatus.FAILED && (
        <Message type="error" showIcon>
          This webhook has failed
        </Message>
      )}
      <h2 className="text-title mb-2">{title}</h2>
      <Form className="c-webhook-form" error={error}>
        <Input field="name" label="Name" value={name} onChange={setName} placeholder="My awesome webhook" />
        <Select
          label="Type"
          field="type"
          defaultValue={type}
          options={[
            { label: "New Post", value: WebhookType.NEW_POST },
            { label: "New Comment", value: WebhookType.NEW_COMMENT },
            { label: "Change Status", value: WebhookType.CHANGE_STATUS },
            { label: "Delete Post", value: WebhookType.DELETE_POST },
          ]}
          onChange={setType}
        />
        <Field label="Status">
          <Toggle active={status === WebhookStatus.ENABLED} onToggle={setStatus} />
          {status === WebhookStatus.FAILED && <p className="text-muted mt-1">This webhook was disabled due to a trigger failure</p>}
        </Field>
        <Input
          field="url"
          label="URL"
          afterLabel={<HoverInfo text="You can use Go template formatting with many properties here" onClick={showModal} />}
          value={url}
          onChange={setUrl}
          placeholder="https://webhook.site/..."
        />
        <TextArea
          className="c-webhook-form__content"
          field="content"
          label="Content"
          afterLabel={<HoverInfo text="You can use Go template formatting with many properties here" onClick={showModal} />}
          value={content}
          onChange={setContent}
          placeholder="Request body"
        />
        <Input field="http_method" label="HTTP Method" value={httpMethod} onChange={setHttpMethod} placeholder="POST" />
        <Field label="HTTP Headers" afterLabel={<HoverInfo text="Those headers are sent in the request when the webhook is triggered" />}>
          <VStack spacing={2}>
            {Object.entries(httpHeaders).map(([header, value]) => (
              <HttpHeader key={header} header={header} value={value} onEdit={setHttpHeader} onRemove={removeHttpHeader} />
            ))}
            <HttpHeader onEdit={setHttpHeader} allHeaders={allHeaders} />
          </VStack>
        </Field>
        {(url || content) && (
          <Field label="Preview" className="c-webhook-form__preview">
            {preview === null ? (
              <p className="text-muted">Failed to load preview</p>
            ) : preview === undefined ? (
              <Loader className="text-center" text="Loading preview" />
            ) : (
              <VStack className="bg-gray-50 rounded-md p-2" spacing={2}>
                {url && (
                  <div>
                    <h3 className="text-bold mb-1">URL</h3>
                    <pre>{preview.url.value ? preview.url.value : preview.url.error}</pre>
                    {preview.url.message && <p className="text-muted">{preview.url.message}</p>}
                  </div>
                )}
                {content && (
                  <div>
                    <h3 className="text-bold mb-1">Content</h3>
                    <pre>{preview.content.value ? preview.content.value : preview.content.error}</pre>
                    {preview.content.message && <p className="text-muted">{preview.content.message}</p>}
                  </div>
                )}
              </VStack>
            )}
          </Field>
        )}
        <HStack>
          <Button variant="primary" onClick={handleSave}>
            Save
          </Button>
          <Button onClick={handleCancel} variant="tertiary">
            Cancel
          </Button>
        </HStack>
      </Form>
      <WebhookTemplateInfoModal type={type} isModalOpen={isModalOpen} onModalClose={hideModal} />
    </>
  )
}
