import * as React from 'react';

interface SuggesionBoxProps {
    top: number;
    left: number;
    shown: boolean;
    data: string[];
    itemSelected?: number;
}

export class SuggestionBox extends React.Component<SuggesionBoxProps, {}> {

    constructor(props : SuggesionBoxProps) {
        super(props);
    }
    
    public render(){
        return <div
        style={{
            position: "absolute",
            width: "200px",
            borderRadius: "6px",
            background: "white",
            boxShadow: "rgba(0, 0, 0, 0.4) 0px 1px 4px",
                         
            display: this.props.shown ? "block" : "none",
            top: this.props.top,
            left: this.props.left,
          }}
          >

          {
            this.props.data.map((user, index) => (
              <div
                style={{
                  padding: '5px 5px',
                  background: index === this.props.itemSelected ? '#eee' : ''
                }}
                key = {user}
              >
                { user }
              </div>
            ))  
          }
        </div>
    }



}