import React from "react"
import { ErrorPageWrapper } from "./components/ErrorPageWrapper"

interface MaintenanceProps {
  message: string
  until?: string
}

const Maintenance = (props: MaintenanceProps) => {
  return (
    <ErrorPageWrapper id="p-maintenance" showHomeLink={false}>
      <h1 className="text-display">UNDER MAINTENANCE</h1>
      <p>{props.message}</p>
      {props.until ? (
        <p>
          We&apos;ll be back at <strong>{props.until}</strong>.
        </p>
      ) : (
        <p>We&apos;ll be back soon.</p>
      )}
    </ErrorPageWrapper>
  )
}

export default Maintenance
