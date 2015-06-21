import Ember from 'ember';

export default Ember.Route.extend({
  model: function() {
    return this.get('store').rawRequest({url: 'images?detailed=false'});
  }
});
