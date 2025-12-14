import React, { useEffect, useState } from "react"
import { ReactionCount } from "@fider/models"
import { Icon } from "@fider/components"
import ReactionAdd from "@fider/assets/images/reaction-add.svg"
import { HStack } from "@fider/components/layout"
import { classSet } from "@fider/services"
import { useFider } from "@fider/hooks"
import "./Reactions.scss"

interface ReactionsProps {
  emojiSelectorRef: React.RefObject<HTMLDivElement>
  toggleReaction: (emoji: string) => void
  reactions?: ReactionCount[]
}

const availableEmojis = ["👍", "👎", "😄", "🎉", "😕", "❤️", "🚀", "👀"]

export const Reactions: React.FC<ReactionsProps> = ({ emojiSelectorRef, toggleReaction, reactions }) => {
  const fider = useFider()
  const [isEmojiSelectorOpen, setIsEmojiSelectorOpen] = useState(false)

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (emojiSelectorRef.current && !emojiSelectorRef.current.contains(event.target as Node)) {
        setIsEmojiSelectorOpen(false)
      }
    }

    document.addEventListener("click", handleClickOutside)
    return () => {
      document.removeEventListener("click", handleClickOutside)
    }
  }, [])

  return (
    <div ref={emojiSelectorRef}>
      <HStack spacing={2} align="center" className="mt-2 c-reactions relative">
        {fider.session.isAuthenticated && (
          <>
            <span
              onClick={() => setIsEmojiSelectorOpen(!isEmojiSelectorOpen)}
              className="c-reactions-add-reaction relative text-gray-600 clickable inline-flex items-center px-1 py-1 rounded-full text-xs bg-gray-100 hover:bg-gray-200"
            >
              <Icon width="18" height="18" sprite={ReactionAdd} className="" />
            </span>
            {isEmojiSelectorOpen && (
              <div className="c-reactions-emojis p-2 absolute bg-white border rounded shadow-lg">
                {availableEmojis.map((emoji) => (
                  <a
                    key={emoji}
                    className="clickable p-2 hover:bg-gray-100"
                    onClick={() => {
                      toggleReaction(emoji)
                      setIsEmojiSelectorOpen(false)
                    }}
                  >
                    {emoji}
                  </a>
                ))}
              </div>
            )}
          </>
        )}
        {reactions !== undefined && (
          <>
            {reactions.map((reaction) => (
              <span
                key={reaction.emoji}
                {...(fider.session.isAuthenticated && { onClick: () => toggleReaction(reaction.emoji) })}
                className={classSet({
                  "inline-flex items-center px-2 py-1 rounded-full text-xs": true,
                  "bg-blue-100": reaction.includesMe,
                  "bg-gray-100": !reaction.includesMe,
                  "clickable hover:bg-blue-200": fider.session.isAuthenticated && reaction.includesMe,
                  "clickable hover:bg-gray-200": fider.session.isAuthenticated && !reaction.includesMe,
                })}
              >
                {reaction.emoji} <span className="ml-1 text-semibold">{reaction.count}</span>
              </span>
            ))}
          </>
        )}
      </HStack>
    </div>
  )
}
