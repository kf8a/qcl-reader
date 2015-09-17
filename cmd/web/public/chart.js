// var d3Chart = require('d3Chart');

var Chart = React.createClass({displayName: "Chart",
  propTypes: {
    data: React.PropTypes.array
  },

  componentDidMount: function() {
    var el = this.getDOMNode();

    this.d3Chart = d3Chart.create(el, this.getChartState());
  },

  componentDidUpdate: function() {
    var el = this.getDOMNode();
    this.d3Chart.update(el, this.getChartState());
  },

  getChartState: function() {
    var el = this.getDOMNode();
    return {
      data: this.props.data,
      now: this.props.now,
      slope: this.props.slope,
      intercept: this.props.intercept,
      recording: this.props.recording
    };
  },

  componentWillUnmount: function() {
    var el = this.getDOMNode();
    this.d3Chart.destroy(el);
  },

  render: function() {
    return (
      React.createElement("div", {className: "Chart"}
      )
    );
  }
});

// module.exports = Chart;
