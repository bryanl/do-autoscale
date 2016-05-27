import Model from 'ember-data/model';
import attr from 'ember-data/attr';

export default Model.extend({
  regions: attr(),
  sizes: attr(),
  keys: attr()
});
