var Chart = React.createClass({displayName: "Chart",
  render: function() {
    return( React.createElement("div", null, "I am a chart"));
  }
});

var QCL = React.createClass({displayName: "QCL",
  render: function() {
    return(
      React.createElement("div", {id: "app"}, 
      React.createElement(Chart, null), 
      React.createElement(Chart, null)
        )
    )
  }
})
  React.render(React.createElement(QCL, null), document.getElementById("example"));
