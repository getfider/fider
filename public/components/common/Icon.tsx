import React from "react"

interface IconProps {
  sprite: SpriteSymbol | string
  height?: string
  width?: string
  className?: string
  onClick?: () => void
}

export const Icon = (props: IconProps) => {
  if (typeof props.sprite === "string") {
    const styles = { height: props.height && `${props.height}px`, width: props.width && `${props.width}px` }
    return <img style={styles} height={props.height} width={props.width} className={props.className} src={props.sprite} />
  }

  return (
    <svg onClick={props.onClick} height={props.height} width={props.width} className={props.className} viewBox={props.sprite.viewBox}>
      <use href={props.sprite.url} />
    </svg>
  )
}
