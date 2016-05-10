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
    <option value='F1-SS-1'>F1-SS1</option>
    <option value='F1-SQ-1'>F1-SS1</option>
    <option value='F9-SS-1'>F1-SS1</option>
    <option value='F9-SQ-1'>F1-SS1</option>
    <option value='F1-SS-2'>F1-SS2</option>
    <option value='F1-SQ-2'>F1-SS2</option>
    <option value='F9-SS-2'>F1-SS2</option>
    <option value='F9-SQ-2'>F1-SS2</option>
    <option value='F1-SS-3'>F1-SS3</option>
    <option value='F1-SQ-3'>F1-SS3</option>
    <option value='F9-SS-3'>F1-SS3</option>
    <option value='F9-SQ-3'>F1-SS3</option>
    <option value='F1-SS-4'>F1-SS4</option>
    <option value='F1-SQ-4'>F1-SS4</option>
    <option value='F9-SS-4'>F1-SS4</option>
    <option value='F9-SQ-4'>F1-SS4</option>
    <option value='F1-SS-5'>F1-SS5</option>
    <option value='F1-SQ-5'>F1-SS5</option>
    <option value='F9-SS-5'>F1-SS5</option>
    <option value='F9-SQ-5'>F1-SS5</option>'
    <option value="F80-s-2d">F80-s-2d</option>
    <option value="F80-s-14d">F80-s-14d</option>
    <option value="F80-s-28d">F80-s-28d</option>
    <option value="F80-s-0d">F80-s-0d</option>
    <option value="F80-d-2d">F80-d-2d</option>
    <option value="F80-d-14d">F80-d-14d</option>
    <option value="F80-d-28d">F80-d-28d</option>
    <option value="F80-d-0d">F80-d-0d</option>
    <option value="F83-s-2d">F83-s-2d</option>
    <option value="F83-s-14d">F83-s-14d</option>
    <option value="F83-s-28d">F83-s-28d</option>
    <option value="F83-s-0d">F83-s-0d</option>
    <option value="F83-d-2d">F83-d-2d</option>
    <option value="F83-d-14d">F83-d-14d</option>
    <option value="F83-d-28d">F83-d-28d</option>
    <option value="F83-d-0d">F83-d-0d</option>
    <option value="F85-s-2d">F85-s-2d</option>
    <option value="F85-s-14d">F85-s-14d</option>
    <option value="F85-s-28d">F85-s-28d</option>
    <option value="F85-s-0d">F85-s-0d</option>
    <option value="F85-d-2d">F85-d-2d</option>
    <option value="F85-d-14d">F85-d-14d</option>
    <option value="F85-d-28d">F85-d-28d</option>
    <option value="F85-d-0d">F85-d-0d</option>
    <option value="F85-d-0d">F85-d-0d</option>
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
