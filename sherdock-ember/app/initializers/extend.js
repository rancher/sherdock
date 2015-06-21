import Ember from 'ember';

export function initialize(/*container, application*/) {
  Ember.TextField.reopen({
    attributeBindings: ['style'],
  });
}

export default {
  name: 'extend',
  initialize: initialize
};
