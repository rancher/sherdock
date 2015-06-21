import Ember from 'ember';

export default Ember.Route.extend({
  actions: {
    refresh: function() {
      this.refresh();
    },
  },

  model: function() {
    return this.get('store').rawRequest({url: 'config'});
  },

  setupController: function(controller, model) {
    controller.set('model', model);
    controller.set('gcRegex', (model.ImagesToNotGC||[]).map((val) => {
      return {value: val};
    }));
    controller.set('images', (model.ImagesToPull||[]).map((val) => {
      return {value: val};
    }));
    controller.gcRegexDidChange();
    controller.imagesDidChange();
  },
});
