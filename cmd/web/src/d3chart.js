var d3Chart = {};

d3Chart.elements = {};

d3Chart.create = function(el, state) {
  var boundingBox = d3.select(el).node().getBoundingClientRect();
  console.log(boundingBox)
  var height = boundingBox.height
  var width =  boundingBox.width
  var data  =  state.data
  var margin = {top: 20, right: 20, bottom: 20, left: 40}

  var svg = d3.select(el).append("svg")
    .attr("width", width + margin.left + margin.right)
    .attr("height", height + margin.top + margin.bottom)
    .attr("preserveAspectRatio", "xMinYMin meet")
    .append("g")
    .attr("transform", "translate(" + margin.left + "," + margin.top + ")");
  this.elements.svg = svg;

  svg.append("defs").append("clipPath")
    .attr("id", "clip")
    .append("rect")
    .attr("width", width)
    .attr("height", height);

  var x = d3.time.scale()
      .domain([d3.min(data, function(d) {return d.time}), 
               d3.max(data, function(d) {return d.time})])
      .range([0, width]);
  this.elements.x = x;

  var y = d3.scale.linear()
    .domain([2.2, 4.5])
    .range([height, 0]);
  this.elements.y = y

  var line = d3.svg.line()
    .x(function(d, i) { return x(d.time); })
    .y(function(d, i) { return y(d.value); });
  this.elements.line = line;

  xAxis = d3.svg.axis().scale(x).orient("top");

  svg.append("g")
    .attr("class", "x axis")
    .call(xAxis);

  yAxis = d3.svg.axis().scale(y).orient("left");

  svg.append("g")
  .attr("class", "y axis")
  .call(yAxis);

  svg.append("g")
    .attr("class", "circles")
    .selectAll("circle")
    .data(data)
    .enter()
    .append("circle")
    .attr("r", 3.5)
    .attr("cx", x(function(d) {return(d.time)}))
    .attr("cy", y(function(d) {return(d.value)}))

    this.elements.path = svg.append("g")
      .attr("clip-path", "url(#clip)")
      .append("path")
      .datum(state.data)
      .attr("class", "line")
      .attr("d", line);

  svg.append("g")
    .attr("class", "trendlines")

  return(this);
};

d3Chart.update = function(el, state) {

  var data = state.data;
  // if (data.length == 0) { return }

  var svg = d3.select(el).select('svg');

  // Re-compute the elements, and render the data points
  var x = this.elements.x;
  var min_x = d3.min(data, function(d) {return d.time});
  var max_x = d3.max(data, function(d) {return d.time});

  x.domain([min_x, max_x]);
  svg.select(".x.axis")
    .transition().duration(200).call(xAxis);

  var y = this.elements.y;
  var min_y = d3.min(data, function(d) {return d.value});
  var max_y = d3.max(data, function(d) {return d.value});
  y.domain([min_y, max_y]);
  svg.select(".y.axis")
    .transition().duration(200).call(yAxis);

  // redraw the line, and slide it to the left
  var path = svg.select("path.line"); 
  var line = this.elements.line;

  var tr = d3.min(data, function(d) {return d.time});

  svg.select(".circles").selectAll("circle")
    .data(data)
    .transition()
    .duration(0)
    .attr("r",2)
    .attr("cx", function(d) {return x(d.time)})
    .attr("cy", function(d) { return y(d.value)})

  svg.select(".circles").selectAll("circle")
    .data(data)
    .enter()
    .append("circle")
    .attr("r", 2)
    .attr("cx", function(d) { return x(d.time)})
    .attr("cy", function(d) { return y(d.value)})

  svg.select(".circles").selectAll("circle")
    .data(data)
    .exit()
    .remove()

  if (data.length == 0 ) {
    path.datum(data)
  } else {
    path
    .attr("d", line)
    .attr("transform", null)
    .transition()
    .duration(0)
    .ease("linear")
    .attr("transform", "translate(" + x(tr-1) + ",0)");
  }

  if (state.recording) {
    var slope = state.slope
    var intercept = state.intercept
    var base = state.now
    var y1 = slope * (min_x - base) + intercept;
    var y2 = slope * (max_x - base) + intercept;
    var trendData = [[min_x, y1, max_x, y2]];
  } else {
    var trendData = []
  }

    var trendline = svg.select(".trendlines")
    .selectAll(".trendline")
    .data(trendData);

    trendline.enter()
    .append("line")
    .attr("class", "trendline")
    .attr("x1", function(d) { return x(d[0]); })
    .attr("y1", function(d) { return y(d[1]); })
    .attr("x2", function(d) { return x(d[2]); })
    .attr("y2", function(d) { return y(d[3]); })
    .attr("stroke", "black")
    .attr("stroke-width", 2);

    trendline.transition()
    .attr("x1", function(d) { return x(d[0]); })
    .attr("y1", function(d) { return y(d[1]); })
    .attr("x2", function(d) { return x(d[2]); })
    .attr("y2", function(d) { return y(d[3]); })
    .attr("stroke", "black")
    .attr("stroke-width", 2);

    trendline.exit().remove();
};

d3Chart.destroy = function(el) {
  // Any clean-up would go here
};

// module.exports = d3Chart;
