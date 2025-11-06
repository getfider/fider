import { Trans } from "@lingui/react/macro"
import React from "react"
import { HStack } from "../layout"
import { Button } from "./Button"

interface PaginationProps {
  currentPage: number
  totalPages: number
  onPageChange: (page: number) => void
}

export const Pagination = ({ currentPage, totalPages, onPageChange }: PaginationProps) => {
  if (totalPages <= 1) {
    return null
  }

  const getVisiblePages = () => {
    const visiblePages: (number | string)[] = []

    if (totalPages <= 7) {
      // Show all pages if 7 or fewer
      for (let i = 1; i <= totalPages; i++) {
        visiblePages.push(i)
      }
    } else {
      // Always show first page
      visiblePages.push(1)

      if (currentPage <= 4) {
        // Show first 5 pages + ellipsis + last page
        for (let i = 2; i <= 5; i++) {
          visiblePages.push(i)
        }
        visiblePages.push("...")
        visiblePages.push(totalPages)
      } else if (currentPage >= totalPages - 3) {
        // Show first page + ellipsis + last 5 pages
        visiblePages.push("...")
        for (let i = totalPages - 4; i <= totalPages; i++) {
          visiblePages.push(i)
        }
      } else {
        // Show first page + ellipsis + current-1, current, current+1 + ellipsis + last page
        visiblePages.push("...")
        for (let i = currentPage - 1; i <= currentPage + 1; i++) {
          visiblePages.push(i)
        }
        visiblePages.push("...")
        visiblePages.push(totalPages)
      }
    }

    return visiblePages
  }

  const handlePageClick = (page: number | string) => {
    if (typeof page === "number" && page !== currentPage) {
      onPageChange(page)
    }
  }

  const handlePrevClick = () => {
    if (currentPage > 1) {
      onPageChange(currentPage - 1)
    }
  }

  const handleNextClick = () => {
    if (currentPage < totalPages) {
      onPageChange(currentPage + 1)
    }
  }

  const visiblePages = getVisiblePages()

  return (
    <HStack justify="between">
      <Button className="no-focus" size="small" variant="secondary" onClick={handlePrevClick} disabled={currentPage === 1}>
        <Trans id="pagination.prev">Previous</Trans>
      </Button>

      {/* Page Numbers */}
      <div className="flex gap-1">
        {visiblePages.map((page, index) => (
          <React.Fragment key={index}>
            {page === "..." ? (
              <span className="px-3 align-self-end text-sm text-gray-500">...</span>
            ) : (
              <Button className="no-focus" size="small" variant="secondary" onClick={() => handlePageClick(page)}>
                {page}
              </Button>
            )}
          </React.Fragment>
        ))}
      </div>

      {/* Next Button */}
      <Button className="no-focus" onClick={handleNextClick} disabled={currentPage === totalPages} variant="secondary" size="small">
        <Trans id="pagination.next">Next</Trans>
      </Button>
    </HStack>
  )
}
