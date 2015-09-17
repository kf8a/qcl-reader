'use strict';

var proto = TimeSeries.prototype;

function TimeSeries(el, opts) {
  if (!(this instanceof TimeSeries)) return new Timeseris(el, opts);

  this._el = el;
  this._initGraph();
  this._render();
}

proto.update = function(data) {
  var self = this;
  
}

proto.render = function() {

}
