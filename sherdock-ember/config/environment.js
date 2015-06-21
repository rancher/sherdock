/* jshint node: true */

module.exports = function(environment) {
  var ENV = {
    modulePrefix: 'sherdock-ember',
    environment: environment,
    baseURL: '/',
    locationType: 'auto',
    EmberENV: {
      FEATURES: {
        // Here you can enable experimental features on an ember canary build
        // e.g. 'with-controller': true
      }
    },

    APP: {
      endpoint: 'http://localhost:8080',
      // Here you can pass flags/options to your application instance
      // when it is created
    }
  };

  if (environment === 'development') {
    // ENV.APP.LOG_RESOLVER = true;
    // ENV.APP.LOG_ACTIVE_GENERATION = true;
    // ENV.APP.LOG_TRANSITIONS = true;
    // ENV.APP.LOG_TRANSITIONS_INTERNAL = true;
    // ENV.APP.LOG_VIEW_LOOKUPS = true;
  }

  if (environment === 'test') {
    // Testem prefers this...
    ENV.baseURL = '/';
    ENV.locationType = 'none';

    // keep test console output quieter
    ENV.APP.LOG_ACTIVE_GENERATION = false;
    ENV.APP.LOG_VIEW_LOOKUPS = false;

    ENV.APP.rootElement = '#ember-testing';
  }

  if (environment === 'production') {

  }

  // Override the endpoint with environment var
  var endpoint = process.env.SHERDOCK_ENDPOINT;
  if ( endpoint )
  {
    // variable can be an ip "1.2.3.4" -> http://1.2.3.4:8080
    // or a URL+port
    if ( endpoint.indexOf('http') !== 0 )
    {
      if ( endpoint.indexOf(':') === -1 )
      {
        endpoint = 'http://' + endpoint + ':8080';
      }
      else
      {
        endpoint = 'http://' + endpoint;
      }
    }

    ENV.APP.endpoint = endpoint;
  }

  return ENV;
};
