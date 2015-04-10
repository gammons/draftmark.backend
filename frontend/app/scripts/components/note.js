var NoteView = React.createClass({
  getInitialState: function() {
      return({content: ""});
  },
  componentDidMount: function() {
    this.getData(function(data) {
      this.setState({content: data});
    }.bind(this));

    //var channel = this.props.pusher.subscribe('updates');
    // channel.bind('update', function(data) {
    //   this.getData(function(data) {
    //     this.setState({content: data.content});
    //   }.bind(this));
    // }.bind(this));
  },
  getData: function(successFn) {
    $.ajax({
      url: "/content"+this.props.path,
      success: successFn,
      error: function(xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }.bind(this)
    });
  },
  render: function() {
    var rawMarkup = marked(this.state.content.toString());
    return(
      <div className="row">
        <div className="large-12 columns">
          <a href="#">Back to notes list</a>
          <span dangerouslySetInnerHTML={{__html: rawMarkup}} />
        </div>
      </div>
    );
  }
});
