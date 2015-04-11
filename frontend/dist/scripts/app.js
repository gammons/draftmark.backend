'use strict';

var IndexView = React.createClass({
  displayName: 'IndexView',

  getInitialState: function getInitialState() {
    return { notes: [], router: null, view: 'index' };
  },

  componentDidMount: function componentDidMount() {
    this.getData((function (data) {
      this.setState({ notes: data });
      this.setupRouter();
    }).bind(this));
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

  getData: function getData(successFn) {
    $.ajax({
      url: '/notes.json',
      dataType: 'json',
      success: successFn,
      error: function error(xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }
    });
  },

  setupRouter: function setupRouter() {
    var router = new Router().init();
    router.on(/notes\/?((\w|.)*)/, this.viewNoteView);
    router.on(/asdf/, function () {
      console.log('got asdf');
    });
    router.on('/', this.indexView);

    this.setState({ notes: this.state.notes, router: router });
  },
  indexView: function indexView() {
    this.setState({ view: 'index' });
  },
  viewNoteView: function viewNoteView(path) {
    console.log('view note view', path);
    this.setState({ view: 'note', selectedNotePath: path });
  },

  render: function render() {
    if (this.state.view == 'index') {
      return this.renderIndex();
    } else if (this.state.view == 'note') {
      return this.renderNote();
    }
  },

  renderIndex: function renderIndex() {
    var showNote = this.showNote;
    var notes = _.map(this.state.notes, function (note) {
      return React.createElement(
        'li',
        { key: note.id },
        React.createElement(NoteCardView, { className: 'noteCardView', path: note.path, title: note.title, content: note.content, clickNote: showNote })
      );
    });

    return React.createElement(
      'div',
      null,
      React.createElement(
        'h1',
        null,
        'Notes'
      ),
      React.createElement(
        'ul',
        { className: 'small-block-grid-2 medium-block-grid-4 large-block-grid-6' },
        notes
      )
    );
  },

  renderNote: function renderNote() {
    var selectedNote = this.getNote(this.state.selectedNotePath);
    return React.createElement(
      'div',
      { className: 'view-note' },
      React.createElement(NoteView, { path: selectedNote.path, title: selectedNote.title, pusher: this.state.pusher })
    );
  },
  getNote: function getNote(path) {
    path = '/' + path;
    return _.find(this.state.notes, function (note) {
      return note.path == path;
    });
  }

});
"use strict";

var NoteView = React.createClass({
  displayName: "NoteView",

  getInitialState: function getInitialState() {
    return { content: "" };
  },
  componentDidMount: function componentDidMount() {
    this.getData((function (data) {
      this.setState({ content: data });
    }).bind(this));

    //var channel = this.props.pusher.subscribe('updates');
    // channel.bind('update', function(data) {
    //   this.getData(function(data) {
    //     this.setState({content: data.content});
    //   }.bind(this));
    // }.bind(this));
  },
  getData: function getData(successFn) {
    $.ajax({
      url: "/content" + this.props.path,
      success: successFn,
      error: (function (xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }).bind(this)
    });
  },
  render: function render() {
    var rawMarkup = marked(this.state.content.toString());
    return React.createElement(
      "div",
      { className: "row" },
      React.createElement(
        "div",
        { className: "large-12 columns" },
        React.createElement(
          "a",
          { href: "#" },
          "Back to notes list"
        ),
        React.createElement("span", { dangerouslySetInnerHTML: { __html: rawMarkup } })
      )
    );
  }
});
"use strict";

var NoteCardView = React.createClass({
  displayName: "NoteCardView",

  render: function render() {
    var href = "#/notes" + this.props.path;
    return React.createElement(
      "div",
      { className: "note" },
      React.createElement(
        "a",
        { href: href },
        React.createElement(
          "h5",
          null,
          this.props.title
        )
      )
    );
  }
});
'use strict';

React.render(React.createElement(IndexView, null), document.getElementById('index'));