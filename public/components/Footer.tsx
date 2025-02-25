import React from "react"

const Footer: React.FC = () => {
  const currentYear = new Date().getFullYear()

  return (
    <footer className="site-footer">
      <nav>
        <ul>
          <li>
            <a href="/terms">Terms & Conditions</a>
          </li>
          <li>
            <a href="/privacy">Privacy Policy</a>
          </li>
          <li>
            <a href="mailto:contact@tarkov.community">Contact</a>
          </li>
          <li>
            <a href="https://discord.gg/escapefromtarkovofficial" target="_blank" rel="noreferrer">
              Official Tarkov Discord
            </a>
          </li>
          <li>
            <a href="https://discord.com/invite/7ZeEyfU" target="_blank" rel="noreferrer">
              Tarkov Wiki Discord
            </a>
          </li>
          <li>
            <a href="https://www.escapefromtarkov.com/support" target="_blank" rel="noreferrer">
              Official Tarkov Support
            </a>
          </li>
        </ul>
      </nav>
      <p>© {currentYear} - This website is not affiliated with BSG (Battlestate Games Ltd).</p>
    </footer>
  )
}

export default Footer
