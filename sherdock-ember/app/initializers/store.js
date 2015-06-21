import Ember from 'ember';

function ajaxPromise(opt) {
  var promise = new Ember.RSVP.Promise(function(resolve,reject) {
    Ember.$.ajax(opt).then(success,fail);

    function success(body, textStatus, xhr) {
      Ember.run(function() {
        resolve(body, 'AJAX Response: '+ opt.url + '(' + xhr.status + ')');
      });
    }

    function fail(xhr, textStatus, err) {
      Ember.run(function() {
        reject({xhr: xhr, textStatus: textStatus, err: err}, 'AJAX Error:' + opt.url + '(' + xhr.status + ')');
      });
    }
  },'Raw AJAX Request: '+ opt.url);

  return promise;
}

var Store = Ember.Object.extend({
  //baseUrl: 'http://192.168.59.104:8080/api',
  baseUrl: '/api',

  rawRequest: function(opt) {
    var url = opt.url;
    if ( url.indexOf('http') !== 0 && url.indexOf('/') !== 0 )
    {
      url = this.get('baseUrl').replace(/\/\+$/,'') + '/' + url;
    }

    opt.url = url;

    if ( opt.data )
    {
      opt.contentType = 'application/json';
      opt.data = JSON.stringify(opt.data);
    }

    var promise = ajaxPromise(opt);
    return promise;
  },
});

export function initialize(container, application) {
  var store = Store.create({});

  container.register('store:main',   store,  {instantiate: false});
  application.inject('controller',  'store', 'store:main');
  application.inject('route',       'store', 'store:main');
}

export default {
  name: 'store',
  initialize: initialize
};
