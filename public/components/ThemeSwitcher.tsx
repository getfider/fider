import IconMoon from "@fider/assets/images/heroicons-moon.svg"
import IconSun from "@fider/assets/images/heroicons-sun.svg"
import React, { useEffect, useState } from "react"
import { Icon } from "./common"
import "./ThemeSwitcher.scss"

export const ThemeSwitcher = () => {
  const [currentTheme, setCurrentTheme] = useState(localStorage.getItem("theme") || "light")

  const toggleTheme = () => {
    const currentTheme = localStorage.getItem("theme") || "light"
    const newTheme = currentTheme === "light" ? "dark" : "light"
    setCurrentTheme(newTheme)
  }

  useEffect(() => {
    localStorage.setItem("theme", currentTheme)
    document.body.setAttribute("data-theme", currentTheme)
  }, [currentTheme])

  const icon = currentTheme === "light" ? <Icon sprite={IconMoon} className="h-6" /> : <Icon sprite={IconSun} className="h-6" />

  return (
    <button onClick={toggleTheme} aria-label="Toggle theme" className="c-themeswitcher">
      {icon}
    </button>
  )
}
