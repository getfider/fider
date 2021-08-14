import "./WebhookTemplateInfoModal.scss"

import { Button, Loader, Modal } from "@fider/components"
import { WebhookProperties } from "@fider/pages/Administration/components/webhook/WebhookProperties"
import React, { useEffect, useState } from "react"
import { WebhookType } from "@fider/models"
import { actions, StringObject } from "@fider/services"
import { HStack, VStack } from "@fider/components/layout"
import { HoverInfo } from "@fider/components/common/HoverInfo"

interface WebhookTemplateInfoProps {
  type: WebhookType
  isModalOpen: boolean
  onModalClose: () => void
}

interface FunctionSpecification {
  params: {
    type: string
    desc: string
    info?: string
    link?: string
  }[]
  description: string
  info?: string
  link?: string
}

const functions: StringObject<FunctionSpecification> = {
  stripHtml: {
    params: [{ type: "string", desc: "The input string containing HTML to strip" }],
    description: "Strip HTML tags from the input data",
    info: "It removes all tags to keep only content",
  },
  md5: {
    params: [{ type: "string", desc: "The input string to hash" }],
    description: "Hash text using md5 algorithm",
    info: "What is MD5?",
    link: "https://en.wikipedia.org/wiki/MD5",
  },
  lower: {
    params: [{ type: "string", desc: "The input string to lowercase" }],
    description: "Lowercase text",
  },
  upper: {
    params: [{ type: "string", desc: "The input string to uppercase" }],
    description: "Uppercase text",
  },
  markdown: {
    params: [{ type: "string", desc: "The input string containing Markdown to parse" }],
    description: "Parse Markdown to HTML from the input data",
    info: "When parsing, input is sanitized from HTML to prevent XSS attacks",
  },
  format: {
    params: [
      {
        type: "string",
        desc: "The date format according to Go specifications",
        info: "See Go time format",
        link: "https://yourbasic.org/golang/format-parse-string-time-date-example/#standard-time-and-date-formats",
      },
      { type: "time", desc: "The time to format" },
    ],
    description: "Format the given date and time",
  },
  quote: {
    params: [{ type: "string", desc: "The input string to quote" }],
    description: "Enquote a string and escape inner special characters",
    info: "You should use this function when using a value as a JSON field",
  },
  escape: {
    params: [{ type: "string", desc: "The input string to escape" }],
    description: "Escape inner special characters of a string, without enquoting it",
    info: "You should use this function when using a value within an enquoted string",
  },
  urlquery: {
    params: [{ type: "string", desc: "The input string to encode as URL query value" }],
    description: "Encode a string into a valid URL query element",
    info: "You should use this function when using a value in URL",
  },
}
const textExample = 'A new post entitled "{{ .post_title }}" has been created by {{ .author_name }}.'
const jsonExample = `{
  "title": "New post: {{ escape .post_title }}",
  "content": {{ quote .post_description }},
  "user": {{ quote .author_name }},
  "date": {{ format "2006-01-02T15:04:05-0700" .post_created_at | quote }}
}`

export const WebhookTemplateInfoModal = (props: WebhookTemplateInfoProps) => {
  const [properties, setProperties] = useState<StringObject | null>()

  useEffect(() => {
    let mounted = true
    setProperties(undefined)
    actions.getWebhookHelp(props.type).then((result) => mounted && setProperties(result.ok ? result.data : null))
    return () => {
      mounted = false
    }
  }, [props.type])

  return (
    <Modal.Window className="c-webhook-templateinfo" isOpen={props.isModalOpen} onClose={props.onModalClose} size="large">
      <Modal.Header>Webhook template formatting help</Modal.Header>
      <Modal.Content>
        <VStack spacing={4} divide>
          <div>
            <h3 className="text-title mb-1">What is a template?</h3>
            <p>
              The template engine used is the native Go <code>text/template</code> package. The simpliest way to create a template is to write your text, and
              insert the property name, prefixed by a dot, enclosed in double braces with spaces within (wierd? complete!). Example:
            </p>
            <pre className="text-left">{textExample}</pre>
            <p>
              When using a structured text format such as JSON, you may need to wrap those properties into quotes. But be careful: you must escape values by
              yourself by calling the appropriate function depending of your situation. Example:
            </p>
            <pre className="text-left">{jsonExample}</pre>
            <Button href="https://pkg.go.dev/text/template" target="_blank" variant="primary">
              Official Go templates documentation
            </Button>
          </div>
          {properties === undefined ? (
            <p className="text-muted">Failed to load help data</p>
          ) : properties === null ? (
            <Loader text="Loading help data" />
          ) : (
            <>
              <div>
                <h3 className="text-title mb-1">Properties</h3>
                <WebhookProperties properties={properties} propsName="Property name" valueName="Example value" />
              </div>
              <div>
                <h3 className="text-title mb-1">Functions</h3>
                <VStack spacing={2} divide>
                  <HStack className="c-webhook-templateinfo__header flex-wrap" spacing={0}>
                    <div className="c-webhook-templateinfo__header-func">Function</div>
                    <div className="c-webhook-templateinfo__header-desc">Description</div>
                    <VStack className="c-webhook-templateinfo__params">
                      <div className="c-webhook-templateinfo__header-params">Parameters</div>
                      <HStack>
                        <div className="c-webhook-templateinfo__header-param">Type</div>
                        <div className="c-webhook-templateinfo__header-param-desc">Description</div>
                      </HStack>
                    </VStack>
                  </HStack>
                  {Object.entries(functions).map(([func, spec]) => (
                    <HStack key={func} className="flex-wrap" spacing={0}>
                      <div className="c-webhook-templateinfo__func">{func}</div>
                      <div className="c-webhook-templateinfo__desc">
                        {spec.description}
                        {spec.info && <HoverInfo text={spec.info} href={spec.link} target="_blank" />}
                      </div>
                      <VStack className="c-webhook-templateinfo__params" spacing={2} divide>
                        {spec.params.map((param, j) => (
                          <HStack key={j} className="flex-wrap" spacing={0}>
                            <div className="c-webhook-templateinfo__param">{param.type}</div>
                            <div className="c-webhook-templateinfo__param-desc">
                              {param.desc}
                              {param.info && <HoverInfo text={param.info} href={param.link} target="_blank" />}
                            </div>
                          </HStack>
                        ))}
                      </VStack>
                    </HStack>
                  ))}
                </VStack>
              </div>
            </>
          )}
        </VStack>
      </Modal.Content>
      <Modal.Footer>
        <Button variant="tertiary" onClick={props.onModalClose}>
          Close
        </Button>
      </Modal.Footer>
    </Modal.Window>
  )
}
