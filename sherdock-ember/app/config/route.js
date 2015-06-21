import Ember from 'ember';

export default Ember.Route.extend({
  actions: {
    refresh: function() {
      this.refresh();
    },
  },

  model: function() {
    return this.get('store').rawRequest({url: 'config'});
  }
});
