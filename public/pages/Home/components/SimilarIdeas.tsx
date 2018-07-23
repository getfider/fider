import * as React from "react";
import { Idea, Tag, CurrentUser } from "@fider/models";
import { Heading, Loader } from "@fider/components";
import { ListIdeas } from "./ListIdeas";
import { actions } from "@fider/services";

interface SimilarIdeasProps {
  title: string;
  tags: Tag[];
  user?: CurrentUser;
}

interface SimilarIdeasState {
  title: string;
  ideas: Idea[];
  loading: boolean;
}

export class SimilarIdeas extends React.Component<SimilarIdeasProps, SimilarIdeasState> {
  constructor(props: SimilarIdeasProps) {
    super(props);
    this.state = {
      title: props.title,
      loading: !!props.title,
      ideas: []
    };
  }

  public static getDerivedStateFromProps(nextProps: SimilarIdeasProps, prevState: SimilarIdeasState) {
    if (nextProps.title !== prevState.title) {
      return {
        loading: true,
        title: nextProps.title
      };
    }
    return null;
  }
  public componentDidMount() {
    this.loadSimilarIdeas();
  }

  private timer?: number;
  public componentDidUpdate() {
    window.clearTimeout(this.timer);
    this.timer = window.setTimeout(this.loadSimilarIdeas, 200);
  }

  private loadSimilarIdeas = () => {
    if (this.state.loading) {
      actions.searchIdeas({ query: this.state.title }).then(x => {
        this.setState({ loading: false, ideas: x.data });
      });
    }
  };

  public render() {
    return (
      <>
        <Heading
          title="Similar ideas"
          subtitle="Consider voting on existing ideas instead of posting a new one."
          icon="lightbulb outline"
          size="small"
          dividing={true}
        />
        {this.state.loading ? (
          <Loader />
        ) : (
          <ListIdeas
            ideas={this.state.ideas}
            tags={this.props.tags}
            emptyText={`No similar ideas matched '${this.props.title}'.`}
          />
        )}
      </>
    );
  }
}
