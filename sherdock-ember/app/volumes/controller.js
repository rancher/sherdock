import Ember from 'ember';

export default Ember.ArrayController.extend({
  actions: {
    removeAll: function() {
      this.get('store').rawRequest({
        url: 'volumes',
        method: 'DELETE',
      }).then(() => {
        this.send('refresh');
      });
    },

    remove: function(volumeId) {
      this.get('store').rawRequest({
        url: 'volumes/'+volumeId,
        method: 'DELETE',
      }).then(() => {
        this.send('refresh');
      });
    }
  }
});
