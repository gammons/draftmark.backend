var NoteCardView = React.createClass({
  render: function() {
    var href = "#/notes/" + this.props.id;
    return(
      <div className="note">
        <a href={href}>
          <h5>{this.props.title}</h5>
        </a>
      </div>
    );
  }
});

