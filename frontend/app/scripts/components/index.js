var IndexView = React.createClass({
  getInitialState: function() {
    return {notes: [], router:  null, view: 'index' }
  },

  componentDidMount: function() {
    this.getData(function(data) {
      this.setState({notes: data});
      this.setupRouter();
    }.bind(this));
    //this.setupPusher();
  },

  // setupPusher: function() {
  //   this.setState({pusher: new Pusher('2cdc6bc2a2113ae973d8') })
  //   var channel = this.state.pusher.subscribe('updates');
  //   channel.bind('update', function(data) {
  //     this.getData(function(data) {
  //       this.setState({notes: data});
  //     }.bind(this));
  //   }.bind(this));
  // },

  getData: function(successFn) {
    $.ajax({
      url: "/notes.json",
      dataType: 'json',
      success: successFn,
      error: function(xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }
    });
  },

  setupRouter: function() {
    var router = new Router().init();
    router.on(/notes\/?((\w|.)*)/, this.viewNoteView)
    router.on(/asdf/, function() { console.log("got asdf"); });
    router.on("/", this.indexView)

    this.setState({notes: this.state.notes, router: router});
  },
  indexView: function() {
    this.setState({view: "index"});
  },
  viewNoteView: function(path) {
    console.log("view note view", path)
    this.setState({view: "note", selectedNotePath: path});
  },

  render: function() {
    if (this.state.view == 'index') {
      return this.renderIndex();
    } else if (this.state.view == 'note') {
      return this.renderNote();
    }
  },

  renderIndex: function() {
    var showNote = this.showNote;
    var notes = _.map(this.state.notes, function(note) {
      return(
        <li key={note.id}>
          <NoteCardView className="noteCardView" path={note.path} title={note.title} content={note.content} clickNote={showNote} />
        </li>
      );
    });

    return(
      <div>
        <h1>Notes</h1>
        <ul className="small-block-grid-2 medium-block-grid-4 large-block-grid-6">
          {notes}
        </ul>
      </div>
    );
  },

  renderNote: function() {
    var selectedNote = this.getNote(this.state.selectedNotePath);
    return (
      <div className="view-note">
      <NoteView path={selectedNote.path} title={selectedNote.title} pusher={this.state.pusher} />
      </div>
    );
  },
  getNote: function(path) {
    path = "/"+path;
    return _.find(this.state.notes, function(note) { return note.path == path });
  }

});
