var LocationSelect = React.createClass({
  handleChange: function(event) {
    this.setState({value: event.target.value});
    this.props.updatePlot(event.target.value);
  },

  render: function() {
    return(
    <div className="form-group">
    <label for="location">Location</label>
    <select className="location form-control input-lg" value={ this.props.value } onChange={this.handleChange} >

<option value="1r1-2d">"1r1-2d"</option>
<option value="G1r1-7d">"G1r1-7d"</option>
<option value="G1r1-14d">"G1r1-14d"</option>
<option value="G1r1-28d">"G1r1-28d"</option>
<option value="G1r2-2d">"G1r2-2d"</option>
<option value="G1r2-7d">"G1r2-7d"</option>
<option value="G1r2-14d">"G1r2-14d"</option>
<option value="G1r2-28d">"G1r2-28d"</option>
<option value="G1r3-2d">"G1r3-2d"</option>
<option value="G1r3-7d">"G1r3-7d"</option>
<option value="G1r3-14d">"G1r3-14d"</option>
<option value="G1r3-28d">"G1r3-28d"</option>
<option value="G1r4-2d">"G1r4-2d"</option>
<option value="G1r4-7d">"G1r4-7d"</option>
<option value="G1r4-14d">"G1r4-14d"</option>
<option value="G1r4-28d">"G1r4-28d"</option>
<option value="G1r5-2d">"G1r5-2d"</option>
<option value="G1r5-7d">"G1r5-7d"</option>
<option value="G1r5-14d">"G1r5-14d"</option>
<option value="G1r5-28d">"G1r5-28d"</option>
<option value="G5r1-2d">"G5r1-2d"</option>
<option value="G5r1-7d">"G5r1-7d"</option>
<option value="G5r1-14d">"G5r1-14d"</option>
<option value="G5r1-28d">"G5r1-28d"</option>
<option value="G5r2-2d">"G5r2-2d"</option>
<option value="G5r2-7d">"G5r2-7d"</option>
<option value="G5r2-14d">"G5r2-14d"</option>
<option value="G5r2-28d">"G5r2-28d"</option>
<option value="G5r3-2d">"G5r3-2d"</option>
<option value="G5r3-7d">"G5r3-7d"</option>
<option value="G5r3-14d">"G5r3-14d"</option>
<option value="G5r3-28d">"G5r3-28d"</option>
<option value="G5r4-2d">"G5r4-2d"</option>
<option value="G5r4-7d">"G5r4-7d"</option>
<option value="G5r4-14d">"G5r4-14d"</option>
<option value="G5r4-28d">"G5r4-28d"</option>
<option value="G5r5-2d">"G5r5-2d"</option>
<option value="G5r5-7d">"G5r5-7d"</option>
<option value="G5r5-14d">"G5r5-14d"</option>
<option value="G6r5-28d">"G5r5-28d"</option>
    </select>
    </div>
    )
  }
});

var Result =  React.createClass({

  render: function() {
    if (this.props.recording){
      return(
        <div className='results'>
        <div className="row">
        <div className="col-md-3, col-sm-3 col-xs-3">
          <h2>{ this.props.plot}</h2>
        </div>
        <div className="col-md-3 col-sm-3 col-xs-3">
          <h2>{this.props.height} cm</h2>
        </div>
        </div>
        <div className="row">
        <div className="col-md-6 col-sm-6 col-xs-6">
        <ul className='flux list-unstyled'>
        <li>CO2: { this.props.co2 }</li>
        <li>N2O: { this.props.n2o }</li>
        <li>CH4: { this.props.ch4 }</li>
        <li>N2O Doubling time: {this.props.n2o_intercept/this.props.n2o/60} min</li>
        </ul>
        </div>
        <div className="col-md-3 col-sm-3 col-xs-3">
        <button type='button' className="btn btn-primary btn-block btn-lg" onClick={this.props.handleSave}>Save</button>
        </div>
        <div className="col-md-3 col-sm-3 col-xs-3">
        <button className='btn btn-default btn-block btn-lg' type='button' onClick={this.props.handleCancel}>Cancel</button>
        </div>
        </div>
        </div>
      )
    } else {
      return(
        <div className='results'>
        <form className='form'>
        <LocationSelect value={this.props.plot} updatePlot={this.props.updatePlot} />
          <div className="form-group">
          <label for="height">Height</label>
          <input id="height" type='number' min="0" max="50" step="0.1" className="form-control input-default" value={this.props.height} onChange={this.props.updateHeight}/> cm
        </div>
        </form>
        <button type='button' className="btn btn-primary btn-default" onClick={this.props.handleRecord}>Record</button>
        </div>
      )
  }
  }
})
