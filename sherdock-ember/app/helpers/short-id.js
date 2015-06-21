import Ember from 'ember';

export default Ember.Handlebars.makeBoundHelper(function(str) {
  return (str||'').substr(0,8);
});
