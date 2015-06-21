import Ember from 'ember';

export default Ember.Controller.extend({
  gcRegex: null,
  images: null,

  actions: {
    save: function() {
      this.get('store').rawRequest({
        url: 'config',
        method: 'POST',
        data:  {
          GCIntervalMinutes: parseInt(this.get('model.GCIntervalMinutes'),10),
          PullIntervalMinutes: parseInt(this.get('model.PullIntervalMinutes'),10),
          ImagesToNotGC: this.get('gcRegex').map((entry) => {return entry.value; }).filter((entry) => {return !!entry;}),
          ImagesToPull: this.get('images').map((entry) => {return entry.value; }).filter((entry) => {return !!entry;}),
        }
      });
    }
  },

  gcRegexDidChange: function() {
    var vals = this.get('gcRegex');
    for ( var i = 0 ; i < vals.get('length')-1 ; i++ )
    {
      var val = (vals.objectAt(i).value||'').trim();
      if ( !val )
      {
        vals.removeAt(i);
        i--;
      }
    }

    var last = vals.objectAt(vals.get('length')-1);
    if ( !last || (last && last.value) )
    {
      vals.pushObject({value: ''});
    }
  }.observes('gcRegex.@each.value'),

  imagesDidChange: function() {
    var vals = this.get('images');
    for ( var i = 0 ; i < vals.get('length')-1 ; i++ )
    {
      var val = (vals.objectAt(i).value||'').trim();
      if ( !val )
      {
        vals.removeAt(i);
        i--;
      }
    }

    var last = vals.objectAt(vals.get('length')-1);
    if ( !last || (last && last.value) )
    {
      vals.pushObject({value: ''});
    }
  }.observes('images.@each.value'),

});
