var Result =  React.createClass({displayName: "Result",

  render: function() {
    if (this.props.recording){
    return(
      React.createElement("div", {className: "results"}, 
        React.createElement("div", {className: "row"}, 
        React.createElement("div", {className: "col-md-6"}, 
          React.createElement("ul", {className: "flux list-unstyled"}, 
          React.createElement("li", null, "N2O: ",  this.props.n2o, " "), 
          React.createElement("li", null, "CO2: ",  this.props.co2, " "), 
          React.createElement("li", null, "CH4: ",  this.props.ch4, " ")
          )
        ), 
          React.createElement("div", {className: "col-md-3"}, 
          React.createElement("button", {type: "button", className: "btn btn-primary btn-block btn-lg", onClick: this.props.handleSave}, "Save")
          ), 
          React.createElement("div", {className: "col-md-3"}, 
          React.createElement("button", {className: "btn btn-default btn-block btn-lg", type: "button", onClick: this.props.handleCancel}, "Cancel")
          )
        )
      )
    )
    } else {
      return(
      React.createElement("div", {className: "results"}, 
        React.createElement("form", {className: "form"}, 
        React.createElement("div", {className: "form-group"}, 
        React.createElement("label", {for: "location"}, "Location"), 
          React.createElement("select", {id: "location"}, 
            React.createElement("option", {value: "T1R1"}, "T1R1"), 
            React.createElement("option", {value: "T1R2"}, "T1R2")
          )
          ), 
          React.createElement("div", {className: "form-group"}, 
          React.createElement("label", {for: "height"}, "Height"), 
          React.createElement("input", {id: "height", type: "number"}), " cm"
        )
        ), 
        React.createElement("button", {type: "button", className: "btn btn-primary btn-lg", onClick: this.props.handleRecord}, "Record")
        )
      )
  }
  }
})
