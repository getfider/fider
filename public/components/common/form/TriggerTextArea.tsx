import * as React from "react";
import getCaretCoordinates from "textarea-caret";
import { TextArea, TextAreaProps } from "./TextArea";

export interface TextAreaTriggerStart {
  top: number;
  left: number;
  key: string;
}

export interface TriggerTextAreaProps extends TextAreaProps {
  onTriggerStart?: (e: TextAreaTriggerStart) => void;
  onTriggerEnd?: () => void;
  onTriggerChange?: (text: string) => void;
  onTriggerArrow?: (arrow: string) => void;
  onTriggerSelected?: (triggerPosition: number, cursorPosition: number) => void;
}

interface TriggerTextAreaState {
  suggestionTriggered: boolean;
  suggestionTriggerPosition: number;
}

export class TriggerTextArea extends React.Component<TriggerTextAreaProps, TriggerTextAreaState> {
  private onKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    const { onTriggerEnd, onTriggerChange, onTriggerStart, onTriggerArrow, onTriggerSelected } = this.props;
    const { suggestionTriggered, suggestionTriggerPosition } = this.state;
    const { key, which, currentTarget } = e;

    if (onTriggerStart) {
      const { selectionStart } = e.currentTarget;

      if (suggestionTriggered) {
        // Backspace handling
        if (which === 8 && selectionStart <= suggestionTriggerPosition) {
          this.setState({
            suggestionTriggered: false,
            suggestionTriggerPosition: 0
          });

          if (onTriggerEnd) {
            setTimeout(() => onTriggerEnd(), 0);
          }
        }
        // Down arrow handling
        else if (which === 40 && onTriggerArrow) {
          e.preventDefault();
          setTimeout(() => {
            onTriggerArrow("down");
          }, 0);
        }
        // Up arrow handling
        else if (which === 38 && onTriggerArrow) {
          e.preventDefault();
          setTimeout(() => {
            onTriggerArrow("up");
          }, 0);
        }
        // Enter handling
        else if (which === 13 && onTriggerSelected) {
          const capturedText = currentTarget.value.substring(suggestionTriggerPosition, selectionStart + 1);
          e.preventDefault();
          setTimeout(() => {
            onTriggerSelected(suggestionTriggerPosition, selectionStart);
          }, 0);
          this.setState({
            suggestionTriggered: false,
            suggestionTriggerPosition: 0
          });
        } else {
          if (onTriggerChange) {
            setTimeout(() => {
              const capturedText = currentTarget.value.substring(suggestionTriggerPosition, selectionStart + 1);
              onTriggerChange(capturedText);
            });
          }
        }
      } else {
        if (key === "@") {
          const caretCoordinates = getCaretCoordinates(e.currentTarget, selectionStart);

          this.setState({
            suggestionTriggered: true,
            suggestionTriggerPosition: selectionStart + 1
          });

          setTimeout(() => {
            onTriggerStart({
              key,
              left: caretCoordinates.left,
              top: caretCoordinates.top + caretCoordinates.height + 5
            });
          });
        }
      }
    }
  };

  public componentDidMount = () => {
    this.setState({
      suggestionTriggered: false
    });
  };

  public render() {
    return <TextArea {...this.props} onKeyDown={this.onKeyDown} />;
  }
}
