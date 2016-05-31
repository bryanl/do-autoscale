import DS from 'ember-data';
import Inflector from 'ember-inflector';
import Ember from 'ember';

export default DS.RESTAdapter.extend({
  namespace: 'api',
  pathForType: function(type) {
    const inflector = Inflector.inflector;

    var withUnderscore = Ember.String.underscore(type);
    return inflector.pluralize(withUnderscore);
  }
});
