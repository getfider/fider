import * as React from 'react';
import { Tag } from '@fider/models';

interface TagsFilterProps {
  tags: Tag[];
  selectionChanged: (selected: number[]) => void;
}

interface TagsFilterState {
  selected: number[];
}

export class TagsFilter extends React.Component<TagsFilterProps, TagsFilterState> {
  private element: HTMLDivElement;

  constructor(props: TagsFilterProps) {
    super(props);
    this.state = {
      selected: []
    };
  }

  public componentDidMount() {
    $(this.element).dropdown({
      onAdd: (value: string) => {
        const id = parseInt(value, 10);
        const selected = this.state.selected.concat(id);
        this.setState({ selected });
        this.props.selectionChanged(selected);
      },
      onRemove: (value: string) => {
        const id = parseInt(value, 10);
        const idx = this.state.selected.indexOf(id);
        const selected = this.state.selected.splice(idx, 1) && this.state.selected;
        this.setState({ selected });
        this.props.selectionChanged(selected);
      }
    });
  }

  public render() {
    return (
      <div className={`tags-filter ${this.state.selected.length > 0 ? 'has-selection' : ''}`}>
        <div className={`ui multiple dropdown `} ref={(e) => this.element = e!}>
          <i className="filter icon"/>
          <span className="text">Filter by tag...</span>
          <div className="menu">
            <div className="ui icon search input">
              <i className="search icon"/>
              <input type="text" placeholder="Search tags..."/>
            </div>
            <div className="divider"/>
            <div className="header">
              <i className="tags icon"/>
              Tag
            </div>
            <div className="scrolling menu">
              {
                this.props.tags.map((t) => (
                  <div key={t.id} className="item" data-value={t.id}>
                    <div className="ui empty circular label" style={{backgroundColor: `#${t.color}`}} />
                    {t.name}
                  </div>
                ))
              }
            </div>
          </div>
        </div>
      </div>
    );
  }
}
