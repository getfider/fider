import * as React from "react";

import {TriggerTextArea, TextAreaTriggerStart, TextAreaProps} from "@fider/components/common";
import {SuggestionBox} from "@fider/components"
import {mentions} from "@fider/services";
import {User} from "@fider/models";

interface MentionableTextAreaProps extends TextAreaProps {
  onMention: (start: number, end:number, mention:string) => void;
} 

interface MentionableTextAreaState {
  suggestionBoxVisible: boolean;
  suggestionBoxTop: number;
  suggestionBoxLeft: number;
  suggestedUsers: User[];
  suggestionSelected: number;
}

export class MentionableTextArea extends React.Component<MentionableTextAreaProps,MentionableTextAreaState>{


  constructor(props : MentionableTextAreaProps){
    super(props);
    this.state = {
      suggestionBoxVisible : false,
      suggestionBoxTop : 0,
      suggestionBoxLeft: 0,
      suggestedUsers: [],
      suggestionSelected: 0
    };
  }



  private mentionStart = (e : TextAreaTriggerStart) => {
    this.setState({
      suggestionBoxLeft : e.left,
      suggestionBoxTop  : e.top,
      suggestionBoxVisible: true
    })
  }
  
  private mentionChange = (text : string) => {
    if (text.length > 0){
      mentions.get(text).then((response) => {
        this.setState({
          suggestedUsers : response.data
        });
        if (this.state.suggestionSelected > response.data.length) {
          this.setState({
            suggestionSelected : response.data.length
          });
        }

      })
    } else {
      this.setState({
        suggestionSelected: 0,
      });
    }
    
  }
  
  private mentionEnd = () => {
    this.setState({
      suggestionBoxVisible: false,
      suggestedUsers : []
    })
  }
  
  private suggestionSelected = (mentionStartIndex : number, cursorPosition: number) =>{
    if (this.props.onMention){
      const mention = this.state.suggestedUsers[this.state.suggestionSelected].name;
      this.props.onMention(mentionStartIndex, cursorPosition, mention)
    }
    this.setState({
      suggestionBoxVisible: false,
      suggestedUsers : []
    })
  }
  
  private mentionArrow = (selectEvent : string) => {
    switch (selectEvent){
      case ("up"): {
        const selected = (this.state.suggestionSelected - 1) % this.state.suggestedUsers.length;
        this.setState({suggestionSelected : selected})
        break;
      }
      case ("down"): {
        const selected = (this.state.suggestionSelected + 1) % this.state.suggestedUsers.length;
        this.setState({suggestionSelected : selected})
        break;
      }
    }
  };
  
  
  public render(){
    return <TriggerTextArea
    onTriggerStart={this.mentionStart}
    onTriggerChange={this.mentionChange}
    onTriggerEnd={this.mentionEnd}
    onTriggerSelected={this.suggestionSelected}
    onTriggerArrow={this.mentionArrow}
    {...this.props}
    >
      <SuggestionBox
      top={this.state.suggestionBoxTop}
      left={this.state.suggestionBoxLeft}
      shown={this.state.suggestionBoxVisible}
      data={this.state.suggestedUsers.map(u => u.name)}
      itemSelected={this.state.suggestionSelected}
      />
    </TriggerTextArea>
  }
}