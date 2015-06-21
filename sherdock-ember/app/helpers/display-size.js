import Ember from 'ember';

var suffixes = ['B','KiB','MiB','GiB','TiB'];

export default Ember.Handlebars.makeBoundHelper(function(bytes) {
  var i = 0;
  var val = bytes;

  while ( val > 1024 && i < suffixes.length )
  {
    val /= 1024;
    i++;
  }

  return Math.round(val*10)/10 +' '+ suffixes[i];
});
