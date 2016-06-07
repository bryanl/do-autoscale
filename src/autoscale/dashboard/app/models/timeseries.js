import attr from 'ember-data/attr';
import Fragment from 'model-fragments/fragment';

export default Fragment.extend({
  value: attr('number'),
  timestamp: attr('date')
});
