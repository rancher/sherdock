import Ember from 'ember';

export default Ember.View.extend({
  didInsertElement: function() {
    this._super();
    this.initGraph();
  },

  initGraph: function() {
    var edges = [];
    var nodes = [];
    var seenIds = {};
    var volumes = this.get('context');

    volumes.forEach(function(volume) {
      var tag = (volume.RepoTags||[])[0];
      if ( tag === '<none>:<none>' )
      {
        tag = volume.Id.substr(0,8)+'...';
      }

      if ( !volume.Id || !volume.ParentId )
      {
        return;
      }

      console.log('Adding', tag, volume.Id, volume.ParentId);
      nodes.push({ data: { id: volume.Id, name: tag } });
      seenIds[volume.Id] = true;
    });

    volumes.forEach(function(volume) {
      if ( seenIds[volume.ParentId] && seenIds[volume.Id] )
      {
        console.log('Edge from', volume.ParentId.substr(0,8),'to',volume.Id.substr(0,8),'for',volume.RepoTags);
        edges.push({ data: { source: volume.ParentId, target: volume.Id } });
      }
    });

    $('#graph').cytoscape({
      style: cytoscape.stylesheet()
        .selector('node')
          .css({
            'content': 'data(name)',
            'text-valign': 'center',
            'color': 'white',
            'text-outline-width': 2,
            'text-outline-color': '#888'
          })
        .selector('edge')
          .css({
            'target-arrow-shape': 'triangle'
          })
        .selector(':selected')
          .css({
            'background-color': 'black',
            'line-color': 'black',
            'target-arrow-color': 'black',
            'source-arrow-color': 'black'
          })
        .selector('.faded')
          .css({
            'opacity': 0.25,
            'text-opacity': 0
          }),
      elements: {
        nodes: nodes,
        edges: edges,
      },

      layout: {
        name: 'grid',
        padding: 10
      },

      // on graph initial layout done (could be async depending on layout...)
      ready: function(){
        window.cy = this;

        // giddy up...
        cy.elements().unselectify();

        cy.on('tap', 'node', function(e){
          var node = e.cyTarget; 
          var neighborhood = node.neighborhood().add(node);

          cy.elements().addClass('faded');
          neighborhood.removeClass('faded');
        });

        cy.on('tap', function(e){
          if( e.cyTarget === cy ){
            cy.elements().removeClass('faded');
          }
        });
      }
    });
  },
});
