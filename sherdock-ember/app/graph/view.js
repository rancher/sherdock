import Ember from 'ember';

export default Ember.View.extend({
  didInsertElement: function() {
    this._super();
    this.initGraph();
  },

  initGraph: function() {
    var outer = d3.select("#graph svg");
    var inner = outer.select("g");
    var zoom = d3.behavior.zoom().on("zoom", function() {
       inner.attr("transform", "translate(" + d3.event.translate + ")" +
                                   "scale(" + d3.event.scale + ")");
    });

    outer.call(zoom);

    var g = new dagreD3.graphlib.Graph().setGraph({
      rankdir: "LR",
      nodesep: 20,
      ranksep: 50,
      marginx: 10,
      marginy: 10
    });

    var render = new dagreD3.render();

    this.setProperties({
      graphZoom: zoom,
      graphOuter: outer,
      graphInner: inner,
      graphRender: render,
      graph: g
    });

    this.updateGraph();
  },

  updateGraph: function() {
    var g = this.get('graph');

    var edgeOpts = {
      arrowhead: 'vee',
      lineInterpolate: 'linear',
    };

    var seenIds = {};
    var volumes = this.get('context');
    volumes.forEach(function(volume) {
      var tag = (volume.RepoTags||[])[0];
      if ( tag === '<none>:<none>' )
      {
        tag = volume.Id.substr(0,8)+'...';
      }

      g.setNode(volume.Id, {
        label: tag,
      });

      seenIds[volume.Id] = true;
    });

    volumes.forEach(function(volume) {
      if ( seenIds[volume.ParentId] && seenIds[volume.Id] )
      {
        g.setEdge(volume.ParentId, volume.Id, edgeOpts);
        console.log('Edge from', volume.ParentId.substr(0,8),'to',volume.Id.substr(0,8));
      }
    });

    this.renderGraph();
  },

  renderGraph: function() {
    var zoom = this.get('graphZoom');
    var render = this.get('graphRender');
    var inner = this.get('graphInner');
    var outer = this.get('graphOuter');
    var g = this.get('graph');

    inner.call(render, g);

    // Zoom and scale to fit
    var zoomScale = zoom.scale();
    var graphWidth = g.graph().width;
    var graphHeight = g.graph().height;
    var width = $('#svg').width();
    var height = $('#svg').height();
    zoomScale = Math.min(2.0, Math.min(width / graphWidth, height / graphHeight));
    var translate = [(width/2) - ((graphWidth*zoomScale)/2), (height/2) - ((graphHeight*zoomScale)/2)];
    zoom.translate(translate);
    zoom.scale(zoomScale);
    zoom.event(outer);
  },
});
