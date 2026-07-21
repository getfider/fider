"use strict"

// Local ESLint rule: flags JSX `href` attributes set to root-relative paths
// without going through the <Link> component or resolveHref(). Catches both
// string literals (href="/admin") and template literals (href={`/posts/${id}`}).
//
// Broken under sub-path hosting: when BASE_URL has a path component
// (e.g. https://example.com/feedback), bare /-prefixed hrefs resolve to the
// domain root and skip the sub-path. Using <Link> or resolveHref() routes
// the href through basePath() and produces correct URLs at every host mode.
//
// Safe elements: <Link> (auto-resolves), <Button> (auto-resolves), and
// <Dropdown.ListItem> (auto-resolves) are exempt. Anything else using a
// root-relative href is flagged.
module.exports = {
  meta: {
    type: "problem",
    docs: {
      description: "Disallow bare root-relative hrefs in JSX; use <Link> or resolveHref()",
    },
    messages: {
      bareHref: "Avoid hardcoded root-relative href on <{{element}}>; use the <Link> component or resolveHref() so sub-path hosting is respected.",
    },
    schema: [],
  },
  create(context) {
    const SAFE_ELEMENTS = new Set(["Link", "Button"])

    function getElementName(jsxOpeningElement) {
      const name = jsxOpeningElement.name
      if (name.type === "JSXIdentifier") return name.name
      if (name.type === "JSXMemberExpression") {
        // Dropdown.ListItem — walk to the tail identifier
        let node = name
        while (node.type === "JSXMemberExpression") node = node.property
        return node.name
      }
      return ""
    }

    function startsWithSlash(node) {
      if (node.type === "Literal" && typeof node.value === "string") {
        return node.value.startsWith("/")
      }
      if (node.type === "TemplateLiteral" && node.quasis.length > 0) {
        const first = node.quasis[0].value.cooked
        return typeof first === "string" && first.startsWith("/")
      }
      return false
    }

    return {
      JSXAttribute(node) {
        if (node.name.name !== "href" || !node.value) return

        const opening = node.parent
        const element = getElementName(opening)
        if (SAFE_ELEMENTS.has(element)) return
        // Dropdown.ListItem auto-resolves too
        if (element === "ListItem" && opening.name.type === "JSXMemberExpression") return

        let expr = null
        if (node.value.type === "Literal") {
          expr = node.value
        } else if (node.value.type === "JSXExpressionContainer") {
          expr = node.value.expression
        }
        if (!expr) return

        if (startsWithSlash(expr)) {
          context.report({
            node,
            messageId: "bareHref",
            data: { element },
          })
        }
      },
    }
  },
}
