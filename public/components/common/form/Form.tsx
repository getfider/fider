import "./Form.scss"

import React from "react"
import { Failure, classSet } from "@fider/services"
import { DisplayError } from "@fider/components"

interface ValidationContext {
  error?: Failure
  clearError?: (field: string) => void
}

interface FormProps {
  children?: React.ReactNode
  className?: string
  error?: Failure
  autoComplete?: string
  onSubmit?: (e: React.FormEvent) => void
}

export const ValidationContext = React.createContext<ValidationContext>({
  clearError: (field: string) => {
    /* Default implementation does nothing with field: ${field} */
    void field // Prevent unused parameter warning
  },
})

export const Form: React.FunctionComponent<FormProps> = (props) => {
  const className = classSet({
    "c-form": true,
    [props.className || ""]: props.className,
  })

  const [formError, setFormError] = React.useState<Failure | undefined>(props.error)

  // Update formError when props.error changes
  React.useEffect(() => {
    setFormError(props.error)
  }, [props.error])

  // Function to clear error for a specific field
  const clearError = (field: string) => {
    if (formError && formError.errors) {
      const newErrors = formError.errors.filter((err) => err.field !== field)
      if (newErrors.length !== formError.errors.length) {
        setFormError({ errors: newErrors.length > 0 ? newErrors : undefined })
      }
    }
  }

  const handleSubmit = (e: React.FormEvent) => {
    if (props.onSubmit) {
      e.preventDefault()
      props.onSubmit(e)
    }
  }

  return (
    <form autoComplete={props.autoComplete || "off"} className={className} onSubmit={handleSubmit}>
      <DisplayError error={formError} />
      <ValidationContext.Provider value={{ error: formError, clearError }}>{props.children}</ValidationContext.Provider>
    </form>
  )
}
