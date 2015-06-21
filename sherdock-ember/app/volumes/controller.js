import Ember from 'ember';

export default Ember.ArrayController.extend({
  actions: {
    removeAll: function() {
      this.get('store').rawRequest({
        url: 'images',
        method: 'DELETE',
      }).then(() => {
        this.send('refresh');
      });
    },

    remove: function(volumeId) {
      this.get('store').rawRequest({
        url: 'images/'+volumeId,
        method: 'DELETE',
      }).then(() => {
        this.send('refresh');
      });
    }
  }
});
