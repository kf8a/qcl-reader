'use strict';

// var React = require('react');
// var Chart = require('chart');

var App = React.createClass({
  getInitialState: function() {
    return { 
      ch4: [],
      n2o: [],
      co2: [],
      recording: false,
      n2o_flux: 0,
      co2_flux: 0,
      ch4_flux: 0,
      n2o_intercept: 0,
      co2_intercept: 0,
      ch4_intercept: 0,
      now: new Date()
    };
  },

  componentDidMount: function() {
    this.requestData();
  },

  render: function() {
    return (
      <div className="flux">
      <div className="row">
      <Chart 
        now={ this.state.now}
        slope={this.state.co2_flux}
        intercept={this.state.co2_intercept}
        recording={this.state.recording}
        data={this.state.co2} />
        </div>
        <div className="row">
      <Chart 
        now={this.state.now} 
        slope={this.state.n2o_flux}
        intercept={this.state.n2o_intercept}
        recording={this.state.recording}
        data={this.state.n2o} />
        </div>
        <div className="row">
      <Chart 
        now={this.state.now }
        slope={this.state.ch4_flux}
        intercept={this.state.ch4_intercept}
        recording={this.state.recording}
        data={this.state.ch4} />
        </div>
        <div className="row">
      <Result 
        handleRecord={this.handleRecord}
        handleSave={this.handleSave}
        handleCancel={this.handleCancel}
        recording={this.state.recording}
        n2o={this.state.n2o_flux} 
        co2={this.state.co2_flux} 
        ch4={this.state.ch4_flux}/>
        </div>
      </div>
    );
  },

  resetData: function() {
    this.setState({co2:[], ch4:[], n2o:[]});
  },

  handleCancel: function(e) {
    e.preventDefault();
    this.resetData();
    this.setState({n2o_flux: null,
                  co2_flux: null,
                  ch4_flux: null,
                  recording: false});
  },

  handleSave: function(e) {
    e.preventDefault();
    // send data back to server
    jQuery.ajax({
      type: "POST",
      url: "/save",
      data: JSON.stringify({"co2": this.state.co2, "ch4": this.state.ch4, "n2o": this.state.n2o}),
      dataType: 'json'
    });

    this.resetData();
    this.setState({n2o_flux: null,
                  co2_flux: null,
                  ch4_flux: null,
                  recording: false})
  },

  handleRecord: function(e) {
    e.preventDefault();
    this.resetData();
    jQuery.ajax({
      type: "POST",
      url: "/record",
      data: "",
      dataType: 'json'
    });
    this.setState({recording: true,
                  now: new Date()})
  },

  prepareData: function(datum) {
      // this.state.data.push(datum);
      // var myco2 = this.state.co2.concat({ time: datum.obs_time, value: datum.CO2_ppm})
      this.state.co2.push({time: datum.obs_time, value: datum.CO2_ppm})
      this.state.n2o.push({time: datum.obs_time, value: datum.N2O_dry_ppm})
      this.state.ch4.push({time: datum.obs_time, value: datum.CH4_dry_ppm})

      if (this.state.co2.length > 200) {
        this.state.co2.shift();
        this.state.n2o.shift();
        this.state.ch4.shift()
      };

      this.setState({ co2: this.state.co2,
                    ch4: this.state.ch4, 
                    n2o: this.state.n2o});
  },

  computeFluxes: function() {
    if (this.state.n2o.length < 2) {return; };
    var n2o_flux = this.computeFlux(this.state.n2o);
    var co2_flux = this.computeFlux(this.state.co2);
    var ch4_flux = this.computeFlux(this.state.ch4);
    this.setState({n2o_flux: n2o_flux['slope'], 
                  n2o_intercept: n2o_flux['intercept'],
                  co2_flux: co2_flux['slope'], 
                  co2_intercept: co2_flux['intercept'],
                  ch4_flux: ch4_flux['slope'],
                  ch4_intercept: ch4_flux['intercept']});
  },

  computeFlux: function(data) {
    var x = data.map(function(e, now) {
      return (e.time.valueOf() - this.state.now.valueOf());
    }, this);
    var y = data.map(function(e) {
      return e.value;
    });
    var lr = this.linearRegression(x,y);
    return(lr);
  },

  linearRegression: function(x,y) {
    var lr = {};
    var n = y.length;
    var sum_x = 0;
    var sum_y = 0;
    var sum_xy = 0;
    var sum_xx = 0;
    var sum_yy = 0;

    for (var i = 0; i < y.length; i++) {

      sum_x += x[i];
      sum_y += y[i];
      sum_xy += (x[i]*y[i]);
      sum_xx += (x[i]*x[i]);
      sum_yy += (y[i]*y[i]);
    } 

    lr['slope'] = (n * sum_xy - sum_x * sum_y) / (n*sum_xx - sum_x * sum_x);
    lr['intercept'] = (sum_y - lr.slope * sum_x)/n;
    lr['r2'] = Math.pow((n*sum_xy - sum_x*sum_y)/Math.sqrt((n*sum_xx-sum_x*sum_x)*(n*sum_yy-sum_y*sum_y)),2);

    return lr;
  },

  requestData: function() {
    var socket = this.props.socketService;

    socket.onDataReceived(function (datum) {
      datum.obs_time = new Date(datum.obs_time);
      datum.time = new Date(datum.time);

      this.prepareData(datum);
      if (this.state.recording) {
        this.computeFluxes();
      }
    }.bind(this));
  }

});

var socketService = new SocketService({
    path: '/ws'
});

React.render(
  < App socketService={socketService} />, 
  document.getElementById('example')
)

